package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"os"
)

func get_db_connection() *sql.DB {

	dbuser := os.Getenv("DBUSER")
	dbname := os.Getenv("DBNAME")
	dbpass := os.Getenv("DBPASS")
	dbhost := os.Getenv("DBHOST")
	dbport := os.Getenv("DBPORT")

	dbconn, err := sql.Open("postgres", "dbname="+dbname+" host="+dbhost+" port="+dbport+" user="+dbuser+" password="+dbpass+" sslmode=disable")

	err = dbconn.Ping()

	if err != nil {
		panic(err.Error())
	}

	return dbconn

}

var dbpool *sql.DB //global db pool connection

func main() {

	dbpool = get_db_connection()

	m := martini.Classic()
	m.Map(dbpool)
	m.Get("/movelist", GetMoveList)
	m.Get("/getuser/:username", GetUser)
	m.Get("/getmoveid", GetMoveID)
	m.Post("/registeruser", RegisterUser)
	m.Get("/move/:moveid", GetMove)
	m.Post("/move", AddMove)
	m.Delete("/move/:moveid", DeleteMove)
	m.Put("/move/:moveid", AmendMove)
	m.Run()
}
