// ** React Imports
import { SyntheticEvent, useState } from "react";

// ** MUI Imports
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import TabList from "@mui/lab/TabList";
import TabPanel from "@mui/lab/TabPanel";
import TabContext from "@mui/lab/TabContext";
import { styled } from "@mui/material/styles";
import MuiTab, { TabProps } from "@mui/material/Tab";

// ** Icons Imports
import MagnifyPlus from "mdi-material-ui/MagnifyPlus";
import CogOutline from "mdi-material-ui/CogOutline";

// ** Demo Tabs Imports
import TabLogs from "src/views/feed-details/TabLogs";
import TabQualityChecks from "src/views/feed-details/TabQualityChecks";


// ** Third Party Styles Imports
import "react-datepicker/dist/react-datepicker.css";

const Tab = styled(MuiTab)<TabProps>(({ theme }) => ({
  [theme.breakpoints.down("md")]: {
    minWidth: 100,
  },
  [theme.breakpoints.down("sm")]: {
    minWidth: 67,
  },
}));

const TabName = styled("span")(({ theme }) => ({
  lineHeight: 1.71,
  fontSize: "0.875rem",
  marginLeft: theme.spacing(2.4),
  [theme.breakpoints.down("md")]: {
    display: "none",
  },
}));

const AccountSettings = () => {
  // ** State
  const [value, setValue] = useState<string>("quality-checks");

  const handleChange = (event: SyntheticEvent, newValue: string) => {
    setValue(newValue);
  };

  return (
    <Card>
      <TabContext value={value}>
        <TabList
          onChange={handleChange}
          aria-label="account-settings tabs"
          sx={{ borderBottom: (theme) => `1px solid ${theme.palette.divider}` }}
        >
          <Tab
            value="quality-checks"
            label={
              <Box sx={{ display: "flex", alignItems: "center" }}>
                <MagnifyPlus />
                <TabName>Quality Checks</TabName>
              </Box>
            }
          />
          <Tab
            value="logs"
            label={
              <Box sx={{ display: "flex", alignItems: "center" }}>
                <CogOutline />
                <TabName>Logs</TabName>
              </Box>
            }
          />
        </TabList>

        <TabPanel sx={{ p: 0 }} value="quality-checks">
          <TabQualityChecks />
        </TabPanel>
        <TabPanel sx={{ p: 0 }} value="logs">
          <TabLogs />
        </TabPanel>
      </TabContext>
    </Card>
  );
};

export default AccountSettings;
