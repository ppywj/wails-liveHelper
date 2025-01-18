package systemtray

import (
	"context"
	_ "embed"
	"os/exec"
	"runtime"

	"github.com/energye/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed LH.ico
var icon []byte

type SystemTray struct {
	ctx context.Context
}

func NewSystemTray() *SystemTray {
	return &SystemTray{}
}

func openUrl(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // Linux 等其他 Unix-like 系统
		cmd = "xdg-open"
	}

	args = append(args, url)
	exec.Command(cmd, args...).Start()

}
func (a *SystemTray) systemTray() {
	systray.SetIcon(icon)
	aboutItem := systray.AddMenuItem("关于liveHelper", "关于liveHelper")
	exitItem := systray.AddMenuItem("退出", "退出程序")
	exitItem.Click(func() { wailsRuntime.Quit(a.ctx) })
	aboutItem.Click(func() { openUrl("https://github.com/ppywj/wails-liveHelper") })
	systray.SetOnClick(func(menu systray.IMenu) { wailsRuntime.WindowShow(a.ctx) })
	systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })
	systray.SetTooltip("liveHelper")
}

func (a *SystemTray) Startup(ctx context.Context) {
	a.ctx = ctx
	systray.Run(a.systemTray, func() {})
}
