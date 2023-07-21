package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// 1.定义响应路由的页面
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is Index page!")
}

// 1.定义响应路由的页面
func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the about page.")
}

// 3.处理表单提交
func submit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uname := r.Form.Get("username")
	fmt.Fprintf(w, "usename : %s", uname)

}

// 4.使用模板引擎
type User struct {
	Username string
	Password string
}

// 4.使用模板引擎
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		uname := r.Form.Get("uname")
		pword := r.Form.Get("pass")

		user := User{uname, pword}
		fmt.Println(user)
		tmpl := template.Must(template.ParseFiles("login.html"))
		tmpl.Execute(w, user)
	} else {
		http.ServeFile(w, r, "login.html")
	}
}

// 5.上传文件
func uphtml(w http.ResponseWriter, r *http.Request) {
	html := `<html>
	<body>
	<form action="/upfile" method="post" enctype="multipart/form-data">
		<input type="file" name="upfile" />
		<input type="submit" value="submit file"/>
	</body></html>
	`
	w.Write([]byte(html))
}

// 5.上传文件
func upload(w http.ResponseWriter, r *http.Request) {
	//读取请求文件
	file, _, err := r.FormFile("upfile")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	//将文件存入到服务器指定位置
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = ioutil.WriteFile("./uploads/updata.zip", data, 0666)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("file up success!"))
	// fmt.Printf("upfile ok")

}

// 6.文件下载
func download(w http.ResponseWriter, r *http.Request) {
	//设置文件头
	w.Header().Add("Content-Disposition", "attachment;filename=updata.zip")
	w.Header().Add("Content-Type", "application/octet-stream")
	//读取文件内容
	file, err := ioutil.ReadFile("./uploads/updata.zip")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	//返回文件内容
	io.WriteString(w, string(file))

}

func main() {
	//1.设置url对应的路由
	http.HandleFunc("/", home)
	http.HandleFunc("/index", home)
	http.HandleFunc("/about", about)

	//2.处理静态文件
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	//3.处理表单提交
	http.HandleFunc("/submit", submit) //http://localhost:84/login?username=aaa

	//4.使用模板引擎
	http.HandleFunc("/login", login)

	//5.上传文件
	http.HandleFunc("/uphtml", uphtml)
	http.HandleFunc("/upfile", upload)

	//6.文件下载
	http.HandleFunc("/download", download)

	log.Fatal(http.ListenAndServe(":84", nil))
}
