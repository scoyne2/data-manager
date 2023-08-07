// ** MUI Imports
import * as React from 'react';
import Button from "@mui/material/Button";
import FileAlert from "mdi-material-ui/FileAlert";

// ** Types Imports
import { EMRLogsProps } from "src/@core/components/emr-logs/types";

function openInNewTab(url: string) {
    window.open(url, '_blank', 'noopener,noreferrer');
}

const EMRLogs = (props: EMRLogsProps) => {
    const { log_url } = props;
    if (log_url === null || log_url === undefined || log_url === "") {
        return (
            <>
            </>
        );
    }
    return (
        <>
        <Button onClick={() => openInNewTab(log_url)}><FileAlert/></Button>
        </>
    );
};

export default EMRLogs;
