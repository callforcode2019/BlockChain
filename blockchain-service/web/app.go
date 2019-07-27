package web

import (
	"fmt"
	"github.com/BlockChain/blockchain-service/web/controllers"
	"net/http"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/home", app.HomeHandler_1)
	http.HandleFunc("/request", app.RequestHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/login",app.LoginHandler)
	http.HandleFunc("/register",app.RegisterHandler)
	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}