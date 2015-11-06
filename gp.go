package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
	"strings"
)

func CheckNO(msg string) string {
	if len(msg) == 6 && (strings.HasPrefix(msg, "60") || strings.HasPrefix(msg, "51") || strings.HasPrefix(msg, "00") || strings.HasPrefix(msg, "30")) {
		return GuPiao(msg)
	}
	return ""
}

func GuPiao(no string) string {
	client := &http.Client{}
	url := ""
	if strings.HasPrefix(no, "60") || strings.HasPrefix(no, "51") {
		url = "http://hq.sinajs.cn/list=sh" + no
	} else if strings.HasPrefix(no, "00") || strings.HasPrefix(no, "30") {
		url = "http://hq.sinajs.cn/list=sz" + no
	}

	reqest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "代碼錯誤"
	}
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("Content-Type", "text/html; charset=utf-8")
	resp, err := client.Do(reqest)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "代碼錯誤"
	}
	dec := mahonia.NewDecoder("GBK")
	s := dec.ConvertString(string(data))
	s = s[strings.Index(s, "\"")+1 : strings.LastIndex(s, "\"")]
	s = GuPiaoFormat(s)
	fmt.Println(s)
	return s
}

func GuPiaoFormat(data string) string {
	FormatString := ""
	if len(data) > 1 {
		arr := strings.Split(data, ",")
		FormatString += fmt.Sprintln("股票名称", arr[0])
		FormatString += fmt.Sprintln("昨日收盘价", arr[2])
		FormatString += fmt.Sprintln("今日开盘价", arr[1])
		FormatString += fmt.Sprintln("当前价格", arr[3])
		FormatString += fmt.Sprintln("今日最高价", arr[4])
		FormatString += fmt.Sprintln("今日最低价", arr[5])
		FormatString += fmt.Sprintln("竞买价(“买一”报价)", arr[6])
		FormatString += fmt.Sprintln("竞卖价(“卖一”报价)", arr[7])
		FormatString += fmt.Sprintln("时间", arr[30], arr[31])
	}
	return FormatString
}
