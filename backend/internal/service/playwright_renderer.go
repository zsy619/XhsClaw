package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

type PlaywrightRenderer struct {
	pw      *playwright.Playwright
	browser playwright.Browser
}

var pwRenderer *PlaywrightRenderer

func InitPlaywright() error {
	var err error

	homeDir, _ := os.UserHomeDir()
	browserInstallDir := filepath.Join(homeDir, ".playwright", "browsers")
	os.Setenv("PLAYWRIGHT_BROWSERS_PATH", browserInstallDir)

	log.Println("正在安装Playwright浏览器...")
	err = playwright.Install()
	if err != nil {
		log.Printf("Playwright浏览器安装失败: %v", err)
		return fmt.Errorf("安装Playwright浏览器失败: %v", err)
	}

	pwRenderer = &PlaywrightRenderer{}
	pwRenderer.pw, err = playwright.Run()
	if err != nil {
		return fmt.Errorf("启动Playwright失败: %v", err)
	}

	pwRenderer.browser, err = pwRenderer.pw.Chromium.Launch()
	if err != nil {
		return fmt.Errorf("启动Chromium失败: %v", err)
	}

	log.Println("Playwright初始化成功")
	return nil
}

func (s *PlaywrightRenderer) Close() error {
	if s.browser != nil {
		s.browser.Close()
	}
	if s.pw != nil {
		return s.pw.Stop()
	}
	return nil
}

func boolPtr(b bool) *bool {
	return &b
}

func (s *PlaywrightRenderer) renderHTMLToImageWithPlaywright(htmlContent, outputPrefix, suffix string, width, height int) (string, error) {
	maxRetries := 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err := s.doRenderWithPlaywright(htmlContent, outputPrefix, suffix, width, height)
		if err == nil {
			return result, nil
		}
		lastErr = err
		log.Printf("Playwright渲染尝试 %d/%d 失败: %v", attempt, maxRetries, err)

		if attempt < maxRetries {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return "", fmt.Errorf("Playwright渲染失败（已重试%d次）: %v", maxRetries, lastErr)
}

func (s *PlaywrightRenderer) doRenderWithPlaywright(htmlContent, outputPrefix, suffix string, width, height int) (string, error) {
	dpr := 2.0
	viewportWidth := int(float64(width) * dpr)
	viewportHeight := int(float64(height) * dpr)

	page, err := s.browser.NewPage()
	if err != nil {
		return "", fmt.Errorf("创建页面失败: %v", err)
	}
	defer page.Close()

	err = page.SetViewportSize(viewportWidth, viewportHeight)
	if err != nil {
		return "", fmt.Errorf("设置视口大小失败: %v", err)
	}

	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("playwright_render_%d.html", time.Now().UnixNano()))
	if err := os.WriteFile(tempFile, []byte(htmlContent), 0644); err != nil {
		return "", fmt.Errorf("创建临时HTML文件失败: %v", err)
	}
	defer os.Remove(tempFile)

	_, err = page.Goto(fmt.Sprintf("file://%s", tempFile))
	if err != nil {
		return "", fmt.Errorf("导航到HTML文件失败: %v", err)
	}

	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		return "", fmt.Errorf("等待页面加载失败: %v", err)
	}

	time.Sleep(1 * time.Second)

	screenshot, err := page.Screenshot(playwright.PageScreenshotOptions{
		Type:     playwright.ScreenshotTypePng,
		FullPage: boolPtr(true),
	})
	if err != nil {
		return "", fmt.Errorf("截图失败: %v", err)
	}

	if len(screenshot) < 1000 {
		return "", fmt.Errorf("截图数据无效，大小: %d bytes", len(screenshot))
	}

	filename := fmt.Sprintf("%s_%s_%d.png", outputPrefix, suffix, time.Now().UnixNano())
	imagesDir := "assets/images"

	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return "", fmt.Errorf("创建图片目录失败: %v", err)
	}

	fullPath := filepath.Join(imagesDir, filename)
	if err := os.WriteFile(fullPath, screenshot, 0644); err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	if info, err := os.Stat(fullPath); err != nil || info.Size() < 1000 {
		return "", fmt.Errorf("图片文件保存失败或文件过小")
	}

	log.Printf("图片渲染成功: %s", filename)

	return "/xhsclaw/image/" + filename, nil
}

func GetPlaywrightRenderer() *PlaywrightRenderer {
	return pwRenderer
}
