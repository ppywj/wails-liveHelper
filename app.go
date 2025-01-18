package main

import (
	"context"
	"fmt"

	"github.com/martinlindhe/notify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{ctx: context.Background()}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// 系统通知

func (a *App) Notify(title string,
	message string,
	text string,
	icon string) {
	notify.Notify("Title", "Message", "123456", "")
}

func (a *App) Log(msg string) {
	runtime.LogInfo(a.ctx, msg)
}
