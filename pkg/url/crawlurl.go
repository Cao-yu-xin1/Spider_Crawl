package url

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func CrawlUrl() *rod.Browser {
	url := launcher.New().
		Headless(false).
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	return browser
}
