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

type diannao struct {
	shuxing map[string]string
}

//var d = new(diannao)
//d.chaxun("92:4f:9a:ec:78:fd")
//d.genggai("92:4f:9a:ec:78:fd")
//d.shanchu("92:4f:9a:ec:78:fd")
//d.charu(map[string]string)

// 查询数据
func (s *diannao) chaxun(mac string) (ok bool) {
	var t string = `select * from COMPUTERA where MAC='` + mac + `'`
	row := db.QueryRow(t)
	err := row.Scan(&p["id"], &p["user"], &p["dep"], &p["type"], &p["sys"], &p["ip"], &p["mac"], &p["disk"])
	if err != nil {
		return false
	}
	return true
}
func (s *diannao) genggai(mac string) (ok bool) {

	stmt, err := db.Prepare("update COMPUTERA set ID=?,USER=?,DEP=?,TYPE=?,SYS=?,IP=?,MAC=?,DISK=? where MAC=?")
	if err != nil {
		return false
	}
	stmt.Exec(s.shuxing["id"], s.shuxing["user"], s.shuxing["dep"], s.shuxing["type"], s.shuxing["sys"], s.shuxing["ip"], s.shuxing["mac"], s.shuxing["disk"], mac)
	return true

}
func (s *diannao) shanchu(mac string) (ok bool) {
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
func (s *diannao) charu() (ok bool) {
	sql := `INSERT INTO COMPUTERA VALUES ('` + s.shuxing["id"] + `','` + s.shuxing["user"] + `','` + s.shuxing["dep"] + `','` + s.shuxing["type"] + `','` + s.shuxing["sys"] + `','` + s.shuxing["ip"] + `','` + s.shuxing["mac"] + `','` + s.shuxing["disk"] + `');`
	db.Exec(sql)
	return true
}

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
	if r.URL.Path == "/query" {
		queryPage(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/index.html", "templates/table.html", "templates/head.html", "templates/tail.html")
	*hd.Msg = "欢迎使用办公设备信息查询系统"
	//	data := struct {
	//		Msg  string
	//		Data []string
	//	}{
	//		Msg:  "欢迎使用办公设备信息查询系统",
	//		Data: []string{},
	//	}
	t.ExecuteTemplate(w, "index", *hd)

}

func queryPage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/index.html", "templates/table.html", "templates/head.html", "templates/tail.html")

	v := r.FormValue("MACID")

	if v != "" {
		if m, _ := regexp.MatchString("^([A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2}", v); !m {
			//			data := struct {
			//				Msg  string
			//				Data []string
			//			}{
			//				Msg:  "请输入正确的mac格式...cmd命令行中输入ipconfig即可看到网卡的正确mac信息",
			//				Data: []string{},
			//			}
			*hd.Msg = "请输入正确的mac格式...cmd命令行中输入ipconfig即可看到网卡的正确mac信息"
			t.ExecuteTemplate(w, "index", *hd)
		} else {
			// 调用数据库查询 v 返回对应数据到 Data 里面
			//			data := struct {
			//				Msg  string
			//				Data []string
			//			}{
			//				Msg:  "查询结果如下",
			//				Data: []string{},
			//			}
			*hd.Msg = "查询结果如下"

			/*
					<tr>
				        <td>使用者姓名</td>
				        <td>张三</td>
				    </tr>
				    <tr class="alt">
				        <td>所属部门</td>
				        <td>生产安全部</td>
				    </tr>
			*/
			t.ExecuteTemplate(w, "index", *hd)
		}

		return

	} else {
		//		data := struct {
		//			Msg  string
		//			Data []string
		//		}{
		//			Msg:  "请输入查询内容",
		//			Data: []string{},
		//		}
		*hd.Msg = "请输入查询内容"
		t.ExecuteTemplate(w, "index", hd)
	}

}

func editPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/edit.html")
	t.ExecuteTemplate(w, "edit", "")
	return
}

func main() {

	// 特别说明: 要使用80端口需要使用管理员身份运行程序.
	mux := &MyMux{}
	fmt.Println("now Listening port...")
	err := http.ListenAndServe(":9999", mux)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}
