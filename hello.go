package main

import (
	"net/http"
	"time"

	mvc "github.com/ohayao/gomvc"
)

type hello struct{}

func (*hello) Get_vp1(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	//使用文件模版测试(嵌套模版)
	out.HtmlView("taaa", []string{"taaat"}, nil, nil)
	return out
}
func (*hello) Get_vf1(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	//使用视图文件夹模版测试（单个文件模版）
	out.HtmlView("views/home/index", nil, nil, nil)
	return out
}
func (*hello) Get_vf2(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	//使用视图文件模版测试(嵌套模版)
	out.HtmlView("views/home/about", []string{"views/template/header", "views/template/footer"}, nil, nil)
	return out
}
func (*hello) Get_vf3(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	//使用文本视图
	out.HtmlView("mtxt", nil, nil, nil)
	return out
}
func (*hello) Get_vf4(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	//使用混合视图
	fns := map[string]interface{}{
		"Test": func() string {
			return "Test ....."
		},
	}
	data := map[string]interface{}{
		"Time":  time.Now(),
		"Times": time.Now().Unix(),
	}
	out.CookieAdd(&http.Cookie{})
	out.HtmlView("views/home/contact", []string{"views/template/header", "views/template/footer", "taaat", "mtxt"}, fns, data)
	return out
}
