//go:build windows

// Package windowctrl 通过 Win32 消息控制 Envly 应用窗口（最小化/最大化/关闭/无边框）。
package windowctrl

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procEnumWindows                = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId   = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowTextW             = user32.NewProc("GetWindowTextW")
	procPostMessageW               = user32.NewProc("PostMessageW")
	procIsZoomed                   = user32.NewProc("IsZoomed")
	procGetWindowLong              = user32.NewProc("GetWindowLongW")
	procSetWindowLong              = user32.NewProc("SetWindowLongW")
	procSetWindowPos               = user32.NewProc("SetWindowPos")
	procOpenProcess                = kernel32.NewProc("OpenProcess")
	procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandle                = kernel32.NewProc("CloseHandle")
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
	swpNoMove       = 0x0002
	swpNoSize       = 0x0001
	swpNoZOrder     = 0x0004

	processQueryLimitedInformation = 0x1000
)

var gwlStyleValue = int(-16)

// Controller 控制 Envly 应用窗口。
type Controller struct{}

func (Controller) Action(action string) error {
	hwnd, err := findWindow()
	if err != nil {
		return err
	}
	switch action {
	case "frameless":
		return stripNonClient(hwnd)
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

// EnsureFrameless 移除窗口标题栏等非客户区样式。
func EnsureFrameless() error {
	hwnd, err := findWindow()
	if err != nil {
		return err
	}
	return stripNonClient(hwnd)
}

// findWindow 枚举顶层窗口，返回标题为 Envly 且属于 Envly/pake-envly 进程的窗口。
// 按进程名过滤可跳过任务栏缩略图（由 explorer.exe 托管的同名窗口）。
func findWindow() (uintptr, error) {
	var found uintptr
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		var buf [512]uint16
		n, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if n == 0 || syscall.UTF16ToString(buf[:n]) != "Envly" {
			return 1 // continue
		}
		var pid uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
		if pid == 0 {
			return 1
		}
		hProcess, _, _ := procOpenProcess.Call(processQueryLimitedInformation, 0, uintptr(pid))
		if hProcess == 0 {
			return 1
		}
		defer procCloseHandle.Call(hProcess)
		var pathBuf [1024]uint16
		size := uint32(len(pathBuf))
		ret, _, _ := procQueryFullProcessImageNameW.Call(hProcess, 0, uintptr(unsafe.Pointer(&pathBuf[0])), uintptr(unsafe.Pointer(&size)))
		if ret == 0 {
			return 1
		}
		name := filepath.Base(syscall.UTF16ToString(pathBuf[:size]))
		if name == "Envly.exe" || name == "pake-envly.exe" {
			found = hwnd
			return 0 // stop
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	if found == 0 {
		return 0, fmt.Errorf("Envly window not found")
	}
	return found, nil
}

// stripNonClient 清除标题栏、边框、系统菜单、最小化/最大化按钮等非客户区样式位。
func stripNonClient(hwnd uintptr) error {
	style, _, _ := procGetWindowLong.Call(hwnd, uintptr(gwlStyleValue))
	if style == 0 {
		return fmt.Errorf("failed to read window style")
	}
	nonClient := uintptr(wsCaption | wsThickFrame | wsSysMenu | wsMinimizeBox | wsMaximizeBox)
	newStyle := style &^ nonClient
	if newStyle == style {
		return nil // 已经是无边框
	}
	procSetWindowLong.Call(hwnd, uintptr(gwlStyleValue), newStyle)
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, uintptr(swpFramechanged|swpNoMove|swpNoSize|swpNoZOrder))
	return nil
}
