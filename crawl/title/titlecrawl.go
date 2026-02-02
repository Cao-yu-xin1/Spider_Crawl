package title

import (
	"Spider_Crawl/pkg/url"
	"fmt"
	"time"
)

func CrawlTitle(pageUrl string, selector string) {
	browser := url.CrawlUrl()
	defer browser.MustClose()
	page := browser.MustPage(pageUrl)

	page.MustWaitLoad()

	err := page.Timeout(10 * time.Minute).MustElement(selector).WaitVisible()
	if err != nil {
		println("等待商品标题可见超时", err.Error())
		return
	}
	//title := page.MustElement(".offer-title.ellipsis").MustText()
	//println("商品标题:", title)

	time.Sleep(2 * time.Second)

	elements := page.MustElements(selector)

	// 打印所有商品标题
	fmt.Printf("找到 %d 个商品标题:\n", len(elements))

	for i, el := range elements {
		if el.MustVisible() {
			text := el.MustText()
			fmt.Printf("%d. %s\n", i+1, text)
		}
	}
}

//func main() {
//	CrawlTitle("https://sale.1688.com/factory/u0vjcc4j.html?spm=a260k.home2025.centralDoor.ddoor.66333597BBbHgE&topOfferIds=1005591171200",
//		".offerTitle")
//}
