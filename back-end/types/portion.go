package types

import (
	"fmt"
	"log"
)

type PortionType int

const (
	TURN PortionType = iota
	STRAIGHT
)

type Portion struct {
	Id          string          // Portion ID
	Difficulty  int             // Difficulty of the portion in [0,5]
	DriversOn   []*DriverInRace // Drivers on the portion
	Type        PortionType     // Type of the portion
	IsDRSZone   bool            // True if is a DRS Zone. -> increases chances of overtaking
	NextPortion *Portion        // Next Portion in the circuit
}

func NewPortion(id string, difficulty int, driversOn []*DriverInRace, isDRSZone bool) *Portion {

	d := make([]*DriverInRace, len(driversOn))
	copy(d, driversOn)

	var t PortionType
	if len(id) < len("turn") {
	} else if id[:len("turn")] == "turn" {
		t = TURN
	} else {
		t = STRAIGHT
	}
	return &Portion{
		Id:         id,
		Difficulty: difficulty,
		DriversOn:  d,
		Type:       t,
		IsDRSZone:  isDRSZone,
	}
}

func (p *Portion) AddDriverOn(driver *DriverInRace) {
	p.DriversOn = append(p.DriversOn, driver)
}

func (p *Portion) RemoveDriverOn(driver *DriverInRace) {
	for i, d := range p.DriversOn {
		if d == driver {
			p.DriversOn = append(p.DriversOn[:i], p.DriversOn[i+1:]...)
			return
		}
	}
}

func (p *Portion) SwapDrivers(driver1 *DriverInRace, driver2 *DriverInRace) error {
	var i1, i2 int
	for i, d := range p.DriversOn {
		if d == driver1 {
			i1 = i
		}
		if d == driver2 {
			i2 = i
		}
	}
	if i1 != 0 && i2 != 0 {
		p.DriversOn[i1] = driver2
		p.DriversOn[i2] = driver1
		return nil
	}
	return fmt.Errorf("Driver %s or %s not found on portion %s", driver1.Driver.Id, driver2.Driver.Id, p.Id)
}

func (p *Portion) DisplayDriversOn() {
	log.Print("Drivers on portion", p.Id, " [ ")
	for i, driver := range p.DriversOn {
		log.Printf("%d : %s, ", len(p.DriversOn)-i, driver.Driver.Lastname)
	}
	log.Println("]")
}
