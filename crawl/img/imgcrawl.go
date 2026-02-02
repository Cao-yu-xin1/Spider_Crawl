package img

import (
	"Spider_Crawl/pkg/url"
	"fmt"
	"time"
)

func CrawlImg(pageUrl string, selector string) {
	browser := url.CrawlUrl()
	defer browser.MustClose()
	page := browser.MustPage(pageUrl)

	page.MustWaitLoad()

	err := page.Timeout(10 * time.Minute).MustElement(selector).WaitVisible()
	if err != nil {
		println("等待商品图片可见超时", err.Error())
		return
	}

	time.Sleep(2 * time.Second)
	/*
		elements := page.MustElements(selector)

		// 打印所有商品标题
		fmt.Printf("找到 %d 个商品图片:\n", len(elements))

		for i, el := range elements {
			if el.MustVisible() {
				src, err := el.Attribute("src")
				if err == nil && src != nil {
					fmt.Printf("%d. https://%s\n", i+1, *src)
				}
			}
		}
	// 抓取所有商品图片元素（包含滚动后动态加载的）*/

	elements := page.MustElements(selector)

	fmt.Printf("找到 %d 个商品图片:\n", len(elements))

	/*// 遍历并提取图片src
	for i, el := range elements {
		if el.MustVisible() {
			src, err := el.Attribute("src")
			if err != nil && src != nil {
				fmt.Printf("%d. %s\n", i+1, *src)
			}
		}
	}*/
	for i, el := range elements {
		// 确保元素可见
		if el.MustVisible() {
			src, err := el.Attribute("src")
			if err == nil && src != nil {
				fmt.Printf("%d. https://%s\n", i+1, *src)
			}
		}
	}
}

//func main() {
//	CrawlImg("https://sale.1688.com/factory/u0vjcc4j.html?spm=a260k.home2025.centralDoor.ddoor.66333597BBbHgE&topOfferIds=1005591171200",
//		".offerImg")
//}
