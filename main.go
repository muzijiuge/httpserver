package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectDB() error {
	var err error
	db, err = sql.Open("mysql", "root:123456@/testdb")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("数据库连接成功")
	return nil
}

func HandlerAdd(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")
	_, err := db.Exec("INSERT INTO users(name, age) VALUES (?, ?)", name, age)
	if err != nil {
		fmt.Fprintf(w, "添加用户失败：%v\n", err)
		return
	}

	fmt.Fprintf(w, "成功添加用户\n")
}

func HandlerDel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		fmt.Fprintf(w, "删除用户失败：%v\n", err)
		return
	}

	fmt.Fprintf(w, "成功删除用户\n")
}

func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	_, err := db.Exec("UPDATE users SET name=?, age=? WHERE id=?", name, age, id)
	if err != nil {
		fmt.Fprintf(w, "更新用户失败：%v\n", err)
		return
	}
	fmt.Fprintf(w, "成功更新用户\n")
}

func HandlerFind(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var name string
	var age int
	err := db.QueryRow("SELECT name, age FROM users WHERE id = ?", id).Scan(&name, &age)
	if err != nil {
		fmt.Fprintf(w, "查询用户失败：%v\n", err)
		return
	}

	fmt.Fprintf(w, "用户姓名：%s:年龄:%d\n", name, age)
}

func main() {
	err := connectDB()
	if err != nil {
		fmt.Printf("数据库连接失败：%v\n", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/add", HandlerAdd)
	http.HandleFunc("/del", HandlerDel)
	http.HandleFunc("/update", HandlerUpdate)
	http.HandleFunc("/find", HandlerFind)

	fmt.Println("已开启服务并监听8080端口")
	http.ListenAndServe(":8080", nil)
}
