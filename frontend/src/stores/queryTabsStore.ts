import { defineStore } from "pinia";
import { computed, ref } from "vue";
import type { QueryTab, QueryResult } from "@/types/queryTab";

export const useQueryTabsStore = defineStore("queryTabs", () => {
  const tabs = ref<QueryTab[]>([]);

  const activeTabId = ref<string | null>(null);

  const activeTab = computed(() => {
    return tabs.value.find((tab) => tab.id === activeTabId.value) ?? null;
  });

  function createTab(connectionId?: string) {
    const id = crypto.randomUUID();

    const tab: QueryTab = {
      id,
      title: `console_${tabs.value.length + 1}.sql`,
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
    activeTabId.value = id;
  }

  function closeTab(id: string) {
    const index = tabs.value.findIndex((tab) => tab.id === id);

    if (index === -1) {
      return;
    }

    const wasActive = activeTabId.value === id;

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

  return {
    tabs,
    activeTabId,
    activeTab,

    createTab,
    selectTab,
    closeTab,
    updateSql,
    setResult,
    setLoading,
    setError,
  };
});
