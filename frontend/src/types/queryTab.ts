export interface QueryColumn {
  Name: string;
  Type?: string;
  Nullable?: boolean;
  DefaultValue?: string | number;
}


export interface QueryResult {
  Duration: number;
  Columns: QueryColumn[];
  Rows: Array<Array<string | number>>;
}

export type QueryTabType = "query" | "result"
export interface QueryTab {
  id: string;
  title: string;
  type: QueryTabType;
  
  sql: string;
  result: QueryResult | null;

  loading: boolean;
  error: string | null;
  dirty: boolean;

  createdAt: number;

  connectionId?: string;
}