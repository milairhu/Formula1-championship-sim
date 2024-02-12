import React from "react";
import { Pie } from "react-chartjs-2";

function PieChart({ title, chartData }) {
  return (
      <Pie
        data={chartData}
        options={{
          plugins: {
            title: {
              display: true,
              text: title
            },
            legend: {
              display: false
            }
          }
        }}
      />
  );
}
export default PieChart;