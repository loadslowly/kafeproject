package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	Id, StatusReg, StatusAdmin uint16
	Username, Email, Password  string
}

type PostInfo struct {
	Id                  uint16
	Title, Beat, Poster string
}

var posts = []PostInfo{}
var accounts = []UserInfo{}
var idAccount int

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "mysql:@tcp(127.0.0.1:3306)/beatstore")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT username, password, email, statusreg, statusadmin FROM `users`")
	if err != nil {
		panic(err)
	}

	accounts = []UserInfo{}
	for res.Next() {
		var account UserInfo
		err := res.Scan(&account.Username, &account.Password, &account.Email, &account.StatusReg, &account.StatusAdmin)
		if err != nil {
			panic(err)
		}

		accounts = append(accounts, account)
	}

	t.ExecuteTemplate(w, "index", nil)
}

func sign(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/sign.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "sign", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/register.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "register", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password_conf := r.FormValue("password_conf")
	statusadmin := "0"
	statusreg := "0"

	if username == "" || email == "" || password == "" || password_conf == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else if password != password_conf {
		fmt.Fprintf(w, "Пароли не совпдают")
	} else {
		db, err := sql.Open("mysql", "mysql:@tcp(127.0.0.1:3306)/beatstore")
		if err != nil {
			panic(err)
		}

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `users` (`username`, `email`, `password`, `statusreg`, `statusadmin`) VALUES('%s', '%s', '%s', '%s', '%s')", username, email, password, statusreg, statusadmin))
		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/sign", http.StatusSeeOther)

		defer db.Close()
	}
}

func save_feed(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db, err := sql.Open("mysql", "mysql:@tcp(127.0.0.1:3306)/beatstore")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT username, email, password, statusreg, statusadmin FROM `users`")
	if err != nil {
		panic(err)
	}

	accounts = []UserInfo{}
	for res.Next() {
		var account UserInfo
		err := res.Scan(&account.Username, &account.Email, &account.Password, &account.StatusReg, &account.StatusAdmin)
		if err != nil {
			panic(err)
		}

		accounts = append(accounts, account)
	}

	for i := 0; i < len(accounts); i++ {
		if username == accounts[i].Username && password == accounts[i].Password {
			idAccount = i
			accounts[idAccount].StatusReg = 1
			if accounts[i].Username == "loadslowly" {
				accounts[idAccount].StatusAdmin = 1
			}
		}
	}

	if accounts[idAccount].StatusReg == 1 {
		defer res.Close()
		http.Redirect(w, r, "/feed", http.StatusSeeOther)

		defer db.Close()
	} else if username == "" || password == "" {
		fmt.Fprint(w, "Вы не заполнили все поля")
	} else {
		fmt.Fprint(w, "Неверный логин или пароль")
	}
}

func feed(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/feed.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "feed", accounts[idAccount])
}

func profile(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/profile.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "profile", accounts[idAccount])
}

func contact(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/contact.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "contact", accounts[idAccount])
}

func recipes(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/recipes.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "recipes", accounts[idAccount])
}

func save_post(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	beat := r.FormValue("beat")
	poster := r.FormValue("poster")

	db, err := sql.Open("mysql", "mysql:@tcp(127.0.0.1:3306)/beatstore")
	if err != nil {
		panic(err)
	}

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `post` (`title`, `beat`, `poster`) VALUES('%s', '%s', '%s')", title, beat, poster))
	if err != nil {
		panic(err)
	}

	defer insert.Close()

	http.Redirect(w, r, "/feed", http.StatusSeeOther)

	defer db.Close()
}

func newpost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/newpost.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "newpost", accounts[idAccount])
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/contact", contact).Methods("GET")
	rtr.HandleFunc("/recipes", recipes).Methods("GET")
	rtr.HandleFunc("/sign", sign).Methods("GET")
	rtr.HandleFunc("/register", register).Methods("GET")
	rtr.HandleFunc("/newpost", newpost).Methods("GET")

	rtr.HandleFunc("/feed", feed).Methods("GET")
	rtr.HandleFunc("/profile/{key}", profile).Methods("GET")

	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/save_feed", save_feed).Methods("POST")
	rtr.HandleFunc("/save_post", save_post).Methods("POST")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", rtr)
	http.ListenAndServe(":9090", nil)
}

func main() {
	handleFunc()
}
