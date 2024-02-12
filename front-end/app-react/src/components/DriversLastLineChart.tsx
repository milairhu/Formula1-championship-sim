import "../App.css";
import {useEffect, useState} from 'react';
import React from "react";
import LineChart from "../utils/LineChart.js";
import { teamColorsMap } from "../utils/teamsColor.js";


  

const DriversLastLineChart = ({statistics }) => {

    const [lineChart, setLineChart] = useState<any>({
      labels: [],
      datasets: [
        {
          
        }
      ],
    });

    const [datasets, setDatasets] = useState<any[]>([]);

    
    useEffect(()=>{
      if (statistics){
        updateDatasets(statistics)
      }
  },[statistics])


    function updateDatasets(data) {
      //Etape 1 : mise à jour du dataset
      let res = datasets;
        if (res.length === 0) {
          data.lastChampionshipStatistics.driversTotalPoints.forEach((driver) => {
            res.push({
              label: "Total de points " + driver.driver,
              data: [driver.totalPoints],
              backgroundColor: teamColorsMap[driver.driver],
              borderColor: teamColorsMap[driver.driver],
              borderWidth: 2,
              fill : false
            });
          })
        }
        else {
          //Obligé de tester si lastChampionship est vide car sinon on risque d'ajouter des 0 inutiles
          if (data.lastChampionship !== ""){
            res.forEach((d) => {
            data.lastChampionshipStatistics.driversTotalPoints.forEach((driver) => {
              if (d.label === "Total de points " + driver.driver) {
                d.data.push(driver.totalPoints)
              }
            })
          })
          } else {
            //on a réinitiliaisé le championnat, on réinitialise les datasets
            res.forEach((d) => {
              data.lastChampionshipStatistics.driversTotalPoints.forEach((driver) => {
                if (d.label === "Total de points " + driver.team) {
                  d.data = [driver.totalPoints]
                  d.backgroundColor = teamColorsMap[driver.team]
                  d.borderColor = teamColorsMap[driver.team]
                }
              })
            })
          }
        }
      setDatasets(res)

      //Etape 2 : mise à jour du lineChart
      setLineChart(() => {
        let labels = lineChart.labels;
        if (data.lastChampionship == "" ) {
          labels = []
        } 
        else{
          labels = lineChart.labels.concat(data.lastChampionship)
        }
        return {
          labels: labels,
          datasets: res
          
        };
      })
      }



    return (
      <LineChart title={"Points des pilotes dans chaque championnat"} chartData={lineChart} />
    );
};

export default DriversLastLineChart