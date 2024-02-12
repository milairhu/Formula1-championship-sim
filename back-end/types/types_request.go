package types

type DriverRank struct {
	Rank        int            `json:"rank"`
	FirstName   string         `json:"firstname"`
	Lastname    string         `json:"lastname"`
	Points      int            `json:"points"`
	Personality map[string]int `json:"personality"`
}

func NewDriverRank(rank int, firstname, lastname string, pts int, personality map[string]int) *DriverRank {
	return &DriverRank{Rank: rank, FirstName: firstname, Lastname: lastname, Points: pts, Personality: personality}
}

type PersonalityInfo struct {
	IdDriver    string         `json:"idDriver"`
	Lastname    string         `json:"lastname"`
	Personality map[string]int `json:"personality"`
}

type UpdatePersonalityInfo struct {
	IdDriver    string         `json:"idDriver"`
	Personality map[string]int `json:"personality"`
}

type SimulateChampionship struct {
	LastChampionship           string                     `json:"lastChampionship"`
	NbSimulations              int                        `json:"nbSimulations"`
	TotalStatistics            TotalStatistics            `json:"totalStatistics"`
	LastChampionshipStatistics LastChampionshipStatistics `json:"lastChampionshipStatistics"`
}

type SimulateRace struct {
	Championship           string                 `json:"championship"`
	Race                   string                 `json:"race"`
	ChampionshipStatistics ChampionshipStatistics `json:"championshipStatistics"`
	RaceStatistics         RaceStatistics         `json:"raceStatistics"`
	Highlights             []RaceHighlight        `json:"highlights"`
	IsLastRace             bool                   `json:"isLastRace"`
}

type ChampionshipStatistics struct {
	DriversTotalPoints       []*DriverTotalPoints        `json:"driversTotalPoints"`
	TeamsTotalPoints         []*TeamTotalPoints          `json:"teamsTotalPoints"`
	PersonalityAveragePoints []*PersonalityAveragePoints `json:"personalityAveragePoints"`
	PersonalityAverage       map[string]map[int]float64  `json:"personalityAverage"`
	NbCrashsPersonnality     []*NbCrashsPersonnality     `json:"nbCrashsPersonnality"`
}

type RaceStatistics struct {
	DriversTotalPoints       []*DriverTotalPoints        `json:"driversTotalPoints"`
	TeamsTotalPoints         []*TeamTotalPoints          `json:"teamsTotalPoints"`
	PersonalityAveragePoints []*PersonalityAveragePoints `json:"personalityAveragePoints"`
	PersonalityAverage       map[string]map[int]float64  `json:"personalityAverage"`
	NbCrashsPersonnality     []*NbCrashsPersonnality     `json:"nbCrashsPersonnality"`
}

type TotalStatistics struct {
	DriversTotalPoints       []*DriverTotalPoints        `json:"driversTotalPoints"`
	TeamsTotalPoints         []*TeamTotalPoints          `json:"teamsTotalPoints"`
	PersonalityAveragePoints []*PersonalityAveragePoints `json:"personalityAveragePoints"`
	PersonalityAverage       map[string]map[int]float64  `json:"personalityAverage"`
	NbCrashsPersonnality     []*NbCrashsPersonnality     `json:"nbCrashsPersonnality"`
}
type LastChampionshipStatistics struct {
	DriversTotalPoints       []*DriverTotalPoints        `json:"driversTotalPoints"`
	TeamsTotalPoints         []*TeamTotalPoints          `json:"teamsTotalPoints"`
	PersonalityAveragePoints []*PersonalityAveragePoints `json:"personalityAveragePoints"`
	PersonalityAverage       map[string]map[int]float64  `json:"personalityAverage"` //pour réaliser les statistiques avec le script Python
	NbCrashsPersonnality     []*NbCrashsPersonnality     `json:"nbCrashsPersonnality"`
}

type RaceHighlight struct {
	Description string        // Describe the highlight
	Type        HighlightType // Type of highlight
}

func NewRaceHighlight(desc string, highlightType HighlightType) RaceHighlight {
	return RaceHighlight{Description: desc, Type: highlightType}
}

type DriverTotalPoints struct {
	Driver      string `json:"driver"`
	TotalPoints int    `json:"totalPoints"`
}

type TeamTotalPoints struct {
	Team        string `json:"team"`
	TotalPoints int    `json:"totalPoints"`
}

type PersonalityAveragePoints struct {
	Personality   map[string]int `json:"personality"`
	AveragePoints float64        `json:"averagePoints"`
	NbDrivers     int            `json:"nbDrivers"`
}

type NbCrashsPersonnality struct {
	Personality map[string]int `json:"personality"`
	NbCrash     int            `json:"nbCrash"`
}

func NewDriverTotalPoints(lastname string, pts int) *DriverTotalPoints {
	return &DriverTotalPoints{Driver: lastname, TotalPoints: pts}
}

func NewTeamTotalPoints(team string, pts int) *TeamTotalPoints {
	return &TeamTotalPoints{Team: team, TotalPoints: pts}
}

func NewPersonalityAveragePoints(personality map[string]int, pts int, nbDrivers int) *PersonalityAveragePoints {
	return &PersonalityAveragePoints{Personality: personality, AveragePoints: float64(pts), NbDrivers: nbDrivers}
}

func NewLastChampionshipStatistics(driversTotalPoints []*DriverTotalPoints, teamTotalPoints []*TeamTotalPoints, personalityTotalPoints []*PersonalityAveragePoints, personnalityAverage map[string]map[int]float64, nbCrashsPersonnality []*NbCrashsPersonnality) *LastChampionshipStatistics {
	return &LastChampionshipStatistics{DriversTotalPoints: driversTotalPoints, TeamsTotalPoints: teamTotalPoints, PersonalityAveragePoints: personalityTotalPoints, PersonalityAverage: personnalityAverage, NbCrashsPersonnality: nbCrashsPersonnality}
}

func NewSimulateChampionship(lastChampionship string, nbSim int, totalStatistics TotalStatistics, lastChampionshipStatistics LastChampionshipStatistics) *SimulateChampionship {
	return &SimulateChampionship{LastChampionship: lastChampionship, TotalStatistics: totalStatistics, LastChampionshipStatistics: lastChampionshipStatistics, NbSimulations: nbSim}
}

func NewRaceStatistics(driversTotalPoints []*DriverTotalPoints) *RaceStatistics {
	return &RaceStatistics{DriversTotalPoints: driversTotalPoints}
}

func NewChampionshipStatistics(driversTotalPoints []*DriverTotalPoints, teamTotalPoints []*TeamTotalPoints, personalityTotalPoints []*PersonalityAveragePoints, personnalityAverage map[string]map[int]float64, nbCrashsPersonnality []*NbCrashsPersonnality) *LastChampionshipStatistics {
	return &LastChampionshipStatistics{DriversTotalPoints: driversTotalPoints, TeamsTotalPoints: teamTotalPoints, PersonalityAveragePoints: personalityTotalPoints, PersonalityAverage: personnalityAverage, NbCrashsPersonnality: nbCrashsPersonnality}
}
