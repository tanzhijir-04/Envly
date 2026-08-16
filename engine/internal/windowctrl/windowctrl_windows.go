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
	procGetWindowLong = user32.NewProc("GetWindowLongW")
	procSetWindowLong = user32.NewProc("SetWindowLongW")
	procSetWindowPos  = user32.NewProc("SetWindowPos")
)

const (
	wmSysCommand = 0x0112
	scMinimize   = 0xF020
	scMaximize   = 0xF030
	scRestore    = 0xF120
	scClose      = 0xF060

	wsCaption     = 0x00C00000
	wsThickFrame  = 0x00040000
	wsSysMenu     = 0x00080000
	wsMinimizeBox = 0x00020000
	wsMaximizeBox = 0x00010000
	swpFramechanged = 0x0020
	swpNoMove     = 0x0002
	swpNoSize     = 0x0001
	swpNoZOrder   = 0x0004
)

var gwlStyleValue = int(-16)

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
	case "frameless":
		return stripCaption(hwnd)
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

func findWindow() (uintptr, error) {
	title, err := syscall.UTF16PtrFromString("Envly")
	if err != nil {
		return 0, err
	}
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(title)))
	if hwnd == 0 {
		return 0, fmt.Errorf("Envly window not found")
	}
	return hwnd, nil
}

// EnsureFrameless 移除窗口标题栏样式（保留可调整大小的边框）。
func EnsureFrameless() error {
	hwnd, err := findWindow()
	if err != nil {
		return err
	}
	return stripCaption(hwnd)
}

func stripCaption(hwnd uintptr) error {
	style, _, _ := procGetWindowLong.Call(hwnd, uintptr(gwlStyleValue))
	if style == 0 {
		return fmt.Errorf("failed to read window style")
	}
	// 清除所有非客户区样式位：标题栏、边框、系统菜单、最小化/最大化按钮
	nonClient := uintptr(wsCaption | wsThickFrame | wsSysMenu | wsMinimizeBox | wsMaximizeBox)
	newStyle := style &^ nonClient
	if newStyle == style {
		return nil // 已经是无边框
	}
	procSetWindowLong.Call(hwnd, uintptr(gwlStyleValue), newStyle)
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, uintptr(swpFramechanged|swpNoMove|swpNoSize|swpNoZOrder))
	return nil
}
