import "../App.css";
import {useEffect, useRef, useState} from 'react';
import React from "react";
import PieChart from "../utils/PieChart";
import { generateRGB } from "../utils/generateRGB.js";


  

const PersonalityPieChart = ({personalityStatistics }) => {

    const [pieChart, setPieChart] = useState<any>({
      labels: [],
      datasets: [
        {
          
        }
      ],
    });

    function updatePieChart(data) {
      setPieChart(() => {
        return {
          labels: data.map((perso) => {
            return "Ag: " + perso.personality.Aggressivity + " ; Conc: " + perso.personality.Concentration  + " ; Conf: " + perso.personality.Confidence + " ; D: "+perso.personality.Docility  
          }),
          datasets: [
            {
              label: "Nombre d'apparitions",
              data: data.map((perso) => {
                return perso.nbDrivers
              }),
              backgroundColor: data.map((perso) => {
                return generateRGB(perso)
                
              }),
              borderColor: data.map((perso) => {
                return generateRGB(perso)
              }),
              borderWidth: 2,
              
            },
          ],
        };
      })
    }

    useEffect(()=>{
        if (personalityStatistics){
            updatePieChart(personalityStatistics)
        }
    },[personalityStatistics])

    return (
      <PieChart title={"Fréquences des profils"} chartData={pieChart} />
    );
};

export default PersonalityPieChart