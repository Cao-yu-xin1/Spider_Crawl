package price

import (
	"fmt"
	"github.com/Cao-yu-xin1/Spider_Crawl/pkg/url"
	"time"
)

func CrawlPrice(pageUrl string, selector string) {
	browser := url.CrawlUrl()
	defer browser.MustClose()

	page := browser.MustPage(pageUrl)

	page.MustWaitLoad()

	err := page.Timeout(10 * time.Minute).MustElement(selector).WaitVisible()
	if err != nil {
		println("等待商品价格可见超时", err.Error())
		return
	}

	time.Sleep(2 * time.Second)

	elements := page.MustElements(selector)

	// 打印所有商品标题
	fmt.Printf("找到 %d 个商品价格:\n", len(elements))

	for i, el := range elements {
		if el.MustVisible() {
			text := el.MustText()
			fmt.Printf("%d. %s\n", i+1, text)
		}
	}
}

//func main() {
//	CrawlPrice("https://air.1688.com/kapp/channel-fe/cps-4c-pc/gphg?type=0&offerIds=1007909658416,774041976238,1007285119301",
//		".offer-price")
//}
