import "../App.css";
import { useEffect, useState } from "react";
import React from "react";
import BarChart from "../utils/BarChart";
import { generateRGB } from "../utils/generateRGB.js";

const PersonalityBarChart = ({ personalityStatistics, title }) => {
  const [barChart, setBarChart] = useState<any>({
    labels: [],
    datasets: [{}],
  });

  function updateBarChart(data) {
    setBarChart(() => {
      return {
        labels: data.map((perso) => {
          return (
            "Ag: " +
            perso.personality.Aggressivity +
            " ; Conc: " +
            perso.personality.Concentration +
            " ; Conf: " +
            perso.personality.Confidence +
            " ; D: " +
            perso.personality.Docility
          );
        }),
        datasets: [
          {
            label: "Moyenne par champ.",
            data: data.map((perso) => {
              return perso.averagePoints.toFixed(2);
            }),
            backgroundColor: data.map((perso) => {
              return generateRGB(perso);
            }),
            borderColor: data.map((perso) => {
              return generateRGB(perso);
            }),
            borderWidth: 2,
          },
        ],
        option: {
          scales: {
            xAxes: [
              {
                display: false, // Masquer les labels sous l'axe des x
              },
            ],
            yAxes: [
              {
                display: false,
              },
            ],
          },
        },
      };
    });
  }

  useEffect(() => {
    if (personalityStatistics) {
      updateBarChart(personalityStatistics);
    }
  }, [personalityStatistics]);

  return <BarChart title={title} chartData={barChart} />;
};

export default PersonalityBarChart;
