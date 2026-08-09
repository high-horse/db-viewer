<template>
    <div
        class="h-screen w-screen flex flex-col bg-[#0c0b09] text-[#94a3b8] overflow-hidden select-none"
    >
        <!-- Header -->
        <WorkspaceHeader
            @disconnect="handleDisconnect"
        />

        <!-- Workspace -->
        <div class="grow flex relative">
            <q-splitter
                v-model="sidebarWidth"
                :limits="[15, 35]"
                class="absolute-full"
            >
                <!-- LEFT SIDEBAR -->
                <template #before>
                    <SchemaSidebar />
                </template>

                <!-- MAIN -->
                <template #after>
                    <q-splitter
                        v-model="editorHeight"
                        horizontal
                        :limits="[20, 80]"
                        class="h-full"
                    >
                        <!-- QUERY CONSOLE -->
                        <template #before>
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
                                @execute="executeQuery"
                            />
                        </template>

                        <!-- RESULTS -->
                        <template #after>
                            <ResultGrid
                                :result="
                                    queryTabsStore.activeTab
                                        ?.result ?? null
                                "
                                :loading="
                                    queryTabsStore.activeTab
                                        ?.loading ?? false
                                "
                                :error="
                                    queryTabsStore.activeTab
                                        ?.error ?? null
                                "
                            />
                        </template>
                    </q-splitter>
                </template>
            </q-splitter>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

import WorkspaceHeader from "@/components/workspace/WorkspaceHeader.vue";
import SchemaSidebar from "@/components/workspace/SchemaSidebar.vue";
import QueryConsole from "@/components/workspace/QueryConsole.vue";
import ResultGrid from "@/components/workspace/ResultGrid.vue";

import { useQueryTabsStore } from "@/stores/queryTabsStore";

const sidebarWidth = ref(20);
const editorHeight = ref(45);

const queryTabsStore = useQueryTabsStore();

async function executeQuery(id: string) {
    const tab = queryTabsStore.tabs.find(
        (tab) => tab.id === id,
    );

    if (!tab || !tab.sql.trim()) {
        return;
    }

    queryTabsStore.setLoading(id, true);

    try {
        // Temporary dummy execution.
        // Replace this with DbService later.

        await new Promise((resolve) =>
            setTimeout(resolve, 500),
        );

        queryTabsStore.setResult(id, {
            Duration: 8,

            Columns: [
                { Name: "actor_id" },
                { Name: "first_name" },
                { Name: "last_name" },
            ],

            Rows: [
                [1, "PENELOPE", "GUINESS"],
                [2, "NICK", "WAHLBERG"],
                [3, "ED", "CHASE"],
            ],
        });
    } catch (error) {
        queryTabsStore.setError(
            id,
            error instanceof Error
                ? error.message
                : "Query execution failed",
        );
    }
}

function handleDisconnect() {
    console.log("Disconnect");
}
</script>