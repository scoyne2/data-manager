// ** MUI Imports
import Alert from "@mui/material/Alert";
import { AlertColor } from "@mui/material";
import Card from "@mui/material/Card";
import Grid from "@mui/material/Grid";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Select from "@mui/material/Select";
import CardHeader from "@mui/material/CardHeader";
import CardContent from "@mui/material/CardContent";
import InputAdornment from "@mui/material/InputAdornment";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import InputLabel from "@mui/material/InputLabel";

// ** Icons Imports
import DownloadBoxOutline from "mdi-material-ui/DownloadBoxOutline";
import OfficeBuilding from "mdi-material-ui/OfficeBuilding";

// ** React Imports
import { useState } from "react";

import { useMutation, gql } from "@apollo/client";
// import { FormEvent } from "react";

const ADD_FEED = gql`
  mutation AddFeed (
    $vendor: String!
    $feedName: String!
    $feedMethod: String!
  ) {
    addFeed(vendor: $vendor, feedName: $feedName, feedMethod: $feedMethod) {
      id
      vendor
      feedName
      feedMethod
    }
  }
`;



const FormNewFeed = () => {
  let message = "";
  let severity: AlertColor = "success";

  const [showAlert, setShowAlert] = useState(true);
  const [formState, setFormState] = useState({
    vendor: '',
    feedName: '',
    feedMethod: 'sftp'
  });

  const [addFeed, { data, error }] = useMutation(ADD_FEED, {
    onError: () => {},
    variables: {
      vendor: formState.vendor,
      feedName: formState.feedName,
      feedMethod: formState.feedMethod
    }
  });

  if (error) {
    message = `An Error Occured: ${error}`
    severity = "error";
  }

  if (data) {
    message = "Feed Added!";
    severity = "success";
  }

  let alertMessage: JSX.Element = <></>;
  if (showAlert) {
    alertMessage = <Alert severity={severity} onClose={() => setShowAlert(false)}>{message}</Alert>
  }

  return (
  <>
    {alertMessage}
    <Card>
      <CardHeader
        title="Add New Feed"
        titleTypographyProps={{ variant: "h6" }}
      />
      <CardContent>
        <form onSubmit={(e) => {
        e.preventDefault();
        addFeed();
        setShowAlert(true);
      }}>
          <Grid container spacing={5}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                id="vendor"
                label="Vendor"
                onChange={(e) =>
                  setFormState({
                    ...formState,
                    vendor: e.target.value
                  })
                }
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <OfficeBuilding />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                id="feedName"
                label="Feed Name"
                onChange={(e) =>
                  setFormState({
                    ...formState,
                    feedName: e.target.value
                  })
                }
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <DownloadBoxOutline />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Feed Type</InputLabel>
                <Select
                label="Feed Type"
                defaultValue="sftp"
                id="feedMethod"
                onChange={(e) =>
                  setFormState({
                    ...formState,
                    feedMethod: e.target.value
                  })
                }>
                  <MenuItem value="sftp">SFTP</MenuItem>
                  <MenuItem value="s3">S3</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <Button type="submit" variant="contained" size="large">
                Submit
              </Button>
            </Grid>
          </Grid>
        </form>
      </CardContent>
    </Card>
    </>
  );

};

export default FormNewFeed;
