package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

type user struct {
	Id       string
	Password string
	Email    string
}

type get_user_resp struct {
	User   user
	Result string
}

func RegisterUser(res http.ResponseWriter, req *http.Request, db *sql.DB) {

	newuser := user{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&newuser)
	if err != nil {
		panic(err.Error())
	}

	var insert_statement string = "insert into user_accounts (username, password_hash, email) values ($1, $2, $3)"

	stmt, err := db.Prepare(insert_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Exec(newuser.Id, newuser.Password, newuser.Email)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	stmt.Close()
}

func GetUser(res http.ResponseWriter, req *http.Request, parms martini.Params, db *sql.DB) {

	username := parms["username"]

	row := db.QueryRow("SELECT trim(email), password_hash FROM user_accounts WHERE username = $1", username)

	var result, email, password_hash string

	var user_info user

	err := row.Scan(&email, &password_hash)

	if err != nil {
		fmt.Printf("rows.Scan error: %v\n", err)
		result = "no user found"
	} else {
		user_info = user{username, password_hash, email}
		result = "user found"
	}

	resp := new(get_user_resp)
	resp.User = user_info
	resp.Result = result

	//convert resp struct to json
	jsonResp, encError := json.Marshal(resp)
	if encError != nil {
		fmt.Println(encError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	res.Write(jsonResp)

}

type move struct {
	Name       string
	Category   string
	Tags       string
	Video_link string
	Username   string
}

type get_move_list_resp struct {
	Moves []move
}

func GetMoveList(res http.ResponseWriter, req *http.Request, db *sql.DB) {

	username := req.FormValue("username")

	var move_list []move

	rows, statementError := db.Query("SELECT trim(name), trim(category), trim(tags), trim(video_link) FROM moves WHERE username = $1", username)

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

		current_move := move{name, category, tags, video_link, username}

		move_list = append(move_list, current_move)

	}

	resp := new(get_move_list_resp)
	resp.Moves = move_list

	//convert resp struct to json
	jsonResp, encError := json.Marshal(resp)
	if encError != nil {
		fmt.Println(encError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	res.Write(jsonResp)
}

func AddMove(res http.ResponseWriter, req *http.Request, db *sql.DB) {

	decoder := json.NewDecoder(req.Body)

	newmove := move{}

	err := decoder.Decode(&newmove)
	if err != nil {
		panic(err.Error())
	}

	var insert_statement string = "insert into moves (name, category, tags, video_link, username) values ($1, $2, $3, $4, $5)"

	stmt, err := db.Prepare(insert_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Exec(newmove.Name, newmove.Category, newmove.Tags, newmove.Video_link, newmove.Username)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	stmt.Close()

	res.Write([]byte("Move added"))
}

func get_db_connection() *sql.DB {

	dbuser := os.Getenv("DBUSER")
	dbname := os.Getenv("DBNAME")
	dbpass := os.Getenv("DBPASS")

	dbconn, err := sql.Open("postgres", "dbname="+dbname+" host=localhost port=5432 user="+dbuser+" password="+dbpass+" sslmode=disable")

	err = dbconn.Ping()

	if err != nil {
		panic(err.Error())
	}

	return dbconn

}

// DB Returns a martini.Handler
//func DB() martini.Handler {
//	session, err := get_db_connection()
//	if err != nil {
//		panic(err)
//	}

//	return func(c martini.Context) {
//		s := session.Clone()
//		c.Map(s.DB("advent"))
//		defer s.Close()
//		c.Next()
//	}
//}

var dbpool *sql.DB //global db pool connection

func main() {

	dbpool = get_db_connection()

	m := martini.Classic()
	m.Map(dbpool)
	m.Get("/movelist", GetMoveList)
	m.Get("/getuser/:username", GetUser)
	m.Post("/addmove", AddMove)
	m.Post("/registeruser", RegisterUser)
	m.Run()
}
