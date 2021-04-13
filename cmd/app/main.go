package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/weather-api-service/pkg/persistence"
	. "github.com/weather-api-service/pkg/router"
	"github.com/weather-api-service/pkg/util"
	"net/http"
)

var dbClient = getDBClient()

func getDBClient() *sql.DB {
	fmt.Print("here in get")
	database, err := sql.Open("sqlite3", "./weather-api.db?cache=shared")
	if err != nil {
		panic("Unable to start, failed to connect to db. Error " + err.Error())
	}
	database.SetMaxOpenConns(1)
	return database
}

func init() {
	if dbClient != nil {
		statement1, _ := dbClient.Prepare("DROP TABLE IF EXISTS WEATHER")
		_, _ = statement1.Exec()
		statement, _ := dbClient.Prepare("CREATE TABLE IF NOT EXISTS WEATHER (location TEXT PRIMARY KEY, last_accessed TEXT, temperature TEXT)")
		_, _ = statement.Exec()
	}
}

func main() {

	fmt.Print("started in main")
	conf := util.InitApp("configs/")
	r := mux.NewRouter()
	wr := &WeatherResource{
		Conf:       conf,
		Repository: persistence.JDBCRepository{DBClient: dbClient},
	}
	r.HandleFunc("/temperature/{location}", wr.GetWeather)
	http.Handle("/", r)
	_ = http.ListenAndServe(":8080", nil)
}
