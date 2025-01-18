package main

import (
	"context"
	"embed"

	systemtray "wails-liveHelper/backend/systray"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	//创建系统托盘
	systp := systemtray.NewSystemTray()

	// Create application with options
	app := NewApp()
	//关闭任务栏不退出程序

	// Create application with options
	err := wails.Run(&options.App{
		Title:       "liveHelper", //窗口标题
		Width:       800,
		Height:      600,
		Frameless:   true, // 无边框模式
		AlwaysOnTop: true, //窗口置顶
		Windows: &windows.Options{
			WebviewIsTransparent:              true, // 网页透明
			WindowIsTranslucent:               true, // 窗口透明
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
		},

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Fullscreen: false,
		//BackgroundColour:  &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		HideWindowOnClose: true, //关闭时隐藏窗口而不退出
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			systp.Startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			runtime.LogInfo(ctx, "应用程序退出")
			return
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
