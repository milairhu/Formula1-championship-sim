import "../App.css";
import {useEffect, useState} from 'react';
import React from "react";
import LineChart from "../utils/LineChart.js";
import { teamColorsMap } from "../utils/teamsColor.js";


  

const TeamsLastLineChart = ({statistics }) => {
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
        data.lastChampionshipStatistics.teamsTotalPoints.forEach((team) => {
          res.push({
            label: "Total de points " + team.team,
            data: [team.totalPoints],
            backgroundColor: teamColorsMap[team.team],
            borderColor: teamColorsMap[team.team],
            borderWidth: 2,
            fill : false
          });
        })
      }
      else {
        //Obligé de tester si lastChampionship est vide car sinon on risque d'ajouter des 0 inutiles
        if (data.lastChampionship !== ""){
          res.forEach((d) => {
          data.lastChampionshipStatistics.teamsTotalPoints.forEach((team) => {
            if (d.label === "Total de points " + team.team) {
              d.data.push(team.totalPoints)
            }
          })
        })
        } else {
          //on a réinitiliaisé le championnat, on réinitialise les datasets
          res.forEach((d) => {
            data.lastChampionshipStatistics.teamsTotalPoints.forEach((team) => {
              if (d.label === "Total de points " + team.team) {
                d.data = [team.totalPoints]
                d.backgroundColor = teamColorsMap[team.team]
                d.borderColor = teamColorsMap[team.team]
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
      <LineChart title={"Points des équipes par championnat"} chartData={lineChart} />
    );
};

export default TeamsLastLineChart