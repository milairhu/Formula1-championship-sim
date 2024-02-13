import "../App.css";
import {FunctionComponent, useEffect, useState} from 'react';
import API from '../api.js';
import React from "react";
import { ENDPOINT_PERSONNALITIES } from "../endpoints/endpoints.ts";
import BarChart from "../utils/BarChart.js";


const EditPersonnalities: FunctionComponent<{}> = ({}) => {
  
const [drivers, setDrivers] = useState<any[]>([]);



useEffect(() => {
    API.get(ENDPOINT_PERSONNALITIES)
        .then((res) => {
          setDrivers(res.data);
          updatePersonalitiesChart(res.data);
        })
        .catch((e: Error) => {
            console.log("ERROR GET ", ENDPOINT_PERSONNALITIES);
            console.log(e);
        })
}, []);

function handleSubmit(event) {
  event.preventDefault();
  API.put(ENDPOINT_PERSONNALITIES, drivers)
  .then((res) => {
    setDrivers(drivers.map((driver) => {
      const matchingDriver = res.data.find((resDriver) => resDriver.idDriver === driver.idDriver);
      return {
        ...driver,
        personality: {
          ...matchingDriver.personality,
        },
      };
    }
    ));
    updatePersonalitiesChart(res.data);
  }
  )
  .catch((e: Error) => {
    console.log("ERROR PUT ", ENDPOINT_PERSONNALITIES);
    console.log(e);
  })
}

// Handle traits modification
const handleTraitChange = (driverIndex, traitName, value) => {
  setDrivers(() =>  
    drivers.map((driver, index) => {
      if (index === driverIndex) {
        return {
          ...driver,
          personality: {
            ...driver.personality,
            [traitName]: value,
          },
        };
      }
      return driver;
    })
  );
};

//== Chart //
const [personalitiesChart, setPersonalitiesChart] = useState<any>({
  labels: [],
  datasets: [
    {
      
    }
  ],
});

function updatePersonalitiesChart(driversData : any[]) {
  setPersonalitiesChart( () => {

    //Create map
    const personalityCountMap = new Map();
    driversData.forEach((driver) => {
        const configString = JSON.stringify(driver.personality);

        if (personalityCountMap.has(configString)) {
          personalityCountMap.set(configString, personalityCountMap.get(configString) + 1);
        } else {
          personalityCountMap.set(configString, 1);
        }
      }
    );
    const keys : string []= Array.from(personalityCountMap.keys());
    keys.sort();
    let labels = []
    keys.forEach((key) => {
      const personality = JSON.parse(key);
      labels.push(`Ag:${personality.Aggressivity} Conc:${personality.Concentration} Conf:${personality.Confidence} D:${personality.Docility}`);
    });
    const data = []
    for (let i = 0; i < keys.length; i++) {
      data.push(personalityCountMap.get(keys[i]));
    }
    return{
      labels: labels,
      datasets: [
        {
          label: "Number of pilotes",
          data: data,
          backgroundColor: [
            "red",
          ],
          borderColor: "red",
          borderWidth: 2,
          fill : false
        }
      ],
    }


  }
  );
}

  return (
    <>
      <div className="w-100 text-center bg-gray-700 text-white text-xl border rounded-lg p-2">
        Modifier les personnalités
      </div>
      <div className="w-100 flex justify-around h-full">
        <div className="w-4/12">
          <form onSubmit={handleSubmit} className=" mx-auto rounded-xl">
            <table className="w-full h-full overflow-auto my-8 table-auto border border-gray-300 p-2 ">
              <thead>
                <tr>
                  <th className="text-left pl-1">Pilote</th>
                  <th className="text-center">Agressivité</th>
                  <th className="text-center">Concentration</th>
                  <th className="text-center">Confiance</th>
                  <th className="text-center">Docilité</th>
                </tr>
              </thead>
              <tbody>
                {drivers.map((driver, index) => (
                  <tr key={index} className={(index % 2 === 0 ? 'bg-white' : 'bg-gray-100')+ " border border-gray-300 p-2  h-full"}>
                    <td className="text-left border border-gray-300 pl-1">{driver.lastname}</td>
                    <td className="text-center border border-gray-300 py-1">
                      <input
                        type="number"
                        min="1"
                        max="5"
                        className={(index % 2 === 0 ? 'bg-white' : 'bg-gray-100')+ " border border-gray-300 px-1 py-1 rounded  h-full"}
                        value={driver.personality.Aggressivity}
                        onChange={(e) => handleTraitChange(index, "Aggressivity", parseInt(e.target.value))}
                      />
                    </td>
                    <td className="text-center border border-gray-300 py-1">
                      <input
                        type="number"
                        min="1"
                        max="5"
                        className={(index % 2 === 0 ? 'bg-white' : 'bg-gray-100')+  " border border-gray-300 px-1 py-1 rounded h-full"}
                        value={driver.personality.Concentration}
                        onChange={(e) => handleTraitChange(index, "Concentration", parseInt(e.target.value))}
                      />
                    </td>
                    <td className="text-center border border-gray-300 py-1">
                      <input
                        type="number"
                        min="1"
                        max="5"
                        className={(index % 2 === 0 ? 'bg-white' : 'bg-gray-100')+ " border border-gray-300 px-1 py-1 rounded  h-full"}
                        value={driver.personality.Confidence}
                        onChange={(e) => handleTraitChange(index, 'Confidence', parseInt(e.target.value))}
                      />
                    </td>
                    <td className="text-center py-1">
                      <input
                        type="number"
                        min="1"
                        max="5"
                        className={(index % 2 === 0 ? 'bg-white' : 'bg-gray-100') + " border border-gray-300 px-1 py-1 rounded  h-full"}
                        value={driver.personality.Docility  }
                        onChange={(e) => handleTraitChange(index, 'Docility', parseInt(e.target.value))}
                      />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          <button className="mt-4 bg-blue-500 text-white px-4 py-2 rounded" type="submit">Enregistrer</button>
        </form>
      </div>
      <div className="w-7/12 flex items-center justify-center">
        <BarChart title={"Répartition des personnalités"} chartData={personalitiesChart}/>
      </div>
      
    
    </div>
  </>
  );
};

export default EditPersonnalities;
