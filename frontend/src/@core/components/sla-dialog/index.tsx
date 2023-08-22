import * as React from 'react';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import { CheckCircle, CloseCircle, MinusCircle } from 'mdi-material-ui';

import { ThemeColor } from "src/@core/layouts/types";
import { SLADialogProps } from "src/@core/components/sla-dialog/types";

import { useMutation, gql } from "@apollo/client";

const UPDATE_SLA = gql`
  mutation UpdateSLA (
    $feedID: Int!
    $schedule: String!
  ) {
    updateSLA(feedID: $feedID, schedule: $schedule)
  }
`;


const getIconButtonForStatus = (status: string, handleClickOpen: any) => {
    let icon = <MinusCircle />
    let iconColor = ("info" as ThemeColor)
    
    switch (status) {
      case "Met":
        iconColor = ("success" as ThemeColor)
        icon = <CheckCircle />
        break;
      case "Missed":
        iconColor = ("error" as ThemeColor)
        icon = <CloseCircle />
        break;
    }

    return ( 
        <IconButton color={iconColor} onClick={handleClickOpen} >
         {icon}
        </IconButton>
    )
  };

const SLADialog = (props: SLADialogProps) => {
  const { feed_id, sla_status, schedule } = props;
  const [open, setOpen] = React.useState(false);

  const [scheduleSelection, setSecheduleSelection] = React.useState('');

  const handleChange = (event: SelectChangeEvent) => {
    setSecheduleSelection(event.target.value as string);
  };

  const handleClickOpen = () => {
    let currentSchedule = "none"
    if(schedule == "" || schedule == "none"){
      currentSchedule = schedule
    }
    setSecheduleSelection(currentSchedule)
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };


  const [updateSLA, { error }] = useMutation(UPDATE_SLA, {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    onError: () => {console.log(error)},
    variables: {
      feedID: feed_id,
      schedule: scheduleSelection,
    }
  });

  const handleSave = () => {
    // Make call to API to update sla status for feed
    console.log(feed_id)
    console.log(scheduleSelection)
    updateSLA({ variables: { feedID: feed_id, schedule: scheduleSelection } })
    setOpen(false);
  };

  return (
    <>
      {getIconButtonForStatus(sla_status, handleClickOpen)}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>SLA Settings</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Select the appropriate SLA for the feed.
          </DialogContentText>
          <Select
            labelId="sla-setting"
            id="sla-setting"
            value={scheduleSelection}
            label="sla-setting"
            onChange={handleChange}
          >
            <MenuItem value={"daily"}>Daily</MenuItem>
            <MenuItem value={"weekly"}>Weekly</MenuItem>
            <MenuItem value={"monthly"}>Monthly</MenuItem>
            <MenuItem value={"none"}>None</MenuItem>
          </Select>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleSave}>Save</Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
export default SLADialog;