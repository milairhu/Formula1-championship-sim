package restserver

import (
	"log"
	"net/http"
	"sync"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

type RestServer struct {
	sync.Mutex
	addr              string
	pointTabCircuit   []*types.Circuit //circuits
	pointTabTeam      []*types.Team    //current teams
	initPersonalities map[string]types.Personality
}

func NewRestServer(addr string, pointTabCircuit []*types.Circuit, pointTabTeam []*types.Team, personalities map[string]types.Personality) *RestServer {
	return &RestServer{addr: addr, pointTabCircuit: pointTabCircuit, pointTabTeam: pointTabTeam, initPersonalities: personalities}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rsa *RestServer) Start() {
	//initialise statistics
	statistics = &types.SimulateChampionship{}
	for _, team := range rsa.pointTabTeam {
		for _, driver := range team.Drivers {
			statistics.TotalStatistics.DriversTotalPoints = append(statistics.TotalStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
			statistics.LastChampionshipStatistics.DriversTotalPoints = append(statistics.LastChampionshipStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
		}
		statistics.TotalStatistics.TeamsTotalPoints = append(statistics.TotalStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
		statistics.LastChampionshipStatistics.TeamsTotalPoints = append(statistics.LastChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
	}

	//Initialise personnalités dans les statistiques
	statistics.TotalStatistics.PersonalityAveragePoints = make([]*types.PersonalityAveragePoints, 0)
	statistics.LastChampionshipStatistics.PersonalityAveragePoints = make([]*types.PersonalityAveragePoints, 0)
	for indeTeam := range rsa.pointTabTeam {
		for indDriv := range rsa.pointTabTeam[indeTeam].Drivers {
			d := rsa.pointTabTeam[indeTeam].Drivers[indDriv].Id
			var perso types.Personality
			perso.TraitsValue = make(map[string]int)
			perso.TraitsValue["Confidence"] = rsa.initPersonalities[d].TraitsValue["Confidence"]
			perso.TraitsValue["Aggressivity"] = rsa.initPersonalities[d].TraitsValue["Aggressivity"]
			perso.TraitsValue["Docility"] = rsa.initPersonalities[d].TraitsValue["Docility"]
			perso.TraitsValue["Concentration"] = rsa.initPersonalities[d].TraitsValue["Concentration"]
			statistics.TotalStatistics.PersonalityAveragePoints = append(statistics.TotalStatistics.PersonalityAveragePoints, &types.PersonalityAveragePoints{Personality: perso.TraitsValue, AveragePoints: 0, NbDrivers: 0})
			statistics.LastChampionshipStatistics.PersonalityAveragePoints = append(statistics.LastChampionshipStatistics.PersonalityAveragePoints, &types.PersonalityAveragePoints{Personality: perso.TraitsValue, AveragePoints: 0, NbDrivers: 0})

		}

	}

	//Idem pour raceStatistics
	raceStatistics = &types.SimulateRace{}
	raceStatistics.IsLastRace = false
	for _, team := range rsa.pointTabTeam {
		for _, driver := range team.Drivers {
			raceStatistics.RaceStatistics.DriversTotalPoints = append(raceStatistics.RaceStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
			raceStatistics.ChampionshipStatistics.DriversTotalPoints = append(raceStatistics.ChampionshipStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
		}
		raceStatistics.RaceStatistics.TeamsTotalPoints = append(raceStatistics.ChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
		raceStatistics.ChampionshipStatistics.TeamsTotalPoints = append(raceStatistics.ChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
	}

	raceStatistics.RaceStatistics.PersonalityAveragePoints = make([]*types.PersonalityAveragePoints, 0)
	raceStatistics.ChampionshipStatistics.PersonalityAveragePoints = make([]*types.PersonalityAveragePoints, 0)
	for indeTeam := range rsa.pointTabTeam {
		for indDriv := range rsa.pointTabTeam[indeTeam].Drivers {
			d := rsa.pointTabTeam[indeTeam].Drivers[indDriv].Id
			var perso types.Personality
			perso.TraitsValue = make(map[string]int)
			perso.TraitsValue["Confidence"] = rsa.initPersonalities[d].TraitsValue["Confidence"]
			perso.TraitsValue["Aggressivity"] = rsa.initPersonalities[d].TraitsValue["Aggressivity"]
			perso.TraitsValue["Docility"] = rsa.initPersonalities[d].TraitsValue["Docility"]
			perso.TraitsValue["Concentration"] = rsa.initPersonalities[d].TraitsValue["Concentration"]
			raceStatistics.RaceStatistics.PersonalityAveragePoints = append(raceStatistics.RaceStatistics.PersonalityAveragePoints, &types.PersonalityAveragePoints{Personality: perso.TraitsValue, AveragePoints: 0, NbDrivers: 0})
			raceStatistics.ChampionshipStatistics.PersonalityAveragePoints = append(raceStatistics.ChampionshipStatistics.PersonalityAveragePoints, &types.PersonalityAveragePoints{Personality: perso.TraitsValue, AveragePoints: 0, NbDrivers: 0})

		}

	}

	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/simulateRace", rsa.startRaceSimulation)
	mux.HandleFunc("/resetSimulateRace", rsa.resetRaceSimulation)
	mux.HandleFunc("/simulateChampionship", rsa.startSimulation)
	mux.HandleFunc("/simulate100Championships", rsa.start100Simulations)
	mux.HandleFunc("/personalities", rsa.getAndUpdatePersonalities)
	mux.HandleFunc("/statisticsChampionship", rsa.statisticsChampionship)
	mux.HandleFunc("/statisticsRace", rsa.statisticsRace)
	mux.HandleFunc("/reset", rsa.reset)
	corsHandler := corsMiddleware(mux)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        corsHandler,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())

}
