package main

import (
	"database/sql"
	"github.com/Mgeorg1/go-todo-list/app"
	"github.com/Mgeorg1/go-todo-list/db"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	conn, err := sql.Open("postgres", "postgresql://root:Ds12345!@localhost:5432/to_do_list?sslmode=disable")
	if err != nil {
		log.Fatalln(err.Error())
	}
	dbStore := db.NewStore(conn)
	server, err := app.CreateServer(&dbStore)
	if err != nil {
		log.Fatalln(err.Error())
	}
	http.HandleFunc("/", server.ListTasks)
	http.HandleFunc("/add", server.AddTodo)
	http.HandleFunc("/delete/", server.DeleteTask)
	http.HandleFunc("/done/", server.SetDone)
	http.HandleFunc("/search/", server.Search)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	log.Println("Starting server on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
