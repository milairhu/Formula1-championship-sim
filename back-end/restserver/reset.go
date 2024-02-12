package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

func (rsa *RestServer) reset(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /reset")
	// reset des variables globales
	nextChampionship = "2023/2024"
	nbSimulation = 0
	statistics = &types.SimulateChampionship{}
	//Remet les pilotes à 0
	for indTeam := range rsa.pointTabTeam {
		for indDriv := range rsa.pointTabTeam[indTeam].Drivers {
			driver := rsa.pointTabTeam[indTeam].Drivers[indDriv]
			statistics.TotalStatistics.DriversTotalPoints = append(statistics.TotalStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
			statistics.LastChampionshipStatistics.DriversTotalPoints = append(statistics.LastChampionshipStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
		}
		statistics.TotalStatistics.TeamsTotalPoints = append(statistics.TotalStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: rsa.pointTabTeam[indTeam].Name, TotalPoints: 0})
		statistics.LastChampionshipStatistics.TeamsTotalPoints = append(statistics.LastChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: rsa.pointTabTeam[indTeam].Name, TotalPoints: 0})
	}
	//On remet les personnalités à 0 dans le tableau d'équipes
	for indTeam := range rsa.pointTabTeam {
		for indDriver := 0; indDriver < 2; indDriver++ {
			var d string = rsa.pointTabTeam[indTeam].Drivers[indDriver].Id
			var perso types.Personality
			//Obligé de faire ça, sinon effet de bord sur initPersonalities
			perso.TraitsValue = make(map[string]int)
			perso.TraitsValue["Confidence"] = rsa.initPersonalities[d].TraitsValue["Confidence"]
			perso.TraitsValue["Aggressivity"] = rsa.initPersonalities[d].TraitsValue["Aggressivity"]
			perso.TraitsValue["Docility"] = rsa.initPersonalities[d].TraitsValue["Docility"]
			perso.TraitsValue["Concentration"] = rsa.initPersonalities[d].TraitsValue["Concentration"]
			rsa.pointTabTeam[indTeam].Drivers[indDriver].Personality = perso
		}
	}
	// On remet les personnalité à 0 dans les statistiques
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
	serial, err := json.Marshal(statistics) //statistics is defined in simulateChampionship.go
	if err != nil {
		panic("Error /reset : can't marshal statistics" + err.Error())
	}
	w.Write(serial)
}
