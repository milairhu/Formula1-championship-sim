import React from "react";
import { Line } from "react-chartjs-2";

function LineChart({title,  chartData }) {
  return (
      <Line
      
        data={chartData}
        options={{
          plugins: {
            title: {
              display: true,
              text: title,
            },
            legend: {
              display: false
            }
          }
        }}
      />
  );
}
export default LineChart;