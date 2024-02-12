import "./App.css";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import { useState } from "react";
import LaunchSimulation from "./pages/LaunchSimulation.tsx"
import EditPersonnalities from "./pages/EditPersonnalities.tsx";
import RaceSimulation from "./pages/RaceSimulation.tsx";

const CustomTabPanel = ({ index, value, children }) => {
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
};

const App = () => {
  const [value, setValue] = useState(0);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  return (
    
    
    <Box sx={{ width: "100%" }}>
      <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
        <Tabs
          value={value}
          onChange={handleChange}
          aria-label="basic tabs example"
        >
          <Tab label="Simulations" />
          {/*<Tab label="???" />*/}
          <Tab label="Simulations course par course" />
          <Tab label="Modification personnalités" />
          {/*<Tab label="Stats générales" />*/}
        </Tabs>
      </Box>
      <CustomTabPanel value={value} index={0}>
        <LaunchSimulation/>
      </CustomTabPanel>
      <CustomTabPanel value={value} index={1}>
        <RaceSimulation/>
      </CustomTabPanel>
      <CustomTabPanel value={value} index={2}>
        <EditPersonnalities />
      </CustomTabPanel>
      
      {/*<CustomTabPanel value={value} index={2}>
        <SimulationResults />
  </CustomTabPanel>*/}
    </Box>
  );
};

export default App;
