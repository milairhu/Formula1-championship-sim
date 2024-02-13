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
			desc = fmt.Sprintf("CRASH on lap %d: Multiple drivers crashed: %s and %s during an overtaking attempt", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[1].Driver.Lastname)
		} else if len(drivers) == 1 {
			desc = fmt.Sprintf("CRASH on lap %d: Driver %s crashed during an overtaking attempt", drivers[0].NbLaps, drivers[0].Driver.Lastname)
		} else {
			return nil, fmt.Errorf("CRASH highlight must include 1 or 2 drivers")
		}
	case CRASHPORTION:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("CRASHPORTION highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("CRASH on lap %d: Driver %s crashed on portion %s", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[0].Position.Id)
		}
	case OVERTAKE:
		if len(drivers) != 2 {
			return nil, fmt.Errorf("OVERTAKE highlight must have 2 drivers")
		}
		desc = fmt.Sprintf("OVERTAKE on lap %d: Driver %s successfully overtook %s", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[1].Driver.Lastname)
	case FINISH:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("FINISH highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("FINISH: Driver %s finished the race!", drivers[0].Driver.Lastname)
		}
	case DRIVER_PITSTOP:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("DRIVER_PITSTOP highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("PITSTOP: Driver %s entered the pit stop", drivers[0].Driver.Lastname)
		}
	case DRIVER_PITSTOP_CHANGETYRE:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("DRIVER_PITSTOP_CHANGETYRE highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("PITSTOP: Driver %s entered the pit stop (changed tire from %s to %s)", drivers[0].Driver.Lastname, tyreTypes[drivers[0].PrevTyre], tyreTypes[drivers[0].CurrentTyre])
		}
	case CREVAISON:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("CREVAISON highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("TIRE BLOWOUT: Driver %s experienced a tire blowout on lap %d", drivers[0].Driver.Lastname, drivers[0].NbLaps)
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
