package simulator

import (
	"log"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

type Simulator struct {
	Championships []types.Championship //Contient les championnats à simuler. Nous n'en passons qu'un en se servant d'un simulator différent pour chaque championnat
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

	log.Println("Lancement d'une nouvelle simulation...")
	for _, championship := range s.Championships {
		//On simule chaque championnat
		log.Printf("Lancement d'un nouveau championnat : %s...\n", championship.Name)
		for i, circuit := range championship.Circuits {
			//On simule chaque course
			//Etape 1 : on crée la course
			var id = circuit.Name + " " + championship.Name

			var date = time.Now()
			if i != 0 {
				date = championship.Races[i-1].Date.AddDate(0, 0, 14)
			}
			var meteo = circuit.GenerateMeteo()
			new_Race := types.NewRace(id, circuit, date, championship.Teams, meteo)

			//Etape 2 (la principale) : on joue la course
			pointsMap, err := new_Race.SimulateRace()
			if err != nil {
				log.Printf("Erreur simulation cours %s : %s\n", new_Race.Id, err.Error())
			}

			//On enregistre les points gagnés par chaque pilote
			for indT := range championship.Teams {
				for indD := range championship.Teams[indT].Drivers {
					championship.Teams[indT].Drivers[indD].ChampionshipPoints += pointsMap[championship.Teams[indT].Drivers[indD].Id]
				}
			}
			//Etape 3 : on ajoute la course au championnat
			championship.Races[i] = *new_Race
		}
		//On affiche le classement du championnat
		log.Printf("\n\n===== Classements du championnat %s =====\n", championship.Name)
		teamTotalPoints = championship.DisplayTeamRank()
		driverTotalPoints, personalityAveragePoints, personnalityAverage = championship.DisplayDriverRank()

	}
	return driverTotalPoints, teamTotalPoints, personalityAveragePoints, personnalityAverage
}
