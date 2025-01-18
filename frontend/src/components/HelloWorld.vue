<template>
  <div id="app" :class="mainWnd">
    <div v-if="!isPlaying">
      <h1>未开启摄像头</h1>
    </div>
    <div class="video-container" v-if="isPlaying">
      <video ref="videoEl" autoplay class="videoStyle"></video>
      <canvas ref="overlayEl" class="overlay"></canvas>
    </div>
    <div v-if="showSettingsMenu" class="settings-menu">
      <el-switch v-model="isCircle" active-text="圆形边框" inactive-text="矩形边框" @change="toggleShape" />
      <el-switch v-model="isPlaying" active-text="开启摄像头" inactive-text="停止摄像头" @change="toggleVideo" />
      <el-switch v-model="notMirror" active-text="镜像关闭" inactive-text="镜像开启" @change="toggleMirror" />
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
import videojs from 'video.js';
// 假设 Greet, Notify, Log 和 WindowSetPosition 是全局可用的函数或通过其他方式引入
const borderStyle = ref('5px'); // 默认边框样式
const videoEl = ref(null);
const overlayEl = ref(null);
const isPlaying = ref(false);
const isCircle = ref(false);
const showSettingsMenu = ref(false);
const notMirror = ref(false); // 是否开启镜像
const brightnessLevel = ref(1); // 默认亮度为100%
const selectedFilter = ref('none'); // 默认无滤镜
let stream = null;
const offset = ref({ x: 0, y: 0 });
const isDragging = ref(false);
let intervalId = null;
let canvas;
//动态样式对象
const mainWnd = ref({
  'rectWnd': true,
  'circleWnd': false
})



const toggleSettingsMenu = (event) => {
  if (event.ctrlKey && event.key.toLowerCase() === 'h') {
    event.preventDefault();
    showSettingsMenu.value = !showSettingsMenu.value;
  }
};

//设置边框样式
const toggleShape = () => {
  if (isCircle.value) {
    mainWnd.value = {
      "rectWnd":false,
      "circleWnd": true
    };
    console.log('当前形状: 圆形', mainWnd.value);
  } else {
    mainWnd.value = {
      "rectWnd": true,
      "circleWnd": false
    };
    console.log('当前形状: 矩形', mainWnd.value);
  }
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
    videoEl.value.style.transform = notMirror.value ? 'scaleX(-1)' : '';
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
    videoEl.value.style.transform = notMirror.value ? 'scaleX(-1)' : '';
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
  pointer-events: none;
  /* 确保点击事件不会被 Canvas 阻挡 */
  height: 100%; 
  width: 100%;  
  overflow: hidden;
}

.rectWnd {
  /* 默认背景透明，允许鼠标拖动*/
  background-color: transparent;
  --wails-draggable: drag;
  /* 设置背景为透明 */
  overflow: hidden;
  /* 确保内容不会超出矩形边界 */
  border-radius: 5px;
  height: 100%; 
  width: 100%;  
  min-width: 100px;
  min-height: 100px;
}

.circleWnd {
  background-color: transparent;
  --wails-draggable: drag;
  /* 设置背景为透明 */
  overflow: hidden;
  /* 确保内容不会超出矩形边界 */
  border-radius: 50%;
  height: 100%; 
  width: 100%;  
  min-width: 100px;
  min-height: 100px;
}




.video-container {
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: transparent;
  height: 100%; 
  width: 100%;  
  overflow: hidden;
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
  object-fit: cover;
  height: 100%; 
  width: 100%;  
}

</style>