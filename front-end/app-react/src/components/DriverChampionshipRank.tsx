import "../App.css";
import {useEffect, useState} from 'react';
import React from "react";


  

const DriverChampionshipRank = ({driversStatistics, lastChamp }) => {

    const [rows, setRows] = useState<any[]>([])
    
    useEffect(()=>{
        if (driversStatistics){
            setRows(()=>{
                let res = []
                driversStatistics.forEach((driver) =>{
                  let lastPoints = 0
                    lastChamp.forEach((lastDriver) => {
                        if (driver.driver === lastDriver.driver){
                            lastPoints = lastDriver.totalPoints
                        }
                    })
                    res.push({name : driver.driver, points: driver.totalPoints, lastPoints: lastPoints })
                })
                //trie les pilotes par points
                res.sort((a,b) => {
                    return b.points - a.points
                })
                return res
            })
        }
        
    },[driversStatistics])

    return (
        <table className="w-full h-full table-auto border border-gray-300 p-2">
        <thead>
          <tr>
            <th className="text-left pl-1 ">Rank</th>
            <th className="text-center">Driver</th>
            <th className="text-center">Tot. Points</th>
            <th className="text-center">Points at last champ.</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={index} className={(index % 2 === 0 ? 'bg-white' : 'bg-gray-100')+ " border border-gray-300 p-2"}>
              <td className="text-left border border-gray-300 pl-1">{index+1}</td>
              <td className="text-center border border-gray-300 py-1">
                {row.name}
              </td>
              <td className="text-center border border-gray-300 py-1">
                {row.points}
              </td>
              <td className="text-center border border-gray-300 py-1">
                {row.lastPoints}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    );
};

export default DriverChampionshipRank