import "../App.css";
import {useEffect, useState} from 'react';
import React from "react";


  

const PersonalityChampionshipRank = ({personalityStatistics }) => {

    const [rows, setRows] = useState<any[]>([])
    const [totPers, setTotPers] = useState<number>()
    
    useEffect(()=>{
        if (personalityStatistics){
            setRows(()=>{
                let res = []
                personalityStatistics.forEach((perso) =>{
                  let shortName = "Ag: " + perso.personality.Aggressivity + " ; Conc: " + perso.personality.Concentration  + " ; Conf: " + perso.personality.Confidence + " ; D: "+perso.personality.Docility  
                  
                  res.push({name : shortName, points: perso.averagePoints.toFixed(2), nbDrivers: perso.nbDrivers })
                })
                //trie les personnalités par points
                res.sort((a,b) => {
                    return b.points - a.points
                })
                return res
            })
        }
        
    },[personalityStatistics])

    useEffect(()=>{
      if (rows){
        setTotPers(()=>{
          let tot = 0
          rows.forEach((row)=>{
            tot += row.nbDrivers
          })
          return tot
        })
      }
    },[rows])

    return (
        <table className="w-full h-full table-auto border border-gray-300 p-2">
        <thead>
          <tr>
            <th className="text-left pl-1 ">Rang</th>
            <th className="text-center">Personnalité</th>
            <th className="text-center">Moyenne par champ.</th>
            <th className="text-center">Nb. Pilotes</th>
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
                {row.nbDrivers}
              </td>
            </tr>
          ))}
        </tbody>
        {totPers}
      </table>
     
    );
};

export default PersonalityChampionshipRank