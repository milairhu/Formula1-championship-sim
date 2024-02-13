package simulator

import (
	"log"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

type Simulator struct {
	Championships []types.Championship // championships to simulate. Each championship contains a list of circuits and teams
}

func NewSimulator(championships []types.Championship) *Simulator {
	c := make([]types.Championship, len(championships))
	copy(c, championships)

	return &Simulator{
		Championships: c,
	}
}

func (s *Simulator) LaunchSimulation() ([]*types.DriverTotalPoints, []*types.TeamTotalPoints, []*types.PersonalityAveragePoints, map[string]map[int]float64) {

	var driverTotalPoints []*types.DriverTotalPoints
	var teamTotalPoints []*types.TeamTotalPoints
	var personalityAveragePoints []*types.PersonalityAveragePoints
	var personnalityAverage map[string]map[int]float64

	log.Println("Launching new simulation...")
	for _, championship := range s.Championships {
		// Simulate each championship
		log.Printf("Launching new championship : %s...\n", championship.Name)
		for i, circuit := range championship.Circuits {
			// Simulate each race
			// Step 1 : create race
			var id = circuit.Name + " " + championship.Name

			var date = time.Now()
			if i != 0 {
				date = championship.Races[i-1].Date.AddDate(0, 0, 14)
			}
			var meteo = circuit.GenerateMeteo()
			new_Race := types.NewRace(id, circuit, date, championship.Teams, meteo)

			//Step 2 (main step) : simulate
			pointsMap, err := new_Race.SimulateRace()
			if err != nil {
				log.Printf("Error simulation race %s : %s\n", new_Race.Id, err.Error())
			}

			// Record points earned by each driver
			for indT := range championship.Teams {
				for indD := range championship.Teams[indT].Drivers {
					championship.Teams[indT].Drivers[indD].ChampionshipPoints += pointsMap[championship.Teams[indT].Drivers[indD].Id]
				}
			}
			//Step 3 : add the race to the championship
			championship.Races[i] = *new_Race
		}
		// Display championship rank
		log.Printf("\n\n===== Championship ranking %s =====\n", championship.Name)
		teamTotalPoints = championship.DisplayTeamRank()
		driverTotalPoints, personalityAveragePoints, personnalityAverage = championship.DisplayDriverRank()

	}
	return driverTotalPoints, teamTotalPoints, personalityAveragePoints, personnalityAverage
}
