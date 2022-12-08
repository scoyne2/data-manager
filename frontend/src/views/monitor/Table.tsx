// ** MUI Imports
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Chip from "@mui/material/Chip";
import Table from "@mui/material/Table";
import TableRow from "@mui/material/TableRow";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import Typography from "@mui/material/Typography";
import TableContainer from "@mui/material/TableContainer";
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

// const rows = getFeedStatus();
const headerValues = getColumns();

const header: JSX.Element[] = [];
headerValues.forEach((value) => {
  header.push(<TableCell>{value}</TableCell>);
});

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


  data?.feedstatuses.forEach((row: FeedStatusType) =>
    tableBody.push(
      <TableRow
        hover
        key={row.id}
        sx={{ "&:last-of-type td, &:last-of-type th": { border: 0 } }}
      >
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
      </TableRow>
    )
  );

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
