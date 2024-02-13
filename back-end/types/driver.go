package types

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

type Driver struct {
	Id                 string      `json:"id"`                 // Driver ID
	Firstname          string      `json:"firstname"`          // Firstname
	Lastname           string      `json:"lastname"`           // Lastname
	Level              int         `json:"level"`              // Level of the driver, in [1, 10]
	Country            string      `json:"country"`            // Country
	Personality        Personality `json:"personality"`        // Personality
	ChampionshipPoints int         `json:"championshipPoints"` // Points in the current champonship
}
type DriverInRace struct {
	Driver         *Driver      // The driver himself
	Position       *Portion     // Position of the driver
	NbLaps         int          // Number of laps completed
	Status         DriverStatus // Status of the driver
	IsPitStop      bool         // PitStop --> true if the driver is in pitstop
	TimeWoPitStop  int          // Time without pitstop --> increments at each step
	ChanEnv        chan Action  // channel to receive and send actions & the environment
	PitstopSteps   int          // Number of steps blocked in pitstop
	PrevTyre       Tyre
	CurrentTyre    Tyre         // Current type of tyre
	AvailableTyres map[Tyre]int // Quantity of available tyres by type
	TyreTypeCount  int          // Number of different types of tyres used during the race
	UsedTyreTypes  []Tyre
	CurrentRank    int // Current ranking of the driver, to see if we implement it
	Speed          int // Speed of the driver
}

// Actions of a driver

type Action int

const (
	TRY_OVERTAKE Action = iota
	NOOP
	CONTINUE
	ACCIDENTPNEUS
)

type DriverStatus int

const (
	RACING DriverStatus = iota
	CRASHED
	ARRIVED
	PITSTOP
	PITSTOP_CHANGETYRE
)

type Tyre int

const (
	WET  Tyre = 0
	SOFT Tyre = iota
	MEDIUM
	HARD
)

func NewDriver(id string, firstname string, lastname string, level int, country string, team *Team, personality Personality) *Driver {

	return &Driver{
		Id:                 id,
		Firstname:          firstname,
		Lastname:           lastname,
		Level:              level,
		Country:            country,
		Personality:        personality,
		ChampionshipPoints: 0,
	}
}

func NewDriverInRace(driver *Driver, position *Portion, channel chan Action, meteoCondition Meteo) *DriverInRace {
	// WET tyre if RAINY weather
	if meteoCondition == RAINY {
		return &DriverInRace{
			Driver:        driver,
			Position:      position,
			NbLaps:        0,
			ChanEnv:       channel,
			Status:        RACING,
			TimeWoPitStop: 0,
			PitstopSteps:  0,
			CurrentTyre:   WET,
			TyreTypeCount: 1,
			Speed:         1,
		}
	} else {
		return &DriverInRace{
			Driver:        driver,
			Position:      position,
			NbLaps:        0,
			ChanEnv:       channel,
			Status:        RACING,
			TimeWoPitStop: 0,
			PitstopSteps:  0,
			CurrentTyre:   SOFT,
			TyreTypeCount: 1,
			Speed:         1,
		}
	}
}

// Function to test if a driver successfully completes a portion without crashing
func (d *DriverInRace) PortionSuccess(penalty int) bool {
	// For now we take into account the driver's level, the portion's difficulty, and tire wear
	portion := d.Position

	// If the speed is low, we consider the portion is successfully completed without difficulty
	if d.Speed <= 3 {
		return true
	}

	successProbability := 995
	successProbability += d.Driver.Level * 20
	successProbability -= portion.Difficulty * 18
	successProbability -= d.TimeWoPitStop
	successProbability -= penalty
	successProbability -= d.Speed * 5
	successProbability += d.Driver.Personality.TraitsValue["Concentration"] * 10

	var dice int = rand.Intn(999) + 1

	if dice <= 10 && d.Speed < 10 {
		d.Speed += int(d.CurrentTyre)
	} else if dice >= 990 && d.Speed > 1 {
		d.Speed--
	}

	return dice <= successProbability
}

func MakeSliceOfDriversInRace(teams []*Team, startingPortion *Portion, mapChan sync.Map, weatherCondition Meteo) ([]*DriverInRace, error) {
	res := make([]*DriverInRace, 0)
	for _, team := range teams {
		for _, driver := range team.Drivers {
			dtamp := driver // Necessary, otherwise only the address of a team member is used
			c, ok := mapChan.Load(dtamp.Id)
			if !ok {
				return nil, fmt.Errorf("error while creating driver in race : %s", driver.Id)
			}
			d := NewDriverInRace(&dtamp, startingPortion, c.(chan Action), weatherCondition)
			res = append(res, d)
		}
	}
	return res, nil
}

func ShuffleDrivers(drivers []*DriverInRace) []*DriverInRace {
	rand.Shuffle(len(drivers), func(i, j int) {
		drivers[i], drivers[j] = drivers[j], drivers[i]
	})
	return drivers
}

func (d *DriverInRace) PitStop() bool {

	if d.Status == PITSTOP || d.Status == PITSTOP_CHANGETYRE || d.Position.Id != "straight_1" {
		return false
	}

	pitStopProbability := 0

	// Check if a pit stop is necessary
	pitStopProbability += (d.TimeWoPitStop * 10) / 2

	var dice int = rand.Intn(999) + 1

	if dice <= pitStopProbability {
		d.PitstopSteps = 3
		d.TimeWoPitStop = 0
		return true
	}

	return false
}

func (d *DriverInRace) ChangeTyreType() bool {
	// Keep the WET tire if it's raining
	if d.CurrentTyre == WET {
		return false
	}

	// Check if we change the tire type, more likely if later in the race and fewer tire types used < 2
	changeTyreTypeProbability := 0
	tyreTypeCount := len(d.UsedTyreTypes)

	changeTyreTypeProbability = d.NbLaps + (3-tyreTypeCount)*20

	var dice = rand.Intn(99) + 1

	if dice < changeTyreTypeProbability {
		switch d.CurrentTyre {
		case SOFT:
			d.PrevTyre = SOFT
			d.CurrentTyre = MEDIUM
			d.UsedTyreTypes = append(d.UsedTyreTypes, MEDIUM)
			return true
		case MEDIUM:
			d.PrevTyre = MEDIUM
			d.CurrentTyre = SOFT
			d.UsedTyreTypes = append(d.UsedTyreTypes, SOFT)
			return true
		case HARD:
			d.PrevTyre = HARD
			d.CurrentTyre = MEDIUM
			d.UsedTyreTypes = append(d.UsedTyreTypes, MEDIUM)
			return true
		}
	}

	return false
}

func (d *DriverInRace) TestTyres() bool {

	if d.Status == PITSTOP || d.Status == PITSTOP_CHANGETYRE {
		return true
	}

	tyresProbability := 0

	// Check if the tires are fine
	tyresProbability += d.TimeWoPitStop - (int(d.CurrentTyre)*10 + 100)

	var dice int = rand.Intn(999) + 1

	return dice > tyresProbability
}

// Function to overtake another driver
func (d *DriverInRace) Overtake(otherDriver *DriverInRace) (success bool, crashedDrivers []*DriverInRace) {

	var overtakingProbability int

	// If the other driver is in pitstop, overtaking is guaranteed
	if otherDriver.Status == PITSTOP || otherDriver.Status == PITSTOP_CHANGETYRE {
		return true, []*DriverInRace{}
	}

	if d.Status == PITSTOP || d.Status == PITSTOP_CHANGETYRE {
		return false, []*DriverInRace{}
	}

	// Attempting an overtake, tire wear increases
	d.TimeWoPitStop += 10

	// Depending on the confidence and concentration traits values of the driver, the probability of a successful overtake varies
	if d.Driver.Personality.TraitsValue["Confidence"] > 3 && d.Driver.Personality.TraitsValue["Concentration"] > 3 {
		overtakingProbability = 75
	} else if d.Driver.Personality.TraitsValue["Confidence"] <= 3 && d.Driver.Personality.TraitsValue["Concentration"] >= 3 {
		overtakingProbability = 65
	} else {
		overtakingProbability = 50
	}

	if d.Driver.Level > otherDriver.Driver.Level {
		overtakingProbability += 10
	} else if d.Driver.Level < otherDriver.Driver.Level {
		overtakingProbability -= 10
	}

	portion := d.Position

	// For now we consider the drivers' levels and the "difficulty" of the portion
	overtakingProbability -= portion.Difficulty * 2

	overtakingProbability *= 10

	var dice int = rand.Intn(999) + 1

	// Depending on the dice roll, it affects the driver's speed
	if dice <= 10 && d.Speed < 10 {
		d.Speed += int(d.CurrentTyre)
	} else if dice >= 990 && d.Speed > 1 {
		d.Speed--
	}

	// If the dice roll is below overtakingProbability, the overtake is successful and the driver's confidence increases
	if dice <= overtakingProbability {
		if d.Driver.Personality.TraitsValue["Confidence"] < 5 {
			d.Driver.Personality.TraitsValue["Confidence"] += 1
			// If the driver's confidence is already maxed out and the overtake is successful, they become less docile
		} else if d.Driver.Personality.TraitsValue["Confidence"] == 5 && d.Driver.Personality.TraitsValue["Docility"] > 1 {
			d.Driver.Personality.TraitsValue["Docility"] -= 1
		}
		return true, []*DriverInRace{}
	}

	// Otherwise, check if a crash occurs

	// Here we have a critical failure, both drivers crash
	if dice >= 995 {
		return false, []*DriverInRace{d, otherDriver}
	}

	// Here, only one driver crashes, randomly determined
	if dice >= 990 {
		if dice%2 == 0 {
			return false, []*DriverInRace{d}
		} else {
			return false, []*DriverInRace{otherDriver}
		}
	}

	// In the default case, the overtake fails but no crash occurs
	// Reset the driver's speed
	d.Speed = 1
	return false, []*DriverInRace{}

}

func (d *DriverInRace) DriverToOvertake() (*DriverInRace, error) {

	if d.Status == PITSTOP || d.Status == PITSTOP_CHANGETYRE {
		return nil, nil
	}

	p := d.Position
	for i := range p.DriversOn {
		if p.DriversOn[i] == d {
			if len(p.DriversOn) > i+1 && p.DriversOn[i+1] != nil {
				return p.DriversOn[i+1], nil
			} else {
				return nil, nil
			}
		}
	}

	return nil, fmt.Errorf("Driver %s (%s, crashed if =1 : %d) who want to overtake is not found on portion %s", d.Driver.Id, d.Driver.Lastname, d.Status, p.Id)
}

// Function to decide whether to attempt an overtake or not
func (d *DriverInRace) OvertakeDecision(driverToOvertake *DriverInRace) (bool, error) {

	toOvertake, err := d.DriverToOvertake()
	if err != nil {
		return false, err
	}

	p := d.Position
	overtakeProbability := 0

	// Aggressiveness and confidence of the driver impact overtaking attempts
	if d.Driver.Personality.TraitsValue["Aggressivity"] > 3 || (d.Driver.Personality.TraitsValue["Aggressivity"] >= 3 && d.Driver.Personality.TraitsValue["Confidence"] >= 3) {
		overtakeProbability += d.Driver.Personality.TraitsValue["Aggressivity"] * 2
	} else {
		overtakeProbability -= (d.Driver.Personality.TraitsValue["Docility"]) * 2
	}

	if p.Difficulty != 0 {
		overtakeProbability += 20 / p.Difficulty
	} else {
		overtakeProbability = 0
	}

	if toOvertake != nil {
		// If the driver is in pitstop, choose to overtake
		if driverToOvertake.Status == PITSTOP || driverToOvertake.Status == PITSTOP_CHANGETYRE {
			return true, nil
		}

		var dice = rand.Intn(99) + 1
		return dice <= overtakeProbability, nil
	}
	return false, nil
}

// Change the driver's speed based on their personality
func (d *DriverInRace) ChangeSpeed() {

	if d.Driver.Personality.TraitsValue["Confidence"] > 3 && d.Driver.Personality.TraitsValue["Concentration"] > 3 {
		d.Speed += 2
	}
	if d.Driver.Personality.TraitsValue["Confidence"] <= 3 && d.Driver.Personality.TraitsValue["Concentration"] >= 3 {
		d.Speed += 1
	}
	if d.Driver.Personality.TraitsValue["Confidence"] <= 3 && d.Driver.Personality.TraitsValue["Concentration"] <= 3 {
		d.Speed -= 1
	}
	if d.Driver.Personality.TraitsValue["Confidence"] > 3 && d.Driver.Personality.TraitsValue["Concentration"] <= 3 {
		d.Speed -= 2
	}
	if d.Driver.Personality.TraitsValue["Aggressivity"] > 3 {
		d.Speed += 2
	}
	if d.Driver.Personality.TraitsValue["Aggressivity"] <= 3 {
		d.Speed -= 2
	}

	if d.Speed > 10 {
		d.Speed = 10
	} else if d.Speed < 1 {
		d.Speed = 1
	}

}

func (d *DriverInRace) Start(position *Portion, nbLaps int) {
	log.Printf("		Starting driver %s %s...\n", d.Driver.Firstname, d.Driver.Lastname)

	for {
		if d.Status == ARRIVED || d.Status == CRASHED {
			return
		}
		// Wait for environment to signal decision time
		<-d.ChanEnv
		if d.Status == ARRIVED || d.Status == CRASHED {
			return
		}
		// Make decision

		// Change the driver's speed based on their personality
		d.ChangeSpeed()

		// Check if the tires are in good condition
		if !d.TestTyres() {
			d.ChanEnv <- ACCIDENTPNEUS
			continue
		}

		// Check if a pitstop is needed
		pitstop := false

		if d.Status != PITSTOP && d.Status != PITSTOP_CHANGETYRE {
			d.TimeWoPitStop++
			pitstop = d.PitStop()
		}

		if pitstop {
			// Check if tire type needs to be changed
			changeTyre := d.ChangeTyreType()
			if changeTyre {
				d.Status = PITSTOP_CHANGETYRE
			} else {
				d.Status = PITSTOP
			}
			d.ChanEnv <- NOOP
			continue
		} else if d.Status == PITSTOP || d.Status == PITSTOP_CHANGETYRE {
			d.PitstopSteps--
			d.ChanEnv <- NOOP
			continue
		}

		// Check if overtaking is possible
		toOvertake, err := d.DriverToOvertake()
		if err != nil {
			log.Printf("Error while getting the driver to overtake: %s\n", err)
		}
		if toOvertake != nil {
			// Decide whether to attempt an overtake
			decision, err := d.OvertakeDecision(toOvertake)
			if err != nil {
				log.Printf("Error while getting the decision to overtake: %s\n", err)
			}
			if decision {
				// Signal the decision to the environment
				d.ChanEnv <- TRY_OVERTAKE
			} else {
				d.ChanEnv <- CONTINUE
			}

		} else {
			// If no overtaking opportunity, do nothing
			d.ChanEnv <- CONTINUE
		}
		// Check if race is finished
		if d.NbLaps == nbLaps {
			return
		}
	}
}
