package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB

const (
	DBUser     = "root"
	DBPassword = "0000"
	DBName     = "user_db"
)

func InitDB() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@/%s", DBUser, DBPassword, DBName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var err error
	db, err = InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/signup", SignUpPage)
	http.HandleFunc("/users", GetUsers)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server started on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/login.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == http.MethodPost {
        // Ambil data formulir dari r.PostForm
        username := r.PostFormValue("username")
        password := r.PostFormValue("password")

        // Validasi login
        if isValidLogin(username, password) {
            // Login berhasil, redirect atau berikan respons sesuai
            w.Write([]byte("Login berhasil!"))
        } else {
            // Login gagal, berikan respons sesuai
            w.Write([]byte("Login gagal. Cek username dan password."))
        }
    }
}

func isValidLogin(username, password string) bool {
    // Lakukan query ke database untuk memeriksa apakah username dan password valid
    row := db.QueryRow("SELECT username FROM users WHERE username=? AND password=?", username, password)

    var validUsername string
    if err := row.Scan(&validUsername); err != nil {
        // Kesalahan atau username dan password tidak valid
        return false
    }

    // Username dan password valid
    return true
}




func SignUpPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/signup.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == http.MethodPost {
        // Ambil data formulir dari r.PostForm
        username := r.PostFormValue("username")
        password := r.PostFormValue("password")
        confirmPassword := r.PostFormValue("confirm-password")

        // Validasi pendaftaran
        if password == confirmPassword {
            // Pendaftaran berhasil, simpan data ke database
            if err := RegisterUser(username, password); err != nil {
                http.Error(w, "Gagal mendaftarkan pengguna", http.StatusInternalServerError)
                log.Fatal(err)
                return
            }

            w.Write([]byte("Pendaftaran berhasil!"))
        } else {
            // Pendaftaran gagal, berikan respons sesuai
            w.Write([]byte("Pendaftaran gagal. Password tidak sesuai."))
        }
    }
}

func RegisterUser(username, password string) error {
    // Lakukan penyisipan data ke database
    _, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
    return err
}






func GetUsers(w http.ResponseWriter, r *http.Request) {
    // Fetch user data from the database
    users, err := RetrieveUsers()
    if err != nil {
        http.Error(w, "Gagal mengambil data pengguna", http.StatusInternalServerError)
        log.Fatal(err)
        return
    }

    // Send user data as JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func RetrieveUsers() ([]User, error) {
    // Lakukan query ke database
    rows, err := db.Query("SELECT id, username, password FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, &user.Password)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}