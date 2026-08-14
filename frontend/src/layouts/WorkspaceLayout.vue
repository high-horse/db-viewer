<template>
    <div
        class="h-screen w-screen flex flex-col bg-[#0c0b09] text-[#94a3b8] overflow-hidden select-none"
    >
        <WorkspaceHeader @disconnect="handleDisconnect" />

        <div class="grow flex relative min-w-0 min-h-0 overflow-hidden">
            <q-splitter
                v-model="sidebarWidth"
                :limits="[15, 35]"
                class="absolute-full min-w-0 min-h-0"
            >
                <!-- LEFT SIDEBAR -->
                <template #before>
                    <div class="h-full min-w-0 overflow-hidden">
                        <SchemaSidebar />
                    </div>
                </template>

                <!-- MAIN -->
                <template #after>
                    <div class="h-full min-w-0 min-h-0 overflow-hidden">
                        <q-splitter
                            v-if="showResults"
                            v-model="editorHeight"
                            horizontal
                            :limits="[20, 80]"
                            class="h-full min-w-0 min-h-0"
                        >
                            <!-- QUERY CONSOLE -->
                            <template #before>
                                <div class="h-full min-w-0 min-h-0 overflow-hidden">
                                    <QueryConsole
                                        :tabs="queryTabsStore.tabs"
                                        :active-tab-id="
                                            queryTabsStore.activeTabId
                                        "
                                        :active-tab="
                                            queryTabsStore.activeTab
                                        "
                                        @create-tab="
                                            queryTabsStore.createTab()
                                        "
                                        @select-tab="
                                            queryTabsStore.selectTab
                                        "
                                        @close-tab="
                                            queryTabsStore.closeTab
                                        "
                                        @update-sql="
                                            queryTabsStore.updateSql
                                        "
                                        :on-execute="executeQuery"
                                        @toggle-result-tab="handleToggleResultTab"
                                    />
                                </div>
                            </template>

                            <!-- RESULTS -->
                            <template #after>
                                <div class="h-full min-w-0 min-h-0 overflow-hidden">
                                    <ResultGrid
                                        :result="
                                            queryTabsStore.activeTab?.result ??
                                            null
                                        "
                                        :loading="
                                            queryTabsStore.activeTab?.loading ??
                                            false
                                        "
                                        :error="
                                            queryTabsStore.activeTab?.error ??
                                            null
                                        "
                                        @close="hideResults"
                                    />
                                </div>
                            </template>
                        </q-splitter>

                        <!-- Query only -->
                        <div
                            v-else
                            class="h-full w-full min-w-0 min-h-0 overflow-hidden"
                        >
                            <QueryConsole
                                :tabs="queryTabsStore.tabs"
                                :active-tab-id="
                                    queryTabsStore.activeTabId
                                "
                                :active-tab="
                                    queryTabsStore.activeTab
                                "
                                @create-tab="
                                    queryTabsStore.createTab()
                                "
                                @select-tab="
                                    queryTabsStore.selectTab
                                "
                                @close-tab="
                                    queryTabsStore.closeTab
                                "
                                @update-sql="
                                    queryTabsStore.updateSql
                                "
                                :on-execute="executeQuery"
                                @toggle-result-tab="handleToggleResultTab"
                            />
                        </div>
                    </div>
                </template>
            </q-splitter>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import { DbService } from "@bindings/db-viewer/internal/app";
import type { QueryResult } from "@/types/queryTab";
import WorkspaceHeader from "@/components/workspace/WorkspaceHeader.vue";
import SchemaSidebar from "@/components/workspace/SchemaSidebar.vue";
import QueryConsole from "@/components/workspace/QueryConsole.vue";
import ResultGrid from "@/components/workspace/ResultGrid.vue";


const $router = useRouter();
import { useQueryTabsStore } from "@/stores/queryTabsStore";
import { useConnectionStore } from "@/stores/connectionStore";

const sidebarWidth = ref(20);
const editorHeight = ref(45);

const showResults = ref(true);
const queryTabsStore = useQueryTabsStore();
const connectionStore = useConnectionStore();
const { activeConnection, activeConnectionMetadata } = storeToRefs(connectionStore);

function hideResults() {
    showResults.value = false;
}

function handleToggleResultTab() {
    showResults.value = !showResults.value;
}

async function executeQuery(id: string, sql: string) {
    console.log("executing sql query", sql);

    const tab = queryTabsStore.tabs.find((tab) => tab.id === id);

    if (!tab || !tab.sql.trim()) {
        return;
    }
    showResults.value = true;

    queryTabsStore.setLoading(id, true);

    try {
        const response = await DbService.ExecuteQuery(sql);

        console.log("query response", response);

        if (!response) {
            throw new Error("Query returned no response");
        }

        const result: QueryResult = {
            Duration: response.duration ?? 0,

            Columns: (response.columns ?? []).map((column) => ({
              Name: column.name,
              Type: column.databaseType,
              Nullable: column.nullable,
              DefaultValue: column.defaultValue,
            })),

            Rows: (response.rows ?? []).map(
                (row) =>
                    (row ?? []).map(
                        (value) => value ?? "",
                    ),
            ),
        };

        queryTabsStore.setResult(id, result);
    } catch (error) {
        queryTabsStore.setError(
            id,
            error instanceof Error
                ? error.message
                : "Query execution failed",
        );
    } finally {
        queryTabsStore.setLoading(id, false);
    }
}

async function handleDisconnect() {
  const session = connectionStore.getActiveSession();
  console.log("disconnecting", session, connectionStore.getActiveSession());
  if (session) {
    await DbService.Disconnect(String(session.id));
    connectionStore.clearActiveSession();
    $router.push({name: 'Welcome'});
  }
}


async function setActiveSession() {
  try {
    const [session, exist] = await DbService.GetActiveConnectionObject();
    console.log("active session", session, exist);
    if (exist && session) {
      connectionStore.setActiveSession(session);
      
      if (!activeConnectionMetadata.value) {
        connectionStore.setActiveConnectionMetadata();
      }

    } else {
      console.log("no active session");
    }
  } catch (error) {
    console.error("failed to set active session", error);
  }
}

onMounted(() => {
  if (!activeConnection.value) {
    setActiveSession();
  };
  
});

</script>


<style scoped>
:deep(.q-splitter) {
    min-width: 0;
    min-height: 0;
}

:deep(.q-splitter__panel) {
    min-width: 0;
    min-height: 0;
    overflow: hidden;
}
</style>