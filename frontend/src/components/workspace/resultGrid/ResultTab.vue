<template>
    <div class="h-full relative overflow-hidden bg-[#100e0c]">
        <q-table
            flat
            square
            dense
            dark
            :rows="mappedRows"
            :columns="mappedColumns"
            row-key="id"
            :pagination="{
                rowsPerPage: 0,
            }"
            :virtual-scroll-item-size="28"
            class="absolute-full data-explorer-grid"
        >
            <!-- Header -->
            <template #header="props">
                <q-tr :props="props">
                    <q-th
                        v-for="col in props.cols"
                        :key="col.name"
                        :props="props"
                        class="column-header"
                    >
                        {{ col.label }}

                        <span
                            v-if="col.name !== 'sn'"
                            class="text-grey text-[8px] ml-1"
                        >
                            {{ col.Type }}
                        </span>

                        <q-tooltip
                            v-if="col.name !== 'sn'"
                            anchor="bottom middle"
                            self="top middle"
                            :offset="[0, 6]"
                            class="column-tooltip"
                        >
                            <div class="font-mono text-[10px]">
                                <div class="mt-1">
                                    {{ col.Type }}

                                    |

                                    {{
                                        col.Nullable
                                            ? "Nullable"
                                            : "Not Nullable"
                                    }}

                                    {{
                                        col.DefaultValue
                                            ? ` | Default: ${col.DefaultValue}`
                                            : ""
                                    }}
                                </div>
                            </div>
                        </q-tooltip>
                    </q-th>
                </q-tr>
            </template>

            <!-- Body -->
            <template #body="props">
                <q-tr :props="props">
                    <q-td
                        v-for="col in props.cols"
                        :key="col.name"
                        :props="props"
                        :class="{
                            'sn-cell': col.name === 'sn',
                        }"
                        align="left"
                    >
                        <template v-if="col.name === 'sn'">
                            {{ props.rowIndex + 1 }}
                        </template>

                        <template v-else>
                            {{ props.row[col.field] }}
                        </template>
                    </q-td>
                </q-tr>
            </template>
        </q-table>
    </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { QTableColumn } from "quasar";
import type { QueryResult } from "@/types/queryTab";

const props = defineProps<{
    result: QueryResult;
}>();

const mappedColumns = computed<QTableColumn[]>(() => {
    return [
        {
            name: "sn",
            label: "",
            field: "sn",
            align: "left",
            sortable: false,
        },

        ...props.result.Columns.map((column) => ({
            name: column.Name,
            label: column.Name,
            field: column.Name,
            align: "left" as const,
            sortable: true,

            Type: column.Type,
            Nullable: column.Nullable,
            DefaultValue: column.DefaultValue,
        })),
    ];
});

const mappedRows = computed(() => {
    return props.result.Rows.map((row, rowIndex) => {
        const rowObject = {
            id: rowIndex,
        } as Record<string, unknown>;

        props.result.Columns.forEach((column, columnIndex) => {
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

.data-explorer-grid :deep(.q-table__container) {
    background: transparent !important;
}

/* Header */

.data-explorer-grid :deep(thead tr th) {
    position: sticky;
    top: 0;
    z-index: 1;

    background-color: #161310;

    border-bottom: 2px solid #292521;

    color: #f59e0b;

    font-size: 11px;
    font-family: monospace;
    font-weight: bold;

    white-space: nowrap;
}

.data-explorer-grid :deep(thead tr th:first-child) {
    width: 45px;
    min-width: 45px;
    max-width: 45px;

    padding-left: 8px;
    padding-right: 4px;

    color: #4b5563 !important;

    text-align: left !important;
    font-weight: normal !important;
}

/* Body */

.data-explorer-grid :deep(tbody tr) {
    background-color: #100e0c;
}

.data-explorer-grid :deep(tbody tr:hover) {
    background-color: #231f1a !important;
}

.data-explorer-grid :deep(tbody tr:nth-child(even)) {
    background-color: #13110f;
}

.data-explorer-grid :deep(td) {
    border-bottom: 1px solid #292521;

    font-family: monospace;
    font-size: 11px;

    color: #d1d5db;

    white-space: nowrap;
}

/* SN */

.data-explorer-grid :deep(.sn-cell),
.data-explorer-grid :deep(.sn-cell.text-left),
.data-explorer-grid :deep(.sn-cell.text-right),
.data-explorer-grid :deep(.sn-cell.text-center) {
    width: 45px;
    min-width: 45px;
    max-width: 45px;

    padding-left: 8px;
    padding-right: 4px;

    color: #4b5563 !important;

    text-align: left !important;

    font-family: monospace;
    font-size: 11px;
    font-weight: normal !important;
}

/* Column spacing */

.data-explorer-grid :deep(th:not(:first-child)) {
    padding-left: 8px;
    padding-right: 8px;
}

.data-explorer-grid :deep(td:not(:first-child)) {
    padding-left: 8px;
    padding-right: 8px;
}

/* Scrollbar */

.data-explorer-grid :deep(.q-table__middle) {
    scrollbar-width: thin;
    scrollbar-color: #292521 #100e0c;
}

.data-explorer-grid :deep(.q-table__middle::-webkit-scrollbar) {
    width: 8px;
    height: 8px;
}

.data-explorer-grid :deep(.q-table__middle::-webkit-scrollbar-track) {
    background: #100e0c;
}

.data-explorer-grid :deep(.q-table__middle::-webkit-scrollbar-thumb) {
    background: #292521;
    border-radius: 4px;
}

.data-explorer-grid :deep(.q-table__middle::-webkit-scrollbar-thumb:hover) {
    background: #3a342e;
}

/* Tooltip */

.data-explorer-grid :deep(.column-tooltip) {
    background: #1c1916 !important;
    border: 1px solid #292521;

    color: #d1d5db;

    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.data-explorer-grid :deep(.q-table__bottom) {
    background: #161310;
    border-top: 1px solid #292521;
}
</style>