package types

import "math/rand"

type Circuit struct {
	Id                string    // Circuit ID
	Name              string    // Circuit name
	Country           string    // Circuit country
	Portions          []Portion // Portions
	MeteoDistribution []int     // Distribution of meteo conditions
	NbLaps            int       //Number of Laps
}

func NewCircuit(id string, name string, country string, portions []Portion, meteoDistribution []int) *Circuit {

	p := make([]Portion, len(portions))
	copy(p, portions)

	return &Circuit{
		Id:                id,
		Name:              name,
		Country:           country,
		Portions:          p,
		MeteoDistribution: meteoDistribution,
	}
}

/*MeteoDistribution is built as follow : [40,70].
* Means that in 40% of the cases, weather is RAINY, 30% DRY, 30% HEAT
 */
func (c *Circuit) GenerateMeteo() Meteo {
	var dice int = rand.Intn(100) //generate number between 0 and 99
	if dice < c.MeteoDistribution[0] {
		return RAINY
	}
	if dice < c.MeteoDistribution[1] {
		return DRY
	}
	return HEAT
}
