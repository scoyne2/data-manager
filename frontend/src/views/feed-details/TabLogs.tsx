
// ** MUI Imports
import CardContent from "@mui/material/CardContent";
import Table from "@mui/material/Table";
import TableRow from "@mui/material/TableRow";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";

import {
  FeedLogsType,
  FeedLogType,
  statusObj,
} from "../../@core/api/FeedsAPI";

import { useQuery, gql } from "@apollo/client";

const GET_LOGS = gql`
  query GetFeedLogs {
    feedlogs {
      id
      status
      feed_name
      vendor
      process_date
      content
    }
  }
`;

const headerValues: string[] = [
  "Date",
  "Feed",
  "Vendor",
  "Content"
];

const header: JSX.Element[] = [];
headerValues.forEach((value) => {
  header.push(<TableCell>{value}</TableCell>);
});


const TabLogs = () => {
  const { loading, error, data } = useQuery<FeedLogsType>(GET_LOGS);
  const tableBody: JSX.Element[] = [];

  // if (error) {
  //   console.log(error);
  //   tableBody.push(<>Error: {error.message}</>);
  // }

  // START TESTING ONLY
  if (error) {
    
  tableBody.push(<>
  <TableRow hover key="1">
    <TableCell>2022-01-01 12:00:33</TableCell>
    <TableCell>Orders</TableCell>
    <TableCell>Coyne Enterprises</TableCell>
    <TableCell>Row 2 empty</TableCell>
  </TableRow>
  <TableRow hover key="2">
    <TableCell>2022-01-01 12:00:33</TableCell>
    <TableCell>Orders</TableCell>
    <TableCell>Coyne Enterprises</TableCell>
    <TableCell>Row 221 malformed</TableCell>
  </TableRow>
  <TableRow hover key="3">
    <TableCell>2022-01-01 12:00:33</TableCell>
    <TableCell>Orders</TableCell>
    <TableCell>Coyne Enterprises</TableCell>
    <TableCell>Row 52 malformed</TableCell>
  </TableRow>
  </>
  );
  }
  // END TESTING ONLY


  if (loading) {
    tableBody.push(
      <TableRow>
        <TableCell>Loading...</TableCell>
      </TableRow>
    );
  }

  data?.feedlogs.forEach((row: FeedLogType) => 
    tableBody.push(
      <TableRow hover key={row.id}>
        <TableCell>{row.process_date}</TableCell>
        <TableCell>{row.feed_name}</TableCell>
        <TableCell>{row.vendor}</TableCell>
        <TableCell>{row.content}</TableCell>
      </TableRow>
    )
  );

  return (
    <CardContent>
      <TableContainer>
        <Table sx={{ minWidth: 800 }} aria-label="Quality Check Results">
          <TableHead>
            <TableRow>{header}</TableRow>
          </TableHead>
          <TableBody>{tableBody}</TableBody>
        </Table>
      </TableContainer>
    </CardContent>
  );
};

export default TabLogs;
