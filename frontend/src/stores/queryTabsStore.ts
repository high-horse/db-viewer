import { defineStore } from "pinia";
import { computed, ref } from "vue";
import type { QueryTab, QueryResult } from "@/types/queryTab";
import { DbService } from "@bindings/db-viewer/internal/app";
import { QueryExecutionType } from "@bindings/db-viewer/internal/engine/entities";
import type { QueryResult as BackendResult } from "@bindings/db-viewer/internal/engine/entities";

const initialTab: QueryTab = {
  id: "query-1",
  title: "console_1.sql",
  type: "query",
  
  sql: "",

  result: null,

  loading: false,

  error: null,

  dirty: false,

  createdAt: Date.now(),

  connectionId: undefined,
};

export const useQueryTabsStore = defineStore("queryTabs", () => {
  const tabs = ref<QueryTab[]>([initialTab]);

  const activeTabId = ref<string | null>("query-1");

  const activeTab = computed(() => {
    return tabs.value.find((tab) => tab.id === activeTabId.value) ?? null;
  });

  function createTab(connectionId?: string) {
    const id = crypto.randomUUID();

    const tab: QueryTab = {
      id,
      title: `console_${tabs.value.length + 1}.sql`,
      type: "query",
      sql: "",
      result: null,
      loading: false,
      error: null,
      dirty: false,
      createdAt: Date.now(),
      connectionId,
    };

    tabs.value.push(tab);

    activeTabId.value = id;

    return tab;
  }

  function createResultTab(title: string, connectionId?: string) {
    const id = crypto.randomUUID();

    const tab: QueryTab = {
      id,
      title,
      type: "result",
      sql: "",
      result: null,
      loading: false,
      error: null,
      dirty: false,
      createdAt: Date.now(),
      connectionId,
    };

    tabs.value.push(tab);

    activeTabId.value = id;

    return tab;
  }

  function selectTab(id: string) {
    const exists = tabs.value.some((tab) => tab.id === id);

    if (!exists) {
      return;
    }

    activeTabId.value = id;
  }

  function closeTab(id: string) {
    const index = tabs.value.findIndex((tab) => tab.id === id);

    if (index === -1) {
      return;
    }

    const wasActive = activeTabId.value === id;

    releaseCursor(tabs.value[index].cursor);

    tabs.value.splice(index, 1);

    if (!wasActive) {
      return;
    }

    if (tabs.value.length === 0) {
      activeTabId.value = null;
      return;
    }

    const nextIndex = Math.min(index, tabs.value.length - 1);

    activeTabId.value = tabs.value[nextIndex].id;
  }

  function updateSql(id: string, sql: string) {
    const tab = tabs.value.find((tab) => tab.id === id);

    if (!tab) {
      return;
    }

    tab.sql = sql;
    tab.dirty = true;
  }

  function setResult(id: string, result: QueryResult) {
    const tab = tabs.value.find((tab) => tab.id === id);

    if (!tab) {
      return;
    }

    tab.result = result;
    tab.loading = false;
    tab.error = null;
    tab.dirty = false;
  }

  function setLoading(id: string, loading: boolean) {
    const tab = tabs.value.find((tab) => tab.id === id);

    if (tab) {
      tab.loading = loading;
    }
  }

  function setError(id: string, error: string) {
    const tab = tabs.value.find((tab) => tab.id === id);

    if (!tab) {
      return;
    }

    tab.loading = false;
    tab.error = error;
  }

  function releaseCursor(cursor?: string) {
    if (!cursor) return;
    void DbService.ExecuteQuery({ query: "", cursor, type: QueryExecutionType.QueryExecutionClose, pageSize: 0, page: 0 }).catch(() => {});
  }

  function mapResult(response: BackendResult): QueryResult {
    return {
      Duration: response.duration,
      Columns: (response.columns ?? []).map(column => ({ Name: column.name, Type: column.databaseType, Nullable: column.nullable, DefaultValue: column.defaultValue })),
      Rows: (response.rows ?? []).map(row => row ?? []),
      Cursor: response.cursor, HasMore: response.hasMore,
      StartRow: response.startRow, PageSize: response.pageSize, IsQuery: response.isQuery, CanNavigate: response.canNavigate,
    };
  }

  async function execute(id: string, sql: string) {
    const tab = tabs.value.find(tab => tab.id === id);
    if (!tab || tab.loading || !sql.trim()) return;
    const previousCursor = tab.cursor ?? "";
    tab.loading = true;
    tab.error = null;
    tab.pageError = undefined;
    tab.totalRows = undefined;
    tab.fetchingLast = false;
    tab.sortColumn = 0;
    tab.sortDirection = undefined;
    tab.streamNextPage = 2;
    tab.result = null;
    tab.pages = [];
    tab.pageIndex = 0;
    tab.cursor = "";
    tab.executedSql = sql;
    try {
      const response = await DbService.ExecuteQuery({ query: sql, cursor: previousCursor, type: QueryExecutionType.QueryExecutionExecute, pageSize: 100, page: 1 });
      if (!response) throw new Error("Query returned no response");
      if (!tabs.value.includes(tab)) { releaseCursor(response.cursor); return; }
      const result = mapResult(response);
      tab.pages = [result];
      tab.cursor = result.Cursor;
      if (response.totalRows != null) tab.totalRows = response.totalRows;
      else if (!result.HasMore) tab.totalRows = result.StartRow - 1 + result.Rows.length;
      setResult(id, result);
    } catch (error) {
      setError(id, error instanceof Error ? error.message : String(error));
    } finally { tab.loading = false; }
  }

  async function navigate(id: string, direction: "previous" | "next" | "last") {
    const tab = tabs.value.find(tab => tab.id === id);
    if (!tab || tab.loading || !tab.pages?.length) return;
    if (tab.result?.CanNavigate) {
      const page = Math.floor((tab.result.StartRow - 1) / tab.result.PageSize) + 1;
      const lastPage = Math.max(1, Math.ceil((tab.totalRows ?? 0) / tab.result.PageSize));
      const target = direction === "last" ? lastPage : page + (direction === "next" ? 1 : -1);
      if (target < 1 || target > lastPage || target === page) return;
      await fetchReadPage(tab, target);
      return;
    }
    const index = tab.pageIndex ?? 0;
    const target = direction === "last"
      ? tab.pages.length - (tab.cursor ? 0 : 1)
      : index + (direction === "next" ? 1 : -1);
    if (target < 0) return;
    if (target < tab.pages.length) {
      tab.pageIndex = target;
      tab.result = tab.pages[target];
      tab.pageError = undefined;
      return;
    }
    if (!tab.cursor) return;
    tab.loading = true;
    tab.fetchingLast = direction === "last";
    tab.pageError = undefined;
    try {
      do {
        const last = tab.pages[tab.pages.length - 1];
        const response = await DbService.ExecuteQuery({ query: "", cursor: tab.cursor, type: QueryExecutionType.QueryExecutionFetchPaged, pageSize: last.PageSize, page: Math.floor((last.StartRow - 1) / last.PageSize) + 2 });
        if (!response) throw new Error("No result page returned");
        if (!tabs.value.includes(tab)) { releaseCursor(response.cursor); return; }
        const result = mapResult(response);
        tab.cursor = result.Cursor;
        if (!result.HasMore) tab.totalRows = result.StartRow - 1 + result.Rows.length;
        tab.pages.push(result);
        // Keep navigation memory bounded regardless of the database size.
        if (tab.pages.length > 10) tab.pages.shift();
        tab.pageIndex = tab.pages.length - 1;
        tab.result = result;
      } while (tab.fetchingLast && tab.cursor && tabs.value.includes(tab));
    } catch (error) {
      tab.pageError = error instanceof Error ? error.message : String(error);
    } finally { tab.loading = false; tab.fetchingLast = false; }
  }

  async function fetchReadPage(tab: QueryTab, page: number, sort?: { column: number; direction?: "asc" | "desc" }) {
    if (!tab.result || tab.loading) return;
    const size = tab.result.PageSize;
    if (!sort) {
      const cached = tab.pages?.findIndex(result => result.StartRow === (page - 1) * size + 1) ?? -1;
      if (cached >= 0 && tab.pages) {
        tab.pageIndex = cached;
        tab.result = tab.pages[cached];
        tab.pageError = undefined;
        return;
      }
    }
    const useCursor = !sort && !!tab.cursor && page === tab.streamNextPage;
    tab.loading = true;
    tab.pageError = undefined;
    try {
      const response = await DbService.ExecuteQuery({
        query: tab.executedSql ?? "", cursor: tab.cursor ?? "",
        type: useCursor ? QueryExecutionType.QueryExecutionFetchPaged : QueryExecutionType.QueryExecutionNavigate,
        pageSize: size, page,
        sortColumn: sort?.column ?? tab.sortColumn ?? 0,
        sortDirection: sort ? sort.direction ?? "" : tab.sortDirection ?? "",
      });
      if (!response) throw new Error("No result page returned");
      if (!tabs.value.includes(tab)) { releaseCursor(response.cursor); return; }
      const result = mapResult(response);
      tab.cursor = result.Cursor;
      tab.streamNextPage = Math.floor((result.StartRow - 1) / size) + 2;
      if (response.totalRows != null) tab.totalRows = response.totalRows;
      if (sort) {
        tab.sortColumn = sort.column;
        tab.sortDirection = sort.direction;
        tab.pages = [];
      }
      tab.pages ??= [];
      tab.pages = tab.pages.filter(cached => cached.StartRow !== result.StartRow);
      tab.pages.push(result);
      if (tab.pages.length > 10) tab.pages.shift();
      tab.pageIndex = tab.pages.length - 1;
      tab.result = result;
    } catch (error) {
      // A random jump closes the old stream before opening the new one.
      if (!useCursor) tab.cursor = "";
      tab.pageError = error instanceof Error ? error.message : String(error);
    } finally { tab.loading = false; }
  }

  async function sortResult(id: string, column: number, direction?: "asc" | "desc") {
    const tab = tabs.value.find(tab => tab.id === id);
    if (!tab?.result?.CanNavigate || tab.loading) return;
    await fetchReadPage(tab, 1, { column, direction });
  }

  function stopFetching(id: string) {
    const tab = tabs.value.find(tab => tab.id === id);
    // Complete the in-flight batch, preserving a resumable cursor.
    if (tab) tab.fetchingLast = false;
  }

  function clearResults() {
    for (const tab of tabs.value) releaseCursor(tab.cursor);
    tabs.value = [];
    activeTabId.value = null;
  }

  return {
    tabs,
    activeTabId,
    activeTab,

    createTab,
    createResultTab,
    selectTab,
    closeTab,
    updateSql,
    setResult,
    setLoading,
    setError,
    execute,
    navigate,
    stopFetching,
    sortResult,
    clearResults,
  };
});
