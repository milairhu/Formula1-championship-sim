package types

import (
	"log"
	"sort"
)

type Championship struct {
	Id       string     // Championship ID
	Name     string     // Name
	Circuits []*Circuit // Set of circuits that compose the championship. Defined at the creation of the championship
	Races    []Race     // Array of Races, filled during the championship
	Teams    []*Team    // Set of teams
}

func NewChampionship(id string, name string, circuits []*Circuit, teams []*Team) *Championship {

	c := make([]*Circuit, len(circuits))
	copy(c, circuits)

	t := make([]*Team, len(teams))
	copy(t, teams)
	for i := range t {
		for j := range t[i].Drivers {
			t[i].Drivers[j].ChampionshipPoints = 0
		}

	}

	r := make([]Race, len(circuits))

	return &Championship{
		Id:       id,
		Name:     name,
		Circuits: c,
		Races:    r,
		Teams:    t,
	}
}

func NewChampionshipRandom(id string, name string, circuits []*Circuit, teams []*Team) *Championship {

	c := make([]*Circuit, len(circuits))
	copy(c, circuits)

	t := make([]*Team, len(teams))
	copy(t, teams)
	for i := range t {
		for j := range t[i].Drivers {
			t[i].Drivers[j].ChampionshipPoints = 0
			t[i].Drivers[j].Personality = *NewPersonalityRandom()
		}

	}

	r := make([]Race, len(circuits))

	return &Championship{
		Id:       id,
		Name:     name,
		Circuits: c,
		Races:    r,
		Teams:    t,
	}
}

// Note : pointers are used when life cycle is not handled by the object itself

func (c *Championship) CalcTeamRank() []*Team {
	res := make([]*Team, len(c.Teams))
	copy(res, c.Teams)
	sort.Slice(res, func(i, j int) bool {
		return res[i].CalcChampionshipPoints() > res[j].CalcChampionshipPoints()
	})

	return res
}

func (c *Championship) DisplayTeamRank() []*TeamTotalPoints {
	log.Printf("\n\n==== Teams Rankig=ng ====\n")
	teamRank := c.CalcTeamRank()
	teamsRankTab := make([]*TeamTotalPoints, 0)
	for i, team := range teamRank {
		teamRank := NewTeamTotalPoints(team.Name, team.CalcChampionshipPoints())
		log.Printf("%d : %s : %d points\n", i+1, team.Name, team.CalcChampionshipPoints())

		teamsRankTab = append(teamsRankTab, teamRank)
	}

	return teamsRankTab
}

func (c *Championship) CalcDriverRank() []*Driver {

	res := make([]*Driver, 0)
	for indT := range c.Teams {
		for indD := range c.Teams[indT].Drivers {
			res = append(res, &c.Teams[indT].Drivers[indD])
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ChampionshipPoints > res[j].ChampionshipPoints
	})

	return res
}

func (c *Championship) DisplayDriverRank() ([]*DriverTotalPoints, []*PersonalityAveragePoints, map[string]map[int]float64) {
	log.Printf("\n\n==== Drivers Ranking ====\n")
	driversRank := c.CalcDriverRank()
	driversRankTab := make([]*DriverTotalPoints, 0)
	personalityRankTab := make([]*PersonalityAveragePoints, 0)
	personnalityAverage := make(map[string]map[int]float64)
	nb := make(map[string]map[int]int)

	for i, driver := range driversRank {
		driverRank := NewDriverTotalPoints(driver.Lastname, driver.ChampionshipPoints)

		// Add drivers points to the personalities
		for personnality, level := range driver.Personality.TraitsValue {
			if _, ok := personnalityAverage[personnality]; !ok {
				personnalityAverage[personnality] = make(map[int]float64)
				nb[personnality] = make(map[int]int)
			}
			personnalityAverage[personnality][level] += float64(driver.ChampionshipPoints)
			nb[personnality][level] += 1
		}

		// If personality already in the array, add the points. Else, create a new personality in the array
		var found bool
		for indPers := range personalityRankTab {
			if personalityRankTab[indPers].Personality["Aggressivity"] == driver.Personality.TraitsValue["Aggressivity"] &&
				personalityRankTab[indPers].Personality["Concentration"] == driver.Personality.TraitsValue["Concentration"] &&
				personalityRankTab[indPers].Personality["Confidence"] == driver.Personality.TraitsValue["Confidence"] &&
				personalityRankTab[indPers].Personality["Docility"] == driver.Personality.TraitsValue["Docility"] {
				personalityRankTab[indPers].AveragePoints += float64(driver.ChampionshipPoints)
				personalityRankTab[indPers].NbDrivers += 1
				found = true
				break
			}
		}
		if !found {
			var perso Personality
			perso.TraitsValue = make(map[string]int)
			perso.TraitsValue["Confidence"] = driver.Personality.TraitsValue["Confidence"]
			perso.TraitsValue["Aggressivity"] = driver.Personality.TraitsValue["Aggressivity"]
			perso.TraitsValue["Docility"] = driver.Personality.TraitsValue["Docility"]
			perso.TraitsValue["Concentration"] = driver.Personality.TraitsValue["Concentration"]
			personalityRank := NewPersonalityAveragePoints(perso.TraitsValue, driver.ChampionshipPoints, 1)
			personalityRankTab = append(personalityRankTab, personalityRank)
		}

		log.Printf("%d : %s %s : %d points\n", i+1, driver.Firstname, driver.Lastname, driver.ChampionshipPoints)
		log.Printf("%v", driver.Personality.TraitsValue)

		driversRankTab = append(driversRankTab, driverRank)

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

	// Compute personalities means
	for personnality, level := range personnalityAverage {
		for level, points := range level {
			personnalityAverage[personnality][level] = points / float64(nb[personnality][level])
		}
	}

	return driversRankTab, personalityRankTab, personnalityAverage
}

func (c *Championship) DisplayPersonalityRepartition() {
	log.Printf("\n\n==== Distribution of personalities ====\n")
	driverRank := c.CalcDriverRank()

	aggressivity_value_5 := 0
	aggressivity_value_4 := 0
	aggressivity_value_3 := 0
	aggressivity_value_2 := 0
	aggressivity_value_1 := 0
	aggressivity_value_0 := 0
	for i, driver := range driverRank {
		if i < 15 {
			switch driver.Personality.TraitsValue["Aggressivity"] {
			case 0:
				aggressivity_value_0 += 1
			case 1:
				aggressivity_value_1 += 1
			case 2:
				aggressivity_value_2 += 1
			case 3:
				aggressivity_value_3 += 1
			case 4:
				aggressivity_value_4 += 1
			case 5:
				aggressivity_value_5 += 1
			default:
				log.Printf("Value of aggressivity out of range : %d", driver.Personality.TraitsValue["Aggressivity"])
			}
		}
		if i == 4 {
			log.Printf("Distribution of aggressiveness level in the top 5: \n")
			log.Printf("Aggressiveness 5 : %d", aggressivity_value_5)
			log.Printf("Aggressiveness 4 : %d", aggressivity_value_4)
			log.Printf("Aggressiveness 3 : %d", aggressivity_value_3)
			log.Printf("Aggressiveness 2 : %d", aggressivity_value_2)
			log.Printf("Aggressiveness 1 : %d", aggressivity_value_1)
			log.Printf("Aggressiveness 0 : %d", aggressivity_value_0)
		}
		if i == 9 {
			log.Printf("Distribution of aggressiveness level in the top 10: \n")
			log.Printf("Aggressiveness 5 : %d", aggressivity_value_5)
			log.Printf("Aggressiveness 4 : %d", aggressivity_value_4)
			log.Printf("Aggressiveness 3 : %d", aggressivity_value_3)
			log.Printf("Aggressiveness 2 : %d", aggressivity_value_2)
			log.Printf("Aggressiveness 1 : %d", aggressivity_value_1)
			log.Printf("Aggressiveness 0 : %d", aggressivity_value_0)
		}
		if i == 14 {
			log.Printf("Distribution of aggressiveness level in the top 15: \n")
			log.Printf("Aggressiveness 5 : %d", aggressivity_value_5)
			log.Printf("Aggressiveness 4 : %d", aggressivity_value_4)
			log.Printf("Aggressiveness 3 : %d", aggressivity_value_3)
			log.Printf("Aggressiveness 2 : %d", aggressivity_value_2)
			log.Printf("Aggressiveness 1 : %d", aggressivity_value_1)
			log.Printf("Aggressiveness 0 : %d", aggressivity_value_0)
		}
	}
}
