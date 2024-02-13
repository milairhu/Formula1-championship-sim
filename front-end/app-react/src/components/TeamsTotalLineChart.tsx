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
      // Step 1 : update dataset
      let res = datasets;
        if (res.length === 0) {
          data.totalStatistics.teamsTotalPoints.forEach((team) => {
            res.push({
              label: "Total of points " + team.team,
              data: [team.totalPoints],
              backgroundColor: teamColorsMap[team.team],
              borderColor: teamColorsMap[team.team],
              borderWidth: 2,
              fill : false
            });
          })
        }
        else {
          if (data.lastChampionship !== ""){
            res.forEach((team) => {
            data.totalStatistics.teamsTotalPoints.forEach((teamData) => {
              if (team.label === "Total of points " + teamData.team) {
                team.data.push(teamData.totalPoints)
              }
            })
          })
          } else {
            res.forEach((team) => {
              data.totalStatistics.teamsTotalPoints.forEach((teamData) => {
                if (team.label === "Total of points " + teamData.team) {
                  team.data = [teamData.totalPoints]
                  team.backgroundColor = teamColorsMap[teamData.team]
                  team.borderColor = teamColorsMap[teamData.team]
                }
              })
            })
          }
        }
      setDatasets(res)

      //Step 2 : update lineChart
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
      <LineChart title={"Total of points per teams"} chartData={lineChart} />
    );
};

export default TeamsTotalLineChart