package main

import (
	"context"
	"day-10/connection"
	"day-10/middleware"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	ID           int
	Description  string
	Image        string
	Author       string
	ProjectName  string
	Technologies []string
	StartDate    time.Time
	EndDate      time.Time
	Duration     string
	FormatDate   string
	FormatDatee  string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type DataSession struct {
	IsLogin bool
	Name    string
}

var userdata = DataSession{}

func main() {

	connection.DatabaseConnect()

	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	e.Static("/public", "public")
	e.Static("/upload", "upload")

	e.GET("/", home)
	e.GET("/addProject", addProject)
	e.GET("/contactMe", contactMe)
	e.GET("/projeect-detail/:id", projectDetail)
	e.GET("/editProjectCard/:id", editProjects)
	e.GET("/form-regist", formRegist)
	e.GET("/form-login", formLogin)

	e.POST("/logout", logout)
	e.POST("/login", login)
	e.POST("/regist", regist)
	e.POST("/edit-project/:id", middleware.UploadFile(editFormProject))
	e.POST("/delete-project/:id", deleteProject)
	e.POST("/addFormProject", middleware.UploadFile(addFormProject))

	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userdata.IsLogin = false
	} else {
		userdata.IsLogin = sess.Values["isLogin"].(bool)
		userdata.Name = sess.Values["name"].(string)
	}

	var result []Project
	var query string

	if userdata.IsLogin {

		userID := sess.Values["id"].(int)
		query = fmt.Sprintf("SELECT tb_projects.id, description, image, tb_user.name AS author, name_project, technologies, start_date, end_date, durasi FROM tb_projects JOIN tb_user ON tb_projects.author = tb_user.id WHERE tb_user.id = %d ORDER BY tb_projects.id DESC", userID)
	} else {

		query = "SELECT tb_projects.id, description, image, tb_user.name AS author, name_project, technologies, start_date, end_date, durasi FROM tb_projects JOIN tb_user ON tb_projects.author = tb_user.id ORDER BY tb_projects.id DESC"
	}

	item, _ := connection.Conn.Query(context.Background(), query)

	for item.Next() {
		var each Project
		err := item.Scan(&each.ID, &each.Description, &each.Image, &each.Author, &each.ProjectName, &each.Technologies, &each.StartDate, &each.EndDate, &each.Duration)

		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		result = append(result, each)
	}

	if len(result) == 0 {
		projects := map[string]interface{}{
			"projects":    result,
			"SnapStatus":  sess.Values["status"],
			"SnapMessage": sess.Values["message"],
			"DataSession": userdata,
			"IsEmpty":     true,
		}

		delete(sess.Values, "message")
		delete(sess.Values, "status")
		sess.Save(c.Request(), c.Response())

		var tmpl, err = template.ParseFiles("views/index.html")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return tmpl.Execute(c.Response(), projects)
	}

	projects := map[string]interface{}{
		"projects":    result,
		"SnapStatus":  sess.Values["status"],
		"SnapMessage": sess.Values["message"],
		"DataSession": userdata,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func addProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/addProject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userdata.IsLogin = false
	} else {
		userdata.IsLogin = sess.Values["isLogin"].(bool)
		userdata.Name = sess.Values["name"].(string)
	}
	snap := map[string]interface{}{
		"DataSession": userdata,
	}
	return tmpl.Execute(c.Response(), snap)
}

func projectDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	query := `SELECT tb_projects.id, description, image, tb_user.name as author, name_project, technologies, start_date, end_date, durasi FROM tb_projects JOIN tb_user ON tb_projects.author = tb_user.id WHERE tb_projects.id = $1`
	row := connection.Conn.QueryRow(context.Background(), query, id)

	var project = Project{}
	err = row.Scan(&project.ID, &project.Description, &project.Image, &project.Author, &project.ProjectName, &project.Technologies, &project.StartDate, &project.EndDate, &project.Duration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	project.FormatDate = project.StartDate.Format("02/01/2006")
	project.FormatDatee = project.EndDate.Format("02/01/2006")

	templateFile := "views/add-project-detail.html"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userdata.IsLogin = false
	} else {
		userdata.IsLogin = sess.Values["isLogin"].(bool)
		userdata.Name = sess.Values["name"].(string)
	}

	item := map[string]interface{}{
		"Project":     project,
		"DataSession": userdata,
	}

	return tmpl.Execute(c.Response(), item)
}

func contactMe(c echo.Context) error {
	var template, err = template.ParseFiles("views/contact-me.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userdata.IsLogin = false
	} else {
		userdata.IsLogin = sess.Values["isLogin"].(bool)
		userdata.Name = sess.Values["name"].(string)
	}
	snap := map[string]interface{}{
		"DataSession": userdata,
	}
	return template.Execute(c.Response(), snap)
}

func addFormProject(c echo.Context) error {
	projectName := c.FormValue("projectName")
	startDateStr := c.FormValue("startDate")
	endDateStr := c.FormValue("endDate")
	description := c.FormValue("desc")
	technologies := c.Request().Form["technologies"]
	sess, _ := session.Get("session", c)
	author := sess.Values["id"].(int)
	image := c.Get("dataFile").(string)

	// Validasi input
	if projectName == "" || startDateStr == "" || endDateStr == "" || description == "" || len(technologies) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data yang dikirim tidak lengkap"})
	}

	// Parsing tanggal mulai dan tanggal selesai
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Format tanggal mulai tidak valid"})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Format tanggal selesai tidak valid"})
	}

	// Menghitung durasi antara tanggal mulai dan tanggal selesai
	duration := endDate.Sub(startDate)

	// Menghitung durasi dalam satuan hari
	days := int(duration.Hours() / 24)

	// Mengubah durasi menjadi format "X days"
	durationFormatted := fmt.Sprintf("%d days", days)

	// Insert data ke database
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects (description, image, author, name_project, technologies, start_date, end_date, durasi) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", description, image, author, projectName, technologies, startDate, endDate, durationFormatted)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(delete echo.Context) error {
	i, _ := strconv.Atoi(delete.Param("id"))

	fmt.Println("index : ", i)

	_, error := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", i)
	if error != nil {
		return delete.JSON(http.StatusInternalServerError, map[string]string{"message": error.Error()})
	}

	return delete.Redirect(http.StatusMovedPermanently, "/")
}

func editProjects(c echo.Context) error {
	// Mendapatkan ID dari parameter URL
	i, erri := strconv.Atoi(c.Param("id"))
	if erri != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": erri.Error()})
	}
	var editproject = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, description, image, name_project, technologies, start_date, end_date FROM tb_projects WHERE id=$1", i).
		Scan(&editproject.ID, &editproject.Description, &editproject.Image, &editproject.ProjectName, &editproject.Technologies, &editproject.StartDate, &editproject.EndDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userdata.IsLogin = false
	} else {
		userdata.IsLogin = sess.Values["isLogin"].(bool)
		userdata.Name = sess.Values["name"].(string)
	}

	item := map[string]interface{}{
		"Project":     editproject,
		"DataSession": userdata,
	}
	var tmpl, errtemplate = template.ParseFiles("views/editProject.html")
	if errtemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), item)

}

func editFormProject(c echo.Context) error {

	id, erri := strconv.Atoi(c.Param("id"))
	if erri != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": erri.Error()})
	}
	projectName := c.FormValue("projectName")
	startDateStr := c.FormValue("startDate")
	endDateStr := c.FormValue("endDate")
	description := c.FormValue("desc")
	technologies := c.Request().Form["technologies"]
	sess, _ := session.Get("session", c)
	author := sess.Values["id"].(int)
	image := c.Get("dataFile").(string)

	// Validasi input
	if projectName == "" || startDateStr == "" || endDateStr == "" || description == "" || len(technologies) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "isi dengan lengkap"})
	}

	// Parsing startDate and endDate
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Format tanggal mulai tidak valid"})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Format tanggal selesai tidak valid"})
	}

	// Menghitung durasi antara tanggal mulai dan tanggal selesai
	duration := endDate.Sub(startDate)

	// Menghitung durasi dalam satuan hari
	days := int(duration.Hours() / 24)

	// Mengubah durasi menjadi format "X days"
	durationFormatted := fmt.Sprintf("%d days", days)

	// Update data proyek ke tabel database
	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name_project=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6, author=$7, durasi=$8 WHERE id=$9",
		projectName, startDate, endDate, description, technologies, image, author, durationFormatted, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// Redirect ke halaman proyek setelah update berhasil
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func formRegist(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/form-regist.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userdata.IsLogin = false
	} else {
		userdata.IsLogin = sess.Values["isLogin"].(bool)
		userdata.Name = sess.Values["name"].(string)
	}

	snap := map[string]interface{}{
		"SnapStatus":  sess.Values["status"],
		"SnapMessage": sess.Values["message"],
		"DataSession": userdata,
	}
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), snap)
}

func regist(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := c.FormValue("input-name")
	email := c.FormValue("input-email")
	password := c.FormValue("input-pass")

	// Check if email already exists in tb_user
	var count int
	err = connection.Conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tb_user WHERE email = $1", email).Scan(&count)
	if err != nil {
		return redirectWithMessage(c, "regist failed bro", false, "/form-regist")
	}

	if count > 0 {
		return redirectWithMessage(c, "Email already exists", false, "/form-regist")
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		return redirectWithMessage(c, "regist failed bro", false, "/form-regist")
	}

	return redirectWithMessage(c, "Your registration is successful", true, "/form-login")
}

func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)

	snap := map[string]interface{}{
		"SnapStatus":  sess.Values["status"],
		"SnapMessage": sess.Values["message"],
	}
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	tmpl, err := template.ParseFiles("views/form-login.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), snap)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("input-email")
	password := c.FormValue("input-pass")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return redirectWithMessage(c, "wrong email", false, "/form-login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "wrong pass bro", false, "/form-login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800
	sess.Values["message"] = "login success bro!"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	fmt.Println("message", message)
	return c.Redirect(http.StatusMovedPermanently, path)
}
