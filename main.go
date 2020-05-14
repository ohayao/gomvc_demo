package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	mvc "github.com/ohayao/gomvc"
)

func init() {
	cfgStatics()
	cfgView()
	cfgPage()
	cfgController()
	cfgMiddleware()
}

func main() {
	//mvc.ShowRoutes()
	err := mvc.Run(11225)
	fmt.Println(err)
}

func cfgStatics() {
	//正则表达式匹配备份、配置、压缩包等后缀
	reg, _ := regexp.Compile(`\.(bak|zip|gz|rar|ini)$`)
	mvc.RegistStatics(&mvc.Statics{
		AbsolutePath:     "./statics",
		RelativePath:     "/",
		DefaultPage:      "index.html",
		IgnoreFileRegexp: reg,
	})
}

func cfgView() {
	//配置视图文件夹
	mvc.ConfigViewPath("./views", []string{".html"}, time.Second*10)
	mvc.ConfigViewPath("./views2", []string{".html", ".gtpl"}, time.Minute*1)
	//配置单个视图文件
	mvc.ConfigViewFile("taaa", "./temp/aaa.html", time.Second*60)
	mvc.ConfigViewFile("taaat", "./temp/taaat.html", time.Second*60)
	//配置文本视图
	mvc.ConfigViewText("mtxt", `{{define "mtxt"}}<div>Hello this is txt template for "mtxt"</div>{{end}}`)

	//配置视图函数
	mvc.ConfigViewFunc("Hello", func(name string) string {
		return fmt.Sprintf("[GFN]Hello,%s", name)
	})
}

func cfgPage() {
	//注册home首页
	mvc.RegistPage("GET", "/home", home1)
	mvc.RegistPage("GET", "/", home)

	//注册404 500页面
	mvc.RegistCodePage(404, func(in *mvc.Input) *mvc.Output {
		out := mvc.NewOutput()
		return out.SetStatusCode(404).Html([]byte(`<html><head><title>404 NotFound</title></head><body><h1>404 ....</h1></body></html>`))
	})
	mvc.RegistCodePage(500, func(in *mvc.Input) *mvc.Output {
		out := mvc.NewOutput()
		return out.SetStatusCode(500).Html([]byte(`<html><head><title>500 Server Error</title></head><body><h1>500 ....</h1></body></html>`))
	})
}

func cfgController() {
	mvc.RegistController("", "", "admin", &admin{})
	mvc.RegistController("api", "v1", "admin", &admin{})
	mvc.RegistController("api", "v2", "admin", &admin{})
	mvc.RegistController("a/b/c", "v3", "admin", &admin{})
	mvc.RegistController("", "", "hello", &hello{})
}

func cfgMiddleware() {
	logMiddleware := &mvc.Middleware{
		Name: "LOG",
		Handler: func(in *mvc.Input) (*mvc.Output, bool) {
			fmt.Printf("LOG [%s:%s] [UA=%s]\n", in.Method, in.URL, in.GetHeader("user-agent"))
			return nil, false
		},
		Enable:   true,
		Position: mvc.MiddlewareBefore,
		Rules: []*mvc.MiddlewareRule{
			&mvc.MiddlewareRule{
				Method:    []string{"GET", "POST"},
				IsSkip:    true,
				URLRegexp: regexp.MustCompile(`^/admin.*`),
			},
			&mvc.MiddlewareRule{
				Method:    []string{"GET", "POST"},
				IsSkip:    false,
				URLRegexp: regexp.MustCompile(`.*`),
			},
		},
	}
	//跳过/admin开头的请求 注意rules先后顺序
	fmt.Println("Use Middleware LOG ", mvc.UseMiddleware(logMiddleware))

	authMiddleware := &mvc.Middleware{
		Name:   "AUTH",
		Enable: true,
		Handler: func(in *mvc.Input) (*mvc.Output, bool) {
			fmt.Println("Exec AUTH middleware ..", in.URL)
			if ck := in.GetCookie("token"); ck == nil {
				fmt.Println("token is nil, redirect")
				out := mvc.NewOutput()
				out.CookieAdd(&http.Cookie{Name: "token", Value: "123456", Expires: time.Now().Add(time.Minute * 5), Path: "/"})
				out.Redirect("/admin/login?reffer="+in.URL, 302)
				return out, true
			}
			return nil, false
		},
		Position: mvc.MiddlewareBefore,
		Rules: []*mvc.MiddlewareRule{
			&mvc.MiddlewareRule{
				Method:    []string{"GET", "POST"},
				IsSkip:    true,
				URLRegexp: regexp.MustCompile(`/admin/login`),
			},
			&mvc.MiddlewareRule{
				Method:    []string{"GET", "POST"},
				IsSkip:    false,
				URLRegexp: regexp.MustCompile(`/admin.*`),
			},
		},
	}
	//auth 用户认证
	fmt.Println("Use Middleware AUTH ", mvc.UseMiddleware(authMiddleware))
}

func home(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	out.Json([]byte(`{"Hello":"张三","World":"你好啊！"}`))
	return out
}
func home1(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	return out.Redirect("/", 302)
}
