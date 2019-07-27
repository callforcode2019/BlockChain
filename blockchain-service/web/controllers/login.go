package controllers

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"log"
)

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w,r,"index.html",nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		pwd := Sha1maker(password)

		db,err := sql.Open("mysql","root:@/test?charset=utf8")
		if err != nil {
			log.Fatal(err)
		}
		rows,err := db.Query("select * from `user`")
		var name string
		var pass string
		var email string
		for rows.Next() {
			rows.Scan(&name,&pass,&email)
			if username == name {
				if pwd == pass {
					rows.Close()
					cookie := http.Cookie{Name:"username",Value:username,MaxAge:60*60*24}
					http.SetCookie(w,&cookie)
					w.Write([]byte("success"))
					return
				}
			}
		}
		rows.Close()
		w.Write([]byte("no user"))
		return
	}
}

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w,r,"register.html",nil)
	} else if r.Method == "POST" {
		db,err := sql.Open("mysql","root:@/test?charset=utf8")
		if err != nil {
			log.Fatal(err)
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		fmt.Println(username,password,email)
		rows, err := db.Query("SELECT `username` from `user`")
		fmt.Println(err)
		var name string
		for rows.Next() {
			rows.Scan(&name)
			if name == username {
				//username exist
				rows.Close()
				w.Write([]byte("user existed"))
				return
			}
		}
		rows.Close()

		//shal1
		pwd := Sha1maker(password)
		//
		_,err = db.Exec("INSERT into `user` set `username`=?, `password`=?, `email`=?",username,pwd,email)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte("success"))
		return
	}
}

func Sha1maker(pwd string) string {
	checker := sha1.New()
	io.WriteString(checker,pwd)
	hSum := checker.Sum(nil)
	hexString := hex.EncodeToString(hSum)
	return hexString
}