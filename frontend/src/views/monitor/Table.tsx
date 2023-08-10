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

import DataPreviewModal from 'src/@core/components/data-preview';
import EMRLogs from 'src/@core/components/emr-logs';
import DataQualityResults from 'src/@core/components/data-quality';

import {
  getColumns,
  FeedStatusType,
  FeedStatusDetailedType,
  FeedStatusesDetailsType,
  statusObj,
} from "../../@core/api/FeedsAPI";

import { useQuery, gql } from "@apollo/client";

const GET_FEED_STATUSES = gql`
  query GetFeedStatuses {
    feedstatusesdetailed {
      error_count
      feed_method
      feed_name
      id
      process_date
      record_count
      status
      vendor
      previous_feeds {
        error_count
        feed_method
        feed_name
        file_name
        id
        process_date
        record_count
        status
        vendor
        emr_logs
        data_quality_url
      }
    }
  }
`;

const headerValues = getColumns();

const header: JSX.Element[] = [];
headerValues.forEach((value) => {
  header.push(<TableCell>{value}</TableCell>);
});

function DetailRows(previous_feeds: FeedStatusType[]){
  const detailRow: JSX.Element[] = [];
  previous_feeds.forEach((prow: FeedStatusType) => 
    detailRow.push(
        <TableRow key={prow.id}>
          <TableCell >{prow.file_name}</TableCell>
          <TableCell >{prow.process_date}</TableCell>
          <TableCell >{prow.record_count}</TableCell>
          <TableCell >{prow.error_count}</TableCell>
          <TableCell >
            <Chip
              label={prow.status}
              color={statusObj[prow.status].color}
              sx={{
                height: 24,
                fontSize: "0.75rem",
                textTransform: "capitalize",
                "& .MuiChip-label": { fontWeight: 500 },
              }}
            />
          </TableCell>
          <TableCell>
            <DataPreviewModal vendor={prow.vendor} feed_name={prow.feed_name} file_name={prow.file_name} />
          </TableCell>
          <TableCell>
            <EMRLogs log_url={prow.emr_logs} />
          </TableCell>
          <TableCell>
            <DataQualityResults data_quality_url={prow.data_quality_url} status={prow.status} />
          </TableCell>
        </TableRow>
      )
    )
    return detailRow
}

function FeedRow(row: FeedStatusDetailedType, open: boolean | undefined, setOpen: (key: number, value: boolean) => void) {  
  open = open ?? false;
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
              onClick={() => setOpen(row.id, !open)}
            >
              <UnfoldMoreHorizontal /> 
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
        </TableRow>
        {/* Feed Details */}
        <TableRow style={{ background: "#F8F8F8"}}>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0, background: "#F8F8F8"}} colSpan={4}>
        <Collapse in={open} timeout="auto" unmountOnExit>
            <Box >
              <Table aria-label="details">
              <TableHead>
                <TableRow>
                  <TableCell>File Name</TableCell>
                  <TableCell>Date</TableCell>
                  <TableCell>Rows</TableCell>
                  <TableCell>Errors</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Data Preview</TableCell>
                  <TableCell>EMR Logs</TableCell>
                  <TableCell>Data Quality</TableCell>
                  </TableRow>
              </TableHead>
                <TableBody>
                  { DetailRows(row.previous_feeds) }
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
  const { loading, error, data } = useQuery<FeedStatusesDetailsType>(GET_FEED_STATUSES);
  const tableBody: JSX.Element[] = [];
  const [open, setOpen] = React.useState(new Map<number, boolean>());
  const updateOpen = (key: number, value: boolean) => {
    setOpen(map => new Map(map.set(key, value)));
  }

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
    data?.feedstatusesdetailed.forEach((row: FeedStatusDetailedType) => (
      tableBody.push(FeedRow(row, open.get(row.id), updateOpen))
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
