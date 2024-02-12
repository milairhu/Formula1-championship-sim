import { Bar } from "react-chartjs-2";
const BarChart = ({ title,  chartData }) => {
  return (
      <Bar
        data={chartData}
        options={{
          plugins: {
            title: {
              display: true,
              text: title
            },
            legend: {
              display: false
            },
          },
          scales: {
            x: {
              display: false,
            },
          },
        }}
      />
  );
};

export default BarChart;