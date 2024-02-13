package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/milairhu/Formula1-championship-sim/back-end/types"
)

const (
	// Path to the JSON file containing the circuits
	CIRCUITS_PATH = "instances/circuits/inst-circuits.json"
	TEST_PATH     = "instances/circuits/inst-test.json"
	// Path to the JSON file containing the teams
	TEAMS_PATH = "instances/teams/inst-teams.json"
)

type TeamsJSON struct {
	Name    string
	Country string
	Level   int
	Drivers []DriverJSON
}

type DriverJSON struct {
	FirstName   string
	LastName    string
	Country     string
	Level       int
	Personality PersonalityJSON
}

type PersonalityJSON struct {
	Aggressivity  int
	Confidence    int
	Docility      int
	Concentration int
}

func ReadCircuit() ([]types.Circuit, error) {
	// Ouvrir et lire le fichier JSON
	file, err := os.Open(CIRCUITS_PATH)
	if err != nil {
		log.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	defer file.Close()

	var circuits []types.Circuit

	// Décoder le fichier JSON dans la structure de données
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&circuits); err != nil {
		log.Println("Erreur lors de la lecture du fichier JSON :", err)
		return nil, err
	}

	for i := 0; i < len(circuits); i++ {
		//On donne un Id à la portion
		circuits[i].Id = fmt.Sprintf("circuit-%d", i)
		for j := 0; j < len(circuits[i].Portions); j++ {
			// On spécifie le type de chaque portion
			if len(circuits[i].Portions[j].Id) < len("turn") {
			} else if circuits[i].Portions[j].Id[:len("turn")] == "turn" {
				circuits[i].Portions[j].Type = types.TURN
			} else {
				circuits[i].Portions[j].Type = types.STRAIGHT
			}

			//On spécifie de quelle portion celle-ci est la suivante
			if j == 0 {
				circuits[i].Portions[len(circuits[i].Portions)-1].NextPortion = &(circuits[i].Portions[j])
			} else {
				circuits[i].Portions[j-1].NextPortion = &(circuits[i].Portions[j])
			}
		}

	}

	return circuits, nil
}

func ReadTeams() ([]types.Team, map[string]types.Personality, error) {
	// Ouvrir et lire le fichier JSON
	file, err := os.Open(TEAMS_PATH)
	if err != nil {
		log.Println("Erreur lors de l'ouverture du fichier pour lecture des équipes :", err)
		return nil, nil, err
	}

	teamsJSON := make([]TeamsJSON, 0)
	teams := make([]types.Team, 0)
	// Décoder le fichier JSON pour créer les équipes
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&teams); err != nil {
		log.Println("Erreur lors de la lecture du fichier JSON pour les équipes :", err)
		return nil, nil, err
	}
	file.Close()
	file, err = os.Open(TEAMS_PATH)
	if err != nil {
		log.Println("Erreur lors de l'ouverture du fichier pour lecture des personnalités:", err)
		return nil, nil, err
	}
	decoder = json.NewDecoder(file)
	if err := decoder.Decode(&teamsJSON); err != nil {
		log.Println("Erreur lors de la lecture du fichier JSON pour les personnalités :", err)
		return nil, nil, err
	}

	//Ajout des personnalités aux pilotes
	for i, team := range teams {
		for j := range team.Drivers {
			m := make(map[string]int)
			m["Aggressivity"] = teamsJSON[i].Drivers[j].Personality.Aggressivity
			m["Confidence"] = teamsJSON[i].Drivers[j].Personality.Confidence
			m["Docility"] = teamsJSON[i].Drivers[j].Personality.Docility
			m["Concentration"] = teamsJSON[i].Drivers[j].Personality.Concentration
			teams[i].Drivers[j].Personality = *types.NewPersonality(m)
		}
	}

	//Ajout d'Id aux pilotes et aux team
	for i, team := range teams {
		teams[i].Id = fmt.Sprintf("team-%d", i)
		for j := range team.Drivers {
			teams[i].Drivers[j].Id = fmt.Sprintf("driver-%d-%d", i, j)
		}
	}

	//Stockage des personnalités initiales
	mapPersonality := make(map[string]types.Personality)
	for i, team := range teams {
		for j := range team.Drivers {
			var perso types.Personality
			perso.TraitsValue = make(map[string]int)
			perso.TraitsValue["Aggressivity"] = teamsJSON[i].Drivers[j].Personality.Aggressivity
			perso.TraitsValue["Confidence"] = teamsJSON[i].Drivers[j].Personality.Confidence
			perso.TraitsValue["Docility"] = teamsJSON[i].Drivers[j].Personality.Docility
			perso.TraitsValue["Concentration"] = teamsJSON[i].Drivers[j].Personality.Concentration
			mapPersonality[teams[i].Drivers[j].Id] = perso
		}
	}
	return teams, mapPersonality, nil
}
