// ** Types Imports
import { ThemeColor } from "src/@core/layouts/types";

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
  "Logs",
  "Quality Checks"
];

export function getColumns(): string[] {
  return column_names;
}

export interface StatusesAggregateType {
  feedstatuseaggregates:{ 
    files: number;
    rows: number;
    errors: number;
  }
}

export interface FeedQualityChecksType {
  feedqualitycheck: FeedQualityCheckType[];
}

export interface FeedQualityCheckType {
  id: number;
  feed_name: string;
  vendor: string;
  quality_check_results: QualityCheckResultsType[];
}

export interface QualityCheckResultsType {
  id: number;
  status: string;
  quality_check_name: string;
  quality_check_description: string;
}

export interface FeedLogsType {
  feedlogs: FeedLogType[];
}

export interface FeedLogType {
  id: number;
  status: string;
  feed_name: string;
  vendor: string;
  process_date: string;
  content: string;
}