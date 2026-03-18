package alipay

import (
	"fmt"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/nacos"
	"strconv"
)
import "github.com/smartwalle/alipay/v3"

func AliPay(orderNo string, total float64) string {
	ali := nacos.GlobalConfig.AliPay
	var privateKey = ali.PrivateKey // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var appId = ali.AppId           // 必须，支付宝分配的应用ID
	client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = ali.NotifyURL // 必须，支付成功后，支付宝服务器异步通知商户服务器的地址
	p.ReturnURL = ali.ReturnURL // 必须，支付成功后，用户在商户网站返回的地址
	p.Subject = "标题"
	p.OutTradeNo = orderNo
	p.TotalAmount = strconv.FormatFloat(total, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	fmt.Println(payURL)
	return payURL
}
