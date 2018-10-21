package main

import (
	"database/sql"
	"fmt"
	"html/template" // 生成需要的网页模板
	"log"
	"net/http"
	"regexp" //验证用户输入

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("sqlite3", "data/szdevDB.db")
}

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

var d *Basic_information_of_device = new(Basic_information_of_device)

// 查询(目标,类型)成否
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

//func (s *Basic_information_of_device) Modify(mac string) (ok bool) {

//	stmt, err := db.Prepare("update COMPUTERA set ID=?,USER=?,DEP=?,TYPE=?,SYS=?,IP=?,MAC=?,DISK=? where MAC=?")
//	if err != nil {
//		return false
//	}
//	stmt.Exec(s.Attributes["id"], s.Attributes["user"], s.Attributes["dep"], s.Attributes["type"], s.Attributes["sys"], s.Attributes["ip"], s.Attributes["mac"], s.Attributes["disk"], mac)
//	return true

//}

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

//func (s *Basic_information_of_device) Increase() (ok bool) {
//	sql := `INSERT INTO COMPUTERA VALUES ('` + s.Attributes["id"] + `','` + s.Attributes["user"] + `','` + s.Attributes["dep"] + `','` + s.Attributes["type"] + `','` + s.Attributes["sys"] + `','` + s.Attributes["ip"] + `','` + s.Attributes["mac"] + `','` + s.Attributes["disk"] + `');`
//	db.Exec(sql)
//	return true
//}

type hd struct {
	Msg  string
	Data *Basic_information_of_device
}

var h = new(hd)

// 路由
type MyMux struct{}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		mainPage(w, r)
	case "/edit":
		editPage(w, r)
	case "/query":
		queryPage(w, r)
	default:
		http.NotFound(w, r)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.tmpl", "templates/table.tmpl",
		"templates/index-top.tmpl", "templates/index-bottom.tmpl")
	h.Msg = "无权限限制,请随意更改数据库"
	h.Data = d
	t.ExecuteTemplate(w, "index", h)

}

func queryPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t, _ := template.ParseFiles("templates/index.tmpl", "templates/table.tmpl",
		"templates/index-top.tmpl", "templates/index-bottom.tmpl")
	v := r.FormValue("MACID")
	if v != "" {
		if m, _ := regexp.MatchString("^([A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2}", v); !m {
			h.Data = d
			h.Msg = "请输入正确的mac格式...cmd命令行中输入ipconfig即可看到网卡的正确mac信息"
			t.ExecuteTemplate(w, "index", h)
		} else {
			// 调用数据库查询 v 返回对应数据到 Data 里面
			if ok := d.Inquire(v, "mac"); ok {
				h.Msg = "查询结果如下"
				h.Data = d
				t.ExecuteTemplate(w, "index", h)
			} else {
				h.Msg = "数据库中没有该设备信息"
				t.ExecuteTemplate(w, "index", h)
			}

		}
		return
	} else {
		h.Msg = "请输入查询内容"
		h.Data = d
		t.ExecuteTemplate(w, "index", h)
	}
}

func editPage(w http.ResponseWriter, r *http.Request) {

	display_page := func(w http.ResponseWriter, r *http.Request, h string, mw string) {
		t, _ := template.ParseFiles(h,
			"templates/index-top.tmpl",
			"templates/index-bottom.tmpl")
		t.ExecuteTemplate(w, mw, "")
	}

	r.ParseForm()

	// TODO 这里有一个非常严重的问题,明明获取到了v的值,却在第二次访问页面的时候无法进入case "editpost":或者是下面的case "editerror":
	// 我反复的验证了v这个隐藏域的值是正确的,但是就是无法成功的进入下面的代码.只有首次访问页面的时候才能正确的执行case "": 后面的内容
	// v变量是隐藏域属性,表示上次访问来源
	switch v := r.PostForm.Get("yc"); v {
	case "":
		display_page(w, r, "templates/edit.tmpl", "edit")
	case "editpost":
		fmt.Println("ri")
		dn := new(Basic_information_of_device)
		dn.Id = r.PostForm.Get("id")
		dn.User = r.PostForm.Get("user")
		dn.Dep = r.PostForm.Get("department")
		dn.Ip = r.PostForm.Get("ip")
		dn.Mac = r.PostForm.Get("mac")
		dn.Sys = r.PostForm.Get("system_type")
		dn.Type = r.PostForm.Get("Equipment_type")
		dn.Disk = r.PostForm.Get("diskid")
		switch ok := false; ok {
		case d.Inquire(dn.Mac, "mac"):
			//如果重复交给用户判断是否覆盖
			display_page(w, r, "templates/editerr.tmpl", "editerr")
			return
		case d.Inquire(dn.Mac, "id"):
			display_page(w, r, "templates/editerr.tmpl", "editerr")
			return
		case d.Inquire(dn.Mac, "diskid"):
			display_page(w, r, "templates/editerr.tmpl", "editerr")
			return
		default:
			//没有重复数据,直接添加到服务器.
			log.Println("现在开始往数据库里面添加数据")

			display_page(w, r, "templates/editok.tmpl", "editok")
			return
		}
	case "editerror": //用户已经决定是否覆盖已有数据展示成功页面.
		//获取用户选择是覆盖还是退出
		log.Println("现在开始往数据库里面添加数据")
		//然后返回一个成功页面,给出一个提示后让用户决定是否返回主页: 其实我要是知道怎么出现弹窗的话就更方便了.
		display_page(w, r, "templates/editok.tmpl", "editok")
	}
}

func main() {
	defer db.Close()
	// 特别说明: 要使用80端口需要使用管理员身份运行程序.
	/*
		模板文件说明
		- 主页模板 #每次都需要载入,下面的模板根据需要加载.
		    - 查询页面
		    - 查询结果页面/新增结果页面/提交结果确认
		    - 编辑页面
		    - 提交确认反馈页面
	*/
	mux := &MyMux{}
	log.Println("now Listening port...")
	err := http.ListenAndServe(":9999", mux)
	if err != nil {
		log.Println("ListenAndServe:", err)
	}

}
