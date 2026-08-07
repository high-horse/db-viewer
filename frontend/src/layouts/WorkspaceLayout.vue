<template>
    <div
        class="h-screen w-screen flex flex-col bg-[#0b0f19] text-[#94a3b8] overflow-hidden select-none"
    >
        <!-- Top Global App Navbar -->
        <header
            class="h-12 bg-[#111827] border-b border-[#1f2937] flex items-center justify-between px-4 z-10"
        >
            <div class="flex items-center gap-3">
                <q-icon name="dns" size="20px" class="text-indigo-400" />
                <span class="font-bold text-white tracking-wide text-sm"
                    >SQL Client Pro</span
                >
                <q-badge
                    color="indigo-10/10"
                    text-color="indigo-4"
                    label="Connected"
                    class="text-xs font-medium border border-indigo-9/30"
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
                    class="text-rose-400 hover:bg-rose-500/10"
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
                        class="h-full bg-[#0f172a]/40 border-r border-[#1f2937] flex flex-col font-sans"
                    >
                        <div
                            class="p-3 border-b border-[#1f2937] flex justify-between items-center bg-[#111827]/30"
                        >
                            <span
                                class="text-[11px] font-bold uppercase tracking-wider text-[#6b7280]"
                                >Schema Explorer</span
                            >
                            <q-btn
                                flat
                                dense
                                round
                                icon="refresh"
                                size="xs"
                                class="text-[#6b7280] hover:text-slate-300"
                            />
                        </div>

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
                                        class="flex items-center gap-2 py-0.5 cursor-pointer group w-full"
                                    >
                                        <q-icon
                                            :name="prop.node.icon"
                                            :color="prop.node.iconColor"
                                            size="14px"
                                        />
                                        <span
                                            :class="
                                                prop.node.type === 'table'
                                                    ? 'text-slate-200 font-medium'
                                                    : 'text-slate-400'
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

                <!-- Dynamic Main Workspace (Right Box) -->
                <template #after>
                    <q-splitter
                        v-model="editorHeight"
                        horizontal
                        :limits="[20, 80]"
                        class="h-full"
                    >
                        <!-- ZONE 2: Top Console Query Window -->
                        <template #before>
                            <div class="h-full flex flex-col bg-[#0b0f19]">
                                <!-- Tab strip bar -->
                                <div
                                    class="h-9 bg-[#111827] border-b border-[#1f2937] flex items-center justify-between px-2"
                                >
                                    <div class="flex items-center h-full">
                                        <div
                                            class="h-full px-4 border-r border-[#1f2937] bg-[#0b0f19] text-xs font-mono text-amber-400 flex items-center gap-2 border-b-2 border-b-amber-500"
                                        >
                                            <q-icon name="code" size="12px" />
                                            console_1.sql
                                        </div>
                                    </div>
                                    <q-btn
                                        unelevated
                                        color="emerald-5"
                                        size="sm"
                                        class="font-bold text-[#0b0f19] px-3 tracking-wide"
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

                                <!-- Mono Text Area Field Wrapper -->
                                <textarea
                                    v-model="sqlStatement"
                                    class="flex-grow w-full p-4 bg-[#070a12] text-amber-100/90 font-mono text-xs border-none resize-none focus:outline-none leading-relaxed tracking-wide"
                                    spellcheck="false"
                                ></textarea>
                            </div>
                        </template>

                        <!-- ZONE 3: Bottom Spreadsheet Rows Window Grid -->
                        <template #after>
                            <div class="h-full flex flex-col bg-[#070a12]">
                                <!-- Stats Metadata Strip bar -->
                                <div
                                    class="h-8 bg-[#111827] border-b border-[#1f2937] flex items-center justify-between px-3 text-[11px] font-mono text-[#6b7280]"
                                >
                                    <div class="flex items-center gap-3">
                                        <span>RESULT GRID</span>
                                        <span
                                            v-if="gridData"
                                            class="text-teal-400 font-bold"
                                            >• {{ gridData.Rows.length }} rows
                                            fetched</span
                                        >
                                    </div>
                                    <span v-if="gridData" class="text-slate-500"
                                        >{{ gridData.Duration }}ms execution
                                        latency</span
                                    >
                                </div>

                                <!-- Spread Data View Matrix -->
                                <div
                                    class="flex-grow relative overflow-hidden bg-[#070a12]"
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
                                    <div
                                        v-else
                                        class="h-full flex flex-col justify-center items-center gap-2 text-[#4b5563]"
                                    >
                                        <q-icon
                                            name="table_rows"
                                            size="32px"
                                            class="opacity-40"
                                        />
                                        <span class="text-xs tracking-wider"
                                            >Execute a statement query to
                                            inspect structural rows</span
                                        >
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
import { ref, computed } from 'vue'
import type { QTableColumn } from 'quasar' // Import the native Quasar column type

const sidebarWidth = ref(20)
const editorHeight = ref(45)
const sqlStatement = ref("SELECT * FROM actor WHERE first_name LIKE '%Joe%';")

const gridData = ref({
  Duration: 0,
  Columns: [] as Array<{ Name: string }>,
  Rows: [] as Array<Array<string | number>>
})

const schemaNodes = ref([
  {
    id: 'db-root',
    label: 'dvdrental',
    icon: 'storage',
    iconColor: 'emerald-4',
    children: [
      { id: 't1', label: 'actor', icon: 'table_chart', iconColor: 'indigo-4', type: 'table' },
      { id: 't2', label: 'customer', icon: 'table_chart', iconColor: 'indigo-4', type: 'table' },
      { id: 't3', label: 'film', icon: 'table_chart', iconColor: 'indigo-4', type: 'table' },
      { id: 't4', label: 'inventory', icon: 'table_chart', iconColor: 'indigo-4', type: 'table' },
    ]
  }
])

// FIX: Explicitly type this as an array of QTableColumn
const mappedColumns = computed<QTableColumn[]>(() => {
  if (!gridData.value || !gridData.value.Columns) return []
  return gridData.value.Columns.map(c => ({
    name: c.Name,
    label: c.Name,
    field: c.Name,
    align: 'left', // TypeScript now knows this matches "left" | "right" | "center"
    sortable: true
  }))
})

const mappedRows = computed(() => {
  if (!gridData.value || !gridData.value.Rows || !gridData.value.Columns) return []
  return gridData.value.Rows.map((r, i) => {
    const rowObj = { id: i } as Record<string, string | number>
    gridData.value.Columns.forEach((c, cIdx) => {
      rowObj[c.Name] = r[cIdx]
    })
    return rowObj
  })
})

function mockExecute() {
  gridData.value = {
    Duration: 8,
    Columns: [{ Name: 'actor_id' }, { Name: 'first_name' }, { Name: 'last_name' }],
    Rows: [
      [1, 'PENELOPE', 'GUINESS'],
      [2, 'NICK', 'WAHLBERG'],
      [3, 'ED', 'CHASE']
    ]
  }
}
</script>


<style scoped>
.data-explorer-grid :deep(.q-table__card) {
    background: transparent !important;
}
.data-explorer-grid :deep(thead tr th) {
    position: sticky;
    top: 0;
    background-color: #111827;
    z-index: 1;
    font-weight: bold;
    border-bottom: 2px solid #1f2937;
    color: #94a3b8;
    font-size: 11px;
}
.data-explorer-grid :deep(td) {
    border-bottom: 1px solid #1f2937;
    font-family: monospace;
    font-size: 11px;
    color: #cbd5e1;
}
</style>
