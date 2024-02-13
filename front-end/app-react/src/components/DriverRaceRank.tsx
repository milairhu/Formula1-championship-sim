import "../App.css";
import { useEffect, useState } from "react";
import React from "react";

const DriverRaceRank = ({ raceStatistics, championshipStatistics }) => {
  const [rows, setRows] = useState<any[]>([]);

  useEffect(() => {
    if (raceStatistics) {
      setRows(() => {
        let res = [];
        raceStatistics.forEach((driver) => {
          let champPoints = 0;
          championshipStatistics.forEach((champDriver) => {
            if (driver.driver === champDriver.driver) {
              champPoints = champDriver.totalPoints;
            }
          });
          res.push({
            name: driver.driver,
            points: driver.totalPoints,
            champPoints,
          });
        });
        //sort drivers by points
        res.sort((a, b) => {
          return b.points - a.points;
        });
        return res;
      });
    }
  }, [raceStatistics]);

  return (
    <table className="w-full h-full table-auto border border-gray-300 p-2">
      <thead>
        <tr>
          <th className="text-left pl-1 ">Rank</th>
          <th className="text-center">Driver</th>
          <th className="text-center">Tot. Points</th>
          <th className="text-center">Points at championship</th>
        </tr>
      </thead>
      <tbody>
        {rows.map((row, index) => (
          <tr
            key={index}
            className={
              (index % 2 === 0 ? "bg-white" : "bg-gray-100") +
              " border border-gray-300 p-2"
            }
          >
            <td className="text-left border border-gray-300 pl-1">
              {index + 1}
            </td>
            <td className="text-center border border-gray-300 py-1">
              {row.name}
            </td>
            <td className="text-center border border-gray-300 py-1">
              {row.points}
            </td>
            <td className="text-center border border-gray-300 py-1">
              {row.champPoints}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default DriverRaceRank;
