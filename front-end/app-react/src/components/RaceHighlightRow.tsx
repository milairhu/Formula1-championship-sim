import { Grid, Typography } from "@mui/material";
import React from "react";

const CRASHOVERTAKE = 0;
const CRASHPORTION = 1;
const OVERTAKE = 2;
const FINISH = 3;
const DRIVER_PITSTOP = 4;
const DRIVER_PITSTOP_CHANGETYRE = 5;
const CREVAISON = 6;

const RaceHighlightRow = ({ highlight }) => {
  let backgroundColor;

  switch (highlight.Type) {
    case CRASHOVERTAKE:
      backgroundColor = "#FFB4B4";
      break;
    case CRASHPORTION:
      backgroundColor = "#FFB4B4";
      break;
    case OVERTAKE:
      backgroundColor = "#bfffb4";
      break;
    case FINISH:
      backgroundColor = "#b4fdff";
      break;
    case DRIVER_PITSTOP:
      backgroundColor = "#fff3b4";
      break;
    case DRIVER_PITSTOP_CHANGETYRE:
      backgroundColor = "#fff3b4";
      break;
    case CREVAISON:
      backgroundColor = "#FFB4B4";
      break;
  }

  return (
    <Grid item sx={{ backgroundColor: backgroundColor }} padding={1}>
      <Typography>{highlight.Description}</Typography>
    </Grid>
  );
};

export default RaceHighlightRow;
