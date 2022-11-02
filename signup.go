package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error
var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type Class struct {
	Id     int
	Year   string
	Class  string
	Stream string
}
type Classes struct {
	Message  string
	Username string
	Items    []Class
	Enrolled []Class
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.j", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}
func addclass(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "addclass.html")
		return
	}

	year := req.FormValue("year")
	code := req.FormValue("class-code")
	semester := req.FormValue("semester")
	_, err = db.Exec("INSERT INTO `classes`( `year`, `class`, `stream`) VALUES (?,?,?)", year, code, semester)
	http.Redirect(res, req, "/home", 301)

}
func enroll(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(res, "Forbidden", http.StatusForbidden)
		return
	}
	//username, _ := session.Values["username"].(string)
	user_id, _ := session.Values["user-id"].(int)
	class := req.FormValue("class")

	_, err = db.Exec("INSERT INTO `enroll`( class_id,`student_id`) VALUES (?,?)", class, user_id)
	http.Redirect(res, req, "/student", 301)

}
func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string
	var id int
	err := db.QueryRow("SELECT id,username, password FROM users WHERE username=?", username).Scan(&id, &databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	session, _ := store.Get(req, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["user-id"] = id
	session.Save(req, res)
	http.Redirect(res, req, "/student", 301)
	// res.Write([]byte("Hello" + databaseUsername))

}
func lecturer(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "lecturer.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string
	var id int
	err := db.QueryRow("SELECT id,username, password FROM teachers WHERE username=?", username).Scan(&id, &databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	session, _ := store.Get(req, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["user-id"] = id
	session.Save(req, res)
	http.Redirect(res, req, "/home", 301)
	// res.Write([]byte("Hello" + databaseUsername))

}

type Attendants struct {
	Items []Attendance
}
type Attendance struct {

	// defining struct fields
	Name   string
	Class  int
	Status int
}

func attendance(res http.ResponseWriter, req *http.Request) {
	data := Attendants{}
	class := req.FormValue("class")

	rows, _ := db.Query("SELECT username,class_id, status FROM `enroll` inner join users on enroll.student_id=users.id WHERE enroll.class_id=?", class)
	var username string
	var class_id int
	var status int
	for rows.Next() {
		rows.Scan(
			&username,
			&class_id,
			&status,
		)

		std1 := Attendance{username, class_id, status}
		data.Items = append(data.Items, std1)
	}

	t, _ := template.ParseFiles("attendance.html")

	// standard output to print merged data
	t.Execute(res, data)

	//http.ServeFile(res, req,data, "attendance.html")
}
func homePage(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(res, "Forbidden", http.StatusForbidden)
		return
	}

	data := Classes{}

	rows, _ := db.Query("SELECT id,year,class,stream FROM classes")
	var year string
	var class string
	var stream string
	var id int
	for rows.Next() {
		rows.Scan(
			&id,
			&year,
			&class,
			&stream,
		)

		std1 := Class{id, year, class, stream}
		data.Items = append(data.Items, std1)
	}

	t, _ := template.ParseFiles("index.html")

	// standard output to print merged data
	t.Execute(res, data)
}
func mark(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(res, "Forbidden", http.StatusForbidden)
		return
	}
	enroll_id := req.FormValue("attending")

	// user_id, _ := session.Values["user-id"].(int)

	_, err = db.Exec("Update   enroll set status=1 where id=?", enroll_id)
}
func studentHomePage(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(res, "Forbidden", http.StatusForbidden)
		return
	}
	username, _ := session.Values["username"].(string)
	user_id, _ := session.Values["user-id"].(int)
	data := Classes{}

	rows, _ := db.Query("SELECT id,year,class,stream FROM classes")
	enrolls, _ := db.Query("SELECT enroll.id ,classes.class,classes.year,classes.stream FROM `enroll` INNER JOIN classes on classes.id=enroll.class_id where enroll.student_id=?", user_id)

	var year string
	var class string
	var stream string
	var id int
	for rows.Next() {
		rows.Scan(
			&id,
			&year,
			&class,
			&stream,
		)

		std1 := Class{id, year, class, stream}
		data.Items = append(data.Items, std1)
	}
	for enrolls.Next() {
		enrolls.Scan(
			&id,
			&year,
			&class,
			&stream,
		)
		std1 := Class{id, year, class, stream}
		data.Enrolled = append(data.Enrolled, std1)
	}

	data.Username = username
	t, _ := template.ParseFiles("student.html")

	// standard output to print merged data
	t.Execute(res, data)
}

func main() {
	db, err = sql.Open("mysql", "root:(127.0.0.1:3306)/attendance")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/attendance", attendance)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/addclass", addclass)
	http.HandleFunc("/enroll", enroll)
	http.HandleFunc("/student", studentHomePage)
	http.HandleFunc("/lecturer", lecturer)
	http.HandleFunc("/mark", mark)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.ListenAndServe(":3036", nil)
}
