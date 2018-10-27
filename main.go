package main

import (
	"database/sql"
	"encoding/json" // ajax 返回给网页脚本端的内容
	"html/template" // 生成需要的网页模板
	"io"
	"log"
	"net/http"
	"regexp" //验证用户输入

	_ "github.com/mattn/go-sqlite3" // 数据库
)

var db *sql.DB

func init() {
	db, _ = sql.Open("sqlite3", "data/szdevDB.db")
}

// 设备信息结构体
type Basic_information_of_device struct {
	User string
	Dep  string
	Type string
	Mac  string
	Ip   string
	Sys  string
	Disk string
	Id   string
}

// 后端返回该结构体给前端ajxa
type Echo struct {
	Msg  string
	Info string
}

var d *Basic_information_of_device = new(Basic_information_of_device)

// 查询(目标,字段)成否
func (s *Basic_information_of_device) Inquire(c string, i string) (ok bool) {
	var t string
	switch i {
	case "mac":
		t = `select * from COMPUTERA where MAC='` + c + `'`
	case "diskid":
		t = `select * from COMPUTERA where DISK='` + c + `'`
	case "id":
		t = `select * from COMPUTERA where ID='` + c + `'`

	}
	row := db.QueryRow(t)
	err := row.Scan(&s.Id, &s.User, &s.Dep, &s.Type, &s.Sys, &s.Ip, &s.Mac, &s.Disk)
	if err != nil {
		return false
	}
	return true
}

func (s *Basic_information_of_device) Modify(mac string) (ok bool) {

	stmt, err := db.Prepare("update COMPUTERA set ID=?,USER=?,DEP=?,TYPE=?,SYS=?,IP=?,MAC=?,DISK=? where MAC=?")
	if err != nil {
		return false
	}
	stmt.Exec(s.Id, s.User, s.Dep, s.Type, s.Sys, s.Ip, s.Mac, s.Disk, mac)
	return true

}

// 输入mac删除对应设备在数据库中的记录
func (s *Basic_information_of_device) Delete(mac string) (ok bool) {
	stmt, err := db.Prepare("delete from COMPUTERA where MAC=?")
	if err != nil {
		return false
	}
	_, err = stmt.Exec(mac)
	if err != nil {
		return false
	}
	return true
}

func (s *Basic_information_of_device) Increase() (ok bool) {
	sql := `INSERT INTO COMPUTERA VALUES ('` + s.Id + `','` + s.User + `','` + s.Dep + `','` + s.Type + `','` + s.Sys + `','` + s.Ip + `','` + s.Mac + `','` + s.Disk + `');`
	db.Exec(sql)
	return true
}

// 包含设备信息与一个用于前后端交互的信息头
type Hd struct {
	Msg  string
	Data *Basic_information_of_device
}

// 实例化
var h = new(Hd)

// 路由
type MyMux struct{}

// 根据访问指向不同的处理器
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		mainPage(w, r)
	case "/edit":
		editPage(w, r)
	case "/query":
		queryPage(w, r)
	case "/editprocess":
		editerrProcessPage(w, r)
	default:
		http.NotFound(w, r)
	}
}

// 展示查询页面
func mainPage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/index.tmpl",
		"templates/index-top.tmpl",
		"templates/index-bottom.tmpl")
	//h.Msg = "无权限限制,请随意更改数据库"
	//h.Data = d
	t.ExecuteTemplate(w, "index", "")
}

// 查询页面的ajxa交互函数
func queryPage(w http.ResponseWriter, r *http.Request) {
	// ajax响应函数,将查询到的数据的json转化为string发送到页面脚本,由脚本提供解析.
	r.ParseForm()
	v := r.PostForm.Get("MACID")
	if v != "" {
		if m, _ := regexp.MatchString("^([A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2}", v); !m {
			h.Data = d
			h.Msg = "请输入正确的mac格式...cmd命令行中输入ipconfig即可看到网卡的正确mac信息"
			m, _ := json.Marshal(h)
			io.WriteString(w, string(m))
		} else {
			if ok := d.Inquire(v, "mac"); ok {
				h.Msg = "查询结果如下"
				h.Data = d
				m, _ := json.Marshal(h)
				io.WriteString(w, string(m))
			} else {
				h.Msg = "数据库中没有该设备信息"
				m, _ := json.Marshal(h)
				io.WriteString(w, string(m))
			}
		}
	} else {
		h.Msg = "地址输入框为空，请输入信息"
		h.Data = d
		m, _ := json.Marshal(h)
		io.WriteString(w, string(m))
	}
}

func editerrProcessPage(w http.ResponseWriter, r *http.Request) {

	/*
		经过几天的挣扎(工作忙)后,我突然有了灵感.
		设计一个简单的交互逻辑(之前查询那块的那个叫h的数据类型好像可以直接拿来用). 通过双方发送json数据来完成以下内容.
		前端发送数据由以下内容组成: {{"访问类型":"写数据/查数据"},{"post数据串":"扒拉扒拉"}}
		后端返回数据由以下内容组成: {{"结果类型":"提示错误/提示成功"},{"错误字段名称":"比如(数据库中已有该ip地址,是否更新原有内容?)/(操作成功完成)"}}

		1.前端首先提交查数据的json,包含用户填写的表单数据.
		2.后端获取数据后根据"访问类型字段"决定是进入查数据还是直接写入数据,根据查询结果返回信息.
		3.前端获取到响应内容后根据"结果类型字段",弹出选择框(选择是否更新数据库内容)或者信息框(确认成功),提示用户相关信息,如果是弹出框根据用户选择结果决定是返回还是提交"写数据"的json.

		大部分代码应该都是在js脚本那里吧.知识果然不是线状的,线越长接触到的面积越宽广.
	*/

	r.ParseForm()
	jsonBlob := r.PostForm.Get("pjson")
	// 先解析传递过来的数据,并传递给h对象
	err := json.Unmarshal([]byte(jsonBlob), &h)
	if err != nil {
		log.Println("error:", err)
	}
	d = h.Data
	e := new(Echo)
	// 判定页面操作类型是查询还是写入 q 查 w 写.
	switch h.Msg {
	case "q":
		//考虑到实际,查询只对mac做唯一性限制
		if d.Inquire(h.Data.Mac, "mac") {
			e.Msg = "e"
			e.Info = "错误,数据库中已存在该设备,是否覆盖原有设备信息?"
			m, _ := json.Marshal(e)
			io.WriteString(w, string(m))
		} else {
			// 未发现新设备,直接写入数据库
			if d.Increase() {
				e.Msg = "w"
				e.Info = "写入成功."
				m, _ := json.Marshal(e)
				io.WriteString(w, string(m))
			} else {
				e.Msg = "w"
				e.Info = "服务器出现错误,写入失败,请联系管理员."
				m, _ := json.Marshal(e)
				io.WriteString(w, string(m))
			}
		}
	case "w":
		//注意:此处应该直接假设数据已存在,并更新数据.
		if d.Modify(h.Data.Mac) {
			e.Msg = "w"
			e.Info = "写入成功."
			m, _ := json.Marshal(e)
			io.WriteString(w, string(m))
		} else {
			e.Msg = "w"
			e.Info = "服务器出现错误,写入失败,请联系管理员."
			m, _ := json.Marshal(e)
			io.WriteString(w, string(m))
		}
	}
}

// 展示编辑页面
func editPage(w http.ResponseWriter, r *http.Request) {
	display_page := func(w http.ResponseWriter, r *http.Request, h string, mw string) {
		t, _ := template.ParseFiles(h,
			"templates/index-top.tmpl",
			"templates/index-bottom.tmpl")
		t.ExecuteTemplate(w, mw, "")
	}
	r.ParseForm()
	if r.Method == "GET" {
		display_page(w, r, "templates/edit.tmpl", "edit")
	}

}

func main() {
	defer db.Close()
	mux := &MyMux{}
	log.Println("now Listening port...")
	err := http.ListenAndServe(":9999", mux)
	if err != nil {
		log.Println("ListenAndServe:", err)
	}

}
