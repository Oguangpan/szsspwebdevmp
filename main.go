package main

import (
	"fmt"
	"html/template" // 生成需要的网页模板
	"net/http"
	"regexp" //验证用户输入

	_ "devmp/data"
)

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

	t, _ := template.ParseFiles("templates/index.html")
	t.ExecuteTemplate(w, "index", "")

}

func queryPage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/index.html", "templates/table.html")
	t.ExecuteTemplate(w, "index", "")
	// t.Execute(w, "目前未有查询数据展示...")
	v := r.FormValue("MACID")
	if v != "" {
		// 取消正则表达式中对XX-XX-XX-XX-XX-XX的匹配,因为懒不想做转换匹配数据库中的内容
		if m, _ := regexp.MatchString("^([A-Fa-f0-9]{2}:{5}[A-Fa-f0-9]{2}", v); !m {
			fmt.Println("User input is not MAC")
			t.ExecuteTemplate(w, "index",
				"<h2>请输入正确合规的查询内容.</h2><br><div class=\"prompt\">cmd命令行中输入ipconfig即可看到网卡的正确mac信息.</div>")
			//fmt.Fprintf(w, "")
			return
		}

		//根据用户输入进入数据库查询并且通过模板写出数据
		t.ExecuteTemplate(w, "table", "")
		return

	}
	t.ExecuteTemplate(w, "index", "<h2>请输入查询内容.</h2>")
	// fmt.Fprintf(w, "<h2>请输入查询内容.</h2>")

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
