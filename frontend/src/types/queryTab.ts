export interface QueryColumn {
  Name: string;
  Type?: string;
}


export interface QueryResult {
  Duration: number;
  Columns: QueryColumn[];
  Rows: Array<Array<string | number>>;
}

export interface QueryTab {
  id: string;
  title: string;
  sql: string;

  result: QueryResult | null;

  loading: boolean;
  error: string | null;
  dirty: boolean;

  createdAt: number;

  connectionId?: string;
}