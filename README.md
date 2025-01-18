<div align="center">
  <img src="docs/big.png" alt="Live Helper" width="200">
  <h1>Live Helper</h1>
  <p>使用Go&wails&vue-ts开发的跨平台直播助手桌面工具</p>
</div>




## Preview

<div align="center">
  <!-- <img src="docs/imgs/light.png">
  <img src="docs/imgs/dark.png"> -->
</div>



## Build

1、Build Environment

- Node.js [link](https://nodejs.org/en)

- npm ：`npm i -g pnpm`

- Go [link](https://go.dev/)

- Wails [link](https://wails.io/) ：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

2、Pull and Build

```bash
git clone https://github.com/ppywj/liveHelper.git

cd liveHelper/frontend

npm install

cd ..

wails build
```

## Implemented components

- [x] 实时捕获摄像头画面
- [x] 人脸追踪
- [ ] 圆形/矩形画面
- [ ] 美颜
- [ ] 镜像开关
- [ ] 截图
- [ ] 录屏
- [ ] 支持快捷键
- [ ] 录音
- [ ] 语音转文字
- [ ] 抖音弹幕/礼物爬取
- [ ] 自动感谢礼物


## Stargazers over time
[![Stargazers over time](https://starchart.cc/ppywj/liveHelper.svg?variant=light)](https://starchart.cc/ppywj/liveHelper)