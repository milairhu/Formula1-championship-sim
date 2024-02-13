package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/simulator"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

var nextChampionship = "2023/2024"
var nbSimulation int
var statistics *types.SimulateChampionship = &types.SimulateChampionship{} // global var used in/statisticsChampionship

func addNewStatistsicsToPrevious(lastStats types.LastChampionshipStatistics) {
	//update "statistics" to add correctly latests stats
	if statistics.LastChampionship == "" {
		//Particular case : first championship
		statistics.TotalStatistics = types.TotalStatistics(lastStats)
		statistics.LastChampionshipStatistics = lastStats
		return
	} else {
		// General case
		// Adding drivers points
		statistics.LastChampionshipStatistics = lastStats
		mapScoreDrivers := make(map[string]int)
		for _, driver := range lastStats.DriversTotalPoints {
			mapScoreDrivers[driver.Driver] = driver.TotalPoints
		}
		for i := range statistics.TotalStatistics.DriversTotalPoints {
			statistics.TotalStatistics.DriversTotalPoints[i].TotalPoints += mapScoreDrivers[statistics.TotalStatistics.DriversTotalPoints[i].Driver]
		}

		// Adding teams points
		mapScoreTeams := make(map[string]int)
		for _, team := range lastStats.TeamsTotalPoints {
			mapScoreTeams[team.Team] = team.TotalPoints
		}
		for i := range statistics.TotalStatistics.TeamsTotalPoints {
			statistics.TotalStatistics.TeamsTotalPoints[i].TotalPoints += mapScoreTeams[statistics.TotalStatistics.TeamsTotalPoints[i].Team]
		}

		for personality := range lastStats.PersonalityAverage {
			for i := range statistics.TotalStatistics.PersonalityAverage[personality] {
				statistics.TotalStatistics.PersonalityAverage[personality][i] += lastStats.PersonalityAverage[personality][i]
			}
		}

		// Adding personality points
		for _, personality := range lastStats.PersonalityAveragePoints {
			var found bool //set to true if personnality is found
			for i := range statistics.TotalStatistics.PersonalityAveragePoints {
				if personality.Personality["Concentration"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Concentration"] &&
					personality.Personality["Aggressivity"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Aggressivity"] &&
					personality.Personality["Confidence"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Confidence"] &&
					personality.Personality["Docility"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Docility"] {
					found = true
					// Using sum
					statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints = statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints*float64(statistics.TotalStatistics.PersonalityAveragePoints[i].NbDrivers) + personality.AveragePoints*float64(personality.NbDrivers)
					// Update number of drivers
					statistics.TotalStatistics.PersonalityAveragePoints[i].NbDrivers += personality.NbDrivers
					// Update average
					statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints = statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints / float64(statistics.TotalStatistics.PersonalityAveragePoints[i].NbDrivers)
					break
				}
			}
			if !found {
				//Case of an unexplored personality
				var perso types.Personality
				perso.TraitsValue = make(map[string]int)
				perso.TraitsValue["Confidence"] = personality.Personality["Confidence"]
				perso.TraitsValue["Aggressivity"] = personality.Personality["Aggressivity"]
				perso.TraitsValue["Docility"] = personality.Personality["Docility"]
				perso.TraitsValue["Concentration"] = personality.Personality["Concentration"]
				statistics.TotalStatistics.PersonalityAveragePoints = append(statistics.TotalStatistics.PersonalityAveragePoints, &types.PersonalityAveragePoints{Personality: perso.TraitsValue, AveragePoints: personality.AveragePoints, NbDrivers: personality.NbDrivers})
			}
		}

	}

}

// Generate next championship name depending on current championship
func getNextChampionshipName(currChampionship string) (string, error) {
	years := strings.Split(currChampionship, "/")
	newFirstYear, err := time.Parse("2006", years[0])
	if err != nil {
		return "", err
	}

	newFirstYear = newFirstYear.AddDate(1, 0, 0)
	newLastYear := newFirstYear.AddDate(1, 0, 0)
	return fmt.Sprintf("%d/%d", newFirstYear.Year(), newLastYear.Year()), nil

}

// Launch simulation of a championship and return statistics
func (rsa *RestServer) startSimulation(w http.ResponseWriter, r *http.Request) {

	// check method of the request
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /simulateChampionship")

	championship := types.NewChampionship(nextChampionship, nextChampionship, rsa.pointTabCircuit, rsa.pointTabTeam)
	ch, err := getNextChampionshipName(nextChampionship)
	if err != nil {
		panic("Error /simulateChampionship : can't create new Dates" + err.Error())
	}
	nextChampionship = ch

	s := simulator.NewSimulator([]types.Championship{*championship})
	nbSimulation += 1

	// Launch simmulation
	driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage := s.LaunchSimulation()
	lastChampionshipstatistics := types.NewLastChampionshipStatistics(driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage, nil)

	//Add new statistics
	addNewStatistsicsToPrevious(*lastChampionshipstatistics)
	statistics.LastChampionship = championship.Name
	statistics.NbSimulations = nbSimulation

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(statistics)
	w.Write(serial)
}

// Launch simulation with random personalities
func (rsa *RestServer) startSimulationRandom(w http.ResponseWriter, r *http.Request) {

	// check method of the request
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /simulateChampionship")

	championship := types.NewChampionshipRandom(nextChampionship, nextChampionship, rsa.pointTabCircuit, rsa.pointTabTeam)
	ch, err := getNextChampionshipName(nextChampionship)
	if err != nil {
		panic("Error /simulateChampionship : can't create new Dates" + err.Error())
	}
	nextChampionship = ch

	s := simulator.NewSimulator([]types.Championship{*championship})
	nbSimulation += 1

	// Launch simmulation
	driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage := s.LaunchSimulation()
	lastChampionshipstatistics := types.NewLastChampionshipStatistics(driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage, nil)

	//Add new statistics
	addNewStatistsicsToPrevious(*lastChampionshipstatistics)
	statistics.LastChampionship = championship.Name
	statistics.NbSimulations = nbSimulation

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(statistics)
	w.Write(serial)
}

func (rsa *RestServer) start100Simulations(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 100; i++ {
		rsa.startSimulationRandom(w, r)
	}
}
