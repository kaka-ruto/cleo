package qatools

import (
	"fmt"
	"path/filepath"
	"strings"

	playwright "github.com/playwright-community/playwright-go"
)

type PlaywrightBrowser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	page    playwright.Page
}

func NewPlaywrightBrowser(headless bool, videoDir string) (*PlaywrightBrowser, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("start playwright runtime: %w", err)
	}
	launch := playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(headless)}
	browser, err := pw.Chromium.Launch(launch)
	if err != nil {
		_ = pw.Stop()
		return nil, fmt.Errorf("launch chromium: %w", err)
	}
	ctxOpts := playwright.BrowserNewContextOptions{}
	if strings.TrimSpace(videoDir) != "" {
		ctxOpts.RecordVideo = &playwright.RecordVideo{Dir: filepath.Clean(videoDir)}
	}
	ctx, err := browser.NewContext(ctxOpts)
	if err != nil {
		_ = browser.Close()
		_ = pw.Stop()
		return nil, fmt.Errorf("create browser context: %w", err)
	}
	page, err := ctx.NewPage()
	if err != nil {
		_ = ctx.Close()
		_ = browser.Close()
		_ = pw.Stop()
		return nil, fmt.Errorf("create page: %w", err)
	}
	return &PlaywrightBrowser{pw: pw, browser: browser, context: ctx, page: page}, nil
}

func (b *PlaywrightBrowser) OpenURL(target string) error {
	if _, err := b.page.Goto(strings.TrimSpace(target)); err != nil {
		return fmt.Errorf("goto %q: %w", target, err)
	}
	return nil
}

func (b *PlaywrightBrowser) TextContent(selector string) (string, error) {
	text, err := b.page.TextContent(strings.TrimSpace(selector))
	if err != nil {
		return "", fmt.Errorf("read text content for %q: %w", selector, err)
	}
	return strings.TrimSpace(text), nil
}

func (b *PlaywrightBrowser) Screenshot(path string, fullPage bool) error {
	options := playwright.PageScreenshotOptions{Path: playwright.String(filepath.Clean(path)), FullPage: playwright.Bool(fullPage)}
	if _, err := b.page.Screenshot(options); err != nil {
		return fmt.Errorf("capture screenshot %q: %w", path, err)
	}
	return nil
}

func (b *PlaywrightBrowser) Close() error {
	var first error
	if b.context != nil {
		if err := b.context.Close(); err != nil && first == nil {
			first = err
		}
	}
	if b.browser != nil {
		if err := b.browser.Close(); err != nil && first == nil {
			first = err
		}
	}
	if b.pw != nil {
		if err := b.pw.Stop(); err != nil && first == nil {
			first = err
		}
	}
	if first != nil {
		return fmt.Errorf("close playwright browser: %w", first)
	}
	return nil
}
