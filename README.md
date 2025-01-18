<div align="center">
  <img src="docs/logo.png" alt="Live Helper" width="200">
  <h1>live Helper</h1>
  <p>使用Go&wails&vue-js开发的跨平台摄像头桌面工具</p>
</div>


## Preview

<div align="center">
  <img src="docs/menu.png">
  <img src="docs/rect.png">
  <img src="docs/circle.png">
  <img src="docs/vague.png">
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
- [x] 无边框窗口
- [x] 实时捕获摄像头画面
- [x] 圆形/矩形画面切换
- [x] 灰度/模糊/复古滤镜
- [x] 镜像开关

## use
Ctrl+h:显示/关闭菜单

## Stargazers over time
## Stargazers over time
[![Stargazers over time](https://starchart.cc/ppywj/wails-liveHelper.svg?variant=light)](https://starchart.cc/ppywj/wails-liveHelper)