package types

import (
	"fmt"
)

type Highlight struct {
	Description string          // Describe the highlight
	Drivers     []*DriverInRace // Drivers involved in the highlight
	Type        HighlightType   // Type of highlight
}

type HighlightType int

const (
	CRASHOVERTAKE HighlightType = iota
	CRASHPORTION
	OVERTAKE
	FINISH
	DRIVER_PITSTOP
	DRIVER_PITSTOP_CHANGETYRE
	CREVAISON
)

func NewHighlight(drivers []*DriverInRace, highlightType HighlightType) (*Highlight, error) {

	tyreTypes := [4]string{"WET", "SOFT", "MEDIUM", "HARD"}

	d := make([]*DriverInRace, len(drivers))
	copy(d, drivers)
	var desc string

	switch highlightType {
	case CRASHOVERTAKE:
		if len(drivers) > 1 {
			desc = fmt.Sprintf("CRASH au tour %d: Plusieurs pilotes sont rentrés en accident : %s et %s pendant une tentative de doublement", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[1].Driver.Lastname)
		} else if len(drivers) == 1 {
			desc = fmt.Sprintf("CRASH au tour %d: Le pilote %s s'est crashé pendant une tentative de doublement", drivers[0].NbLaps, drivers[0].Driver.Lastname)
		} else {
			return nil, fmt.Errorf("CRASH highlight must include 1 or 2 drivers")
		}
	case CRASHPORTION:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("CRASHPORTION highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("CRASH au tour %d: Le pilote %s s'est crashé sur la portion %s", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[0].Position.Id)
		}
	case OVERTAKE:
		if len(drivers) != 2 {
			return nil, fmt.Errorf("OVERTAKE highlight must have 2 drivers")
		}
		desc = fmt.Sprintf("DEPASSEMENT au tour %d: Le pilote %s a réussi son dépassement sur %s", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[1].Driver.Lastname)
	case FINISH:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("FINISH highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("ARRIVEE: Le pilote %s est arrivé!", drivers[0].Driver.Lastname)
		}
	case DRIVER_PITSTOP:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("DRIVER_PITSTOP highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("PITSTOP: Le pilote %s est rentré au stand", drivers[0].Driver.Lastname)
		}
	case DRIVER_PITSTOP_CHANGETYRE:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("DRIVER_PITSTOP_CHANGETYRE highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("PITSTOP: Le pilote %s est rentré au stand (changement pneu %s vers %s)", drivers[0].Driver.Lastname, tyreTypes[drivers[0].PrevTyre], tyreTypes[drivers[0].CurrentTyre])
		}
	case CREVAISON:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("CREVAISON highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("CREVAISON: Le pilote %s a crevé au tour %d", drivers[0].Driver.Lastname, drivers[0].NbLaps)
		}
	}
	return &Highlight{
		Description: desc,
		Drivers:     drivers,
		Type:        highlightType,
	}, nil
}

func (h *Highlight) String() string {
	return h.Description
}
