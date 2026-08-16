package installer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tanzhijir-04/Envly/engine/internal/config"
)

// installDownload 按 URL 列表依次尝试下载，成功后执行静默安装。
func (i *Installer) installDownload(ctx context.Context, spec config.InstallSpec) error {
	if len(spec.DownloadURLs) == 0 {
		return fmt.Errorf("no download URLs")
	}
	dir := filepath.Join(os.TempDir(), "envly-dl")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	dst := filepath.Join(dir, filepath.Base(spec.DownloadURLs[0]))
	var lastErr error
	for _, url := range spec.DownloadURLs {
		if err := downloadTo(ctx, url, dst); err != nil {
			lastErr = err
			continue
		}
		return i.runSilent(ctx, dst, spec.SilentArgs)
	}
	return fmt.Errorf("all downloads failed: %w", lastErr)
}

func downloadTo(ctx context.Context, url, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func (i *Installer) runSilent(ctx context.Context, installerPath, silentArgs string) error {
	args := parseArgs(silentArgs)
	full := append([]string{installerPath}, args...)
	_, err := i.run.Run(ctx, full[0], full[1:]...)
	return err
}

func parseArgs(s string) []string {
	var args []string
	var cur strings.Builder
	inQuote := false
	var quote byte
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '"' || ch == '\'':
			if inQuote && ch == quote {
				inQuote = false
			} else if !inQuote {
				inQuote = true
				quote = ch
			} else {
				cur.WriteByte(ch)
			}
		case ch == ' ' && !inQuote:
			if cur.Len() > 0 {
				args = append(args, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(ch)
		}
	}
	if cur.Len() > 0 {
		args = append(args, cur.String())
	}
	return args
}
