package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	TOKEN    = "finder"
	Text     = "text"
	Location = "location"
	Image    = "image"
	Link     = "link"
	Event    = "event"
	Music    = "music"
	News     = "news"
)

type msgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
}

type Request struct {
	XMLName                xml.Name `xml:"xml"`
	msgBase                         // base struct
	Location_X, Location_Y float32
	Scale                  int
	Label                  string
	PicUrl                 string
	MsgId                  int64
	Event                  string
}

type Response struct {
	XMLName xml.Name `xml:"xml"`
	msgBase
	ArticleCount int     `xml:",omitempty"`
	Articles     []*item `xml:"Articles>item,omitempty"`
	FuncFlag     int
}

type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

func weixinEvent(w http.ResponseWriter, r *http.Request) {
	if weixinCheckSignature(w, r) == false {
		fmt.Println("auth failed, attached?")
		return
	}

	fmt.Println("auth success, parse POST")

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(body))
	var wreq *Request
	if wreq, err = DecodeRequest(body); err != nil {
		log.Fatal(err)
		return
	}

	wresp, err := dealwith(wreq)
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := wresp.Encode()
	if err != nil {
		fmt.Printf("error:%v\n", err)
		return
	}

	fmt.Println(string(data))
	fmt.Fprintf(w, string(data))
	return
}

func dealwith(req *Request) (resp *Response, err error) {
	resp = NewResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Text
	fmt.Println(req.MsgType, len(req.MsgType), Event, len(Event), req.Event, len(req.Event), "subscribe", len("subscribe"))
	if req.MsgType == Event {
		if req.Event == "subscribe" {
			resp.Content = fmt.Sprintln("欢迎关注微信订阅号！")
			resp.Content += help()
			return resp, nil
		} else {
			resp.Content = fmt.Sprintln("暂时不支持事件！")
			resp.Content += help()
			return resp, nil
		}
	}
	if req.MsgType == Text {

		//帮助
		if strings.Trim(strings.ToLower(req.Content), " ") == "help" || strings.Trim(strings.ToLower(req.Content), " ") == "?" {
			resp.Content = help()
			return resp, nil
		}

		//获取金价
		if strings.Trim(strings.ToLower(req.Content), " ") == "gold" {
			resp.Content = GoldPrice()
			return resp, nil
		}

		//是否查询股票代码
		msg := CheckNO(req.Content)
		if len(msg) > 0 {
			resp.Content = msg
			return resp, nil
		}

		resp.Content = "亲，已经收到您的消息, 将尽快回复您."
	} else if req.MsgType == Image {
		var a item
		a.Description = "阿Y。。。^_^^_^1024你懂的"
		a.Title = "阿Y图文测试"
		a.PicUrl = "http://p1.qhimg.com/d/inn/84310112/haosou.png"
		a.Url = "http://www.baidu.com"

		resp.MsgType = News
		resp.ArticleCount = 1
		resp.Articles = append(resp.Articles, &a)
		resp.FuncFlag = 1
	} else {
		resp.Content = fmt.Sprintln("暂时还不支持其他的类型！")
		resp.Content += help()
	}
	return resp, nil
}

func weixinAuth(w http.ResponseWriter, r *http.Request) {
	if weixinCheckSignature(w, r) == true {
		var echostr string = strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(w, echostr)
	}
}

func weixinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("GET begin...")
		weixinAuth(w, r)
		fmt.Println("GET END...")
	} else {
		fmt.Println("POST begin...")
		weixinEvent(w, r)
		fmt.Println("POST END...")
	}
}

func main() {
	http.HandleFunc("/check", weixinHandler)
	port := "80"
	println("Listening on port ", port, "...")
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func weixinCheckSignature(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	fmt.Println(r.Form)
	var signature string = strings.Join(r.Form["signature"], "")
	var timestamp string = strings.Join(r.Form["timestamp"], "")
	var nonce string = strings.Join(r.Form["nonce"], "")
	tmps := []string{TOKEN, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2]
	tmp := str2sha1(tmpStr)
	if tmp == signature {
		return true
	}
	return false
}

func DecodeRequest(data []byte) (req *Request, err error) {
	req = &Request{}
	if err = xml.Unmarshal(data, req); err != nil {
		return
	}
	req.CreateTime *= time.Second
	return
}

func NewResponse() (resp *Response) {
	resp = &Response{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func (resp Response) Encode() (data []byte, err error) {
	resp.CreateTime = time.Duration(time.Now().Unix())
	data, err = xml.Marshal(resp)
	return
}

func help() string {
	msg := ""
	msg += fmt.Sprintln("0.输入\"?\"或者\"help\" 显示帮助")
	msg += fmt.Sprintln("1.输入股票代码(如\"600360\")显示当前价")
	msg += fmt.Sprintln("2.输入gold(如\"gold\")显示当前金价")
	return msg
}
