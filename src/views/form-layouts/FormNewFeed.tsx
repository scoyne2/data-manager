// ** MUI Imports
import Card from '@mui/material/Card'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'
import Select from '@mui/material/Select'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import InputAdornment from '@mui/material/InputAdornment'
import MenuItem from '@mui/material/MenuItem'


// ** Icons Imports
import DownloadBoxOutline from 'mdi-material-ui/DownloadBoxOutline'
import OfficeBuilding from 'mdi-material-ui/OfficeBuilding'

const FormNewFeed = () => {
  return (
    <Card>
      <CardHeader title='Add New Feed' titleTypographyProps={{ variant: 'h6' }} />
      <CardContent>
        <form onSubmit={e => e.preventDefault()}>
          <Grid container spacing={5}>
          <Grid item xs={12}>
              <TextField
                fullWidth
                label='Vendor'
                InputProps={{
                  startAdornment: (
                    <InputAdornment position='start'>
                      <OfficeBuilding />
                    </InputAdornment>
                  )
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label='Feed Name'
                InputProps={{
                  startAdornment: (
                    <InputAdornment position='start'>
                      <DownloadBoxOutline />
                    </InputAdornment>
                  )
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <Select
                fullWidth
                label='Feed Type'
              >
                <MenuItem value={1}>SFTP</MenuItem>
                <MenuItem value={2}>S3</MenuItem>
                </Select>
            </Grid>


            <Grid item xs={12}>
              <Button type='submit' variant='contained' size='large'>
                Submit
              </Button>
            </Grid>
          </Grid>
        </form>
      </CardContent>
    </Card>
  )
}

export default FormNewFeed
