import "../App.css";
import {useEffect, useState} from 'react';
import React from "react";
import LineChart from "../utils/LineChart.js";
import { teamColorsMap } from "../utils/teamsColor.js";
import { c } from "tar";


  

const TeamsTotalLineChart = ({statistics }) => {

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
          //console.log("Res vaut 0")
          data.totalStatistics.teamsTotalPoints.forEach((team) => {
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
            //console.log("On a reçu un championnat")
            res.forEach((team) => {
            data.totalStatistics.teamsTotalPoints.forEach((teamData) => {
              if (team.label === "Total de points " + teamData.team) {
                //console.log("On a trouvé une équipe : " + teamData.team )
                team.data.push(teamData.totalPoints)
                //console.log(team.data)
              }
            })
          })
          } else {
            //on a réinitiliaisé le championnat, on réinitialise les datasets
            //console.log("On a un championnat vide")
            res.forEach((team) => {
              data.totalStatistics.teamsTotalPoints.forEach((teamData) => {
                if (team.label === "Total de points " + teamData.team) {
                  team.data = [teamData.totalPoints]
                  team.backgroundColor = teamColorsMap[teamData.team]
                  team.borderColor = teamColorsMap[teamData.team]
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
        //console.log(labels)
        //console.log(res)
        return {
          labels: labels,
          datasets: res
          
        };
      })
      }



    return (
      <LineChart title={"Total de points par équipes"} chartData={lineChart} />
    );
};

export default TeamsTotalLineChart