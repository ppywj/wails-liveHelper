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
		Width:       300,
		Height:      300,
		Frameless:   true, // 无边框模式
		AlwaysOnTop: true, //窗口置顶
		BackgroundColour: &options.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 0,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true, // 网页透明
			WindowIsTranslucent:               true, // 窗口透明
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: true, //去掉边框
			WebviewUserDataPath:               "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Fullscreen:        false,
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
