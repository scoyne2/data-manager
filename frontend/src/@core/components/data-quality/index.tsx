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
    const { data_quality_url } = props;
    if (data_quality_url === null || data_quality_url === undefined || data_quality_url === "") {
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
