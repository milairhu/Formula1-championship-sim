package main

import (
	"gitlab.utc.fr/vaursdam/formule-1-ia04/restserver"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/utils"
)

func main() {
	c, err := utils.ReadCircuit()
	if err != nil {
		panic(err)
	}

	t, initPersonalities, err := utils.ReadTeams()
	if err != nil {
		panic(err)
	}

	// We create pointers to the teams and the circuits
	pointTabCircuit := make([]*types.Circuit, len(c))
	for i, circuit := range c {
		tempCircuit := circuit // without buffer, all elements of the array will contain the same address
		pointTabCircuit[i] = &tempCircuit
	}
	initTeams := make([]types.Team, len(t))
	pointTabTeam := make([]*types.Team, len(t))
	for i, team := range t {
		tempTeam := team // without buffer, all elements of the array will contain the same address
		pointTabTeam[i] = &tempTeam
		initTeams[i] = tempTeam
	}

	// server launch
	server := restserver.NewRestServer(":8080", pointTabCircuit, pointTabTeam, initPersonalities)
	server.Start()
}
