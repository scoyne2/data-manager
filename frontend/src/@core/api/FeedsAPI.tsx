// ** Types Imports
import { ThemeColor } from "src/@core/layouts/types";

export interface FeedStatusesDetailsType {
  feedstatusesdetailed: FeedStatusDetailedType[];
}

export interface FeedStatusType {
  id: number;
  record_count: number;
  feed_name: string;
  process_date: string;
  vendor: string;
  error_count: number;
  status: string;
  sla_status: string;
  schedule: string;
  feed_method: string;
  file_name: string;
  emr_logs: string;
  data_quality_url: string;
}

export interface FeedStatusDetailedType {
  id: number;
  feed_id: number;
  record_count: number;
  feed_name: string;
  process_date: string;
  vendor: string;
  error_count: number;
  sla_status: string;
  schedule: string;
  status: string;
  feed_method: string;
  file_name: string;
  previous_feeds: FeedStatusType[];
}

export interface StatusObj {
  [key: string]: {
    color: ThemeColor;
  };
}

export const statusObj: StatusObj = {
  Failed: { color: "error" },
  Errors: { color: "warning" },
  Success: { color: "success" },
  Received: { color: "info" },
  Processing: { color: "info" },
  Validating: { color: "info" },
};

const column_names: string[] = [
  "Details",
  "Feed",
  "Vendor",
  "Date",
  "Rows",
  "Errors",
  "Status",
  "SLA Status"
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