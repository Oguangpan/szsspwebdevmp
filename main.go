package main

import (
	"database/sql"
	"fmt"
	"html/template" // 生成需要的网页模板
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
	if r.URL.Path == "/" {
		mainPage(w, r)
		return
	}
	if r.URL.Path == "/edit" {
		editPage(w, r)
		return
	}
	if r.URL.Path == "/edit_post" {
		editpostPage(w, r)
		return
	}
	if r.URL.Path == "/query" {
		queryPage(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func editpostPage(w http.ResponseWriter, r *http.Request) {
	// TODO
	// 与用户交互用的页面
	// 这里必须学会使用JavaScript与golang交互之后才能继续写下去了.
	// 获取输入(关键数据验证由页面上的javascript来验证)
	dn := new(Basic_information_of_device)
	dn.Id = r.PostForm.Get("id")
	dn.User = r.PostForm.Get("user")
	dn.Dep = r.PostForm.Get("department")
	dn.Ip = r.PostForm.Get("ip")
	dn.Mac = r.PostForm.Get("mac")
	dn.Sys = r.PostForm.Get("system_type")
	dn.Type = r.PostForm.Get("Equipment_type")
	dn.Disk = r.PostForm.Get("diskid")

	// 在数据库中查询提交的数据是否重合,判断的标准是MAC\ID\DISKID其中之一.
	if ok := d.Inquire(dn.Mac, "mac"); ok {
		// 提示用户数据重复,询问是否修改,如果选择修改就修改,不修改就返回
		return
	}
	if ok := d.Inquire(dn.Mac, "id"); ok {
		return
	}
	if ok := d.Inquire(dn.Mac, "diskid"); ok {
		return
	}
	// 添加新数据
	// 修改老数据
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html", "templates/table.html", "templates/head.html", "templates/tail.html")
	h.Msg = "无权限限制请随意更改数据库"
	h.Data = d
	t.ExecuteTemplate(w, "index", h)

}

func queryPage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/index.html", "templates/table.html", "templates/head.html", "templates/tail.html")

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
	t, _ := template.ParseFiles("templates/edit.html")
	t.ExecuteTemplate(w, "edit", "")
	return
}

func main() {
	defer db.Close()
	// 特别说明: 要使用80端口需要使用管理员身份运行程序.
	mux := &MyMux{}
	fmt.Println("now Listening port...")
	err := http.ListenAndServe(":9999", mux)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}
