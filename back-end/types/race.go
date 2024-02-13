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
	log.Printf("Launching new race: %s...\n", r.Id)

	// Creating the shared map
	mapChan := sync.Map{}
	for _, t := range r.Teams {
		for _, d := range t.Drivers {
			mapChan.Store(d.Id, make(chan Action))
		}
	}
	// Creating instances of drivers in the race

	drivers, err := MakeSliceOfDriversInRace(r.Teams, &(r.Circuit.Portions[0]), mapChan, r.MeteoCondition)
	if err != nil {
		return nil, err
	}
	drivers = ShuffleDrivers(drivers)
	log.Println("\n\nStarting Line:")
	for i := range drivers {
		log.Printf("%d: %s %s\n", len(drivers)-i, drivers[i].Driver.Firstname, drivers[i].Driver.Lastname)
	}

	// Putting all drivers on the starting line
	for _, driver := range drivers {
		driver.Position.AddDriverOn(driver)
	}
	log.Println("Drivers are on the starting line")
	// Starting driver agents
	log.Println("Race begins...")
	for _, driver := range drivers {
		go driver.Start(driver.Position, r.Circuit.NbLaps)
	}
	var nbFinish = 0
	var nbDrivers = len(r.Teams) * 2
	decisionMap := make(map[string]Action, nbDrivers)

	// Simulating until all drivers finish the race
	for nbFinish < nbDrivers {
		// Each driver, in a random order, performs tests on overtaking probability, etc...
		drivers = ShuffleDrivers(drivers)

		for i := range drivers {
			// Releasing racing drivers to make decisions
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue // Necessary to wait only for racing drivers
			}

			// Adding a random penalty for tire wear if it's hot
			dice := rand.Intn(100)
			if r.MeteoCondition == HEAT && dice < 25 {
				drivers[i].TimeWoPitStop += 5
			}

			drivers[i].ChanEnv <- 1
		}
		// Retrieving drivers' decisions
		for i := range drivers {
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue // Necessary to wait only for racing drivers
			}
			decisionMap[drivers[i].Driver.Id] = <-drivers[i].ChanEnv
		}

		// Processing decisions and updating drivers' positions
		for i := range drivers {
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue
			}
			decision := decisionMap[drivers[i].Driver.Id]
			switch decision {
			case TRY_OVERTAKE:
				// Checking if the driver can overtake
				driverToOvertake, err := drivers[i].DriverToOvertake()
				if err != nil {
					log.Printf("Error while getting driver to overtake: %s\n", err)
				}
				if driverToOvertake != nil {
					// Checking if the driver successfully overtakes
					success, crashedDrivers := drivers[i].Overtake(driverToOvertake)

					if len(crashedDrivers) > 0 {
						// Creating a crash highlight
						highlight, err := NewHighlight(crashedDrivers, CRASHOVERTAKE)
						if err != nil {
							log.Printf("Error while creating highlight: %s\n", err)
						}
						r.HighLigths = append(r.HighLigths, *highlight)
						log.Println(highlight.Description)
						// Removing crashed drivers
						for ind := range crashedDrivers {
							crashedDrivers[ind].Status = CRASHED
							r.FinalResult = append(r.FinalResult, crashedDrivers[ind].Driver) // adding to the list
							drivers[i].Position.RemoveDriverOn(crashedDrivers[ind])
							nbFinish++
						}
					}

					if success {
						// Creating an overtaking highlight
						highlight, err := NewHighlight([]*DriverInRace{drivers[i], driverToOvertake}, OVERTAKE)
						if err != nil {
							log.Printf("Error while creating highlight: %s\n", err)
						}
						r.HighLigths = append(r.HighLigths, *highlight)
						log.Println(highlight.Description)
						// Updating positions
						drivers[i].Position.SwapDrivers(drivers[i], driverToOvertake)
					}

				}
			case CONTINUE:
				// Checking if the driver succeeds in passing the portion

				penalty := 0

				if r.MeteoCondition == RAINY {
					penalty = 25
				}

				success := drivers[i].PortionSuccess(penalty)

				if !success {
					// In case of a crash, the driver's confidence and docility are affected
					if drivers[i].Driver.Personality.TraitsValue["Confidence"] > 1 {
						drivers[i].Driver.Personality.TraitsValue["Confidence"] -= 1
					}
					if drivers[i].Driver.Personality.TraitsValue["Docility"] < 5 {
						drivers[i].Driver.Personality.TraitsValue["Docility"] += 1
					}
					// Creating a crash highlight
					highlight, err := NewHighlight([]*DriverInRace{drivers[i]}, CRASHPORTION)
					if err != nil {
						log.Printf("Error while creating highlight: %s\n", err)
					}
					r.HighLigths = append(r.HighLigths, *highlight)
					log.Println(highlight.Description)
					// Removing crashed driver
					drivers[i].Status = CRASHED
					r.FinalResult = append(r.FinalResult, drivers[i].Driver) // adding to the list
					drivers[i].Position.RemoveDriverOn(drivers[i])
					nbFinish++
				}

			case NOOP:
				// Do nothing

				if (drivers[i].Status == PITSTOP || drivers[i].Status == PITSTOP_CHANGETYRE) && drivers[i].PitstopSteps == 3 {
					// Creating a pitstop highlight
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
				// Creating a crash highlight
				highlight, err := NewHighlight([]*DriverInRace{drivers[i]}, CREVAISON)
				if err != nil {
					log.Printf("Error while creating highlight: %s\n", err)
				}
				r.HighLigths = append(r.HighLigths, *highlight)
				log.Println(highlight.Description)

				drivers[i].Status = CRASHED
				r.FinalResult = append(r.FinalResult, drivers[i].Driver) // Add it to the array
				drivers[i].Position.RemoveDriverOn(drivers[i])
				nbFinish++

			}
		}

		// Move all drivers who haven't finished the race, aren't crashed, and aren't in pitstop
		newDriversOnPortion := make([][]*DriverInRace, len(r.Circuit.Portions)) // Store the new positions of the drivers
		for i := range r.Circuit.Portions {
			for _, driver := range r.Circuit.Portions[i].DriversOn {
				previousPortion := driver.Position
				if driver.Status != CRASHED && driver.Status != ARRIVED && (driver.Status != PITSTOP && driver.Status != PITSTOP_CHANGETYRE) {
					// Update the driver's position field
					driver.Position = driver.Position.NextPortion
					if i == len(r.Circuit.Portions)-1 {
						log.Println("Turn ", driver.NbLaps, " for ", driver.Driver.Lastname)
						// If we completed a lap
						driver.NbLaps += 1
						if driver.NbLaps == r.Circuit.NbLaps {
							// If we finished the race, remove the driver from the circuit and put them in the ranking
							// Create a finish highlight
							highlight, err := NewHighlight([]*DriverInRace{driver}, FINISH)
							if err != nil {
								log.Printf("Error while creating highlight: %s\n", err)
							}
							r.HighLigths = append(r.HighLigths, *highlight)
							log.Println(highlight.Description)
							// Signal that the driver finished the race
							driver.Status = ARRIVED
							nbFinish++
							r.FinalResult = append(r.FinalResult, driver.Driver)
						}
					}
				}
				if driver.Status == PITSTOP || driver.Status == PITSTOP_CHANGETYRE {
					newDriversOnPortion[i] = append(newDriversOnPortion[i], driver)
				} else if driver.Speed > 5 && (driver.Status != CRASHED && driver.Status != ARRIVED) {
					// If the driver is currently first on their portion and there is no one on portion i+1
					canJump := r.Circuit.Portions[i].DriversOn[0] == driver && len(r.Circuit.Portions[(i+1)%len(r.Circuit.Portions)].DriversOn) == 0
					dice := rand.Intn(15) + 1
					if canJump && dice < driver.Speed {
						driver.Position = driver.Position.NextPortion
						newDriversOnPortion[(i+2)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+2)%len(r.Circuit.Portions)], driver)
					} else {
						newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+1)%len(r.Circuit.Portions)], driver)
					}

				} else if driver.Speed <= 5 && (driver.Status != CRASHED && driver.Status != ARRIVED) {
					// If the driver is currently last on their portion and not going very fast, they can "stagnate"
					canStay := r.Circuit.Portions[i].DriversOn[len(r.Circuit.Portions[i].DriversOn)-1] == driver
					dice := rand.Intn(6) + 1
					if canStay && dice > driver.Speed { // At speed = 5, chance to stay = 1/6, at speed = 1, chance to stay = 5/6
						driver.Position = previousPortion
						driver.Speed++ // Increase the driver's speed to avoid them staying on the same portion for a long time
						newDriversOnPortion[i] = append(newDriversOnPortion[i], driver)
					} else {
						newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+1)%len(r.Circuit.Portions)], driver)
					}
				} else if driver.Status != CRASHED && driver.Status != ARRIVED {
					// Add the driver to their new position
					newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+1)%len(r.Circuit.Portions)], driver)
				}
			}
		}

		// Update the drivers' positions
		for i := range r.Circuit.Portions {
			r.Circuit.Portions[i].DriversOn = make([]*DriverInRace, len(newDriversOnPortion[i])) // Overwrite the old slice
			copy(r.Circuit.Portions[i].DriversOn, newDriversOnPortion[i])                        // Replace with the new one
		}
	}
	// Display the ranking
	log.Println("\n\nFinal ranking:")
	for i := range r.FinalResult {
		log.Printf("%d : %s %s\n", len(r.FinalResult)-i, r.FinalResult[i].Firstname, r.FinalResult[i].Lastname)
	}
	// Return the ranking and the points awarded
	res := r.CalcDriversPoints()
	return res, nil
}

func (r *Race) CalcDriversPoints() map[string]int {
	var n int = len(r.FinalResult)
	res := make(map[string]int, n)
	for i := 0; i < len(r.FinalResult); i++ {
		res[r.FinalResult[i].Id] = 0
	}
	// The first one gets 25 points
	res[r.FinalResult[n-1].Id] = 25

	// The second one gets 18 points
	res[r.FinalResult[n-2].Id] = 18

	// The third one gets 15 points
	res[r.FinalResult[n-3].Id] = 15

	// The fourth one gets 12 points
	res[r.FinalResult[n-4].Id] = 12

	// The fifth one gets 10 points, and we decrement by 2 until the ninth one
	for i := 1; i <= 5; i++ {
		res[r.FinalResult[n-4-i].Id] = 12 - 2*i
	}

	// The tenth one gets 1 point
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
