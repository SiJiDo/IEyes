package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/kirinlabs/HttpRequest"
)

//获取当前页面的id和公司
func GetPage(auth_token string, name string, page int) ([]string, []string) {
	//开始获取页面的公司和id号
	fmt.Println("[+]开始获取页面的公司")

	split1 := "/span></div></div></div><div class=\"search-item sv-search-company  \""
	split2 := "&card_name="
	split3 := "&card_type=公司&card_id="
	spilit_id1 := "&card_type=公司&card_id="
	spilit_id2 := "&item=公司&"
	var idlist []string
	var infolist []string

	for p := 1; p <= page; p++ {

		url1 := "https://www.tianyancha.com/search/p" + strconv.Itoa(p) + "?key=" + url.QueryEscape(name)
		req := HttpRequest.NewRequest()
		req.SetCookies(map[string]string{
			"auth_token": auth_token,
		})

		resp, err := req.Get(url1)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		body, err := resp.Body()
		r := strings.Split(string(body), split1)
		flag := false

		for s := range r {
			var info string
			var id string
			if flag == false {
				info = strings.Split(r[s], "<div class=\"xcx-qrcode\" tyc-xcx-qrcode></div><div class=\"info\">")[1]
				info = strings.Split(info, "</div><div class=\"bottom\"><span>")[0]
				fmt.Println(info)
				flag = true
			} else {
				info = strings.Split(r[s], split2)[1]
				info = strings.Split(info, split3)[0]
				fmt.Println(info)
			}
			id = strings.Split(r[s], spilit_id1)[1]
			id = strings.Split(id, spilit_id2)[0]

			idlist = append(idlist, id)
			infolist = append(infolist, info)
		}
	}
	fmt.Println()
	return idlist, infolist
}

func GetFirstCompany(auth_token string, name string) ([]string, []string) {

	//开始获取页面的公司和id号
	fmt.Println("[+]开始准备获取当前公司名称")

	split1 := "/span></div></div></div><div class=\"search-item sv-search-company  \""
	split2 := "<div class=\"xcx-qrcode\" tyc-xcx-qrcode></div><div class=\"info\">"
	split3 := "</div><div class=\"bottom\"><span>"
	spilit_id1 := "&card_type=公司&card_id="
	spilit_id2 := "&item=公司&"

	url1 := "https://www.tianyancha.com/search?key=" + url.QueryEscape(name)
	req := HttpRequest.NewRequest()
	req.SetCookies(map[string]string{
		"auth_token": auth_token,
	})

	resp, err := req.Get(url1)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := resp.Body()

	info := strings.Split(string(body), split2)[1]
	info = strings.Split(info, split3)[0]

	id := strings.Split(string(body), split1)[0]
	id = strings.Split(string(id), spilit_id1)[1]
	id = strings.Split(string(id), spilit_id2)[0]
	fmt.Println(info)
	fmt.Println()

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
		req := HttpRequest.NewRequest()
		req.SetCookies(map[string]string{
			"auth_token": auth_token,
		})

		resp, err := req.Get(url1)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		body, err := resp.Body()

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
