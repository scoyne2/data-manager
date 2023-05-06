// ** MUI Imports
import * as React from 'react';
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Modal from '@mui/material/Modal';
import FileFind from "mdi-material-ui/FileFind";
import Table from "@mui/material/Table";
import TableRow from "@mui/material/TableRow";
import TableHead from "@mui/material/TableHead";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";


// ** Types Imports
import { DataPreviewModalProps } from "src/@core/components/data-preview/types";

const DOMAIN_NAME = process.env.NEXT_PUBLIC_DOMAIN_NAME;

const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 800,
    bgcolor: 'background.paper',
    boxShadow: 24,
    p: 4,
  };
  
function createHeader(data: string[]) {
    const header: JSX.Element[] = [];
    if (data.length > 0) {
        data.forEach(element => header.push(<TableCell>{element}</TableCell>));
        return header;
    }
}


function createBody(data: string[][]) {
    const body: JSX.Element[] = [];
    if (data.length > 0) {
        data.forEach(rows => {
            const currentRow: JSX.Element[] = [];
            rows.forEach(element => {
                currentRow.push(<TableCell>{element}</TableCell>);
            });
            body.push(<TableRow>{currentRow}</TableRow>);
        });
        return body;
    }
}

const DataPreviewModal = (props: DataPreviewModalProps) => {
    const { vendor, feed_name, file_name } = props;

    const [open, setOpen] = React.useState(false);
    const [tableData, setTableData] = React.useState({"Header": [], "Rows": []});
    const handleOpen = () => setOpen(true);
    const handleClose = () => setOpen(false);

    // Example: https://api.datamanagertool.com/preview?vendor=Coyne+Enterprises&feedname=Hello&filename=hello.csv
    const urlVendor = vendor.replace(/\s/g, '+')
    const urlFeedName = feed_name.replace(/\s/g, '+')
    const dataPreviewURL = `https://api.${DOMAIN_NAME}/preview?vendor=${urlVendor}&feedname=${urlFeedName}&filename=${file_name}`;
    React.useEffect(() => {
        fetch(dataPreviewURL)
            .then(res => res.json())
            .then(
                (data) => {
                    setTableData(data);
                }
            )
           .catch(error => console.error(error));
      });
      

    
    return (
        <>
        <Button onClick={handleOpen}><FileFind/></Button>
        <Modal
            open={open}
            onClose={handleClose}
            aria-labelledby="data preview modal"
            aria-describedby="a modal for previewing data"
        >
             <Box sx={style}>
                <TableContainer>
                    <Table sx={{ minWidth: 800 }} aria-label="data preview table">
                    <TableHead>
                        <TableRow>
                        { createHeader(tableData["Header"])}
                        </TableRow>
                    </TableHead>
                    <TableBody>
                    { createBody(tableData["Rows"]) }
                    </TableBody>
                    </Table>
                </TableContainer>
            </Box>
        </Modal>
        </>
    );
};

export default DataPreviewModal;
