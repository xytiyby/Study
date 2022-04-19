package controller

import (
	"Study/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Application) AddStudentShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Msg  string
		Flag bool
	}{

		Msg:  "",
		Flag: false,
	}
	ShowView(w, r, "register.html", data)
}
func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "register.html", nil)
}
func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Msg  string
		Flag bool
	}{

		Msg:  "",
		Flag: false,
	}
	ShowView(w, r, "query.html", data)
}
func (app *Application) AddStudent(w http.ResponseWriter, r *http.Request) {

	stu := service.Student{
		Name:      r.FormValue("name"),
		EntityID:  r.FormValue("entityID"),
		StudentID: r.FormValue("studentID"),
		Password:  r.FormValue("password"),
		Major:     r.FormValue("major"),
	}

	app.Setup.SaveStudent(stu)
	r.Form.Set("studentID", stu.StudentID)
	r.Form.Set("name", stu.Name)
	app.FindCertByNoAndName(w, r)
}
func (app *Application) FindCertByNoAndName(w http.ResponseWriter, r *http.Request) {
	studentID := r.FormValue("studentID")
	name := r.FormValue("name")
	result, err := app.Setup.FindByEntityIdAndName(studentID, name)
	var stu = service.Student{}
	json.Unmarshal(result, &stu)

	fmt.Println("根据证书编号与姓名查询信息成功：")
	fmt.Println(stu)

	data := &struct {
		stu  service.Student
		Msg  string
		Flag bool
	}{
		stu:  stu,
		Msg:  "",
		Flag: false,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}
