package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	mvc "github.com/ohayao/gomvc"
	"github.com/shanzhaiwukong/gutil"
)

type admin struct{}

func (*admin) Get_login(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	in.ParseParam(0)
	year := in.GetUrlInt32("year", 2000)
	fmt.Println("收到参数 year=", year)
	if in.GetCookie("token") != nil {
		out.HtmlView("views/area/admin/login", nil, nil, `已登记，<a href="./index">去聊天</a>`)
	} else {
		out.HtmlView("views/area/admin/login", nil, nil, `未登记，<a href="./index">去登记</a>`)
	}
	return out
}
func (*admin) Post_login(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	out.Text([]byte(`<p>张三，你好啊！</p>`))

	return out
}

var groupLock sync.RWMutex
var groupList map[string]map[string]*mvc.Output

func init() {
	groupList = make(map[string]map[string]*mvc.Output)
	broadcast()
}

func (*admin) Get_index(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	out.HtmlView("views2/admin/index", nil, nil, nil)
	fmt.Println("=======================> ......................")
	return out
}

func (*admin) Get_websocket(in *mvc.Input) *mvc.Output {
	out := mvc.NewOutput()
	in.ParseParam(0)
	gid := in.GetUrlString("gid", gutil.NewID())
	uid := in.GetUrlString("uid", gutil.NewID())
	groupLock.Lock()
	if groupList[gid] == nil {
		groupList[gid] = make(map[string]*mvc.Output)
	}
	if _, ok := groupList[gid][uid]; !ok {
		groupList[gid][uid] = out
	}
	groupLock.Unlock()
	go func() {
		for {
			if rd, mt, mv := out.WebsocketRead(); rd && mt > 0 {
				if string(mv) == "sys_ping" {
					out.WebSocketWrite(1, []byte("sys_pong"))
				} else {
					toOther(gid, uid, mv)
				}
			} else {
				out.CloseWebsocket()
				break
			}
		}
		groupLock.Lock()
		delete(groupList[gid], uid)
		if len(groupList[gid]) < 1 {
			delete(groupList, gid)
		}
		groupLock.Unlock()
	}()
	return out.Websocket(0, true)
}

func toOther(gid, uid string, data []byte) {
	groupLock.Lock()
	defer groupLock.Unlock()
	pre := make([][]byte, 0)
	pre = append(pre, []byte(fmt.Sprintf("%s说：", uid)))
	pre = append(pre, data)
	data = bytes.Join(pre, nil)
	for k, v := range groupList[gid] {
		if k != uid {
			if !v.WebSocketWrite(1, data) {
				delete(groupList[gid], uid)
			}
		}
	}
}

func broadcast() {
	groupLock.Lock()
	defer func() {
		groupLock.Unlock()
		time.AfterFunc(time.Minute*5, func() {
			broadcast()
		})
	}()
	for k, v := range groupList {
		for k1, v1 := range v {
			if !v1.WebSocketWrite(1, []byte(fmt.Sprintf("系统消息：群【%s】用户【%s】系统时间:%s", k, k1, time.Now().Format("2006/01/02 15:04:05")))) {
				delete(groupList[k], k1)
			}
		}
	}
}
