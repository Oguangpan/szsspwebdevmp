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

func (s *Basic_information_of_device) Inquire(mac string) (ok bool) {
	var t string = `select * from COMPUTERA where MAC='` + mac + `'`
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

// 路由
type MyMux struct{}

type hd struct {
	Msg  string
	Data *Basic_information_of_device
}

var h = new(hd)

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		mainPage(w, r)
		return
	}
	if r.URL.Path == "/edit" {
		editPage(w, r)
		return
	}
	if r.URL.Path == "/query" {
		queryPage(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/index.html", "templates/table.html", "templates/head.html", "templates/tail.html")
	h.Msg = "无权限限制请随意更改数据库"
	//h.Data = d
	h.Data.User = ""
	h.Data.Dep = ""
	h.Data.Type = ""
	h.Data.Mac = ""
	h.Data.Ip = ""
	h.Data.Sys = ""
	h.Data.Disk = ""
	h.Data.Id = ""
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
			if ok := d.Inquire(v); ok {
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
