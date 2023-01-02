// ** MUI Imports
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import Avatar from "@mui/material/Avatar";
import CardHeader from "@mui/material/CardHeader";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import CardContent from "@mui/material/CardContent";

// ** Icons Imports
import DotsVertical from "mdi-material-ui/DotsVertical";
import TrendingUp from "mdi-material-ui/TrendingUp";
import Alert from "mdi-material-ui/Alert";
import AccountOutline from "mdi-material-ui/AccountOutline";

import { StatusesAggregateType } from "../../@core/api/FeedsAPI";

import { useQuery, gql } from "@apollo/client";

const GET_STATS = gql`
  query GetFeedStatusesAggregate {
    feedstatuseaggregates(startDate: "2000-01-01", endDate: "2099-01-01") {
      files
      rows
      errors
    }
  }
`;

const StatisticsCard = () => {
  const { loading, error, data } = useQuery<StatusesAggregateType>(GET_STATS);
  let statsBody: JSX.Element = <></>;

  if (error) {
    console.log(error);
    statsBody = (
        <p>Error: {error.message}</p>
    );
  }

  if (loading) {
    statsBody = (
        <p>Loading...</p>
    );
  }

  statsBody = ( 
  <>
    <Grid item xs={12} sm={2} key={1}>
      <Box key={1} sx={{ display: "flex", alignItems: "left" }}>
        <Avatar
          variant="rounded"
          sx={{
            mr: 3,
            width: 44,
            height: 44,
            boxShadow: 3,
            color: "common.white",
            backgroundColor: `primary.main`,
          }}
        >
          <TrendingUp sx={{ fontSize: "1.75rem" }} />
        </Avatar>
        <Box sx={{ display: "flex", flexDirection: "column" }}>
          <Typography variant="caption">Rows</Typography>
          <Typography variant="h6">{data?.feedstatuseaggregates.rows}</Typography>
        </Box>
      </Box>
  </Grid>
  <Grid item xs={12} sm={2} key={2}>
      <Box key={2} sx={{ display: "flex", alignItems: "left" }}>
        <Avatar
          variant="rounded"
          sx={{
            mr: 3,
            width: 44,
            height: 44,
            boxShadow: 3,
            color: "common.white",
            backgroundColor: `success.main`,
          }}
        >
          <AccountOutline sx={{ fontSize: "1.75rem" }} />
        </Avatar>
        <Box sx={{ display: "flex", flexDirection: "column" }}>
          <Typography variant="caption">Files</Typography>
          <Typography variant="h6">{data?.feedstatuseaggregates.files}</Typography>
        </Box>
      </Box>
  </Grid>
  <Grid item xs={12} sm={2} key={3}>
    <Box key={3} sx={{ display: "flex", alignItems: "left" }}>
      <Avatar
        variant="rounded"
        sx={{
          mr: 3,
          width: 44,
          height: 44,
          boxShadow: 3,
          color: "common.white",
          backgroundColor: `warning.main`,
        }}
      >
        <Alert sx={{ fontSize: "1.75rem" }} />
      </Avatar>
      <Box sx={{ display: "flex", flexDirection: "column" }}>
        <Typography variant="caption">Errors</Typography>
        <Typography variant="h6">{data?.feedstatuseaggregates.errors}</Typography>
      </Box>
    </Box>
  </Grid>
  </>
  )

  return (
    <Card>
      <CardHeader
        title="Data Manager"
        action={
          <IconButton
            size="small"
            aria-label="settings"
            className="card-more-options"
            sx={{ color: "text.secondary" }}
          >
          <DotsVertical />
          </IconButton>
        }
        titleTypographyProps={{
          sx: {
            mb: 2.5,
            lineHeight: "2rem !important",
            letterSpacing: "0.15px !important",
          },
        }}
      />
      <CardContent sx={{ pt: (theme) => `${theme.spacing(3)} !important` }}>
        <Grid container spacing={[5, 0]}>
          { statsBody }
        </Grid>
      </CardContent>
    </Card>
  );
};

export default StatisticsCard;