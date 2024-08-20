package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"gopkg.in/cas.v2"
)

var casURL = "http://sso.xxu.edu.cn"
var tPort = ":5002"
var corpid = ""
var appscret = ""

type templateBinding struct {
	Username   string
	Attributes cas.UserAttributes
}

func main() {
	url, _ := url.Parse(casURL)
	client := cas.NewClient(&cas.Options{URL: url})

	root := chi.NewRouter()
	root.Use(client.Handler)

	server := &http.Server{
		Addr:    tPort,
		Handler: client.Handle(root),
	}

	//路由处理

	root.HandleFunc("/home", homeHandleFunc)
	//root.HandleFunc("/apply", applyHandleFunc)
	//root.HandleFunc("/dlogin", dloginHandleFunc)
	root.HandleFunc("/", rootHandleFunc)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func rootHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.New("index.html").Parse(home_html)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}
	binding := &templateBinding{
		Username:   cas.Username(r),
		Attributes: cas.Attributes(r),
	}
	html := new(bytes.Buffer)
	if err := tmpl.Execute(html, binding); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}
	html.WriteTo(w)
}
func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.New("home.html").Parse(main_html)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	binding := &templateBinding{
		Username:   cas.Username(r),
		Attributes: cas.Attributes(r),
	}

	html := new(bytes.Buffer)
	if err := tmpl.Execute(html, binding); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	html.WriteTo(w)
}

const index_html = `<script>
document.location = "/home/"
</script>
`
const home_html = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>导航</title>
<style>
  .login-a {
    position: absolute;
    top: 10px;
    right: 10px;
  }
html {
  width: 560px;
  height: 600px;
  margin: auto;
}
  .grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    grid-gap: 10px;
  }
  .item {
    background-color: #f0f0f0;
    padding: 20px;
    text-align: center;
	width: 60px;
  height: 40px;
  }
  a {
text-decoration: none;
}
</style>
</head>
<body>
<div><h1>智慧校园临时页面</h1></div>
<div><h4>本页面仅作故障时临时页面使用，以保证其他业务的正常运行。</h4></div>
<div class="grid">
  <!-- 循环创建16个方格，‌每个方格内有一个文字的导航链接 -->

 <a href="http://wb.xxu.edu.cn/relax/sso/cas/login" class="link"> <div class="item">网上服务大厅</div></a>
 <a href="http://oa.xxu.edu.cn/pcas/" class="link"> <div class="item">OA系统(校内)</div></a>
 <a href="http://125.229.55.10:81/caslogin.aspx" class="link"> <div class="item">财务<br>系统</div></a>
 <a href="http://125.229.55.10:81/caslogin_cx.aspx" class="link"> <div class="item">财务<br>查询</div></a>
 <a href="http://kyxt.system.xxu.edu.cn/loginsso.aspx" class="link"> <div class="item">科研<br>系统</div></a>
 <a href="https://xxu.cwkeji.cn/ermsLogin/ssologin.do" class="link"> <div class="item">大型仪器共享</div></a>
 <a href="http://sjpt.system.xxu.edu.cn/zxcas/" class="link"> <div class="item">毕业设计系统</div></a>
 <a href="https://xxu.cwkeji.cn/ermsLogin/ssologin.do" class="link"> <div class="item">图书文献校外访问</div></a>
 <a href="http://25.system.xxu.edu.cn:38025/sso/" class="link"> <div class="item">教务系统（不可用）</div></a>
 <a href="#" class="link"> <div class="item">空白</div></a>


</div>

<span class="login-a">Welcome {{.Username}} | <a href="/logout">退出登录</a></span>
</body>
</html>
`
const main_html = `<!DOCTYPE html>
<html>
  <head>
    <title>Welcome {{.Username}}</title>
  </head>
  <body>
    <h1>Welcome {{.Username}} <a href="/logout">Logout</a></h1>
    <p>Your attributes are:</p>
    <ul>ul {{range $key, $values := .Attributes}}
      <li>li {{$len := len $values}}{{$key}}:{{if gt $len 1}}
        <ul>ul {{range $values}}
          <li>li {{.}}</li>{{end}}
        </ul>
      {{else}} {{index $values 0}}{{end}}</li>{{end}}
    </ul>
  </body>
</html>
`
const apply_html = `<!DOCTYPE html>
<html>
  <head>
    <title>Welcome {{.Username}}</title>
  </head>
  <body>
    <h1>Welcome {{.Username}} <a href="/logout">Logout</a></h1>
    <p>Your attributes are:</p>
    <ul>ul {{range $key, $values := .Attributes}}
      <li>li {{$len := len $values}}{{$key}}:{{if gt $len 1}}
        <ul>ul {{range $values}}
          <li>li {{.}}</li>{{end}}
        </ul>
      {{else}} {{index $values 0}}{{end}}</li>{{end}}
    </ul>
  </body>
</html>
`
const error_500 = `<!DOCTYPE html>
<html>
  <head>
    <title>Error 500</title>
  </head>
  <body>
    <h1>Error 500</h1>
    <p>%v</p>
  </body>
</html>
`
const error_404 = `<!DOCTYPE html>
<html>
  <head>
    <title>Error 404</title>
  </head>
  <body>
    <h1>Error 404：不要乱看，出问题了吧？！</h1>
    <p>%v</p>
  </body>
</html>
`
