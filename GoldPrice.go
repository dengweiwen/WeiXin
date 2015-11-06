package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func GoldPrice() string {
	strFormat := fmt.Sprintln("上海黄金交易所金价(仅供参考)")
	rep, err := http.Get("http://gold.hexun.com/hjxh/")
	if err != nil {
		return "数据获取失败"
	}
	defer rep.Body.Close()
	data, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "数据获取失败"
	}
	m := mahonia.NewDecoder("GBK")
	s := m.ConvertString(string(data))
	re, _ := regexp.Compile(`<p>上海黄金交易所行情</p>[\s\S]*<span>上海黄金交易所</span>`)
	s = re.FindString(s)
	re, _ = regexp.Compile(`<tbody[^> ]*>[\s\S]*</tbody>`)
	s = re.FindString(s)
	re, _ = regexp.Compile(`<th>[\s\S]*</th>`)
	s = re.ReplaceAllString(s, "")
	s = strings.Replace(s, "<tr class=\"r\">", "<tr>", -1)
	s = strings.Replace(s, "<td class=\"l\"><span>", "<td>", -1)
	s = strings.Replace(s, "</span></td>", "</td>", -1)
	s = strings.Replace(s, "<tbody>", "", -1)
	s = strings.Replace(s, "</tbody>", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "<tr>", "", -1)
	rows := strings.Split(s, "</tr>")
	for i := 0; i < len(rows); i++ {
		if len(rows[i]) > 50 {
			tr := strings.Replace(rows[i], "<td>", "", -1)
			cells := strings.Split(tr, "</td>")
			st := ShowPriceFormat(cells)
			if len(st) > 10 {
				strFormat += fmt.Sprint(st)
			}
		}
	}
	fmt.Println(strFormat)
	return strFormat
}

func ShowPriceFormat(cells []string) string {
	data := ""
	if cells[0] == "Au99.99" {
		data += fmt.Sprintln("黄金:", cells[1], cells[8])
	} else if cells[0] == "Pt99.95" {
		data += fmt.Sprintln("铂金:", cells[1], cells[8])
	}
	return data
}
