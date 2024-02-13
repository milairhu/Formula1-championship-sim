import "../App.css";
import { FunctionComponent, useEffect, useState } from "react";
import API from "../api";
import React from "react";
import {
  ENDPOINT_RESET,
  ENDPOINT_SIMULATE_CHAMPIONSHIP,
  ENDPOINT_STATISTICS_CHAMPIONSHIP,
} from "../endpoints/endpoints.ts";
import Chart, { CategoryScale } from "chart.js/auto";
import DriverChampionshipRank from "../components/DriverChampionshipRank.tsx";
import TeamsChampionshipRank from "../components/TeamChampionshipRank.tsx";
import PersonalityBarChart from "../components/PersonalityBarChart.tsx";
import PersonalityPieChart from "../components/PersonalityPieChart.tsx";
import TeamsLastLineChart from "../components/TeamsLastLineChart.tsx";
import DriversLastLineChart from "../components/DriversLastLineChart.tsx";

Chart.register(CategoryScale);

const LaunchSimulation: FunctionComponent<{}> = ({}) => {
  // ==== UseStates ==== //
  const [simulationSummary, setSimulationSummary] = useState<any>();
  const [isRunning, setIsRunning] = useState<boolean>(false);

  // ===== useEffects ====//
  useEffect(() => {
    API.get(ENDPOINT_STATISTICS_CHAMPIONSHIP)
      .then((res) => {
        setSimulationSummary(res.data);
      })
      .catch((e) => {
        console.log("ERROR GET ", ENDPOINT_STATISTICS_CHAMPIONSHIP);
        console.log(e);
      });
  }, []);

  useEffect(() => {
    let intervalId;

    const fetchSimulationData = () => {
      API.get(ENDPOINT_SIMULATE_CHAMPIONSHIP)
        .then((res) => {
          setSimulationSummary(res.data);
        })
        .catch((e) => {
          console.log("ERROR GET ", ENDPOINT_SIMULATE_CHAMPIONSHIP);
          console.log(e);
        });
    };

    if (isRunning) {
      // Call fetchSimulationData immediately and then each 2 secondes
      fetchSimulationData();
      intervalId = setInterval(fetchSimulationData, 2000);
    }
    return () => clearInterval(intervalId);
  }, [isRunning]);

  //===== Handlers ====
  const handleSimulateAll = () => {
    setIsRunning(!isRunning);
  };

  const handleSimulateOne = () => {
    if (!isRunning) {
      API.get(ENDPOINT_SIMULATE_CHAMPIONSHIP)
        .then((res) => {
          setSimulationSummary(res.data);
        })
        .catch((e: Error) => {
          console.log("ERROR GET ", ENDPOINT_SIMULATE_CHAMPIONSHIP);
          console.log(e);
        });
    }
  };

  const handleReset = () => {
    setIsRunning(false);
    API.get(ENDPOINT_RESET)
      .then((res) => {
        setSimulationSummary(res.data);
      })
      .catch((e: Error) => {
        console.log("ERROR GET ", ENDPOINT_RESET);
        console.log(e);
      });
  };

  return (
    <>
      <div className="w-100 text-center bg-gray-700 text-white text-xl border rounded-lg p-2">
        <span>
          Statistiques sur les{" "}
          {simulationSummary?.nbSimulations
            ? simulationSummary.nbSimulations
            : 0}{" "}
          championnats simulés
        </span>
        <span className="text-sm">
          {" "}
          (dernier :{" "}
          {simulationSummary?.lastChampionship
            ? simulationSummary.lastChampionship
            : "N/A"}
          ) :{" "}
        </span>
      </div>
      <div className="w-full mt-5 mb-2 flex pb-4 border-b-2 border-gray-500">
        <div className=" flex w-8/12  w-full items-center ">
          <div className="flex mr-1 ">
            <button
              className={`${
                !isRunning
                  ? "bg-blue-500 hover:bg-blue-700 "
                  : " bg-red-500 hover:bg-red-700"
              }  text-white font-bold py-2 px-4 rounded`}
              onClick={handleSimulateAll}
            >
              {isRunning ? "Stop simulation" : "Launch  simulation"}
            </button>
          </div>
          <div className="ml-1">
            <button
              className={`${
                !isRunning
                  ? "bg-blue-500 hover:bg-blue-700 "
                  : " bg-gray-500 hover:bg-gray-700"
              }  text-white font-bold py-2 px-4 rounded`}
              onClick={handleSimulateOne}
              disabled={isRunning}
            >
              Simulate a single championship
            </button>
          </div>
        </div>
        <div className="w-4/12 flex ">
          <button
            className={`${
              !isRunning
                ? "bg-green-500 hover:bg-green-700 "
                : " bg-gray-500 hover:bg-gray-700"
            } ml-auto text-white font-bold py-2 px-4 rounded`}
            onClick={handleReset}
            disabled={isRunning}
          >
            Reinitialize all
          </button>
        </div>
      </div>

      <div className="text-lg mb-2">
        <div className="border-2 shadow-lg rounded-xl p-4 mt-4 mb-4">
          <div className="flex items-center justify-center w-full">
            <span className="bg-gray-500 text-white rounded-xl p-1 pl-3 pr-3">
              Drivers
            </span>
          </div>
          <div className="flex justify-around items-center">
            <div className="w-4/12 text-sm">
              <DriverChampionshipRank
                driversStatistics={
                  simulationSummary?.totalStatistics.driversTotalPoints
                }
                lastChamp={
                  simulationSummary?.lastChampionshipStatistics
                    .driversTotalPoints
                }
              />
            </div>
            <div className="w-6/12">
              {/*<div className="mb-4">
                <DriversTotalLineChart statistics={simulationSummary} />
            </div>*/}
              <div>
                <DriversLastLineChart statistics={simulationSummary} />
              </div>
            </div>
          </div>
        </div>
        <div className="border-2 shadow-lg rounded-xl p-4 mt-4 mb-4">
          <div className="flex items-center justify-center w-full">
            <span className=" bg-gray-500 text-white rounded-xl p-1 pl-3 pr-3">
              Teams
            </span>
          </div>
          <div className="flex justify-around items-center">
            <div className="w-6/12">
              {/*<div className="mb-4">
                <TeamsTotalLineChart statistics={simulationSummary} />
          </div>*/}
              <div>
                <TeamsLastLineChart statistics={simulationSummary} />
              </div>
            </div>
            <div className="w-4/12 text-sm">
              <TeamsChampionshipRank
                teamsStatistics={
                  simulationSummary?.totalStatistics.teamsTotalPoints
                }
                lastChamp={
                  simulationSummary?.lastChampionshipStatistics.teamsTotalPoints
                }
              />
            </div>
          </div>
        </div>
        <div className="border-2 shadow-lg rounded-xl p-4 mt-4 mb-4">
          <div className="flex items-center justify-center w-full">
            <span className=" bg-gray-500 text-white rounded-xl p-1 pl-3 pr-3">
              Personalities
            </span>
          </div>
          <div className="flex justify-around items-end">
            <div className="w-7/12">
              <div className="text-sm">
                Number of different profiles explored :{" "}
                {simulationSummary?.lastChampionship
                  ? simulationSummary?.totalStatistics.personalityAveragePoints
                      .length
                  : 0}
              </div>
              <PersonalityBarChart
                title={"Points average per profile per championship"}
                personalityStatistics={
                  simulationSummary?.totalStatistics.personalityAveragePoints
                }
              />
            </div>
            <div className="w-4/12">
              <PersonalityPieChart
                personalityStatistics={
                  simulationSummary?.totalStatistics.personalityAveragePoints
                }
              />
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default LaunchSimulation;
