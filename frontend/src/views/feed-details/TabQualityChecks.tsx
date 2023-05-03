import * as React from "react";
// ** MUI Imports
import CardContent from "@mui/material/CardContent";
import Chip from "@mui/material/Chip";
import Select from "@mui/material/Select";
import Table from "@mui/material/Table";
import TableRow from "@mui/material/TableRow";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import MenuItem from "@mui/material/MenuItem";
import InputLabel from "@mui/material/InputLabel";

// import { DataGrid, GridToolbar } from '@mui/x-data-grid';
// import { useDemoData } from '@mui/x-data-grid-generator';


import {
  FeedQualityCheckType,
  FeedQualityChecksType,
  QualityCheckResultsType,
  statusObj,
} from "../../@core/api/FeedsAPI";

import { useQuery, gql } from "@apollo/client";

const GET_QUALITY_CHECKS = gql`
  query GetFeedQualityChecks {
    feedqualitychecks {
      id
      feed_name
      vendor
      quality_check_results
    }
  }
`;

const headerValues: string[] = [
  "Feed",
  "Vendor",
  "Quality Check Results"
];

const header: JSX.Element[] = [];
headerValues.forEach((value) => {
  header.push(<TableCell>{value}</TableCell>);
});

function openFeedDetails(feed_id: number) {
  // TODO have this expand the feed details
  alert(feed_id);
}

function displayQualityResults(quality_check_results: QualityCheckResultsType[]): JSX.Element[] {
  const qualityResults: JSX.Element[] = [];
  quality_check_results.forEach((qualityCheck: QualityCheckResultsType) =>
    qualityResults.push(
      <Chip
        color={statusObj[qualityCheck.status].color}
        onClick={() => openFeedDetails(qualityCheck.id)}
        sx={{height: 24}}
      />
    )
  )
  return qualityResults;
}

const TabQualityChecks = () => {
  const { loading, error, data } = useQuery<FeedQualityChecksType>(GET_QUALITY_CHECKS);
  const tableBody: JSX.Element[] = [];

  // if (error) {
  //   console.log(error);
  //   tableBody.push(<>Error: {error.message}</>);
  // }

  // START TESTING ONLY
  if (error) {
    const quality_check_quality_check_results: QualityCheckResultsType[] = [
      {
        id: 1,
        status: "Success",
        quality_check_name: "Null Check",
        quality_check_description: "All Records Passed"
      },
      {
        id: 2,
        status: "Failed",
        quality_check_name: "Unique Check",
        quality_check_description: "Column: ID. 27 Duplicates found"
      },
      {
        id: 3,
        status: "Success",
        quality_check_name: "Range Check",
        quality_check_description: "All Records Passed"
      }
    ];
    

  tableBody.push(<TableRow hover key="1">
    <TableCell>Orders</TableCell>
    <TableCell>Coyne Enterprises</TableCell>
    <TableCell>{displayQualityResults(quality_check_quality_check_results)}</TableCell>
  </TableRow>);
  }
  // END TESTING ONLY


  if (loading) {
    tableBody.push(
      <TableRow>
        <TableCell>Loading...</TableCell>
      </TableRow>
    );
  }

  data?.feedqualitycheck.forEach((row: FeedQualityCheckType) => 
    tableBody.push(
      <TableRow hover key={row.id}>
        <TableCell>{row.feed_name}</TableCell>
        <TableCell>{row.vendor}</TableCell>
        <TableCell>{displayQualityResults(row.quality_check_results)}</TableCell>
      </TableRow>
    )
  );


  // const VISIBLE_FIELDS = ['name', 'rating', 'country', 'dateCreated', 'isAdmin'];
  // const ddata = useDemoData({
  //   dataSet: 'Employee',
  //   visibleFields: VISIBLE_FIELDS,
  //   rowLength: 100,
  // });

  return (
    <CardContent>
    {/* <div style={{ height: 400, width: '100%' }}>
      <DataGrid {...ddata.data} components={{ Toolbar: GridToolbar }} />
    </div> */}
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

export default TabQualityChecks;
