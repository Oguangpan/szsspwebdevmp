package devmpdata

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Cmpr struct {
	id   string
	user string
	dep  string
	ty   string
	sys  string
	ip   string
	mac  string
	dkid string
}

var db *sql.DB
var p Cmpr

func init() {
	db, _ = sql.Open("sqlite3", "./szdevDB.db")
}

func Inc(p map[string]string) bool {
	//根据参数增加数据,成功后返回真
	sql := `INSERT INTO COMPUTERA VALUES ('000001','张飞','设备管理部','台式兼容机','Windows xp','33.66.99.88','3a:33:1d:90:3c:21','DEB28471D');`
	db.Exec(sql)
	return true
}

func Del(s string) bool {
	// 根据传入的mac地址删除所有,成功返回真
	stmt, err := db.Prepare("delete from COMPUTERA where MAC=?")
	if err != nil {
		return false
	}
	_, err = stmt.Exec(s)
	if err != nil {
		return false
	}
	return true
}

func Mod(p map[string]string, m string) bool {
	// 根据传入的字典内容修改内容
	stmt, err := db.Prepare("update COMPUTERA set ID=?,USER=?,DEP=?,TYPE=?,SYS=?,IP=?,MAC=?,DISK=? where MAC=?")
	if err != nil {
		return false
	}
	stmt.Exec(p["id"], p["user"], p["dep"], p["type"], p["sys"], p["ip"], p["mac"], p["disk"], m)
	return true
}
func Que(s string) (p Cmpr, b bool) {
	// 查询s提供的数据(mac),并且返回p结构与布尔值(p,ok:=devmpdata.Que("mac地址")).
	var t string = `select * from COMPUTERA where MAC='` + s + `'`
	row := db.QueryRow(t)
	err := row.Scan(&p.id, &p.user, &p.dep, &p.ty, &p.sys, &p.ip, &p.mac, &p.dkid)
	if err != nil {
		return p, false
	}
	return p, true

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
