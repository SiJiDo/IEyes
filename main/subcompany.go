package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kirinlabs/HttpRequest"
)

//获取子公司(集团内)
func GetChildCompany_jt(id string, auth_token string, srate int) ([]string, []string, []string) {

	fmt.Println("[+]开始获取子公司信息")

	page := 1
	var domainlist []string
	var applist []string
	var weixinlist []string

	split := "window.haveEquityPermissionById('"
	split1 := "','pc_businfo_invest_structure')"
	split2 := "','"
	split3 := "%"
	split4 := "<td class=\"\">"

	for page > 0 {
		url1 := "https://www.tianyancha.com/pagination/holdingCompany.xhtml?id=" + id + "&pn=" + strconv.Itoa(page)
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
			break
		}

		list := strings.Split(string(body), split)
		list = list[1:]
		for i := range list {
			c := strings.Split(list[i], split1)[0]
			childid := strings.Split(c, split2)[0]
			childcompany := strings.Split(c, split2)[1]

			rate := strings.Split(list[i], split3)[0]
			rate = strings.Split(rate, split4)[1]
			rate = strings.Split(rate, ".")[0]
			if rate[0] == '-' {
				rate = "-"
			} else {
				int, err := strconv.Atoi(rate)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if srate <= int {
					//fmt.Println(childid + " " + childcompany + " 比例:" + rate + "%")
					domainlist = append(domainlist, GetDomain(childid, childcompany, auth_token)...)
					applist = append(applist, Getapp(childid, childcompany, auth_token)...)
					weixinlist = append(weixinlist, Getweixin(childid, childcompany, auth_token)...)
				}

			}
		}
		page = page + 1
	}

	fmt.Println()
	return domainlist, applist, weixinlist
}

//获取子公司(大厂)
func GetChildCompany_dc(id string, auth_token string, srate int) ([]string, []string, []string) {

	fmt.Println("[+]开始获取子公司信息")

	page := 1
	var domainlist []string
	var applist []string
	var weixinlist []string

	split := "/td><td class=\"\"><a  class=\"link-click\""
	split1 := "</a></td></tr>"
	split2 := ">"
	split3 := "%"
	split4 := "<span class=\"\">"

	for page > 0 {
		url1 := "https://www.tianyancha.com/pagination/companyholding.xhtml?id=" + id + "&pn=" + strconv.Itoa(page)
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
			break
		}

		list := strings.Split(string(body), split)
		list = list[1:]
		for i := range list {
			c := strings.Split(list[i], split1)[0]
			childcompany := strings.Split(c, split2)[1]
			childid := strings.Split(c, "/company/")[1]
			childid = strings.Split(childid, "\"")[0]

			rate := strings.Split(list[i], split3)[0]
			rate = strings.Split(rate, split4)[1]
			rate = strings.Split(rate, ".")[0]
			if rate[0] == '-' {
				rate = "-"
			} else {
				int, err := strconv.Atoi(rate)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if srate <= int {
					//fmt.Println(childid + " " + childcompany + " 比例:" + rate + "%")
					domainlist = append(domainlist, GetDomain(childid, childcompany, auth_token)...)
					applist = append(applist, Getapp(childid, childcompany, auth_token)...)
					weixinlist = append(weixinlist, Getweixin(childid, childcompany, auth_token)...)
				}

			}
		}
		page = page + 1
	}

	fmt.Println()
	return domainlist, applist, weixinlist
}

//获取子公司(国企)，需要递归
func GetChildCompany_gq(id string, auth_token string, srate int, deep int) ([]string, []string, []string) {

	fmt.Println("[+]开始获取子公司信息")

	var idlist []string
	var companylist []string
	var domainlist []string
	var applist []string
	var weixinlist []string
	page := 1

	split := "window.haveEquityPermissionById('"
	split1 := "','"
	split3 := "%</td><td>"
	split4 := "</span></td><td>"

	for page > 0 {
		url1 := "https://www.tianyancha.com/pagination/investV2.xhtml?id=" + id + "&pn=" + strconv.Itoa(page)
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
			break
		}

		if strings.Contains(string(body), "th width=\"162px\">") == true {
			break
		}

		list := strings.Split(string(body), split)
		list = list[1:]
		for i := range list {
			childid := strings.Split(list[i], split1)[0]
			childcompany := strings.Split(list[i], split1)[1]
			childcompany = strings.Split(childcompany, split1)[0]

			rate := strings.Split(list[i], split3)[0]
			rate = strings.Split(rate, split4)[1]
			rate = strings.Split(rate, ".")[0]
			if rate[0] == '-' {
				rate = "-"
			} else {
				int, err := strconv.Atoi(rate)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if srate <= int {
					//fmt.Println(childid + " " + childcompany + " 比例:" + rate + "%")
					domainlist = append(domainlist, GetDomain(childid, childcompany, auth_token)...)
					applist = append(applist, Getapp(childid, childcompany, auth_token)...)
					weixinlist = append(weixinlist, Getweixin(childid, childcompany, auth_token)...)

					idlist = append(idlist, childid)
					companylist = append(companylist, childcompany)
				}

			}
		}
		page = page + 1
	}

	if len(idlist) == 0 {
		return domainlist, applist, weixinlist
	} else {
		if deep <= 1 {
			return domainlist, applist, weixinlist
		} else {
			deep = deep - 1
			for i := range idlist {
				domainlisttmp, applisttmp, weixinlisttmp := GetChildCompany_gq(idlist[i], auth_token, srate, deep)
				domainlist = append(domainlist, domainlisttmp...)
				applist = append(applist, applisttmp...)
				weixinlist = append(weixinlist, weixinlisttmp...)
			}
		}
	}

	fmt.Println()
	return domainlist, applist, weixinlist
}

//获取子公司(银行类支行)
func GetChildCompany_yh(id string, auth_token string, srate int) ([]string, []string, []string) {

	fmt.Println("[+]开始获取子公司信息")

	var domainlist []string
	var applist []string
	var weixinlist []string
	page := 1

	split := "/company/"
	split1 := "</a></td></tr></table>"
	split2 := "\""

	for page > 0 {
		url1 := "https://www.tianyancha.com/pagination/branch.xhtml?id=" + id + "&pn=" + strconv.Itoa(page)
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
			break
		}

		list := strings.Split(string(body), split)
		list = list[1:]
		for i := range list {

			childcompany := strings.Split(list[i], split1)[0]
			childcompany = strings.Split(childcompany, ">")[1]
			childid := strings.Split(list[i], split2)[0]

			fmt.Println(childid + " " + childcompany)
			domainlist = append(domainlist, GetDomain(childid, childcompany, auth_token)...)
			applist = append(applist, Getapp(childid, childcompany, auth_token)...)
			weixinlist = append(weixinlist, Getweixin(childid, childcompany, auth_token)...)

		}
		page = page + 1
	}

	fmt.Println()
	return domainlist, applist, weixinlist
}
