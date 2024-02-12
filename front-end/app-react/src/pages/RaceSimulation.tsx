import React, { useEffect, useState } from "react";
import API from "../api";
import DriverRaceRank from "../components/DriverRaceRank.tsx";
import {
  ENDPOINT_RESET_RACE_SIMULATION,
  ENDPOINT_SIMULATE_RACE,
  ENDPOINT_STATISTICS_CHAMPIONSHIP,
  ENDPOINT_STATISTICS_RACE,
} from "../endpoints/endpoints.ts";
import PersonalityBarChart from "../components/PersonalityBarChart.tsx";
import PersonalityPieChart from "../components/PersonalityPieChart.tsx";
import RaceHighlights from "../components/RaceHighlights.tsx";
import { Pagination } from "@mui/material";

const RaceSimulation = () => {
  // useStates
  const [page, setPage] = useState<any>(0);
  const [simulationData, setSimulationData] = useState<any>();
  const [isRunning, setIsRunning] = useState(false);
  const [simulationHistory, setSimulationHistory] = useState<any[]>([]);
  // useEffects
  useEffect(() => {
    API.get(ENDPOINT_STATISTICS_RACE)
      .then((res) => {
        setSimulationData(res.data);
      })
      .catch((e) => {
        console.log("ERROR GET ", ENDPOINT_STATISTICS_RACE);
        console.log(e);
      });
  }, []);

  const handleReset = () => {
    API.get(ENDPOINT_RESET_RACE_SIMULATION)
      .then((res) => {
        setSimulationData(res.data);
      })
      .catch((e: Error) => {
        console.log("ERROR GET ", ENDPOINT_RESET_RACE_SIMULATION);
        console.log(e);
      });
    setSimulationHistory([]);
  };

  const handleSimulateRace = () => {
    API.get(ENDPOINT_SIMULATE_RACE)
      .then((res) => {
        setSimulationData(res.data);
        setSimulationHistory([...simulationHistory, res.data]);
        setPage(simulationHistory.length + 1);
      })
      .catch((e: Error) => {
        console.log("ERROR GET ", ENDPOINT_SIMULATE_RACE);
        console.log(e);
      });
  };

  const handlePageChange = (e, value) => {
    setPage(value);
  };

  return (
    <>
      <div className="w-100 text-center bg-gray-700 text-white text-xl border rounded-lg p-2">
        <span>Simulations course par course</span>
      </div>
      <div className="w-full mt-5 mb-2 flex pb-4 border-b-2 border-gray-500">
        <div className=" flex w-8/12  w-full items-center ">
          <div className="flex mr-1 ">
            <button
              disabled={simulationData?.isLastRace}
              onClick={handleSimulateRace}
              className={`${
                !isRunning
                  ? "bg-blue-500 hover:bg-blue-700 "
                  : " bg-red-500 hover:bg-red-700"
              }  text-white font-bold py-2 px-4 rounded`}
            >
              {"Lancer la simulation"}
            </button>
          </div>
        </div>
        <div className="w-4/12 flex ">
          <button
            onClick={handleReset}
            className={`${
              !isRunning
                ? "bg-green-500 hover:bg-green-700 "
                : " bg-gray-500 hover:bg-gray-700"
            } ml-auto text-white font-bold py-2 px-4 rounded`}
            disabled={isRunning}
          >
            Réinitialiser tout
          </button>
        </div>
      </div>

      <div className="text-lg mb-2">
        <div className="border-2 shadow-lg rounded-xl p-4 mt-4 mb-4">
          <div className="flex items-center justify-center w-full">
            <span className="bg-gray-500 text-white rounded-xl p-1 pl-3 pr-3">
              Résultats de la simulation : {simulationHistory[page - 1]?.race}{" "}
              {simulationData?.championship}{" "}
              {simulationData?.isLastRace && "(fin du championnat)"}
            </span>
          </div>
          <div className="flex justify-around items-center">
            <div className="w-4/12 text-sm">
              <DriverRaceRank
                raceStatistics={
                  simulationHistory.length === 0
                    ? simulationData?.raceStatistics.driversTotalPoints
                    : simulationHistory[page - 1]?.raceStatistics
                        .driversTotalPoints
                }
                championshipStatistics={
                  simulationHistory.length === 0
                    ? simulationData?.raceStatistics.driversTotalPoints
                    : simulationHistory[page - 1]?.championshipStatistics
                        .driversTotalPoints
                }
              />
            </div>
            <div>
              {simulationHistory.length !== 0 && (
                <RaceHighlights
                  highlights={simulationHistory[page - 1]?.highlights}
                />
              )}
            </div>
          </div>
          {/* Insert pagination here */}
          <div className="flex items-center justify-center w-full">
            <Pagination
              count={simulationHistory.length}
              page={page}
              onChange={handlePageChange}
            />
          </div>
        </div>
        <div className="border-2 shadow-lg rounded-xl p-4 mt-4 mb-4">
          <div className="flex items-center justify-center w-full">
            <span className=" bg-gray-500 text-white rounded-xl p-1 pl-3 pr-3">
              Personnalités
            </span>
          </div>
          <div className="flex justify-around items-end">
            <div className="w-7/12">
              <div className="text-sm">
                Nombre de profils différents étudiés :{" "}
                {
                  simulationHistory[page - 1]?.championshipStatistics
                    .personalityAveragePoints.length
                }
              </div>
              <PersonalityBarChart
                title={"Moyenne de points par profil"}
                personalityStatistics={
                  simulationHistory[page - 1]?.championshipStatistics
                    .personalityAveragePoints
                }
              />
            </div>
            <div className="w-4/12">
              <PersonalityPieChart
                personalityStatistics={
                  simulationHistory[page - 1]?.championshipStatistics
                    .personalityAveragePoints
                }
              />
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default RaceSimulation;
