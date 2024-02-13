package restserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/milairhu/Formula1-championship-sim/back-end/types"
)

var i = 0
var championship *types.Championship
var firstSimulation = true
var raceStatistics *types.SimulateRace = &types.SimulateRace{}

func (rsa *RestServer) resetRaceSimulation(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /resetSimulateRace")

	firstSimulation = true
	i = 0
	raceStatistics = &types.SimulateRace{}

	for _, team := range rsa.pointTabTeam {
		for _, driver := range team.Drivers {
			raceStatistics.RaceStatistics.DriversTotalPoints = append(raceStatistics.RaceStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
			raceStatistics.ChampionshipStatistics.DriversTotalPoints = append(raceStatistics.ChampionshipStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
		}
		raceStatistics.RaceStatistics.TeamsTotalPoints = append(raceStatistics.ChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
		raceStatistics.ChampionshipStatistics.TeamsTotalPoints = append(raceStatistics.ChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
	}
	raceStatistics.IsLastRace = false
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

	serial, err := json.Marshal(raceStatistics) //statistics is defined in simulateChampionship.go
	if err != nil {
		panic("Error /reset : can't marshal statistics" + err.Error())
	}
	w.Write(serial)

	w.WriteHeader(http.StatusOK)
}

func (rsa *RestServer) startRaceSimulation(w http.ResponseWriter, r *http.Request) {
	if firstSimulation { //Initialize championship if this is the first race
		championship = types.NewChampionship(nextChampionship, nextChampionship, rsa.pointTabCircuit, rsa.pointTabTeam)
		firstSimulation = false
	}

	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /simulateRace")

	//Simulation of race i

	if i <= len(championship.Circuits) {

		//Creating race
		var id = championship.Circuits[i].Name + " " + championship.Name
		raceStatistics.Championship = championship.Name
		raceStatistics.Race = championship.Circuits[i].Name

		var date = time.Now()
		if i != 0 {
			date = championship.Races[i-1].Date.AddDate(0, 0, 14)
		}
		var meteo = championship.Circuits[i].GenerateMeteo()
		new_Race := types.NewRace(id, championship.Circuits[i], date, championship.Teams, meteo)

		//Simulation of the race
		pointsMap, err := new_Race.SimulateRace()
		if err != nil {
			log.Printf("Erreur simulation cours %s : %s\n", new_Race.Id, err.Error())
		}

		// Adding points to the championship
		for indT := range championship.Teams {
			for indD := range championship.Teams[indT].Drivers {
				championship.Teams[indT].Drivers[indD].ChampionshipPoints += pointsMap[championship.Teams[indT].Drivers[indD].Id]
			}
		}

		// Points of drivers for the race
		driversRankTab := make([]*types.DriverTotalPoints, 0)

		// personality Average
		personalityRankTab := make([]*types.PersonalityAveragePoints, 0)
		personnalityAverage := make(map[string]map[int]float64)
		nb := make(map[string]map[int]int)

		for _, driver := range new_Race.FinalResult {
			driverRank := types.NewDriverTotalPoints(driver.Lastname, pointsMap[driver.Id])
			driversRankTab = append(driversRankTab, driverRank)

			for personnality, level := range driver.Personality.TraitsValue {
				if _, ok := personnalityAverage[personnality]; !ok {
					personnalityAverage[personnality] = make(map[int]float64)
					nb[personnality] = make(map[int]int)
				}
				personnalityAverage[personnality][level] += float64(pointsMap[driver.Id])
				nb[personnality][level] += 1
			}

			var found bool
			for indPers := range personalityRankTab {
				if personalityRankTab[indPers].Personality["Aggressivity"] == driver.Personality.TraitsValue["Aggressivity"] &&
					personalityRankTab[indPers].Personality["Concentration"] == driver.Personality.TraitsValue["Concentration"] &&
					personalityRankTab[indPers].Personality["Confidence"] == driver.Personality.TraitsValue["Confidence"] &&
					personalityRankTab[indPers].Personality["Docility"] == driver.Personality.TraitsValue["Docility"] {
					personalityRankTab[indPers].AveragePoints += float64(pointsMap[driver.Id])
					personalityRankTab[indPers].NbDrivers += 1
					found = true
					break
				}
			}
			if !found {
				var perso types.Personality
				perso.TraitsValue = make(map[string]int)
				perso.TraitsValue["Confidence"] = driver.Personality.TraitsValue["Confidence"]
				perso.TraitsValue["Aggressivity"] = driver.Personality.TraitsValue["Aggressivity"]
				perso.TraitsValue["Docility"] = driver.Personality.TraitsValue["Docility"]
				perso.TraitsValue["Concentration"] = driver.Personality.TraitsValue["Concentration"]
				personalityRank := types.NewPersonalityAveragePoints(perso.TraitsValue, pointsMap[driver.Id], 1)
				personalityRankTab = append(personalityRankTab, personalityRank)
			}

		}

		// Compute means
		for indPers := range personalityRankTab {
			if personalityRankTab[indPers].NbDrivers > 1 {
				personalityRankTab[indPers].AveragePoints = personalityRankTab[indPers].AveragePoints / float64(personalityRankTab[indPers].NbDrivers)
			}

		}

		for personnality := range personnalityAverage {
			for i := 1; i < 6; i++ {
				if _, ok := personnalityAverage[personnality][i]; !ok {
					personnalityAverage[personnality][i] = 0
					nb[personnality][i] = 1
				}
			}
		}

		//Compute personality means
		for personnality, level := range personnalityAverage {
			for level, points := range level {
				personnalityAverage[personnality][level] = points / float64(nb[personnality][level])
			}
		}

		raceStatistics.RaceStatistics.DriversTotalPoints = driversRankTab
		raceStatistics.RaceStatistics.PersonalityAveragePoints = personalityRankTab
		raceStatistics.RaceStatistics.PersonalityAverage = personnalityAverage

		// Team points for current race
		teamsRankTab := make([]*types.TeamTotalPoints, 0)
		for _, team := range new_Race.Teams {
			teamPoints := 0
			for _, driver := range team.Drivers {
				teamPoints += pointsMap[driver.Id]
			}
			teamRank := types.NewTeamTotalPoints(team.Name, teamPoints)
			teamsRankTab = append(teamsRankTab, teamRank)
		}
		raceStatistics.RaceStatistics.TeamsTotalPoints = teamsRankTab

		// Championship stats
		var champDriverTotalPoints []*types.DriverTotalPoints
		var champTeamTotalPoints []*types.TeamTotalPoints
		var champPersonalityAveragePoints []*types.PersonalityAveragePoints
		var champPersonnalityAverage map[string]map[int]float64

		champTeamTotalPoints = championship.DisplayTeamRank()
		champDriverTotalPoints, champPersonalityAveragePoints, champPersonnalityAverage = championship.DisplayDriverRank()

		raceStatistics.ChampionshipStatistics.DriversTotalPoints = champDriverTotalPoints
		raceStatistics.ChampionshipStatistics.PersonalityAveragePoints = champPersonalityAveragePoints
		raceStatistics.ChampionshipStatistics.TeamsTotalPoints = champTeamTotalPoints
		raceStatistics.ChampionshipStatistics.PersonalityAverage = champPersonnalityAverage

		var raceHighlightsTab []types.RaceHighlight

		for _, v := range new_Race.HighLigths {
			highlight := types.NewRaceHighlight(v.Description, v.Type)
			raceHighlightsTab = append(raceHighlightsTab, highlight)
		}
		raceStatistics.Highlights = raceHighlightsTab

		// Inc. races count
		i++
	}

	if i == len(championship.Circuits) {
		raceStatistics.IsLastRace = true
	}

	w.WriteHeader(http.StatusOK)
	serial, err := json.Marshal(raceStatistics)
	if err != nil {
		panic("Error /simulateRace : can't marshal statistics" + err.Error())
	}
	w.Write(serial)
}
