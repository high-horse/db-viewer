<template>
    <div class="relative h-full w-full overflow-hidden bg-[#100e0c]">
        <!-- =========================================================
             TABLE
             ========================================================= -->

        <q-table
            flat
            square
            dense
            dark
            :rows="paginatedRows"
            :columns="mappedColumns"
            row-key="id"
            :pagination="{
                rowsPerPage: rowsPerPage,
            }"
            :virtual-scroll-item-size="28"
            class="absolute inset-0 data-explorer-grid bg-transparent"
            table-class="table-fixed"
            table-style="min-width: 100%; width: max-content;"
            hide-pagination
        >
            <!-- =====================================================
                 HEADER
                 ===================================================== -->

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

                        <!-- Resize handle -->
                        <span
                            v-if="col.name !== 'sn'"
                            class="column-resizer"
                            @mousedown.stop="startResize($event, col.name)"
                        />

                        <!-- Column tooltip -->
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

                        <!-- Context menu -->
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
                                @click="clearSort"
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

            <!-- =====================================================
                 BODY
                 ===================================================== -->

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
                        <!-- Row number -->
                        <template v-if="col.name === 'sn'">
                            <span class="text-grey-8">
                                {{
                                    (currentPage - 1) * rowsPerPage +
                                    props.rowIndex +
                                    1
                                }}
                            </span>
                        </template>

                        <!-- Cell -->
                        <template v-else>
                            <div
                                class="block w-full min-w-0 overflow-hidden text-ellipsis whitespace-nowrap"
                                :title="getCellTitle(props.row[col.field])"
                            >
                                {{ props.row[col.field] }}
                            </div>
                        </template>
                    </q-td>
                </q-tr>
            </template>
        </q-table>

        <!-- =========================================================
             FIXED BOTTOM TOOLBAR
             ========================================================= -->

        <div
            class="data-grid-footer absolute bottom-0 left-0 right-0 z-[100] flex h-[36px] items-center justify-between border-t border-[#292521] bg-[#161310] px-2"
        >
            <!-- Left side -->
            <div class="flex items-center gap-2">
                <!-- Refresh -->
                <button
                    type="button"
                    class="flex h-[24px] w-[24px] items-center justify-center rounded text-gray-500 transition-colors hover:bg-[#292521] hover:text-amber-400 disabled:cursor-not-allowed disabled:opacity-40"
                    :disabled="isRefreshing"
                    title="Refresh"
                    @click="refresh"
                >
                    <q-icon
                        name="refresh"
                        size="16px"
                        :class="{ 'animate-spin': isRefreshing }"
                    />
                </button>

                <span class="font-mono text-[10px] text-gray-600">
                    {{ totalRows }} rows
                </span>
            </div>

            <!-- Right side pagination -->
            <div class="flex items-center gap-1">
                <!-- First page -->
                <button
                    type="button"
                    class="pagination-button"
                    :disabled="currentPage === 1"
                    title="First page"
                    @click="goToFirstPage"
                >
                    <q-icon name="first_page" size="16px" />
                </button>

                <!-- Previous -->
                <button
                    type="button"
                    class="pagination-button"
                    :disabled="currentPage === 1"
                    title="Previous page"
                    @click="previousPage"
                >
                    <q-icon name="chevron_left" size="16px" />
                </button>

                <!-- Page -->
                <div
                    class="mx-1 flex h-[24px] min-w-[80px] items-center justify-center border border-[#292521] bg-[#100e0c] px-2 font-mono text-[10px] text-gray-400"
                >
                    Page {{ currentPage }} / {{ totalPages }}
                </div>

                <!-- Next -->
                <button
                    type="button"
                    class="pagination-button"
                    :disabled="currentPage >= totalPages"
                    title="Next page"
                    @click="nextPage"
                >
                    <q-icon name="chevron_right" size="16px" />
                </button>

                <!-- Last page -->
                <button
                    type="button"
                    class="pagination-button"
                    :disabled="currentPage >= totalPages"
                    title="Last page"
                    @click="goToLastPage"
                >
                    <q-icon name="last_page" size="16px" />
                </button>

                <!-- Rows per page -->
                <select
                    v-model.number="rowsPerPage"
                    class="ml-2 h-[24px] border border-[#292521] bg-[#100e0c] px-1 font-mono text-[10px] text-gray-400 outline-none"
                >
                    <option :value="25">25</option>
                    <option :value="50">50</option>
                    <option :value="100">100</option>
                    <option :value="250">250</option>
                    <option :value="500">500</option>
                </select>
            </div>
        </div>
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

const emit = defineEmits<{
    refresh: [];
}>();

/* =========================================================
   COLUMN WIDTHS
   ========================================================= */

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

/* =========================================================
   SAFE CELL TITLE
   ========================================================= */

const getCellTitle = (value: unknown): string => {
    if (value === null || value === undefined) {
        return "";
    }

    return String(value);
};

/* =========================================================
   HEADER WIDTH MEASUREMENT
   ========================================================= */

const measureHeaderWidths = async () => {
    await nextTick();

    await new Promise<void>((resolve) => {
        requestAnimationFrame(() => resolve());
    });

    const table = document.querySelector(
        ".data-explorer-grid table",
    );

    if (!table) {
        return;
    }

    const headers = table.querySelectorAll("thead th");

    headers.forEach((header, index) => {
        // Skip # column.
        if (index === 0) {
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

        const contentWidth = Math.max(
            th.scrollWidth,
            content.scrollWidth,
        );

        // Restore.
        th.style.width = oldWidth;
        th.style.minWidth = oldMinWidth;
        th.style.maxWidth = oldMaxWidth;

        tableElement.style.tableLayout = oldTableLayout;

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

/* =========================================================
   COLUMN WATCH
   ========================================================= */

const columnSignature = computed(() =>
    props.result.Columns
        .map(
            (column) =>
                `${column.Name}:${column.Type}`,
        )
        .join("|"),
);

watch(columnSignature, async () => {
    const validColumns = new Set(
        props.result.Columns.map(
            (column) => column.Name,
        ),
    );

    Object.keys(columnWidths).forEach(
        (columnName) => {
            if (!validColumns.has(columnName)) {
                delete columnWidths[columnName];
            }
        },
    );

    currentPage.value = 1;

    await nextTick();

    await measureHeaderWidths();
});

/* =========================================================
   COLUMN RESIZING
   ========================================================= */

let resizingColumn: string | null = null;

let resizeStartX = 0;

let resizeStartWidth = 0;

const startResize = (
    event: MouseEvent,
    columnName: string,
) => {
    event.preventDefault();
    event.stopPropagation();

    resizingColumn = columnName;

    resizeStartX = event.clientX;

    const target =
        event.currentTarget as HTMLElement;

    const header =
        target.closest("th") as HTMLElement | null;

    resizeStartWidth =
        header?.getBoundingClientRect().width ??
        getColumnWidth(columnName);

    columnWidths[columnName] =
        resizeStartWidth;

    document.body.classList.add(
        "resizing-column",
    );

    document.addEventListener(
        "mousemove",
        handleResize,
    );

    document.addEventListener(
        "mouseup",
        stopResize,
    );
};

const handleResize = (event: MouseEvent) => {
    if (!resizingColumn) {
        return;
    }

    const delta =
        event.clientX - resizeStartX;

    const newWidth = Math.max(
        MIN_COLUMN_WIDTH,
        resizeStartWidth + delta,
    );

    columnWidths[resizingColumn] =
        newWidth;
};

const stopResize = () => {
    resizingColumn = null;

    document.body.classList.remove(
        "resizing-column",
    );

    document.removeEventListener(
        "mousemove",
        handleResize,
    );

    document.removeEventListener(
        "mouseup",
        stopResize,
    );
};

onBeforeUnmount(() => {
    stopResize();
});

/* =========================================================
   SORTING
   ========================================================= */

type SortDirection =
    | "asc"
    | "desc"
    | null;

const activeSortColumn =
    ref<string | null>(null);

const sortDirection =
    ref<SortDirection>(null);

const sortColumn = (
    columnName: string,
    direction: "asc" | "desc",
) => {
    activeSortColumn.value =
        columnName;

    sortDirection.value =
        direction;

    currentPage.value = 1;
};

const clearSort = () => {
    activeSortColumn.value = null;

    sortDirection.value = null;

    currentPage.value = 1;
};

/* =========================================================
   ROW MAPPING
   ========================================================= */

const mappedRows = computed(() => {
    const rows =
        props.result.Rows.map(
            (row, rowIndex) => {
                const rowObject =
                    {
                        id: rowIndex,
                    } as Record<
                        string,
                        unknown
                    >;

                props.result.Columns.forEach(
                    (
                        column,
                        columnIndex,
                    ) => {
                        rowObject[
                            column.Name
                        ] =
                            row[
                                columnIndex
                            ];
                    },
                );

                return rowObject;
            },
        );

    if (
        !activeSortColumn.value ||
        !sortDirection.value
    ) {
        return rows;
    }

    const column =
        activeSortColumn.value;

    const direction =
        sortDirection.value === "asc"
            ? 1
            : -1;

    return [...rows].sort(
        (a, b) => {
            const valueA =
                a[column];

            const valueB =
                b[column];

            if (
                valueA === null ||
                valueA === undefined
            ) {
                if (
                    valueB === null ||
                    valueB === undefined
                ) {
                    return 0;
                }

                return -1 * direction;
            }

            if (
                valueB === null ||
                valueB === undefined
            ) {
                return 1 * direction;
            }

            if (
                typeof valueA ===
                    "number" &&
                typeof valueB ===
                    "number"
            ) {
                return (
                    (valueA - valueB) *
                    direction
                );
            }

            const stringA =
                String(valueA);

            const stringB =
                String(valueB);

            const numberA =
                Number(stringA);

            const numberB =
                Number(stringB);

            if (
                stringA.trim() !== "" &&
                stringB.trim() !== "" &&
                Number.isFinite(
                    numberA,
                ) &&
                Number.isFinite(
                    numberB,
                )
            ) {
                return (
                    (numberA - numberB) *
                    direction
                );
            }

            return (
                stringA.localeCompare(
                    stringB,
                    undefined,
                    {
                        numeric: true,
                        sensitivity:
                            "base",
                    },
                ) * direction
            );
        },
    );
});

/* =========================================================
   PAGINATION
   ========================================================= */

const rowsPerPage = ref(100);

const currentPage = ref(1);

const totalRows = computed(
    () => mappedRows.value.length,
);

const totalPages = computed(() =>
    Math.max(
        1,
        Math.ceil(
            totalRows.value /
                rowsPerPage.value,
        ),
    ),
);

const paginatedRows = computed(() => {
    const start =
        (currentPage.value - 1) *
        rowsPerPage.value;

    const end =
        start + rowsPerPage.value;

    return mappedRows.value.slice(
        start,
        end,
    );
});

const goToFirstPage = () => {
    currentPage.value = 1;
};

const goToLastPage = () => {
    currentPage.value =
        totalPages.value;
};

const previousPage = () => {
    if (currentPage.value > 1) {
        currentPage.value--;
    }
};

const nextPage = () => {
    if (
        currentPage.value <
        totalPages.value
    ) {
        currentPage.value++;
    }
};

/*
 * When rows-per-page changes, return
 * to the first page.
 */
watch(rowsPerPage, () => {
    currentPage.value = 1;
});

/*
 * Prevent current page from becoming
 * invalid after result data changes.
 */
watch(totalPages, (pages) => {
    if (currentPage.value > pages) {
        currentPage.value = pages;
    }
});

/* =========================================================
   REFRESH
   ========================================================= */

const isRefreshing = ref(false);

const refresh = async () => {
    if (isRefreshing.value) {
        return;
    }

    isRefreshing.value = true;

    try {
        emit("refresh");
    } finally {
        /*
         * Small delay so the refresh indicator
         * is visible even for very fast requests.
         */
        await new Promise<void>(
            (resolve) => {
                setTimeout(resolve, 300);
            },
        );

        isRefreshing.value = false;
    }
};

/* =========================================================
   COLUMNS
   ========================================================= */

const mappedColumns =
    computed<QTableColumn[]>(() => {
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

            ...props.result.Columns.map(
                (column) => {
                    const width =
                        getColumnWidth(
                            column.Name,
                        );

                    return {
                        name: column.Name,

                        label: column.Name,

                        field: column.Name,

                        align: "left" as const,

                        sortable: false,

                        Type: column.Type,

                        Nullable:
                            column.Nullable,

                        DefaultValue:
                            column.DefaultValue,

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
                },
            ),
        ];
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

.data-explorer-grid :deep(.q-table__middle) {
    background: #100e0c;

    /*
     * Leave room for the fixed footer.
     */
    padding-bottom: 36px;

    overflow: auto;

    scrollbar-width: thin;

    scrollbar-color:
        #292521
        #100e0c;
}

/* =========================================================
   TABLE LAYOUT
   ========================================================= */

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

.data-explorer-grid :deep(thead) {
    position: sticky;

    top: 0;

    z-index: 40;
}

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
   STICKY # BODY
   ========================================================= */

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
   ALTERNATING ROWS
   ========================================================= */

.data-explorer-grid
    :deep(
        tbody tr:nth-child(even)
            td:first-child
    ) {
    background: #13110f;
}

/* =========================================================
   HOVER
   ========================================================= */

.data-explorer-grid
    :deep(tbody tr:hover td:first-child) {
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

.data-explorer-grid
    :deep(tbody td > div) {
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

.data-explorer-grid
    :deep(thead th + th) {
    border-left-color: #332e28;
}

/* =========================================================
   STICKY COLUMN SHADOW
   ========================================================= */

.data-explorer-grid
    :deep(thead th:first-child),
.data-explorer-grid
    :deep(tbody td:first-child) {
    box-shadow:
        2px 0 3px
        rgba(0, 0, 0, 0.18);
}

/* =========================================================
   SCROLLBAR
   ========================================================= */

.data-explorer-grid
    :deep(.q-table__middle::-webkit-scrollbar) {
    width: 8px;

    height: 8px;
}

.data-explorer-grid
    :deep(.q-table__middle::-webkit-scrollbar-track) {
    background: #100e0c;
}

.data-explorer-grid
    :deep(.q-table__middle::-webkit-scrollbar-thumb) {
    background: #292521;

    border-radius: 4px;
}

.data-explorer-grid
    :deep(.q-table__middle::-webkit-scrollbar-thumb:hover) {
    background: #3a342e;
}

/* =========================================================
   PAGINATION FOOTER
   ========================================================= */

.data-grid-footer {
    box-shadow:
        0 -2px 8px
        rgba(0, 0, 0, 0.25);
}

/* =========================================================
   PAGINATION BUTTONS
   ========================================================= */

.pagination-button {
    display: flex;

    height: 24px;
    width: 24px;

    align-items: center;
    justify-content: center;

    border: 1px solid #292521;

    background: #100e0c;

    color: #6b7280;

    cursor: pointer;

    transition:
        background-color 80ms ease,
        color 80ms ease,
        border-color 80ms ease;
}

.pagination-button:hover:not(
        :disabled
    ) {
    background: #292521;

    border-color: #3a342e;

    color: #f59e0b;
}

.pagination-button:disabled {
    cursor: not-allowed;

    opacity: 0.3;
}

/* =========================================================
   SELECT
   ========================================================= */

.data-grid-footer select {
    cursor: pointer;
}

.data-grid-footer select:focus {
    border-color: #92400e;
}

/* =========================================================
   REFRESH ANIMATION
   ========================================================= */

.animate-spin {
    animation: spin 0.8s linear infinite;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}
</style>