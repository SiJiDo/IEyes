package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/kirinlabs/HttpRequest"
)

func GetPage(auth_token string, name string, page int) ([]string, []string) {
	var idlist []string
	var infolist []string

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[-]爬取异常,程序退出")
			os.Exit(0)
		}
	}()

	for p := 1; p <= page; p++ {

		url1 := "https://www.tianyancha.com/search/p" + strconv.Itoa(p) + "?key=" + url.QueryEscape(name)

		//初始化爬虫
		c := colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:101.0) Gecko/20100101 Firefox/101.0"))

		c.OnRequest(func(r *colly.Request) {
			// Request头部设定
			r.Headers.Set("Connection", "keep-alive")
			r.Headers.Set("Accept", "*/*")
			r.Headers.Set("Cookie", "auth_token="+auth_token)
			r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
			r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")

			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("a", func(e *colly.HTMLElement) {
			if e.Attr("class") == "index_alink__zcia5 link-click" {

				info := string(e.Text)
				link := strings.Split(e.Attr("href"), "/")
				id := link[len(link)-1]
				idlist = append(idlist, id)
				infolist = append(infolist, info)
			}
		})

		c.Visit(url1)
	}

	return idlist, infolist
}

func GetFirstCompany(auth_token string, name string) ([]string, []string) {

	//开始获取页面的公司和id号
	url1 := "https://www.tianyancha.com/search?key=" + url.QueryEscape(name)
	id := ""
	info := ""

	fmt.Println("[+]开始准备获取当前公司名称")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[-]爬取异常,程序退出")
			os.Exit(0)
		}
	}()

	//初始化爬虫
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:101.0) Gecko/20100101 Firefox/101.0"))

	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Cookie", "auth_token="+auth_token)
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")

		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		if e.Attr("class") == "index_alink__zcia5 link-click" {
			if id == "" && info == "" {
				info = string(e.Text)
				link := strings.Split(e.Attr("href"), "/")
				id = link[len(link)-1]
			}
		}
	})

	c.Visit(url1)
	//fmt.Println(id, info)

	return []string{id}, []string{info}

}

func GetDomain(id string, company string, auth_token string) []string {

	//开始获取备案域名
	fmt.Println("[+]开始获取备案信息:" + company)

	split_icp1 := "</td><td class=\"left-col\"><span>"
	split_icp2 := "-"
	split_domain := "</span></td><td class=\"left-col\"><span>"
	split_domain1 := "<td class=\"left-col\">"
	split_domain2 := "</td>"
	var domainlist []string

	page := 1
	flag := false
	var icp string
	// var site_result []string

	for page > 0 {
		url1 := "https://www.tianyancha.com/pagination/icp.xhtml?id=" + id + "&pn=" + strconv.Itoa(page)
		//fmt.Println(url1)
		req := HttpRequest.NewRequest()
		req.SetCookies(map[string]string{
			"auth_token": auth_token,
		})
		req.SetHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		})

		resp, err := req.Get(url1)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		body, err := resp.Body()
		//fmt.Println(string(body))

		if strings.Contains(string(body), "<div") == false {
			return domainlist
		}
		//进行domain和icp的提取
		icp = strings.Split(string(body), split_icp1)[2]
		icp = strings.Split(string(icp), split_icp2)[0]
		if flag == false {
			fmt.Println("备案号: " + icp)
			flag = true
		}

		domain := strings.Split(string(body), split_domain)
		l := len(domain)
		domain = domain[1:l]

		//输出备案信息
		for i := range domain {
			websiteinfo := strings.Split(domain[i], "</span>")[0]
			if strings.Contains(domain[i], split_domain1) == true {
				d := strings.Split(domain[i], split_domain1)[2]
				d = strings.Split(d, split_domain2)[0]
				domainlist = append(domainlist, d)
				fmt.Print(d + "\t")
				if websiteinfo != "-" {
					fmt.Println("(" + websiteinfo + ")")
				} else {
					fmt.Println()
				}
			}
		}

		page = page + 1
	}

	fmt.Println()
	return domainlist
}

func Getweixin(id string, company string, auth_token string) []string {

	//开始获取备案域名
	fmt.Println("[+]开始获微信公众号信息:" + company)

	split_name := "</div></td><td class=\"\"><span>"
	split_name2 := "</span></td></tr></table></td><td class=\"left-col\"><span>"
	split_name3 := "</span></td></tr></table></td><td class=\"left-col\"><span>"
	split_name4 := "</span></td><td>"
	var weixinlist []string

	page := 1
	// var site_result []string

	for page > 0 {
		url1 := "https://www.tianyancha.com/pagination/wechat.xhtml?id=" + id + "&pn=" + strconv.Itoa(page)
		req := HttpRequest.NewRequest()
		req.SetCookies(map[string]string{
			"auth_token": auth_token,
		})
		req.SetHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		})

		resp, err := req.Get(url1)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		body, err := resp.Body()

		if strings.Contains(string(body), split_name2) == false {
			return weixinlist
		}

		weixin := strings.Split(string(body), split_name)
		l := len(weixin)
		weixin = weixin[1:l]

		//输出备案信息
		for i := range weixin {
			weininfo := strings.Split(weixin[i], split_name2)[0]
			weixinid := strings.Split(weixin[i], split_name3)[1]
			weixinid = strings.Split(weixinid, split_name4)[0]

			fmt.Println(weininfo + "(微信号:" + weixinid + ")")
			weixinlist = append(weixinlist, weininfo+"(微信号:"+weixinid+")")
		}

		page = page + 1
	}

	fmt.Println()
	return weixinlist
}

func Getapp(id string, company string, auth_token string) []string {

	//开始获取备案域名
	fmt.Println("[+]开始获取app信息:" + company)

	split_name := "</span></td></tr></table></td><td class=\"left-col\"><span>"
	split_name2 := "</span></td><td><span>"
	var applist []string

	page := 1
	// var site_result []string

	for page > 0 {
		url1 := "https://www.tianyancha.com/pagination/product.xhtml?id=" + id + "&pn=" + strconv.Itoa(page)
		req := HttpRequest.NewRequest()
		req.SetCookies(map[string]string{
			"auth_token": auth_token,
		})
		req.SetHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		})

		resp, err := req.Get(url1)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		body, err := resp.Body()

		if strings.Contains(string(body), "<div") == false {
			return applist
		}

		app := strings.Split(string(body), split_name)
		l := len(app)
		app = app[1:l]

		//输出备案信息
		for i := range app {
			appinfo := strings.Split(app[i], split_name2)[0]
			applist = append(applist, appinfo)
			fmt.Println(appinfo)
		}

		page = page + 1
	}

	fmt.Println()
	return applist
}
