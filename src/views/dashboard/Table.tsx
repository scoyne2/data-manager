// ** MUI Imports
import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import Chip from '@mui/material/Chip'
import Table from '@mui/material/Table'
import TableRow from '@mui/material/TableRow'
import TableHead from '@mui/material/TableHead'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import Typography from '@mui/material/Typography'
import TableContainer from '@mui/material/TableContainer'

// ** Types Imports
import { ThemeColor } from 'src/@core/layouts/types'

interface RowType {
  age: number
  name: string
  date: string
  email: string
  salary: string
  status: string
  designation: string
}

interface StatusObj {
  [key: string]: {
    color: ThemeColor
  }
}

const rows: RowType[] = [
  {
    errors: 27,
    status: 'success',
    date: '09/27/2018',
    name: 'Daily Dog Treats',
    rows: '1,958,623',
    vendor: 'Mochi Dog Fyle Systems',
    designation: 'SFTP'
  },
  {
    errors: 61,
    date: '09/23/2016',
    rows: '2,389,635',
    status: 'success',
    name: 'Dog Bones',
    vendor: 'Mochi Dog Fyle Systems',
    designation: 'S3'
  },
  {
    errors: 59,
    date: '10/15/2017',
    name: 'Physicianss',
    status: 'failed',
    rows: 0,
    vendor: 'HCA',
    designation: 'SFTP'
  },
  {
    errors: 1225212,
    date: '06/12/2018',
    status: 'failed',
    rows: 0,
    name: 'Specialty Mapping',
    vendor: 'HCA',
    designation: 'SFTP'
  },
  {
    errors: 42,
    status: 'errors',
    date: '03/24/2018',
    rows: 13076,
    name: 'Dental Claims',
    designation: 'S3',
    vendor: 'Coyne Enterprises LLC'
  },
  {
    errors: 0,
    date: '08/25/2017',
    rows: 1090952,
    name: 'Procedures',
    status: 'success',
    vendor: 'Coyne Enterprises LLC',
    designation: 'SFTP'
  },
  {
    errors: 2,
    status: 'errors',
    date: '06/01/2017',
    rows: 17803,
    name: 'Specialty Mapping',
    designation: 'SFTP',
    vendor: 'Mochi Dog Fyle Systems',
  },
  {
    errors: 0,
    date: '12/03/2017',
    rows: 1233617,
    name: 'Physicians',
    status: 'success',
    designation: 'S3',
    vendor: 'ABC'
  }
]

const statusObj: StatusObj = {
  failed: { color: 'error' },
  errors: { color: 'warning' },
  success: { color: 'success' }
}

const DashboardTable = () => {
  return (
    <Card>
      <TableContainer>
        <Table sx={{ minWidth: 800 }} aria-label='table in dashboard'>
          <TableHead>
            <TableRow>
              <TableCell>Feed</TableCell>
              <TableCell>Vendor</TableCell>
              <TableCell>Date</TableCell>
              <TableCell>Rows</TableCell>
              <TableCell>Errors</TableCell>
              <TableCell>Status</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row: RowType) => (
              <TableRow hover key={row.name} sx={{ '&:last-of-type td, &:last-of-type th': { border: 0 } }}>
                <TableCell sx={{ py: theme => `${theme.spacing(0.5)} !important` }}>
                  <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                    <Typography sx={{ fontWeight: 500, fontSize: '0.875rem !important' }}>{row.name}</Typography>
                    <Typography variant='caption'>{row.designation}</Typography>
                  </Box>
                </TableCell>
                <TableCell>{row.vendor}</TableCell>
                <TableCell>{row.date}</TableCell>
                <TableCell>{row.rows}</TableCell>
                <TableCell>{row.errors}</TableCell>
                <TableCell>
                  <Chip
                    label={row.status}
                    color={statusObj[row.status].color}
                    sx={{
                      height: 24,
                      fontSize: '0.75rem',
                      textTransform: 'capitalize',
                      '& .MuiChip-label': { fontWeight: 500 }
                    }}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Card>
  )
}

export default DashboardTable
