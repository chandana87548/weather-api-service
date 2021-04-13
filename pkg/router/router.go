package router

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-akka/configuration"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/weather-api-service/pkg/clients"
	"github.com/weather-api-service/pkg/model"
	"github.com/weather-api-service/pkg/persistence"
	"github.com/weather-api-service/pkg/util"
	"golang.org/x/net/html/charset"
	"net/http"
)

type WeatherResource struct {
	Conf       *configuration.Config
	Repository persistence.JDBCRepository
}

func (wr *WeatherResource) GetWeather(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	//TODO Verify if location is valid string
	currTime := util.GetCurrentTime(0)
	pastTime := util.GetCurrentTime(-300)
	resp := make(map[string]string)
	resp["query_time"] = currTime
	dbValue := wr.Repository.GetByLocation(location, pastTime)
	if dbValue != "" {
		log.Info("data found in db")
		resp["temperature"] = dbValue
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	} else {
		log.Info("Data not found in db, accessing weather api")
		rc := clients.RestHttpClient{
			UserAgent: "go-taxy-client/v0.0.1",
		}
		weatherApiResponse, err := rc.Get(fmt.Sprintf(wr.Conf.GetString("weather-endpoint"), location))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		}
		var weatherResponse model.WeatherResponse
		reader := bytes.NewReader(weatherApiResponse)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		_ = decoder.Decode(&weatherResponse)
		wr.Repository.SetLastAccess(location, currTime, weatherResponse.Temperature)
		resp["temperature"] = weatherResponse.Temperature
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}

}
