package types

import (
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Race struct {
	Id             string      // Race ID
	Circuit        *Circuit    // Circuit
	Date           time.Time   // Date
	Teams          []*Team     // Set of teams
	MeteoCondition Meteo       // Meteo condition
	FinalResult    []*Driver   // Final result, drivers rank from 1st to last
	HighLigths     []Highlight // Containes all what happend during the race

}

func NewRace(id string, circuit *Circuit, date time.Time, teams []*Team, meteo Meteo) *Race {

	d := make([]*Team, len(teams))
	copy(d, teams)

	f := make([]*Driver, 0)

	h := make([]Highlight, 0)

	return &Race{
		Id:             id,
		Circuit:        circuit,
		Date:           date,
		Teams:          d,
		MeteoCondition: meteo,
		FinalResult:    f,
		HighLigths:     h,
	}
}

func (r *Race) SimulateRace() (map[string]int, error) {
	log.Printf("	Lancement d'une nouvelle course : %s...\n", r.Id)

	//Création du map partagé
	mapChan := sync.Map{}
	for _, t := range r.Teams {
		for _, d := range t.Drivers {
			mapChan.Store(d.Id, make(chan Action))
		}
	}
	//On crée les instances des pilotes en course

	drivers, err := MakeSliceOfDriversInRace(r.Teams, &(r.Circuit.Portions[0]), mapChan, r.MeteoCondition)
	if err != nil {
		return nil, err
	}
	drivers = ShuffleDrivers(drivers)
	log.Println("\n\nLigne de départ:")
	for i := range drivers {
		log.Printf("%d : %s %s\n", len(drivers)-i, drivers[i].Driver.Firstname, drivers[i].Driver.Lastname)
	}

	//On met tous les agents sur la ligne de départ :
	for _, driver := range drivers {
		driver.Position.AddDriverOn(driver)
	}
	log.Println("Les pilotes sont sur la ligne de départ")
	//On lance les agents pilotes
	log.Println("Début de la course...")
	for _, driver := range drivers {
		go driver.Start(driver.Position, r.Circuit.NbLaps)
	}
	var nbFinish = 0
	var nbDrivers = len(r.Teams) * 2
	decisionMap := make(map[string]Action, nbDrivers)

	//On simule tant que tous les pilotes n'ont pas fini la course
	for nbFinish < nbDrivers {
		//Chaque pilote, dans un ordre aléatoire, réalise les tests sur la proba de dépasser etc...
		drivers = ShuffleDrivers(drivers)

		for i := range drivers {
			//On débloque les pilotes en course pour qu'ils prennent une décision
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue //Obligatoire car il ne faut attendre que les pilotes qui courent encore
			}

			// On ajoute une pénalité aléatoire sur l'usure des pneus si il fait chaud
			dice := rand.Intn(100)
			if r.MeteoCondition == HEAT && dice < 25 {
				drivers[i].TimeWoPitStop += 5
			}

			drivers[i].ChanEnv <- 1
		}
		// On récupère les décisions des pilotes
		for i := range drivers {
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue //Obligatoire car il ne faut attendre que les pilotes qui courent encore
			}
			decisionMap[drivers[i].Driver.Id] = <-drivers[i].ChanEnv
		}

		//On traite les décisions et on met à jour les positions des pilotes
		for i := range drivers {
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue
			}
			decision := decisionMap[drivers[i].Driver.Id]
			switch decision {
			case TRY_OVERTAKE:
				//On vérifie si le pilote peut bien dépasser
				driverToOvertake, err := drivers[i].DriverToOvertake()
				if err != nil {
					log.Printf("Error while getting driver to overtake: %s\n", err)
				}
				if driverToOvertake != nil {
					//On vérifie si le pilote a réussi son dépassement
					success, crashedDrivers := drivers[i].Overtake(driverToOvertake)

					if len(crashedDrivers) > 0 {
						//On crée un Highlight de crash
						highlight, err := NewHighlight(crashedDrivers, CRASHOVERTAKE)
						if err != nil {
							log.Printf("Error while creating highlight: %s\n", err)
						}
						r.HighLigths = append(r.HighLigths, *highlight)
						log.Println(highlight.Description)
						//On supprime les pilotes crashés
						for ind := range crashedDrivers {
							crashedDrivers[ind].Status = CRASHED
							r.FinalResult = append(r.FinalResult, crashedDrivers[ind].Driver) //on l'ajoute au tableau
							drivers[i].Position.RemoveDriverOn(crashedDrivers[ind])
							nbFinish++
						}
					}

					if success {
						//On crée un Highlight de dépassement
						highlight, err := NewHighlight([]*DriverInRace{drivers[i], driverToOvertake}, OVERTAKE)
						if err != nil {
							log.Printf("Error while creating highlight: %s\n", err)
						}
						r.HighLigths = append(r.HighLigths, *highlight)
						log.Println(highlight.Description)
						//On met à jour les positions
						drivers[i].Position.SwapDrivers(drivers[i], driverToOvertake)
					}

				}
			case CONTINUE:
				//On vérifie juste si le pilote réussi à passer la portion

				pénalité := 0

				if r.MeteoCondition == RAINY {
					pénalité = 25
				}

				success := drivers[i].PortionSuccess(pénalité)

				if !success {
					// En cas de crash la confiance et la docilité du pilote est impacté
					if drivers[i].Driver.Personality.TraitsValue["Confidence"] > 1 {
						drivers[i].Driver.Personality.TraitsValue["Confidence"] -= 1
					}
					if drivers[i].Driver.Personality.TraitsValue["Docility"] < 5 {
						drivers[i].Driver.Personality.TraitsValue["Docility"] += 1
					}
					//On crée un Highlight de crash
					highlight, err := NewHighlight([]*DriverInRace{drivers[i]}, CRASHPORTION)
					if err != nil {
						log.Printf("Error while creating highlight: %s\n", err)
					}
					r.HighLigths = append(r.HighLigths, *highlight)
					log.Println(highlight.Description)
					//On supprime le pilote crashé
					drivers[i].Status = CRASHED
					r.FinalResult = append(r.FinalResult, drivers[i].Driver) //on l'ajoute au tableau
					drivers[i].Position.RemoveDriverOn(drivers[i])
					nbFinish++
				}

			case NOOP:
				//On ne fait rien

				if (drivers[i].Status == PITSTOP || drivers[i].Status == PITSTOP_CHANGETYRE) && drivers[i].PitstopSteps == 3 {
					// On crée un highlight de pitstop
					var highlight *Highlight
					var err error
					switch drivers[i].Status {
					case PITSTOP:
						highlight, err = NewHighlight([]*DriverInRace{drivers[i]}, DRIVER_PITSTOP)
					case PITSTOP_CHANGETYRE:
						highlight, err = NewHighlight([]*DriverInRace{drivers[i]}, DRIVER_PITSTOP_CHANGETYRE)
					}

					if err != nil {
						log.Printf("Error while creating highlight: %s\n", err)
					}
					r.HighLigths = append(r.HighLigths, *highlight)
					log.Println(highlight.Description)
				} else if (drivers[i].Status == PITSTOP || drivers[i].Status == PITSTOP_CHANGETYRE) && drivers[i].PitstopSteps == 0 {
					drivers[i].Status = RACING
				}

			case ACCIDENTPNEUS:
				//On crée un Highlight de crash
				highlight, err := NewHighlight([]*DriverInRace{drivers[i]}, CREVAISON)
				if err != nil {
					log.Printf("Error while creating highlight: %s\n", err)
				}
				r.HighLigths = append(r.HighLigths, *highlight)
				log.Println(highlight.Description)

				drivers[i].Status = CRASHED
				r.FinalResult = append(r.FinalResult, drivers[i].Driver) //on l'ajoute au tableau
				drivers[i].Position.RemoveDriverOn(drivers[i])
				nbFinish++

			}
		}

		//On fait avancer tout les pilotes n'ayant pas fini la course, n'étant pas crashés et n'étant pas en pitstop
		newDriversOnPortion := make([][]*DriverInRace, len(r.Circuit.Portions)) //stocke les nouvelles positions des pilotes
		for i := range r.Circuit.Portions {
			for _, driver := range r.Circuit.Portions[i].DriversOn {
				previousPortion := driver.Position
				if driver.Status != CRASHED && driver.Status != ARRIVED && (driver.Status != PITSTOP && driver.Status != PITSTOP_CHANGETYRE) {
					//On met à jour le champ position du pilote
					driver.Position = driver.Position.NextPortion
					if i == len(r.Circuit.Portions)-1 {
						log.Println("Tour ", driver.NbLaps, " pour ", driver.Driver.Lastname)
						//Si on a fait un tour
						driver.NbLaps += 1
						if driver.NbLaps == r.Circuit.NbLaps {
							//Si on a fini la course, on enlève le pilote du circuit et on le met dans le classement
							//On crée un Highlight d'arrivée
							highlight, err := NewHighlight([]*DriverInRace{driver}, FINISH)
							if err != nil {
								log.Printf("Error while creating highlight: %s\n", err)
							}
							r.HighLigths = append(r.HighLigths, *highlight)
							log.Println(highlight.Description)
							//On signal que le coureur a fini la course
							driver.Status = ARRIVED
							nbFinish++
							r.FinalResult = append(r.FinalResult, driver.Driver)
						}
					}
				}
				if driver.Status == PITSTOP || driver.Status == PITSTOP_CHANGETYRE {
					newDriversOnPortion[i] = append(newDriversOnPortion[i], driver)
				} else if driver.Speed > 5 && (driver.Status != CRASHED && driver.Status != ARRIVED) {
					// Si le pilote est actuellement premier sur sa portion et qu'il n'y a personne sur le portion i+1
					canJump := r.Circuit.Portions[i].DriversOn[0] == driver && len(r.Circuit.Portions[(i+1)%len(r.Circuit.Portions)].DriversOn) == 0
					dice := rand.Intn(15) + 1
					if canJump && dice < driver.Speed {
						driver.Position = driver.Position.NextPortion
						newDriversOnPortion[(i+2)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+2)%len(r.Circuit.Portions)], driver)
					} else {
						newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+1)%len(r.Circuit.Portions)], driver)
					}

				} else if driver.Speed <= 5 && (driver.Status != CRASHED && driver.Status != ARRIVED) {
					// Si le pilote est actuellement dernier sur sa portion et qu'il ne va pas très vite, il peut "stagner"
					canStay := r.Circuit.Portions[i].DriversOn[len(r.Circuit.Portions[i].DriversOn)-1] == driver
					dice := rand.Intn(6) + 1
					if canStay && dice > driver.Speed { // A vitesse = 5, chance de rester = 1/6, à vitesse = 1, chance de rester = 5/6
						driver.Position = previousPortion
						driver.Speed++ // On augmente la vitesse du pilote pour éviter qu'il reste 1000 ans sur la même portion
						newDriversOnPortion[i] = append(newDriversOnPortion[i], driver)
					} else {
						newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+1)%len(r.Circuit.Portions)], driver)
					}
				} else if driver.Status != CRASHED && driver.Status != ARRIVED {
					//On ajoute le pilote à sa nouvelle position
					newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+1)%len(r.Circuit.Portions)], driver)
				}
			}
		}

		//On met à jour les positions des pilotes
		for i := range r.Circuit.Portions {
			r.Circuit.Portions[i].DriversOn = make([]*DriverInRace, len(newDriversOnPortion[i])) //on écrase l'ancien slice
			copy(r.Circuit.Portions[i].DriversOn, newDriversOnPortion[i])                        //on remplace par le nouveau
		}
	}
	//On affiche le classement
	log.Println("\n\nClassement final :")
	for i := range r.FinalResult {
		log.Printf("%d : %s %s\n", len(r.FinalResult)-i, r.FinalResult[i].Firstname, r.FinalResult[i].Lastname)
	}
	//On retourne le classement et les points attribués
	res := r.CalcDriversPoints()
	return res, nil
}

func (r *Race) CalcDriversPoints() map[string]int {
	var n int = len(r.FinalResult)
	res := make(map[string]int, n)
	for i := 0; i < len(r.FinalResult); i++ {
		res[r.FinalResult[i].Id] = 0
	}
	//Le premier obtient 25 points
	res[r.FinalResult[n-1].Id] = 25

	//Le deuxième obtient 18 points
	res[r.FinalResult[n-2].Id] = 18

	//Le troisième obtient 15 points
	res[r.FinalResult[n-3].Id] = 15

	//Le quatrième obtient 12 points
	res[r.FinalResult[n-4].Id] = 12

	// Le cinquième obtient 10 points, et on décremente de 2 jusqu'au neuvième
	for i := 1; i <= 5; i++ {
		res[r.FinalResult[n-4-i].Id] = 12 - 2*i
	}

	//Le dixième obtient 1 point
	res[r.FinalResult[n-10].Id] = 1

	return res
}

func (r *Race) CalcDriverRank() []*Driver {

	res := make([]*Driver, 0)
	for indT := range r.Teams {
		for indD := range r.Teams[indT].Drivers {
			res = append(res, &r.Teams[indT].Drivers[indD])
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ChampionshipPoints > res[j].ChampionshipPoints
	})

	return res
}
