//go:build windows

// Package windowctrl 通过 Win32 消息控制 Envly 应用窗口（最小化/最大化/关闭）。
package windowctrl

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procFindWindowW  = user32.NewProc("FindWindowW")
	procPostMessageW = user32.NewProc("PostMessageW")
	procIsZoomed     = user32.NewProc("IsZoomed")
)

const (
	wmSysCommand = 0x0112
	scMinimize   = 0xF020
	scMaximize   = 0xF030
	scRestore    = 0xF120
	scClose      = 0xF060
)

// Controller 按窗口标题查找 Envly 窗口并发送系统命令。
type Controller struct{}

func (Controller) Action(action string) error {
	title, err := syscall.UTF16PtrFromString("Envly")
	if err != nil {
		return err
	}
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(title)))
	if hwnd == 0 {
		return fmt.Errorf("Envly window not found")
	}
	switch action {
	case "minimize":
		procPostMessageW.Call(hwnd, wmSysCommand, scMinimize, 0)
	case "maximize":
		zoomed, _, _ := procIsZoomed.Call(hwnd)
		cmd := uintptr(scMaximize)
		if zoomed != 0 {
			cmd = scRestore
		}
		procPostMessageW.Call(hwnd, wmSysCommand, cmd, 0)
	case "close":
		procPostMessageW.Call(hwnd, wmSysCommand, scClose, 0)
	default:
		return fmt.Errorf("unknown window action %q", action)
	}
	return nil
}
