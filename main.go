package main

import (
	"database/sql"
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
	case "/editerr":
		editerrPage(w, r)
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
	} else {
		h.Msg = "请输入查询内容"
		h.Data = d
		t.ExecuteTemplate(w, "index", h)
	}
}

func editerrPage(w http.ResponseWriter, r *http.Request) {
	display_page := func(w http.ResponseWriter, r *http.Request, h string, mw string) {
		t, _ := template.ParseFiles(h,
			"templates/index-top.tmpl",
			"templates/index-bottom.tmpl")
		t.ExecuteTemplate(w, mw, "")
	}
	r.ParseForm()
	log.Println(r.PostForm.Get("sel"))
	if r.Method == "GET" {
		display_page(w, r, "templates/editerr.tmpl", "editerr")
	} else {

		if r.PostForm.Get("sel") == "yes" {
			if d.Modify(d.Mac) {
				display_page(w, r, "templates/editok.tmpl", "editok")
			}
		} else {
			log.Println("user select clecot")
			display_page(w, r, "templates/edit.tmpl", "edit")
		}
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
	// 通过判断请求方法是get还是post决定执行
	if r.Method == "GET" {
		display_page(w, r, "templates/edit.tmpl", "edit")
	} else {
		d.Id = r.PostForm.Get("id")
		d.User = r.PostForm.Get("user")
		d.Dep = r.PostForm.Get("department")
		d.Ip = r.PostForm.Get("ip")
		d.Mac = r.PostForm.Get("mac")
		d.Sys = r.PostForm.Get("system_type")
		d.Type = r.PostForm.Get("Equipment_type")
		d.Disk = r.PostForm.Get("diskid")
		if d.Inquire(d.Mac, "mac") {
			log.Println("mac已存在")
			display_page(w, r, "templates/editerr.tmpl", "editerr")
		} else if d.Inquire(d.Id, "id") {
			log.Println("id已存在")
			display_page(w, r, "templates/editerr.tmpl", "editerr")
		} else if d.Inquire(d.Disk, "diskid") {
			log.Println("diskid已存在")
			display_page(w, r, "templates/editerr.tmpl", "editerr")
		} else {
			//没有重复数据,直接添加到服务器.
			if d.Increase() {
				display_page(w, r, "templates/editok.tmpl", "editok")
			}

		}

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
