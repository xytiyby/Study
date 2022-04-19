/**
  @Author : hanxiaodong
*/

package web

import (
	"Study/web/controller"
	"fmt"
	"net/http"
)

// 启动Web服务并指定路由信息
func WebStart(app controller.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)

	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/addEduInfo", app.AddStudentShow) // 显示添加信息页面
	http.HandleFunc("/addStudent", app.AddStudent)     // 提交信息请求

	http.HandleFunc("/queryPage", app.QueryPage)       // 转至根据证书编号与姓名查询信息页面
	http.HandleFunc("/query", app.FindCertByNoAndName) // 根据证书编号与姓名查询信息

	http.HandleFunc("/upload", app.UploadFile)
	http.HandleFunc("/addStudent", app.AddStudent)
	fmt.Println("启动Web服务, 监听端口号为: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v", err)
	}

}
