package main

import (
	"database/sql"
	"fmt"
	"os"
	//"regexp"

	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


var (
	tpl *template.Template
	db *sql.DB )

type userData struct{
	firstname	string
	lastname	string
	email		string
	password	string
}

func main() {
	fmt.Println("project starting")
	createDB()
	handleRequests()
}

func createDB() {
	var err error
	godotenv.Load("/Users/godswill/Documents/GitHub/tutorialEdge/envVars.env")
	DB_PASS := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT") 
	db, err = sql.Open("mysql", "root:"+ DB_PASS + "@tcp(localhost:" + DB_PORT + ")/register")

	if err != nil {
		fmt.Println("Error validating SQL.Open")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting with ping")
		panic(err)
	}

	fmt.Println("Successfully connected to MySQL")
}

func init() {
	 tpl, _ = template.ParseGlob("templates/*.html")	
	//this is for using standard html files and not gohtml
	//tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
func handleRequests() {
	http.HandleFunc("/", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/users", displayUser)
	fmt.Println("Starting server on PORT 8080... ")
	http.ListenAndServe(":8080", nil)
	
}

func signup(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "default.html", nil)
		return
	}
	r.ParseForm()

	var err error 

	user := userData{
		firstname: r.FormValue("fname"),
		lastname: r.FormValue("lname"),
		email: r.FormValue("emailAddress"),
		password: r.FormValue("pass"),
	}
	if user.firstname == "" || user.lastname == "" || user.email == "" || user.password == "" {
		fmt.Println("Empty field spotted")
		tpl.ExecuteTemplate(w, "default.html", nil)
		return 
	}
	
	// if !validateEmail(user.email) {
	// 	panic(err)
	// }
	fmt.Println("before insert")

	var ins *sql.Stmt
	ins, err = db.Prepare("INSERT INTO `register`.`users` (`firstname`,`lastname`,`email`,`password`) VALUES (?,?,?,?);")
	if err != nil {
		panic(err)
	}
	fmt.Println("after insert but before ins.Close() ")
	defer ins.Close()
	res, err := ins.Exec(user.firstname, user.lastname, user.email, user.password)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1{
		fmt.Println("Error Inserting Row:", err)
		tpl.ExecuteTemplate(w, "default.html", "Error inserting data, please check all fields")
		return
	}
}

func login (w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
}

// func validateEmail(email string) bool {
// pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
// emailRegex := regexp.MustCompile(pattern)
// return emailRegex.MatchString(user.email) 
// }

func displayUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****browseHandler running*****")
	stmt := "SELECT * FROM users"
	tpl.ExecuteTemplate(w, "users.html", nil)
	rows, err := db.Query(stmt)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

}