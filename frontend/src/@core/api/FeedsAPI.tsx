// ** Types Imports
import { ThemeColor } from "src/@core/layouts/types";
import { ReactElement } from "react";

// ** Icons Imports
import TrendingUp from "mdi-material-ui/TrendingUp";
import Alert from "mdi-material-ui/Alert";
import AccountOutline from "mdi-material-ui/AccountOutline";

export interface FeedStatusesType {
  feedstatuses: FeedStatusType[];
}

export interface FeedStatusType {
  id: number;
  record_count: number;
  feed_name: string;
  process_date: string;
  vendor: string;
  error_count: number;
  status: string;
  feed_method: string;
}

export interface StatusObj {
  [key: string]: {
    color: ThemeColor;
  };
}

export const statusObj: StatusObj = {
  failed: { color: "error" },
  errors: { color: "warning" },
  success: { color: "success" },
};

const column_names: string[] = [
  "Feed",
  "Vendor",
  "Date",
  "Rows",
  "Errors",
  "Status",
];
export function getColumns(): string[] {
  return column_names;
}

export interface StatsType {
  stats: string;
  title: string;
  color: ThemeColor;
  icon: ReactElement;
}

const stats: StatsType[] = [
  {
    stats: "4.21M",
    title: "Rows",
    color: "primary",
    icon: <TrendingUp sx={{ fontSize: "1.75rem" }} />,
  },
  {
    stats: "215",
    title: "Files",
    color: "success",
    icon: <AccountOutline sx={{ fontSize: "1.75rem" }} />,
  },
  {
    stats: "72",
    color: "warning",
    title: "Errors",
    icon: <Alert sx={{ fontSize: "1.75rem" }} />,
  },
];

export function getStats(): StatsType[] {
  return stats;
}
