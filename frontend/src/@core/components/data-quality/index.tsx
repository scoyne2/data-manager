// ** MUI Imports
import * as React from 'react';
import Button from "@mui/material/Button";
import MagnifyExpand from "mdi-material-ui/MagnifyExpand";

// ** Types Imports
import { DataQualityResultsProps } from "src/@core/components/data-quality/types";

function openInNewTab(url: string) {
    window.open(url, '_blank', 'noopener,noreferrer');
}

const DataQualityResults = (props: DataQualityResultsProps) => {
    const { data_quality_url, status } = props;
    // These statuses indicate data quality check has not yet been completed
    if (status.toLowerCase() === "received" || status.toLowerCase() === "processing" || status.toLowerCase() === "validating")  {
        return (
            <>
            </>
        );
    }
    return (
        <>
        <Button onClick={() => openInNewTab(data_quality_url)}><MagnifyExpand/></Button>
        </>
    );
};

export default DataQualityResults;
