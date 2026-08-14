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
                        <SchemaSidebar @selectTable="handleTableSelect" />
                    </div>
                </template>

                <!-- MAIN -->
                <template #after>
                    <div
                        class="h-full min-w-0 min-h-0 overflow-hidden flex flex-col"
                    >

                        <WorkspaceTabs
                            :tabs="queryTabsStore.tabs"
                            :active-tab-id="queryTabsStore.activeTabId"
                            @create-tab="queryTabsStore.createTab()"
                            @select-tab="queryTabsStore.selectTab"
                            @close-tab="queryTabsStore.closeTab"
                        />

                        <!-- ACTIVE TAB CONTENT -->

                        <div class="flex-grow min-h-0 min-w-0 overflow-hidden">

                            <template
                                v-if="
                                    queryTabsStore.activeTab?.type === 'query'
                                "
                            >
                                <q-splitter
                                    v-if="showResults"
                                    v-model="editorHeight"
                                    horizontal
                                    :limits="[20, 80]"
                                    class="h-full min-w-0 min-h-0"
                                >
                                    <!-- QUERY EDITOR -->
                                    <template #before>
                                        <div
                                            class="h-full min-w-0 min-h-0 overflow-hidden"
                                        >
                                            <QueryConsole
                                                :active-tab="
                                                    queryTabsStore.activeTab
                                                "
                                                :on-execute="executeQuery"
                                                @create-tab="
                                                    queryTabsStore.createTab()
                                                "
                                                @update-sql="
                                                    queryTabsStore.updateSql
                                                "
                                            />
                                        </div>
                                    </template>

                                    <!-- QUERY RESULT -->
                                    <template #after>
                                        <div
                                            class="h-full min-w-0 min-h-0 overflow-hidden"
                                        >
                                            <ResultGrid
                                                :result="
                                                    queryTabsStore.activeTab
                                                        .result ?? null
                                                "
                                                :loading="
                                                    queryTabsStore.activeTab
                                                        .loading
                                                "
                                                :error="
                                                    queryTabsStore.activeTab
                                                        .error
                                                "
                                                @close="hideResults"
                                            />
                                        </div>
                                    </template>
                                </q-splitter>

                                <!-- QUERY WITHOUT RESULTS -->
                                <div
                                    v-else
                                    class="h-full w-full min-w-0 min-h-0 overflow-hidden"
                                >
                                    <QueryConsole
                                        :active-tab="queryTabsStore.activeTab"
                                        :on-execute="executeQuery"
                                        @create-tab="queryTabsStore.createTab()"
                                        @update-sql="queryTabsStore.updateSql"
                                    />
                                </div>
                            </template>

                            <!-- ================= RESULT TAB ================= -->

                            <template
                                v-else-if="
                                    queryTabsStore.activeTab?.type === 'result'
                                "
                            >
                                <div
                                    class="h-full w-full min-w-0 min-h-0 overflow-hidden"
                                >
                                    <ResultGrid
                                        :result="
                                            queryTabsStore.activeTab.result ??
                                            null
                                        "
                                        :loading="
                                            queryTabsStore.activeTab.loading
                                        "
                                        :error="queryTabsStore.activeTab.error"
                                        @close="
                                            queryTabsStore.closeTab(
                                                queryTabsStore.activeTab.id,
                                            )
                                        "
                                    />
                                </div>
                            </template>

                            <!-- ================= NO TAB ================= -->

                            <template v-else>
                                <div
                                    class="h-full flex flex-col items-center justify-center gap-3"
                                >
                                    <q-icon
                                        name="tab"
                                        size="36px"
                                        class="text-[#4b4540]"
                                    />

                                    <span class="text-xs text-[#6b7280]">
                                        No tab selected
                                    </span>

                                    <q-btn
                                        unelevated
                                        color="amber"
                                        text-color="black"
                                        size="sm"
                                        label="New Query"
                                        icon="add"
                                        @click="queryTabsStore.createTab()"
                                    />
                                </div>
                            </template>
                        </div>
                    </div>
                </template>
            </q-splitter>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeMount, onBeforeUnmount } from "vue";

import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";

import { DbService } from "@bindings/db-viewer/internal/app";

import type { QueryResult } from "@/types/queryTab";

import WorkspaceHeader from "@/components/workspace/WorkspaceHeader.vue";
import SchemaSidebar from "@/components/workspace/SchemaSidebar.vue";
import QueryConsole from "@/components/workspace/QueryConsole.vue";
// import ResultGrid from "@/components/workspace/ResultGrid.vue";
import ResultGrid from "@/components/workspace/resultGrid/ResultGrid.vue";

import WorkspaceTabs from "@/components/workspace/WorkspaceTabs.vue";

import { useQueryTabsStore } from "@/stores/queryTabsStore";
import { useConnectionStore } from "@/stores/connectionStore";

const $router = useRouter();

const queryTabsStore = useQueryTabsStore();
const connectionStore = useConnectionStore();

const { activeConnection, activeConnectionMetadata } =
    storeToRefs(connectionStore);

const sidebarWidth = ref(20);
const editorHeight = ref(45);

const showResults = ref(true);

function hideResults() {
    showResults.value = false;
}

async function executeQuery(id: string, sql: string) {
    const tab = queryTabsStore.tabs.find((tab) => tab.id === id);

    if (!tab || !sql.trim()) {
        return;
    }

    showResults.value = true;

    queryTabsStore.setLoading(id, true);

    try {
        const response = await DbService.ExecuteQuery(sql);

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

            Rows: (response.rows ?? []).map((row) =>
                (row ?? []).map((value) => value ?? ""),
            ),
        };

        queryTabsStore.setResult(id, result);
    } catch (error) {
        queryTabsStore.setError(
            id,
            error instanceof Error ? error.message : "Query execution failed",
        );
    } finally {
        queryTabsStore.setLoading(id, false);
    }
}

async function handleTableSelect(node: {
    id: string;
    label: string;
    type?: string;
}) {
    if (node.type !== "table" && node.type !== "view") {
        return;
    }

    const tab = queryTabsStore.createResultTab(node.label);

    await executeTableResult(tab.id, node);
}

async function executeTableResult(
    id: string,
    node: {
        id: string;
        label: string;
        type?: string;
    },
) {
    queryTabsStore.setLoading(id, true);

    try {
        const tableName = node.id
            .split(".")
            .map((part) => `"${part.replace(/"/g, '""')}"`)
            .join(".");

        const sql = `SELECT * FROM ${tableName};`;

        const response = await DbService.ExecuteQuery(sql);

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

            Rows: (response.rows ?? []).map((row) =>
                (row ?? []).map((value) => value ?? ""),
            ),
        };

        queryTabsStore.setResult(id, result);
    } catch (error) {
        queryTabsStore.setError(
            id,
            error instanceof Error ? error.message : "Failed to load table",
        );
    } finally {
        queryTabsStore.setLoading(id, false);
    }
}

async function handleDisconnect() {
    const session = connectionStore.getActiveSession();

    if (session) {
        await DbService.Disconnect(String(session.id));

        connectionStore.clearActiveSession();

        $router.push({
            name: "Welcome",
        });
    }
}

async function setActiveSession() {
    try {
        const [session, exist] = await DbService.GetActiveConnectionObject();

        if (exist && session) {
            connectionStore.setActiveSession(session);

            if (!activeConnectionMetadata.value) {
                connectionStore.setActiveConnectionMetadata();
            }
        }
    } catch (error) {
        console.error("failed to set active session", error);
    }
}

function handleCreateTabShortcut(event: KeyboardEvent) {
    if (event.ctrlKey && event.key.toLowerCase() === "t") {
        event.preventDefault();

        queryTabsStore.createTab();
    }
}

function handleCloseTabShortcut(event: KeyboardEvent) {
    if (event.ctrlKey && event.key.toLowerCase() === "w") {
        event.preventDefault();

        if (queryTabsStore.activeTabId) {
            queryTabsStore.closeTab(queryTabsStore.activeTabId);
        }
    }
}

function handleToggleResultTabShortcut(event: KeyboardEvent) {
    if (event.ctrlKey && event.key.toLowerCase() === "j") {
        event.preventDefault();

        if (queryTabsStore.activeTab?.type === "query") {
            showResults.value = !showResults.value;
        }
    }
}

onBeforeMount(() => {
    window.addEventListener("keydown", handleCreateTabShortcut);

    window.addEventListener("keydown", handleCloseTabShortcut);

    window.addEventListener("keydown", handleToggleResultTabShortcut);
});

onBeforeUnmount(() => {
    window.removeEventListener("keydown", handleCreateTabShortcut);

    window.removeEventListener("keydown", handleCloseTabShortcut);

    window.removeEventListener("keydown", handleToggleResultTabShortcut);
});

onMounted(() => {
    if (!activeConnection.value) {
        setActiveSession();
    }
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
