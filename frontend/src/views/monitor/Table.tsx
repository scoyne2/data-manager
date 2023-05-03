import * as React from 'react';

// ** MUI Imports
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Chip from "@mui/material/Chip";
import Collapse from '@mui/material/Collapse';
import IconButton from '@mui/material/IconButton';
import Table from "@mui/material/Table";
import TableRow from "@mui/material/TableRow";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import Typography from "@mui/material/Typography";
import TableContainer from "@mui/material/TableContainer";
import UnfoldMoreHorizontal from "mdi-material-ui/UnfoldMoreHorizontal";
import UnfoldLessHorizontal from "mdi-material-ui/UnfoldLessHorizontal";


import {
  getColumns,
  FeedStatusType,
  FeedStatusesType,
  statusObj,
} from "../../@core/api/FeedsAPI";

import { useQuery, gql } from "@apollo/client";

const GET_FEED_STATUSES = gql`
  query GetFeedStatuses {
    feedstatuses {
      error_count
      feed_method
      feed_name
      id
      process_date
      record_count
      status
      vendor
    }
  }
`;

const GET_FEED_STATUS_DETAILS = gql`
  query GetFeedStatusDetails($id: Int!) {
    feedstatuses {
      error_count
      feed_method
      feed_name
      id
      process_date
      record_count
      status
      vendor
    }
  }
`;


function openFeedDetails(feed_id: number, location: string) {
  // TODO have this redirect to logs or quality checks
  alert(feed_id + location);
}

const headerValues = getColumns();

const header: JSX.Element[] = [];
headerValues.forEach((value) => {
  header.push(<TableCell>{value}</TableCell>);
});

function FeedDetailRow(feed_id: number){
  const { loading, error, data } = useQuery<FeedStatusesType>(GET_FEED_STATUS_DETAILS, { variables: {id: feed_id}});
  const tableDetailBody: JSX.Element[] = [];

  if (error) {
    console.log(error);
    tableDetailBody.push(
      <TableRow>
        <TableCell>Error: {error.message}</TableCell>
      </TableRow>
    );
  }

  if (loading) {
    tableDetailBody.push(
      <TableRow>
        <TableCell>Loading...</TableCell>
      </TableRow>
    );
  }

  if (!loading && !error) {
  data?.feedstatuses.forEach((row: FeedStatusType) =>
    tableDetailBody.push(
      <TableRow>
        <TableCell></TableCell>
        <TableCell>{row.feed_name}</TableCell>
        <TableCell>{row.vendor}</TableCell>
        <TableCell>{row.process_date}</TableCell>
        <TableCell>{row.record_count}</TableCell>
        <TableCell>{row.error_count}</TableCell>
        <TableCell>
          <Chip
            label={row.status}
            // color={statusObj[row.status].color}
            sx={{
              height: 24,
              fontSize: "0.75rem",
              textTransform: "capitalize",
              "& .MuiChip-label": { fontWeight: 500 },
            }}
          />
        </TableCell>
        <TableCell onClick={() => openFeedDetails(row.id, "logs")}>Link</TableCell>
        <TableCell onClick={() => openFeedDetails(row.id, "quality-checks")}>Link</TableCell>
      </TableRow>
      ));
    }

  return tableDetailBody
}

function FeedRow(row: FeedStatusType) {
  const [open, setOpen] = React.useState(false);
  // const detailHeader = header.shift();
  return (
    <React.Fragment>
      <TableRow
        hover
        key={row.id}
        sx={{ "&:last-of-type td, &:last-of-type th": { border: 0 } }}
      >
        <TableCell>
          <IconButton
              aria-label="expand row"
              size="small"
              onClick={() => setOpen(!open)}
            >
              {open ? <UnfoldMoreHorizontal /> : <UnfoldLessHorizontal />}
          </IconButton>
        </TableCell>
        <TableCell sx={{ py: (theme) => `${theme.spacing(0.5)} !important` }}>
          <Box sx={{ display: "flex", flexDirection: "column" }}>
            <Typography
              sx={{ fontWeight: 500, fontSize: "0.875rem !important" }}
            >
              {row.feed_name}
            </Typography>
            <Typography variant="caption">{row.feed_method}</Typography>
          </Box>
        </TableCell>
        <TableCell>{row.vendor}</TableCell>
        <TableCell>{row.process_date}</TableCell>
        <TableCell>{row.record_count}</TableCell>
        <TableCell>{row.error_count}</TableCell>
        <TableCell>
          <Chip
            label={row.status}
            color={statusObj[row.status].color}
            sx={{
              height: 24,
              fontSize: "0.75rem",
              textTransform: "capitalize",
              "& .MuiChip-label": { fontWeight: 500 },
            }}
          />
        </TableCell>
        <TableCell onClick={() => openFeedDetails(row.id, "logs")}>Link</TableCell>
        <TableCell onClick={() => openFeedDetails(row.id, "quality-checks")}>Link</TableCell>
        </TableRow>

        {/* Feed Details */}
        <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
        <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                Details
              </Typography>
              <Table sx={{ minWidth: 800 }} aria-label="details">
                <TableHead>
                  <TableRow>
                    {/* {{ detailHeader }} */}
                  </TableRow>
                </TableHead>
                <TableBody>
                 {/* { FeedDetailRow(row.id) } */}
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
}


const DashboardTable = () => {
  const { loading, error, data } = useQuery<FeedStatusesType>(GET_FEED_STATUSES);
  const tableBody: JSX.Element[] = [];

  if (error) {
    console.log(error);
    tableBody.push(
      <TableRow>
        <TableCell>Error: {error.message}</TableCell>
      </TableRow>
    );
  }

  if (loading) {
    tableBody.push(
      <TableRow>
        <TableCell>Loading...</TableCell>
      </TableRow>
    );
  }

  if (!loading && !error) {
    data?.feedstatuses.forEach((row: FeedStatusType) =>
      tableBody.push(
        FeedRow(row)
      )
    );
  }

  return (
    <Card>
      <TableContainer>
        <Table sx={{ minWidth: 800 }} aria-label="table in dashboard">
          <TableHead>
            <TableRow>{header}</TableRow>
          </TableHead>
          <TableBody>{tableBody}</TableBody>
        </Table>
      </TableContainer>
    </Card>
  );
};

export default DashboardTable;
