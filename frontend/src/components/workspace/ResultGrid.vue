<template>
    <div class="h-full flex flex-col bg-[#100e0c]">
        <!-- Result Metadata Strip -->
        <div
            class="h-8 bg-[#161310] border-b border-[#292521] flex items-center justify-between px-3 text-[11px] font-mono text-[#6b7280]"
        >
            <div class="flex items-center gap-3">
                <span> RESULT GRID </span>

                <span v-if="result" class="text-teal-400 font-bold">
                    • {{ result.Rows.length }} rows fetched
                </span>
            </div>

            <span v-if="result" class="text-[#6b7280]">
                {{ result.Duration }}ms execution latency
            </span>
        </div>

        <!-- Data Grid -->
        <div class="grow relative overflow-hidden bg-[#100e0c]">
            <q-table
                v-if="result"
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
                <q-icon name="table_rows" size="32px" class="opacity-40" />

                <span class="text-xs tracking-wider">
                    Execute a statement query to inspect structural rows
                </span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { QTableColumn } from "quasar";

interface QueryResult {
    Duration: number;
    Columns: Array<{
        Name: string;
    }>;
    Rows: Array<Array<string | number>>;
}

const props = defineProps<{
    result: QueryResult | null;
}>();

const mappedColumns = computed<QTableColumn[]>(() => {
    if (!props.result?.Columns) {
        return [];
    }

    return props.result.Columns.map((column) => ({
        name: column.Name,
        label: column.Name,
        field: column.Name,
        align: "left",
        sortable: true,
    }));
});

const mappedRows = computed(() => {
    if (!props.result) {
        return [];
    }

    return props.result.Rows.map((row, index) => {
        const rowObject = {
            id: index,
        } as Record<string, string | number>;

        props.result!.Columns.forEach((column, columnIndex) => {
            rowObject[column.Name] = row[columnIndex];
        });

        return rowObject;
    });
});
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
