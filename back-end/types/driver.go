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
	Driver        *Driver      //Pilote lui même
	Position      *Portion     //Position du pilote
	NbLaps        int          //Nombre de tours effectués
	Status        DriverStatus //Status du pilote
	IsPitStop     bool         // PitStop --> true if the driver is in pitstop
	TimeWoPitStop int          // Time without pitstop --> increments at each step
	ChanEnv       chan Action  // channel pour recevoir et envoyer les actions & l'environnement

	PitstopSteps int // Nombre de steps bloqué en pitstop

	PrevTyre       Tyre
	CurrentTyre    Tyre         //Type de pneu actuel
	AvailableTyres map[Tyre]int //Quantity of available tyres by type
	TyreTypeCount  int          //Nombre de différents types de pneus utilisés pendant la course
	UsedTyreTypes  []Tyre
	CurrentRank    int //Classement actuel du pilote, à voir si on l'implémente
	Speed          int //Vitesse du pilote
}

//Actions d'un pilote

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

// Fonction pour tester si un pilote réussit une portion sans se crasher
func (d *DriverInRace) PortionSuccess(pénalité int) bool {
	// Pour le moment on prend en compte le niveau du pilote, la difficulté de la portion et l'usure des pneus
	portion := d.Position

	// Si la vitesse est faible, on considère que la portion est réussie sans difficulté
	if d.Speed <= 3 {
		return true
	}

	probaReussite := 995
	probaReussite += d.Driver.Level * 20
	probaReussite -= portion.Difficulty * 18
	probaReussite -= d.TimeWoPitStop
	probaReussite -= pénalité
	probaReussite -= d.Speed * 5
	probaReussite += d.Driver.Personality.TraitsValue["Concentration"] * 10

	var dice int = rand.Intn(999) + 1

	if dice <= 10 && d.Speed < 10 {
		d.Speed += int(d.CurrentTyre)
	} else if dice >= 990 && d.Speed > 1 {
		d.Speed--
	}

	return dice <= probaReussite
}

func MakeSliceOfDriversInRace(teams []*Team, portionDepart *Portion, mapChan sync.Map, meteoCondition Meteo) ([]*DriverInRace, error) {
	res := make([]*DriverInRace, 0)
	for _, team := range teams {
		for _, driver := range team.Drivers {
			dtamp := driver //nécessaire, sinon n'utilise l'adresse que d'un membre de l'équipe
			c, ok := mapChan.Load(dtamp.Id)
			if !ok {
				return nil, fmt.Errorf("error while creating driver in race : %s", driver.Id)
			}
			d := NewDriverInRace(&dtamp, portionDepart, c.(chan Action), meteoCondition)
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

	probaPitStop := 0

	// On regarde si on doit faire un pitstop
	probaPitStop += (d.TimeWoPitStop * 10) / 2

	var dice int = rand.Intn(999) + 1

	if dice <= probaPitStop {
		d.PitstopSteps = 3
		d.TimeWoPitStop = 0
		return true
	}

	return false
}

func (d *DriverInRace) ChangeTyreType() bool {
	// On garde le pneus WET s'il pleut
	if d.CurrentTyre == WET {
		return false
	}

	// On regarde si on change le type de pneu, plus tendance si plus tard dans la course et types de pneus utilisés < 2
	probaChangeTyreType := 0
	TyreTypeCount := len(d.UsedTyreTypes)

	probaChangeTyreType = d.NbLaps + (3-TyreTypeCount)*20

	var dice = rand.Intn(99) + 1

	if dice < probaChangeTyreType {
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

func (d *DriverInRace) TestPneus() bool {

	if d.Status == PITSTOP || d.Status == PITSTOP_CHANGETYRE {
		return true
	}

	probaPneus := 0

	// On regarde si les pneus vont bien
	probaPneus += d.TimeWoPitStop - (int(d.CurrentTyre)*10 + 100)

	var dice int = rand.Intn(999) + 1

	return dice > probaPneus
}

func (d *DriverInRace) Overtake(otherDriver *DriverInRace) (reussite bool, crashedDrivers []*DriverInRace) {

	var probaDoubler int

	// Si l'autre pilote est en pitstop, on est sûr de doubler
	if otherDriver.Status == PITSTOP || otherDriver.Status == PITSTOP_CHANGETYRE {
		return true, []*DriverInRace{}
	}

	if d.Status == PITSTOP || d.Status == PITSTOP_CHANGETYRE {
		return false, []*DriverInRace{}
	}

	// On tente un doublement, l'usure des pneus augmente donc
	d.TimeWoPitStop += 10

	// En fonction de la valeur des traits de personnalité de confiance et de concentration du pilote, la probabilité de réussir un dépassement varie
	if d.Driver.Personality.TraitsValue["Confidence"] > 3 && d.Driver.Personality.TraitsValue["Concentration"] > 3 {
		probaDoubler = 75
	} else if d.Driver.Personality.TraitsValue["Confidence"] <= 3 && d.Driver.Personality.TraitsValue["Concentration"] >= 3 {
		probaDoubler = 65
	} else {
		probaDoubler = 50
	}

	if d.Driver.Level > otherDriver.Driver.Level {
		probaDoubler += 10
	} else if d.Driver.Level < otherDriver.Driver.Level {
		probaDoubler -= 10
	}

	portion := d.Position

	// Pour le moment on prend en compte le niveaus des pilotes et la "difficulté" de la portion
	probaDoubler -= portion.Difficulty * 2

	probaDoubler *= 10

	var dice int = rand.Intn(999) + 1

	// En fonction du résultat du dé, cela a un impact sur la vitesse du pilote
	if dice <= 10 && d.Speed < 10 {
		d.Speed += int(d.CurrentTyre)
	} else if dice >= 990 && d.Speed > 1 {
		d.Speed--
	}

	// Si on est en dessous de probaDoubler, on double et la confiance du pilote augmente
	if dice <= probaDoubler {
		if d.Driver.Personality.TraitsValue["Confidence"] < 5 {
			d.Driver.Personality.TraitsValue["Confidence"] += 1
			// Si la confiance du pilote est déjà au max et il réussit son dépassement, il devient moins docile
		} else if d.Driver.Personality.TraitsValue["Confidence"] == 5 && d.Driver.Personality.TraitsValue["Docility"] > 1 {
			d.Driver.Personality.TraitsValue["Docility"] -= 1
		}
		return true, []*DriverInRace{}
	}

	// Sinon, on regarde si on crash

	// Ici on a un échec critique, les deux pilotes crashent
	if dice >= 995 {
		return false, []*DriverInRace{d, otherDriver}
	}

	// Ici, un seul pilote crash, on tire au sort lequel
	if dice >= 990 {
		if dice%2 == 0 {
			return false, []*DriverInRace{d}
		} else {
			return false, []*DriverInRace{otherDriver}
		}
	}

	// Dans le cas par défaut, le doublement est échoué mais aucun crash n'a lieu
	// On reset la vitesse du pilote
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

	return nil, fmt.Errorf("Driver %s (%s, crashé si =1 : %d) who want to overtake is not found on portion %s", d.Driver.Id, d.Driver.Lastname, d.Status, p.Id)
}

// Fonction pour décider si on veut ESSAYER de doubler ou non
func (d *DriverInRace) OvertakeDecision(driverToOvertake *DriverInRace) (bool, error) {

	toOvertake, err := d.DriverToOvertake()
	if err != nil {
		return false, err
	}

	p := d.Position
	probaVeutDoubler := 0

	// L'aggressivité et la confiance du pilote impact les tentatives de dépassement
	if d.Driver.Personality.TraitsValue["Aggressivity"] > 3 || (d.Driver.Personality.TraitsValue["Aggressivity"] >= 3 && d.Driver.Personality.TraitsValue["Confidence"] >= 3) {
		probaVeutDoubler += d.Driver.Personality.TraitsValue["Aggressivity"] * 2
	} else {
		probaVeutDoubler -= (d.Driver.Personality.TraitsValue["Docility"]) * 2
	}

	if p.Difficulty != 0 {
		probaVeutDoubler += 20 / p.Difficulty
	} else {
		probaVeutDoubler = 0
	}

	if toOvertake != nil {
		// Si le pilote est en pitstop, on choisit de doubler
		if driverToOvertake.Status == PITSTOP || driverToOvertake.Status == PITSTOP_CHANGETYRE {
			return true, nil
		}

		var dice = rand.Intn(99) + 1
		return dice <= probaVeutDoubler, nil
	}
	return false, nil
}

func (d *DriverInRace) ChangeSpeed() {

	// On change la vitesse du pilote en fonction de sa personnalité
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
	log.Printf("		Lancement du pilote %s %s...\n", d.Driver.Firstname, d.Driver.Lastname)

	for {
		if d.Status == ARRIVED || d.Status == CRASHED {
			return
		}
		//On attend que l'environnement nous dise qu'on peut prendre une décision
		<-d.ChanEnv
		if d.Status == ARRIVED || d.Status == CRASHED {
			return
		}
		//On décide

		// On change la vitesse du pilote en fonction de sa personnalité
		d.ChangeSpeed()

		// On regarde si les pneus vont bien

		if !d.TestPneus() {
			d.ChanEnv <- ACCIDENTPNEUS
			continue
		}

		// On regarde si on doit faire un pitstop

		pitstop := false

		if d.Status != PITSTOP && d.Status != PITSTOP_CHANGETYRE {
			d.TimeWoPitStop++
			pitstop = d.PitStop()
		}

		if pitstop {
			// On regarde si on change le type de pneu
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

		//On regarde si on peut doubler
		toOvertake, err := d.DriverToOvertake()
		if err != nil {
			log.Printf("Error while getting the driver to overtake : %s\n", err)
		}
		if toOvertake != nil {
			//On décide si on veut doubler

			decision, err := d.OvertakeDecision(toOvertake)
			if err != nil {
				log.Printf("Error while getting the decision to overtake : %s\n", err)
			}
			if decision {
				//On envoie la décision à l'environnement
				d.ChanEnv <- TRY_OVERTAKE
			} else {
				d.ChanEnv <- CONTINUE
			}

		} else {
			//Si pas de possibilité de doubler, on ne fait rien
			d.ChanEnv <- CONTINUE
		}
		//On vérifie si on a fini la course
		if d.NbLaps == nbLaps {
			return
		}
	}
}
