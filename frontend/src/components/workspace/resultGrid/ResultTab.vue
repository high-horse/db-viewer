<template>
    <div class="relative h-full w-full overflow-hidden bg-[#100e0c]">
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
            class="absolute inset-0 data-explorer-grid bg-transparent"
            table-class="table-fixed"
            table-style="min-width: 100%; width: max-content;"
        >
            <!--HEADER-->

            <template #header="props">
                <q-tr :props="props" class="bg-[#161310]">
                    <q-th
                        v-for="col in props.cols"
                        :key="col.name"
                        :props="props"
                        class="relative box-border h-[28px] overflow-hidden whitespace-nowrap border-b-2 border-[#292521] px-2 py-0 align-middle font-mono text-[11px] font-bold text-amber-400"
                        :class="{
                            'sticky-col-header pl-2 pr-1 text-left font-normal text-[#4b5563]':
                                col.name === 'sn',
                        }"
                        :style="getColumnStyle(col.name)"
                    >
                        <div
                            class="flex w-full min-w-0 items-center overflow-hidden whitespace-nowrap"
                        >
                            <div
                                class="flex min-w-0 flex-1 items-center overflow-hidden pr-1 whitespace-nowrap"
                            >
                                <!-- Column name -->

                                <span
                                    class="min-w-0 shrink overflow-hidden text-ellipsis whitespace-nowrap"
                                >
                                    {{ col.label }}
                                </span>

                                <!-- Column type -->
                                <span
                                    v-if="col.name !== 'sn'"
                                    class="ml-1 min-w-0 shrink overflow-hidden text-ellipsis whitespace-nowrap text-[8px] font-normal text-gray-500"
                                >
                                    {{ col.Type }}
                                </span>
                            </div>

                            <!-- Sort indicator -->
                            <q-icon
                                v-if="activeSortColumn === col.name"
                                :name="
                                    sortDirection === 'asc'
                                        ? 'arrow_upward'
                                        : 'arrow_downward'
                                "
                                size="12px"
                                class="ml-0.5 shrink-0 text-amber-400"
                            />
                        </div>

                        <span
                            v-if="col.name !== 'sn'"
                            class="column-resizer"
                            @mousedown.stop="startResize($event, col.name)"
                        />

                        <q-tooltip
                            v-if="col.name !== 'sn'"
                            anchor="bottom middle"
                            self="top middle"
                            :offset="[0, 6]"
                            class="bg-[#1c1916] text-gray-300 shadow-lg"
                        >
                            <div class="font-mono text-[10px]">
                                {{ col.Type }}

                                |

                                {{ col.Nullable ? "Nullable" : "Not Nullable" }}

                                {{
                                    col.DefaultValue
                                        ? ` | Default: ${col.DefaultValue}`
                                        : ""
                                }}
                            </div>
                        </q-tooltip>

                        <q-menu
                            v-if="col.name !== 'sn'"
                            context-menu
                            class="min-w-[150px] border border-[#292521] bg-[#1c1916] text-gray-300 shadow-xl"
                        >
                            <!-- Sort Ascending -->
                            <q-item
                                clickable
                                v-close-popup
                                dense
                                class="min-h-[28px] px-2 hover:bg-[#292521]"
                                @click="sortColumn(col.name, 'asc')"
                            >
                                <q-item-section avatar class="min-w-[24px]">
                                    <q-icon
                                        name="arrow_upward"
                                        size="14px"
                                        class="text-gray-400"
                                    />
                                </q-item-section>

                                <q-item-section>
                                    <q-item-label class="font-mono text-[11px]">
                                        Sort Ascending
                                    </q-item-label>
                                </q-item-section>
                            </q-item>

                            <!-- Sort Descending -->
                            <q-item
                                clickable
                                v-close-popup
                                dense
                                class="min-h-[28px] px-2 hover:bg-[#292521]"
                                @click="sortColumn(col.name, 'desc')"
                            >
                                <q-item-section avatar class="min-w-[24px]">
                                    <q-icon
                                        name="arrow_downward"
                                        size="14px"
                                        class="text-gray-400"
                                    />
                                </q-item-section>

                                <q-item-section>
                                    <q-item-label class="font-mono text-[11px]">
                                        Sort Descending
                                    </q-item-label>
                                </q-item-section>
                            </q-item>

                            <q-separator class="my-1 bg-[#292521]" />

                            <!-- Clear Sort -->
                            <q-item
                                v-if="activeSortColumn === col.name"
                                clickable
                                v-close-popup
                                dense
                                class="min-h-[28px] px-2 hover:bg-[#292521]"
                                @click="clearSort()"
                            >
                                <q-item-section avatar class="min-w-[24px]">
                                    <q-icon
                                        name="close"
                                        size="14px"
                                        class="text-gray-500"
                                    />
                                </q-item-section>

                                <q-item-section>
                                    <q-item-label class="font-mono text-[11px]">
                                        Clear Sort
                                    </q-item-label>
                                </q-item-section>
                            </q-item>
                        </q-menu>
                    </q-th>
                </q-tr>
            </template>


            <template #body="props">
                <q-tr
                    :props="props"
                    class="bg-[#100e0c] even:bg-[#13110f] hover:!bg-[#231f1a]"
                >
                    <q-td
                        v-for="col in props.cols"
                        :key="col.name"
                        :props="props"
                        class="box-border overflow-hidden whitespace-nowrap border-b border-[#292521] px-2 py-1.5 font-mono text-[11px] text-gray-300"
                        :class="{
                            'sticky-col-body pl-2 pr-1 text-left text-[#4b5563]':
                                col.name === 'sn',
                        }"
                        :style="getColumnStyle(col.name)"
                    >
                        <template v-if="col.name === 'sn'">
                            <span class="text-grey-8">
                                {{ props.rowIndex + 1 }}
                            </span>
                        </template>

                        <template v-else>
                            <div
                                class="block w-full min-w-0 overflow-hidden text-ellipsis whitespace-nowrap"
                                :title="String(props.row[col.field] ?? '')"
                            >
                                {{ props.row[col.field] }}
                            </div>
                        </template>
                    </q-td>
                </q-tr>
            </template>
        </q-table>
    </div>
</template>

<script setup lang="ts">
import {
    computed,
    nextTick,
    onBeforeUnmount,
    onMounted,
    reactive,
    ref,
    watch,
} from "vue";

import type { QTableColumn } from "quasar";

import type { QueryResult } from "@/types/queryTab";


const MIN_COLUMN_WIDTH = 60;
const SN_COLUMN_WIDTH = 45;
const HEADER_EXTRA_WIDTH = 24;

const props = defineProps<{
    result: QueryResult;
}>();

const columnWidths = reactive<Record<string, number>>({});

const getColumnWidth = (columnName: string): number => {
    if (columnName === "sn") {
        return SN_COLUMN_WIDTH;
    }

    return columnWidths[columnName] ?? MIN_COLUMN_WIDTH;
};

const getColumnStyle = (columnName: string) => {
    const width = getColumnWidth(columnName);

    return {
        width: `${width}px`,
        minWidth: `${width}px`,
        maxWidth: `${width}px`,
    };
};

const measureHeaderWidths = async () => {
    await nextTick();

    await new Promise<void>((resolve) => {
        requestAnimationFrame(() => resolve());
    });

    const table = document.querySelector(".data-explorer-grid table");
    if (!table) {
        return;
    }

    const headers = table.querySelectorAll("thead th");

    headers.forEach((header, index) => {
        if (index === 0) { //Skip # column.
            return;
        }

        const column = props.result.Columns[index - 1];

        if (!column) {
            return;
        }

        if (columnWidths[column.Name] !== undefined) {
            return;
        }

        const th = header as HTMLElement;
        const content = th.querySelector("div") as HTMLElement | null;
        if (!content) {
            return;
        }

        const oldWidth = th.style.width;
        const oldMinWidth = th.style.minWidth;
        const oldMaxWidth = th.style.maxWidth;
        const tableElement = table as HTMLElement;
        const oldTableLayout = tableElement.style.tableLayout;

        th.style.width = "max-content";
        th.style.minWidth = "max-content";
        th.style.maxWidth = "none";
        tableElement.style.tableLayout = "auto";

        const contentWidth = Math.max(th.scrollWidth, content.scrollWidth); // Measure header ONLY.

        // Restore immediately.
        th.style.width = oldWidth;
        th.style.minWidth = oldMinWidth;
        th.style.maxWidth = oldMaxWidth;
        tableElement.style.tableLayout = oldTableLayout;

        // Apply minimum width.
        const width = Math.max(
            MIN_COLUMN_WIDTH,
            contentWidth + HEADER_EXTRA_WIDTH,
        );

        columnWidths[column.Name] = width;
    });
};


onMounted(() => {
    measureHeaderWidths();
});


const columnSignature = computed(() =>
    props.result.Columns.map((column) => `${column.Name}:${column.Type}`).join(
        "|",
    ),
);


watch(columnSignature, async () => {
     // Remove widths for columns no longer exist.
    const validColumns = new Set(
        props.result.Columns.map((column) => column.Name),
    );

    Object.keys(columnWidths).forEach((columnName) => {
        if (!validColumns.has(columnName)) {
            delete columnWidths[columnName];
        }
    });

    // Wait for new columns to render.
    await nextTick();

    await measureHeaderWidths();
});

let resizingColumn: string | null = null;
let resizeStartX = 0;
let resizeStartWidth = 0;

const startResize = (event: MouseEvent, columnName: string) => {
    event.preventDefault();
    event.stopPropagation();
    resizingColumn = columnName;
    resizeStartX = event.clientX;

    // Get the actual rendered header width.
    const target = event.currentTarget as HTMLElement;

    const header = target.closest("th") as HTMLElement | null;

    resizeStartWidth = header?.getBoundingClientRect().width ?? getColumnWidth(columnName);

    // Immediately lock the current width
    columnWidths[columnName] = resizeStartWidth;
    document.body.classList.add("resizing-column");
    document.addEventListener("mousemove", handleResize);
    document.addEventListener("mouseup", stopResize);
};

const handleResize = (event: MouseEvent) => {
    if (!resizingColumn) {
        return;
    }

    const delta = event.clientX - resizeStartX;
    const newWidth = Math.max(MIN_COLUMN_WIDTH, resizeStartWidth + delta);
    columnWidths[resizingColumn] = newWidth;
};

const stopResize = () => {
    resizingColumn = null;
    document.body.classList.remove("resizing-column");
    document.removeEventListener("mousemove", handleResize);
    document.removeEventListener("mouseup", stopResize);
};

onBeforeUnmount(() => {
    stopResize();
});

type SortDirection = "asc" | "desc" | null;
const activeSortColumn = ref<string | null>(null);
const sortDirection = ref<SortDirection>(null);

const sortColumn = (columnName: string, direction: "asc" | "desc") => {
    activeSortColumn.value = columnName;
    sortDirection.value = direction;
};

const clearSort = () => {
    activeSortColumn.value = null;
    sortDirection.value = null;
};

const mappedColumns = computed<QTableColumn[]>(() => {
    return [
        {
            name: "sn",
            label: "#",
            field: "sn",
            align: "left",
            sortable: false,

            style: `
                width: ${SN_COLUMN_WIDTH}px;
                min-width: ${SN_COLUMN_WIDTH}px;
                max-width: ${SN_COLUMN_WIDTH}px;
            `,

            headerStyle: `
                width: ${SN_COLUMN_WIDTH}px;
                min-width: ${SN_COLUMN_WIDTH}px;
                max-width: ${SN_COLUMN_WIDTH}px;
            `,
        },

        ...props.result.Columns.map((column) => {
            const width = getColumnWidth(column.Name);

            return {
                name: column.Name,
                label: column.Name,
                field: column.Name,
                align: "left" as const,
                sortable: false,

                Type: column.Type,
                Nullable: column.Nullable,
                DefaultValue: column.DefaultValue,

                style: `
                    width: ${width}px;
                    min-width: ${width}px;
                    max-width: ${width}px;
                `,

                headerStyle: `
                    width: ${width}px;
                    min-width: ${width}px;
                    max-width: ${width}px;
                `,
            };
        }),
    ];
});
const mappedRows = computed(() => {
    const rows = props.result.Rows.map((row, rowIndex) => {
        const rowObject = {
            id: rowIndex,
        } as Record<string, unknown>;

        props.result.Columns.forEach((column, columnIndex) => {
            rowObject[column.Name] = row[columnIndex];
        });

        return rowObject;
    });

    if (!activeSortColumn.value || !sortDirection.value) {
        return rows;
    }
    const column = activeSortColumn.value;
    const direction = sortDirection.value === "asc" ? 1 : -1;

    return [...rows].sort((a, b) => {
        const valueA = a[column];
        const valueB = b[column];

        if (valueA === null || valueA === undefined) {
            if (valueB === null || valueB === undefined) {
                return 0;
            }

            return -1 * direction;
        }

        if (valueB === null || valueB === undefined) {
            return 1 * direction;
        }


        if (typeof valueA === "number" && typeof valueB === "number") {
            return (valueA - valueB) * direction;
        }


        const stringA = String(valueA);
        const stringB = String(valueB);

        const numberA = Number(stringA);
        const numberB = Number(stringB);

        if (
            stringA.trim() !== "" &&
            stringB.trim() !== "" &&
            Number.isFinite(numberA) &&
            Number.isFinite(numberB)
        ) {
            return (numberA - numberB) * direction;
        }

        return (
            stringA.localeCompare(stringB, undefined, {
                numeric: true,
                sensitivity: "base",
            }) * direction
        );
    });
});
</script>

<style scoped>
/* =========================================================
   RESIZE HANDLE
   ========================================================= */

.column-resizer {
    position: absolute;

    top: 0;
    right: 0;
    bottom: 0;

    width: 8px;

    cursor: col-resize;

    z-index: 70;
}

.column-resizer::after {
    content: "";

    position: absolute;

    top: 3px;
    bottom: 3px;
    left: 3px;

    width: 1px;

    background: transparent;

    transition:
        background-color 80ms ease,
        width 80ms ease;
}

.column-resizer:hover::after {
    width: 2px;

    background: #f59e0b;
}

/* =========================================================
   RESIZING STATE
   ========================================================= */

:global(body.resizing-column) {
    cursor: col-resize !important;

    user-select: none !important;
}

:global(body.resizing-column *) {
    cursor: col-resize !important;

    user-select: none !important;
}

/* =========================================================
   QTABLE CONTAINER
   ========================================================= */

.data-explorer-grid :deep(.q-table__card),
.data-explorer-grid :deep(.q-table__container) {
    background: transparent !important;

    box-shadow: none !important;
}

/*
 * THIS is the scrolling container.
 *
 * Both horizontal and vertical scrolling happen here.
 */

.data-explorer-grid :deep(.q-table__middle) {
    background: #100e0c;

    overflow: auto;

    scrollbar-width: thin;

    scrollbar-color: #292521 #100e0c;
}

/* =========================================================
   TABLE LAYOUT
   ========================================================= */

/*
 * Fixed layout is critical.
 *
 * Body content cannot change the column width.
 */

.data-explorer-grid :deep(table) {
    table-layout: fixed;

    width: max-content;

    min-width: 100%;
}

/* =========================================================
   ALL CELLS
   ========================================================= */

.data-explorer-grid :deep(th),
.data-explorer-grid :deep(td) {
    box-sizing: border-box;

    min-width: 0;

    overflow: hidden;

    white-space: nowrap;

    text-overflow: ellipsis;

    user-select: text;
}

/* =========================================================
   STICKY HEADER
   ========================================================= */

/*
 * Keep the header visible during VERTICAL scrolling.
 *
 * The header has a lower z-index than the top-left
 * sticky # cell.
 */

.data-explorer-grid :deep(thead) {
    position: sticky;

    top: 0;

    z-index: 40;
}

/*
 * Every normal header cell stays at the top.
 */

.data-explorer-grid :deep(thead th) {
    position: sticky;

    top: 0;

    z-index: 40;

    background: #161310;

    overflow: hidden;

    white-space: nowrap;

    text-overflow: ellipsis;
}

/* =========================================================
   STICKY # HEADER
   ========================================================= */

/*
 * The top-left cell needs BOTH:
 *
 * left: 0
 * top: 0
 *
 * Therefore it remains fixed during BOTH horizontal
 * and vertical scrolling.
 */

.data-explorer-grid :deep(thead th:first-child) {
    position: sticky;

    left: 0;

    top: 0;

    z-index: 60;

    width: 45px;

    min-width: 45px;

    max-width: 45px;

    background: #161310;

    border-right: 1px solid #332e28;
}

/* =========================================================
   STICKY # BODY COLUMN
   ========================================================= */

/*
 * The row-number column stays fixed horizontally.
 */

.data-explorer-grid :deep(tbody td:first-child) {
    position: sticky;

    left: 0;

    z-index: 20;

    width: 45px;

    min-width: 45px;

    max-width: 45px;

    overflow: hidden;

    background: #100e0c;

    border-right: 1px solid #292521;
}

/* =========================================================
   STICKY # COLUMN - ALTERNATING ROWS
   ========================================================= */

.data-explorer-grid :deep(tbody tr:nth-child(even) td:first-child) {
    background: #13110f;
}

/* =========================================================
   STICKY # COLUMN - HOVER
   ========================================================= */

.data-explorer-grid :deep(tbody tr:hover td:first-child) {
    background: #231f1a;
}

/* =========================================================
   BODY CELLS
   ========================================================= */

.data-explorer-grid :deep(tbody td) {
    overflow: hidden;

    white-space: nowrap;

    text-overflow: ellipsis;
}

/*
 * Inner cell content also cannot expand the column.
 */

.data-explorer-grid :deep(tbody td > div) {
    display: block;

    width: 100%;

    min-width: 0;

    max-width: 100%;

    overflow: hidden;

    white-space: nowrap;

    text-overflow: ellipsis;
}

/* =========================================================
   COLUMN BORDERS
   ========================================================= */

.data-explorer-grid :deep(th + th),
.data-explorer-grid :deep(td + td) {
    border-left: 1px solid #292521;
}

.data-explorer-grid :deep(thead th + th) {
    border-left-color: #332e28;
}

/* =========================================================
   STICKY COLUMN SHADOW
   ========================================================= */

/*
 * Gives the fixed # column a subtle separation from
 * horizontally scrolling content.
 */

.data-explorer-grid :deep(thead th:first-child),
.data-explorer-grid :deep(tbody td:first-child) {
    box-shadow: 2px 0 3px rgba(0, 0, 0, 0.18);
}

/* =========================================================
   SCROLLBAR
   ========================================================= */

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
</style>
