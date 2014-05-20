package main

import (
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	_ "github.com/lib/pq"
	"net/http"
)

type move struct {
	Name       string
	Category   string
	Tags       string
	Video_link string
}

type get_move_list_resp struct {
	Moves []move
}

type user struct {
	Username      string
	Password_hash string
	Email         string
}

func get_db_connection() *sql.DB {

	dbpool, err := sql.Open("postgres", "dbname=MyDanceMoves user=move_admin password=password sslmode=disable")

	err = dbpool.Ping()

	if err != nil {
		panic(err.Error())
	}

	return dbpool

}

func register_user(w *rest.ResponseWriter, req *rest.Request) {

	newuser := user{}

	err := req.DecodeJsonPayload(&newuser)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var insert_statement string = "insert into user_accounts (username, password_hash, email) values ($1, $2, $3)"

	stmt, err := db.Prepare(insert_statement)

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(newuser.Username, newuser.Password_hash, newuser.Email)

	if err != nil || res == nil {
		panic(err.Error())
	}

	stmt.Close()
}

func get_move_list(w *rest.ResponseWriter, req *rest.Request) {

	var move_list []move

	rows, statementError := db.Query("SELECT trim(name), trim(category), trim(tags), trim(video_link) FROM moves")

	if statementError != nil {
		panic(statementError)
	}

	defer rows.Close()

	move_list = make([]move, 0)

	for rows.Next() {

		var name, category, tags, video_link string

		err := rows.Scan(&name, &category, &tags, &video_link)

		if err != nil {
			fmt.Printf("rows.Scan error: %v\n", err)
		}

		current_move := move{name, category, tags, video_link}

		move_list = append(move_list, current_move)

	}

	resp := new(get_move_list_resp)
	resp.Moves = move_list
	w.WriteJson(&resp)
}

var db *sql.DB //global db pool connection

func main() {

	db = get_db_connection()

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/movelist", get_move_list},
		rest.Route{"POST", "/registeruser", register_user},
	)
	http.ListenAndServe(":8081", &handler)
}
