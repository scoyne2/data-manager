// ** Types Imports
import { ThemeColor } from 'src/@core/layouts/types'
import { ReactElement } from 'react'

// ** Icons Imports
import TrendingUp from 'mdi-material-ui/TrendingUp'
import Alert from 'mdi-material-ui/Alert'
import AccountOutline from 'mdi-material-ui/AccountOutline'

export interface RowType {
    rows: number
    name: string
    date: string
    vendor: string
    errors: number
    status: string
    designation: string
}

export interface StatusObj {
  [key: string]: {
    color: ThemeColor
  }
}

export const statusObj: StatusObj = {
  failed: { color: 'error' },
  errors: { color: 'warning' },
  success: { color: 'success' }
}


const rows: RowType[] = [
    {
      errors: 27,
      status: 'success',
      date: '09/27/2018',
      name: 'Daily Dog Treats',
      rows: 1958623,
      vendor: 'Mochi Dog Fyle Systems',
      designation: 'SFTP'
    },
    {
      errors: 61,
      date: '09/23/2016',
      rows: 2389635,
      status: 'success',
      name: 'Dog Bones',
      vendor: 'Mochi Dog Fyle Systems',
      designation: 'S3'
    },
    {
      errors: 59,
      date: '10/15/2017',
      name: 'Physicians',
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

export function getFeedStatus(): RowType[]  {
    return rows
};


const columns: string[] = [
  'Feed',
  'Vendor',
  'Date',
  'Rows',
  'Errors',
  'Status'

]
export function getColumns(): string[] {
  return columns;
}

export interface StatsType {
  stats: string
  title: string
  color: ThemeColor
  icon: ReactElement
}

const stats: StatsType[] = [
  {
    stats: '4.21M',
    title: 'Rows',
    color: 'primary',
    icon: <TrendingUp sx={{ fontSize: '1.75rem' }} />
  },
  {
    stats: '215',
    title: 'Files',
    color: 'success',
    icon: <AccountOutline sx={{ fontSize: '1.75rem' }} />
  },
  {
    stats: '72',
    color: 'warning',
    title: 'Errors',
    icon: <Alert sx={{ fontSize: '1.75rem' }} />
  }
]

export function getStats(): StatsType[] {
  return stats;
}