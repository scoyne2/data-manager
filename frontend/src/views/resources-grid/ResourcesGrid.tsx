import * as React from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import Grid from "@mui/material/Grid";
import MonitorDashboard from "mdi-material-ui/MonitorDashboard";
import File from "mdi-material-ui/File";

function openInNewTab(url: string) {
  window.open(url, '_blank', 'noopener,noreferrer');
}

// TODO use dynamic resource bucket name
const S3_RESOURCE_BUCKET = process.env.NEXT_PUBLIC_S3_RESOURCE_BUCKET;
const AWS_REGION = process.env.NEXT_PUBLIC_AWS_REGION
const dd_kubernetes_url = "https://us5.datadoghq.com/kubernetes"
const dd_lambda_url ="https://us5.datadoghq.com/functions"
const ge_url = "http://" + S3_RESOURCE_BUCKET + ".s3-website-" + AWS_REGION + ".amazonaws.com/great_expectations/docs/"
const aws_emr_url = "https://s3.console.aws.amazon.com/s3/buckets/" + S3_RESOURCE_BUCKET + "?prefix=logs/"

const ResourcesGrid = () => {
    return (
        <>
          <Grid container spacing={4} >
          <Grid item xs={3}>
            <Card sx={{ minWidth: 50, minHeight: 250, textAlign: "center" }} variant="outlined" >
              <CardContent >
                <MonitorDashboard sx={{ fontSize: 150, color: "primary.light" }}  />
                <Button onClick={() => openInNewTab(dd_kubernetes_url)} size="small">DataDog: Kubernetes Dashboard</Button>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={3}>
            <Card sx={{ minWidth: 50, minHeight: 250, textAlign: "center"  }} variant="outlined">
              <CardContent>
                <MonitorDashboard sx={{ fontSize: 150, color: "secondary.dark" }}  />
                <Button  onClick={() => openInNewTab(dd_lambda_url)} size="small">DataDog: Serverless Dashboard</Button>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={3}>
            <Card sx={{ minWidth: 50, minHeight: 250, textAlign: "center"  }} variant="outlined">
              <CardContent>
                <File sx={{ fontSize: 150, color: "secondary.light" }}  />
                <Button  onClick={() => openInNewTab(ge_url)} size="small">Data Quality Results</Button>
                </CardContent>
            </Card>
          </Grid>
          <Grid item xs={3}>
            <Card sx={{ minWidth: 50, minHeight: 250, textAlign: "center"  }} variant="outlined">
              <CardContent>
                <File sx={{ fontSize: 150, color: "primary.dark" }}  />
                <Button  onClick={() => openInNewTab(aws_emr_url)} size="small">AWS: EMR Spark Logs</Button>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
        </>
        );
      
      };
      
export default ResourcesGrid;