import "../App.css";
import {useEffect, useState} from 'react';
import React from "react";
import LineChart from "../utils/LineChart.js";
import { teamColorsMap } from "../utils/teamsColor.js";
import { c } from "tar";


  

const DriversTotalLineChart = ({statistics }) => {

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
      //Step 1 : update dataset
      let res = datasets;
        if (res.length === 0) {
          data.totalStatistics.driversTotalPoints.forEach((driver) => {
            res.push({
              label: "Total of points " + driver.driver,
              data: [driver.totalPoints],
              backgroundColor: teamColorsMap[driver.driver],
              borderColor: teamColorsMap[driver.driver],
              borderWidth: 2,
              fill : false
            });
          })
        }
        else {
          // Test if lastChampionship is null otherwise too many 0 could be added
          if (data.lastChampionship !== ""){
            res.forEach((d) => {
            data.totalStatistics.driversTotalPoints.forEach((driver) => {
              if (d.label === "Total of points " + driver.driver) {
                d.data.push(driver.totalPoints)
              }
            })
          })
          } else {
            // Championship was reinit, reinitializing datasets
            res.forEach((d) => {
              data.totalStatistics.driversTotalPoints.forEach((driver) => {
                if (d.label === "Total of points " + driver.team) {
                  d.data = [driver.totalPoints]
                  d.backgroundColor = teamColorsMap[driver.team]
                  d.borderColor = teamColorsMap[driver.team]
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
      <LineChart title={"Total of points by driver"} chartData={lineChart} />
    );
};

export default DriversTotalLineChart