import { Box, Grid, Typography } from "@mui/material";
import RaceHighlightRow from "./RaceHighlightRow.tsx";
import React from "react";

const RaceHighlights = ({ highlights }) => {
  return (
    <Box
      sx={{
        border: "1px solid #D1D5DB",
        width: "50rem",
        margin: "1rem",
        alignSelf: "center",
      }}
    >
      <Typography
        variant="h6"
        align="center"
        sx={{ borderBottom: "1px solid #D1D5DB " }}
      >
        {highlights && "Temps fort de la course"}
      </Typography>
      <Grid
        sx={{
          overflowY: "scroll",
        }}
        maxHeight={"45rem"}
        container
        flexDirection={"column"}
        // alignItems={"stretch"}
        flexWrap={"nowrap"}
      >
        {highlights?.map((highlight) => {
          return <RaceHighlightRow highlight={highlight} />;
        })}
      </Grid>
    </Box>
  );
};

export default RaceHighlights;
