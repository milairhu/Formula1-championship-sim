package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
*
* Return statistics about simulations WITHOUT launching a new simulation
*
 */
func (rsa *RestServer) statisticsChampionship(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /statisticsChampionship")
	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(statistics) //statistics is defined in simulateChampionship.go
	if err != nil {
		panic("Error /statisticsChampionship : can't marshal statistics" + err.Error())
	}
	w.Write(serial)
}

func (rsa *RestServer) statisticsRace(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /statisticsRace")
	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(raceStatistics) //statistics is defined in simulateChampionship.go
	if err != nil {
		panic("Error /statisticsRace : can't marshal statistics" + err.Error())
	}
	w.Write(serial)
}
