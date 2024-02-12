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
var statistics *types.SimulateChampionship = &types.SimulateChampionship{} // var globale pour être utilisée dans /statisticsChampionship

func addNewStatistsicsToPrevious(lastStats types.LastChampionshipStatistics) {
	//modifie l'objets "statistics" pour ajouter correctement les dernières stats
	if statistics.LastChampionship == "" {
		//Cas particlier : si on est au premier championnat, le total vaut la dernière simulation
		statistics.TotalStatistics = types.TotalStatistics(lastStats)
		statistics.LastChampionshipStatistics = lastStats
		return
	} else {
		//Cas général
		//Ajout des points des pilotes
		statistics.LastChampionshipStatistics = lastStats
		mapScoreDrivers := make(map[string]int)
		for _, driver := range lastStats.DriversTotalPoints {
			mapScoreDrivers[driver.Driver] = driver.TotalPoints
		}
		for i := range statistics.TotalStatistics.DriversTotalPoints {
			statistics.TotalStatistics.DriversTotalPoints[i].TotalPoints += mapScoreDrivers[statistics.TotalStatistics.DriversTotalPoints[i].Driver]
		}

		//Ajout des points des teams
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

		//Ajout des points des personnalités
		for _, personality := range lastStats.PersonalityAveragePoints {
			var found bool //set to true if personnality is found
			for i := range statistics.TotalStatistics.PersonalityAveragePoints {
				if personality.Personality["Concentration"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Concentration"] &&
					personality.Personality["Aggressivity"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Aggressivity"] &&
					personality.Personality["Confidence"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Confidence"] &&
					personality.Personality["Docility"] == statistics.TotalStatistics.PersonalityAveragePoints[i].Personality["Docility"] {
					found = true
					//On repasse à la somme
					statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints = statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints*float64(statistics.TotalStatistics.PersonalityAveragePoints[i].NbDrivers) + personality.AveragePoints*float64(personality.NbDrivers)
					//Maj du nb de pilotes
					statistics.TotalStatistics.PersonalityAveragePoints[i].NbDrivers += personality.NbDrivers
					//Maj de la moyenne
					statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints = statistics.TotalStatistics.PersonalityAveragePoints[i].AveragePoints / float64(statistics.TotalStatistics.PersonalityAveragePoints[i].NbDrivers)
					break
				}
			}
			if !found {
				//Si la personnalité explorée n'a pas été recensée
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

func getNextChampionshipName(currChampionship string) (string, error) {
	years := strings.Split(currChampionship, "/")
	newFirstYear, err := time.Parse("2006", years[0]) //on souhaite récupérer la première année
	if err != nil {
		return "", err
	}

	newFirstYear = newFirstYear.AddDate(1, 0, 0)
	newLastYear := newFirstYear.AddDate(1, 0, 0)
	return fmt.Sprintf("%d/%d", newFirstYear.Year(), newLastYear.Year()), nil

}

// Lancement d'une simulation d'un championnat et retour de statistiques
func (rsa *RestServer) startSimulation(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête
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

	// Lancement de la simulation
	driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage := s.LaunchSimulation()
	lastChampionshipstatistics := types.NewLastChampionshipStatistics(driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage, nil)

	//Ajoute les nouvelles statistics
	addNewStatistsicsToPrevious(*lastChampionshipstatistics)
	statistics.LastChampionship = championship.Name
	statistics.NbSimulations = nbSimulation

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(statistics)
	w.Write(serial)
}

func (rsa *RestServer) startSimulationRandom(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête
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

	// Lancement de la simulation
	driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage := s.LaunchSimulation()
	lastChampionshipstatistics := types.NewLastChampionshipStatistics(driverLastChampPoints, teamLastChampPoints, personalityLastChampAveragePoints, personnalityAverage, nil)

	//Ajoute les nouvelles statistics
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
