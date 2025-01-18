<template>
  <div id="app" class="mainWnd">
    <div class="video-container" v-if="isPlaying">
      <video ref="videoEl" autoplay class="videoStyle"></video>
      <canvas ref="overlayEl" class="overlay"></canvas>
    </div>
    <div v-if="showSettingsMenu" class="settings-menu">
      <el-switch v-model="isPlaying" active-text="开启摄像头" inactive-text="停止摄像头" @change="toggleVideo" />
      <el-switch v-model="isMirror" active-text="镜像开启" inactive-text="镜像关闭" @change="toggleMirror" />
      <label>
        亮度：
        <input type="range" v-model="brightnessLevel" min="0" max="2" step="0.1" @input="applyBeautyFilters">
      </label>
      <select v-model="selectedFilter" @change="applySelectedFilter">
        <option value="none">无滤镜</option>
        <option value="grayscale">灰度</option>
        <option value="sepia">复古</option>
        <option value="blur">模糊</option>
      </select>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { ElSwitch } from 'element-plus'; // 导入 element-plus 的 Switch 组件
import * as faceapi from "@vladmandic/face-api";
// 假设 Greet, Notify, Log 和 WindowSetPosition 是全局可用的函数或通过其他方式引入

const videoEl = ref(null);
const overlayEl = ref(null);
const isPlaying = ref(false);
const showSettingsMenu = ref(false);
const isMirror = ref(false); // 是否开启镜像
const brightnessLevel = ref(1); // 默认亮度为100%
const selectedFilter = ref('none'); // 默认无滤镜
let stream = null;
const offset = ref({ x: 0, y: 0 });
const isDragging = ref(false);
let intervalId = null;
let canvas;

const drawRoundedRect = (ctx, x, y, width, height, radius, color) => {
  console.log(x, y, width, height, radius, color);
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.lineTo(x + width - radius, y);
  ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
  ctx.lineTo(x + width, y + height - radius);
  ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
  ctx.lineTo(x + radius, y + height);
  ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
  ctx.lineTo(x, y + radius);
  ctx.quadraticCurveTo(x, y, x + radius, y);
  ctx.closePath();
  ctx.lineWidth = 2;
  ctx.strokeStyle = color;
  ctx.stroke();
};
// 加载模型
const loadModels = async () => {
  await faceapi.nets.tinyFaceDetector.loadFromUri('/model');
};

// 检测并框出人脸
// 检测并框出人脸
const detectFaces = async () => {
  if (!videoEl.value || !overlayEl.value) return;

  const displaySize = {
    width: videoEl.value.videoWidth,
    height: videoEl.value.videoHeight
  };

  if (displaySize.width === 0 || displaySize.height === 0) {
    console.warn("视频尺寸无效，跳过检测。");
    return;
  }

  // 确保模型已加载
  if (!faceapi.nets.tinyFaceDetector.params) {
    await loadModels();
  }

  intervalId = setInterval(async () => {
    try {
      // 获取视频元素在视口中的真实尺寸
      const videoRect = videoEl.value.getBoundingClientRect();
      const scale = {
        x: videoEl.value.videoWidth / videoRect.width,
        y: videoEl.value.videoHeight / videoRect.height
      };

      // 检测人脸
      const detections = await faceapi.detectAllFaces(videoEl.value, new faceapi.TinyFaceDetectorOptions());
      const resizedDetections = faceapi.resizeResults(detections, displaySize);

      // 设置 Canvas 的尺寸与视频相同
      overlayEl.value.width = displaySize.width;
      overlayEl.value.height = displaySize.height;

      const ctx = overlayEl.value.getContext('2d');
      if (ctx && Array.isArray(resizedDetections) && resizedDetections.length > 0) {
        ctx.clearRect(0, 0, overlayEl.value.width, overlayEl.value.height);

        // 绘制带圆角的人脸框，并根据视频的显示比例调整检测结果
        resizedDetections.forEach((detection, index) => {
          if (detection && detection.box) {
            const { x, y, width, height } = detection.box;
            const adjustedX = x * scale.x;
            const adjustedY = y * scale.y;
            const adjustedWidth = width * scale.x;
            const adjustedHeight = height * scale.y;

            drawRoundedRect(ctx, adjustedX, adjustedY, adjustedWidth, adjustedHeight, 5, 'rgba(255, 235, 59, 1)'); // 浅黄色边框
            console.log(`绘制矩形框: ${adjustedX}, ${adjustedY}, ${adjustedWidth}, ${adjustedHeight}`);
          } else {
            console.warn(`Detection at index ${index} is missing or malformed.`);
          }
        });
      }
    } catch (error) {
      console.error("人脸检测错误:", error);
    }
  }, 100); // 每1000毫秒更新一次检测
};
const toggleSettingsMenu = (event) => {
  if (event.ctrlKey && event.key.toLowerCase() === 'h') {
    event.preventDefault();
    showSettingsMenu.value = !showSettingsMenu.value;
  }
};

const startDrag = (event) => {
  if (event.button === 0) {
    offset.value.x = event.clientX;
    offset.value.y = event.clientY;
    isDragging.value = true;

    window.addEventListener('mousemove', drag);
    window.addEventListener('mouseup', stopDrag);
  }
};

const drag = (event) => {
  if (isDragging.value) {
    const dx = event.clientX - offset.value.x;
    const dy = event.clientY - offset.value.y;

    window.requestAnimationFrame(() => {
      WindowSetPosition(event.clientX - dx, event.clientY - dy);
    });

    offset.value.x = event.clientX;
    offset.value.y = event.clientY;
  }
};

const stopDrag = () => {
  isDragging.value = false;
  window.removeEventListener('mousemove', drag);
  window.removeEventListener('mouseup', stopDrag);
};

const toggleVideo = async () => {
  if (isPlaying.value) {
    try {
      console.log("开始摄像头");
      stream = await navigator.mediaDevices.getUserMedia({ video: true });
      if (videoEl.value) {
        videoEl.value.srcObject = stream;
        videoEl.value.play();
        if (!overlayEl.value) {
          canvas = document.createElement('canvas');
          canvas.className = 'overlay';
          document.querySelector('.video-container').appendChild(canvas);
          overlayEl.value = canvas;
        }
        // 确保 Canvas 的尺寸与视频一致
        const updateCanvasSize = () => {
          if (overlayEl.value) {
            overlayEl.value.width = videoEl.value.videoWidth;
            overlayEl.value.height = videoEl.value.videoHeight;
          }
        };
        applyBeautyFilters();
        isPlaying.value = true;
        videoEl.value.onloadedmetadata = () => {
          updateCanvasSize();
          detectFaces(); // Call detectFaces when video metadata is loaded
        };
      }
    } catch (error) {
      console.error("访问摄像头时出错:", error);
      console.log(`访问摄像头时出错: ${error}`);
      isPlaying.value = false;
    }
  } else {
    console.log("停止摄像头");
    clearInterval(intervalId);
    if (stream) {
      stream.getTracks().forEach(track => track.stop());
      stream = null;
    }
    isPlaying.value = false;
  }
};

const toggleMirror = () => {
  if (videoEl.value) {
    videoEl.value.style.transform = isMirror.value ? 'scaleX(-1)' : '';
  }
};

const applySelectedFilter = () => {
  if (videoEl.value) {
    let filterString = '';
    switch (selectedFilter.value) {
      case 'grayscale':
        filterString = `grayscale(100%) brightness(${brightnessLevel.value})`;
        break;
      case 'sepia':
        filterString = `sepia(100%) brightness(${brightnessLevel.value})`;
        break;
      case 'blur':
        filterString = `blur(5px) brightness(${brightnessLevel.value})`;
        break;
      default:
        filterString = `brightness(${brightnessLevel.value})`;
        break;
    }
    videoEl.value.style.filter = filterString;
    videoEl.value.style.transform = isMirror.value ? 'scaleX(-1)' : '';
  }
};

const applyBeautyFilters = () => {
  applySelectedFilter();
};

onMounted(() => {
  window.addEventListener('keydown', toggleSettingsMenu);
});

onBeforeUnmount(() => {
  clearInterval(intervalId);
  window.removeEventListener('keydown', toggleSettingsMenu);
});
</script>


<style scoped>
.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none; /* 确保点击事件不会被 Canvas 阻挡 */
}

.mainWnd {
  background-color: transparent;
}

.draggable-window {
  background-color: transparent;
  /* 默认背景透明 */
}

.video-container {
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: transparent;
  /* 设置背景为透明 */
}

.settings-menu {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  flex-direction: column;
  background-color: rgba(255, 255, 255, 0.8);
  /* 半透明背景 */
  padding: 10px;
  border-radius: 10px;
  /* 圆角矩形 */
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.videoStyle {

  /* 圆形 */
 height: 100%;
 Width: 100%;
  object-fit: cover;
}
</style>