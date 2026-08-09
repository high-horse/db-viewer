<template>
    <div
        class="h-screen w-screen flex flex-col bg-[#0c0b09] text-[#94a3b8] overflow-hidden select-none"
    >
        <!-- Global Header -->
        <WorkspaceHeader @disconnect="handleDisconnect" />

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

                <!-- MAIN AREA -->
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
                                v-model="sqlStatement"
                                @execute="executeQuery"
                            />
                        </template>

                        <!-- RESULT GRID -->
                        <template #after>
                            <ResultGrid :result="gridData" />
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

interface QueryResult {
    Duration: number;
    Columns: Array<{
        Name: string;
    }>;
    Rows: Array<Array<string | number>>;
}

const sidebarWidth = ref(20);
const editorHeight = ref(45);

const sqlStatement = ref("SELECT * FROM actor WHERE first_name LIKE '%Joe%';");

const gridData = ref<QueryResult | null>(null);

function executeQuery(query: string) {
    console.log("Executing query:", query);

    // Temporary mock result.
    // Replace this with your DbService query later.

    gridData.value = {
        Duration: 8,

        Columns: [
            {
                Name: "actor_id",
            },
            {
                Name: "first_name",
            },
            {
                Name: "last_name",
            },
        ],

        Rows: [
            [1, "PENELOPE", "GUINESS"],
            [2, "NICK", "WAHLBERG"],
            [3, "ED", "CHASE"],
        ],
    };
}

function handleDisconnect() {
    console.log("Disconnect");

    // Later:
    // 1. Disconnect database
    // 2. Clear active connection
    // 3. Navigate back to welcome page
}
</script>
