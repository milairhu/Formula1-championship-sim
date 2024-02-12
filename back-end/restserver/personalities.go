package restserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

// Décodage de la requête /personalities
func (*RestServer) decodeUpdatePersonalityRequest(r *http.Request) (req []types.UpdatePersonalityInfo, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

// Obtenir les personnalités d'une simulation
func (rsa *RestServer) getAndUpdatePersonalities(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // Obtenir les personnalités
		fmt.Println("GET /personalities")
		driversInfosPersonalities := make([]types.PersonalityInfo, 0)

		for _, team := range rsa.pointTabTeam {
			team := *team
			for _, driver := range team.Drivers {
				driverInfo := types.PersonalityInfo{
					IdDriver:    driver.Id,
					Lastname:    driver.Lastname,
					Personality: driver.Personality.TraitsValue,
				}
				driversInfosPersonalities = append(driversInfosPersonalities, driverInfo)
			}

		}

		serial, _ := json.Marshal(driversInfosPersonalities)
		w.WriteHeader(http.StatusOK)
		w.Write(serial)
		return
	} else if r.Method == "PUT" { // Mettre à jour les personnalités
		fmt.Println("PUT /personalities")
		// décodage de la requête
		req, err := rsa.decodeUpdatePersonalityRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}

		// Réponse à renvoyer
		var resp []types.UpdatePersonalityInfo

		// Parcours des équipes et des pilotes
		for _, team := range rsa.pointTabTeam {
			for i := 0; i < 2; i++ {
				for _, updateReq := range req {
					if updateReq.IdDriver == team.Drivers[i].Id {
						// Test des valeurs des différentes personnalités
						if updateReq.Personality["Aggressivity"] > 5 || updateReq.Personality["Aggressivity"] < 1 ||
							updateReq.Personality["Confidence"] > 5 || updateReq.Personality["Confidence"] < 1 ||
							updateReq.Personality["Docility"] > 5 || updateReq.Personality["Docility"] < 1 ||
							updateReq.Personality["Concentration"] > 5 || updateReq.Personality["Concentration"] < 1 {
							msg := fmt.Sprintf("Un champs de %s de %s n'est pas compris entre 1 et 5", updateReq.Personality, updateReq.IdDriver)
							w.WriteHeader(http.StatusBadRequest)
							serial, _ := json.Marshal(msg)
							w.Write(serial)
							return
						} else {
							// Mise à jour des valeurs de personnalité du pilote
							team.Drivers[i].Personality.TraitsValue = map[string]int{
								"Aggressivity":  updateReq.Personality["Aggressivity"],
								"Confidence":    updateReq.Personality["Confidence"],
								"Docility":      updateReq.Personality["Docility"],
								"Concentration": updateReq.Personality["Concentration"],
							}
						}

						// Remplissage de la réponse
						resp = append(resp, types.UpdatePersonalityInfo{
							IdDriver:    team.Drivers[i].Id,
							Personality: team.Drivers[i].Personality.TraitsValue,
						})
					}
				}
			}
		}

		serial, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusOK)
		w.Write(serial)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed for /personalities", r.Method)
		return
	}
}
