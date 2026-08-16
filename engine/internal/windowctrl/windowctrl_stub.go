//go:build !windows

package windowctrl

import "fmt"

// Controller 在非 Windows 平台返回不支持。
type Controller struct{}

func (Controller) Action(action string) error {
	return fmt.Errorf("window control not supported on this platform")
}
