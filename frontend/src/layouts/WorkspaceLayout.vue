<template>
    <div
        class="h-screen w-screen flex flex-col bg-[#0c0b09] text-[#94a3b8] overflow-hidden select-none"
    >
        <!-- Top Global App Navbar -->
        <header
            class="h-12 bg-[#161310] border-b border-[#292521] flex items-center justify-between px-4 z-10"
        >
            <div class="flex items-center gap-3">
                <q-icon name="dns" size="20px" class="text-amber-400" />

                <span class="font-bold text-white tracking-wide text-sm">
                    SQL Client Pro
                </span>

                <q-badge
                    color="amber-10"
                    text-color="amber-4"
                    label="Connected"
                    class="text-xs font-medium border border-amber-9"
                />
            </div>

            <div class="flex items-center gap-2">
                <q-btn
                    flat
                    dense
                    round
                    icon="tune"
                    size="sm"
                    class="text-[#6b7280] hover:text-white"
                />

                <q-btn
                    flat
                    dense
                    round
                    icon="logout"
                    size="sm"
                    class="text-amber-400 hover:bg-amber-500/10"
                    @click="$emit('disconnect')"
                />
            </div>
        </header>

        <!-- Main Workspace Splitter Container -->
        <div class="flex-grow flex relative">
            <q-splitter
                v-model="sidebarWidth"
                :limits="[15, 35]"
                class="absolute-full"
            >
                <!-- ZONE 1: Left Tree Explorer Sidebar -->
                <template #before>
                    <div
                        class="h-full bg-[#100e0c] border-r border-[#292521] flex flex-col font-sans"
                    >
                        <!-- Sidebar Header -->
                        <div
                            class="p-3 border-b border-[#292521] flex justify-between items-center bg-[#161310]"
                        >
                            <span
                                class="text-[11px] font-bold uppercase tracking-wider text-[#6b7280]"
                            >
                                Schema Explorer
                            </span>

                            <q-btn
                                flat
                                dense
                                round
                                icon="refresh"
                                size="xs"
                                class="text-[#6b7280] hover:text-amber-400"
                            />
                        </div>

                        <!-- Schema Tree -->
                        <q-scroll-area class="flex-grow p-2">
                            <q-tree
                                :nodes="schemaNodes"
                                node-key="id"
                                dark
                                no-connectors
                                class="text-xs"
                            >
                                <template #default-header="prop">
                                    <div
                                        class="flex items-center gap-2 py-1 px-1 rounded cursor-pointer group w-full hover:bg-[#231f1a]"
                                    >
                                        <q-icon
                                            :name="prop.node.icon"
                                            :color="prop.node.iconColor"
                                            size="14px"
                                        />

                                        <span
                                            :class="
                                                prop.node.type === 'table'
                                                    ? 'text-white font-medium'
                                                    : 'text-[#94a3b8]'
                                            "
                                        >
                                            {{ prop.node.label }}
                                        </span>
                                    </div>
                                </template>
                            </q-tree>
                        </q-scroll-area>
                    </div>
                </template>

                <!-- Dynamic Main Workspace -->
                <template #after>
                    <q-splitter
                        v-model="editorHeight"
                        horizontal
                        :limits="[20, 80]"
                        class="h-full"
                    >
                        <!-- ZONE 2: Top Console Query Window -->
                        <template #before>
                            <div class="h-full flex flex-col bg-[#0c0b09]">
                                <!-- Tab Strip -->
                                <div
                                    class="h-9 bg-[#161310] border-b border-[#292521] flex items-center justify-between px-2"
                                >
                                    <div class="flex items-center h-full">
                                        <div
                                            class="h-full px-4 border-r border-[#292521] bg-[#0c0b09] text-xs font-mono text-amber-400 flex items-center gap-2 border-b-2 border-b-amber-500"
                                        >
                                            <q-icon name="code" size="12px" />

                                            console_1.sql
                                        </div>
                                    </div>

                                    <q-btn
                                        unelevated
                                        color="amber"
                                        text-color="black"
                                        size="sm"
                                        class="font-bold px-3 tracking-wide"
                                        @click="mockExecute"
                                    >
                                        <q-icon
                                            name="play_arrow"
                                            size="16px"
                                            class="mr-1"
                                        />

                                        RUN
                                    </q-btn>
                                </div>

                                <!-- SQL Editor -->
                                <textarea
                                    v-model="sqlStatement"
                                    class="flex-grow w-full p-4 bg-[#0c0b09] text-amber-100/90 font-mono text-xs border-none resize-none focus:outline-none leading-relaxed tracking-wide"
                                    spellcheck="false"
                                ></textarea>
                            </div>
                        </template>

                        <!-- ZONE 3: Bottom Results Grid -->
                        <template #after>
                            <div class="h-full flex flex-col bg-[#100e0c]">
                                <!-- Result Metadata Strip -->
                                <div
                                    class="h-8 bg-[#161310] border-b border-[#292521] flex items-center justify-between px-3 text-[11px] font-mono text-[#6b7280]"
                                >
                                    <div class="flex items-center gap-3">
                                        <span> RESULT GRID </span>

                                        <span
                                            v-if="gridData"
                                            class="text-teal-400 font-bold"
                                        >
                                            •
                                            {{ gridData.Rows.length }}
                                            rows fetched
                                        </span>
                                    </div>

                                    <span
                                        v-if="gridData"
                                        class="text-[#6b7280]"
                                    >
                                        {{ gridData.Duration }}ms execution
                                        latency
                                    </span>
                                </div>

                                <!-- Data Grid -->
                                <div
                                    class="flex-grow relative overflow-hidden bg-[#100e0c]"
                                >
                                    <q-table
                                        v-if="gridData"
                                        flat
                                        square
                                        dense
                                        dark
                                        :rows="mappedRows"
                                        :columns="mappedColumns"
                                        row-key="id"
                                        :pagination="{ rowsPerPage: 0 }"
                                        :virtual-scroll-item-size="28"
                                        class="absolute-full data-explorer-grid"
                                    />

                                    <!-- Empty State -->
                                    <div
                                        v-else
                                        class="h-full flex flex-col justify-center items-center gap-2 text-[#4b4540]"
                                    >
                                        <q-icon
                                            name="table_rows"
                                            size="32px"
                                            class="opacity-40"
                                        />

                                        <span class="text-xs tracking-wider">
                                            Execute a statement query to inspect
                                            structural rows
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </template>
                    </q-splitter>
                </template>
            </q-splitter>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import type { QTableColumn } from "quasar";

const sidebarWidth = ref(20);
const editorHeight = ref(45);

const sqlStatement = ref("SELECT * FROM actor WHERE first_name LIKE '%Joe%';");

const gridData = ref({
    Duration: 0,
    Columns: [] as Array<{ Name: string }>,
    Rows: [] as Array<Array<string | number>>,
});

const schemaNodes = ref([
    {
        id: "db-root",
        label: "dvdrental",
        icon: "storage",
        iconColor: "amber-5",
        children: [
            {
                id: "t1",
                label: "actor",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
            {
                id: "t2",
                label: "customer",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
            {
                id: "t3",
                label: "film",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
            {
                id: "t4",
                label: "inventory",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
        ],
    },
]);

const mappedColumns = computed<QTableColumn[]>(() => {
    if (!gridData.value || !gridData.value.Columns) {
        return [];
    }

    return gridData.value.Columns.map((c) => ({
        name: c.Name,
        label: c.Name,
        field: c.Name,
        align: "left",
        sortable: true,
    }));
});

const mappedRows = computed(() => {
    if (!gridData.value || !gridData.value.Rows || !gridData.value.Columns) {
        return [];
    }

    return gridData.value.Rows.map((r, i) => {
        const rowObj = {
            id: i,
        } as Record<string, string | number>;

        gridData.value.Columns.forEach((c, cIdx) => {
            rowObj[c.Name] = r[cIdx];
        });

        return rowObj;
    });
});

function mockExecute() {
    gridData.value = {
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
    };
}
</script>

<style scoped>
.data-explorer-grid :deep(.q-table__card) {
    background: transparent !important;
    box-shadow: none !important;
}

.data-explorer-grid :deep(thead tr th) {
    position: sticky;
    top: 0;
    background-color: #161310;
    z-index: 1;
    font-weight: bold;
    border-bottom: 2px solid #292521;
    color: #f59e0b;
    font-size: 11px;
}

.data-explorer-grid :deep(tbody tr) {
    background-color: #100e0c;
}

.data-explorer-grid :deep(tbody tr:hover) {
    background-color: #231f1a !important;
}

.data-explorer-grid :deep(td) {
    border-bottom: 1px solid #292521;
    font-family: monospace;
    font-size: 11px;
    color: #d1d5db;
}

.data-explorer-grid :deep(tbody tr:nth-child(even)) {
    background-color: #13110f;
}
</style>
