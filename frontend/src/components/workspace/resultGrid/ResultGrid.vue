<template>
    <div class="h-full flex flex-col bg-[#100e0c]">
        <!-- Result Header -->
        <div
            class="h-8 shrink-0 bg-[#161310] border-b border-[#292521] flex items-center justify-between px-3 text-[11px] font-mono text-[#6b7280]"
        >
            <div class="flex items-center gap-3">
                <span>RESULT GRID</span>

                <span
                    v-if="result && !loading"
                    class="text-teal-400 font-bold"
                >
                    • {{ result.Rows.length }} rows in this page
                </span>

                <span
                    v-if="loading"
                    class="text-amber-400 font-bold"
                >
                    • {{ fetchingLast ? 'Fetching remaining pages…' : result ? 'Fetching page…' : 'Executing…' }}
                </span>
            </div>

            <div class="flex items-center gap-3">
                <span
                    v-if="result && !loading"
                    class="text-[#6b7280]"
                >
                    {{ result.Duration }}ms execution latency
                </span>

                 <!-- <q-btn flat dense no-caps
                    type="button"
                    class="w-5 h-5 flex items-center justify-center rounded text-[#4b4540] hover:text-[#d1d5db] hover:bg-[#292521]"
                    title="Hide resultsssssss"
                    @click="emit('close')"
                >
                    <q-icon name="close" size="14px" />
                </q-btn> -->
            </div>
        </div>

        <!-- Loading -->
        <div
            v-if="loading && !result"
            class="flex-grow flex flex-col items-center justify-center gap-3"
        >
            <q-spinner-dots color="amber" size="32px" />

            <span class="text-xs text-[#6b7280] font-mono">
                Executing query...
            </span>
        </div>

        <!-- Error -->
        <div
            v-else-if="error"
            class="flex-grow flex flex-col items-center justify-center gap-3"
        >
            <q-icon
                name="error_outline"
                size="32px"
                class="text-red-400"
            />

            <div
                class="text-xs text-red-300 font-mono max-w-xl text-center px-6"
            >
                {{ error }}
            </div>
        </div>

        <!-- Result Content -->
        <template v-else-if="result">
            <div v-if="pageError" role="alert" class="px-3 py-2 text-xs text-red-300">{{ pageError }}</div>
            <q-tabs
                v-model="activeSubTab"
                dense
                no-caps
                align="left"
                active-color="amber"
                indicator-color="amber"
                class="result-sub-tabs shrink-0"
            >
                <q-tab
                    name="data"
                    label="Data"
                />
                    <!-- icon="table_rows" -->

                <q-tab
                    name="schema"
                    label="Schema"
                />
            </q-tabs>

            <div class="flex-1 min-h-0 min-w-0 overflow-hidden">
                <ResultTab
                    v-if="activeSubTab === 'data'"
                    :result="result"
                    :table-tab-id="tableTabId"
                    :loading="loading"
                    :can-first="canFirst"
                    :can-previous="canPrevious"
                    :can-next="canNext"
                    :can-last="canLast"
                    :total-rows="totalRows"
                    :fetching-last="fetchingLast"
                    :sort-column="sortColumn"
                    :sort-direction="sortDirection"
                    @sort="(column, direction) => emit('sort', column, direction)"
                    @first="emit('first')"
                    @previous="emit('previous')"
                    @next="emit('next')"
                    @last="emit('last')"
                    @stop="emit('stop')"
                    @refresh="emit('refresh')"
                />
                <IndexSchema
                    v-else-if="activeSubTab === 'schema'"
                    :result="result"
                    :table-tab-id="tableTabId"
                />
            </div>


        </template>

        <!-- Empty -->
        <div
            v-else
            class="flex-grow flex flex-col justify-center items-center gap-3 text-[#4b4540]"
        >
            <q-icon
                name="table_rows"
                size="32px"
                class="opacity-40"
            />

            <span class="text-xs tracking-wider">
                Execute a statement query to inspect structural rows
            </span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import type { QueryResult } from "@/types/queryTab";
import ResultTab from "./ResultTab.vue";
import IndexSchema from "./Schema/Index.vue";

const props = defineProps<{
    result: QueryResult | null;
    tableTabId?: string;
    loading?: boolean;
    canFirst?: boolean;
    canPrevious?: boolean;
    canNext?: boolean;
    canLast?: boolean;
    totalRows?: number;
    fetchingLast?: boolean;
    sortColumn?: number;
    sortDirection?: "asc" | "desc";
    pageError?: string;
    error?: string | null;
}>();

const emit = defineEmits<{
    close: [];
    first: [];
    previous: [];
    next: [];
    last: [];
    stop: [];
    refresh: [];
    sort: [column: number, direction?: "asc" | "desc"];
}>();

const activeSubTab = ref<"data" | "schema">("data");

/*
 * Whenever a new result is loaded, go back to DATA.
 */
watch(
    () => props.result,
    () => {
        activeSubTab.value = "data";
    },
);
</script>


<style scoped>
.result-sub-tabs {
    height: 30px;
    min-height: 30px;

    background: #161310;
    border-bottom: 1px solid #292521;
}

.result-sub-tabs :deep(.q-tab) {
    height: 30px;
    min-height: 30px;

    min-width: 70px;
    padding: 0 10px;

    color: #6b7280;

    font-family: monospace;
    font-size: 10px;

    border-right: 1px solid #292521;
}

.result-sub-tabs :deep(.q-tab--active) {
    color: #f59e0b;
    background: #100e0c;
}

/* Important: align icon + text */
.result-sub-tabs :deep(.q-tab__content) {
    height: 100%;

    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;

    gap: 5px;
}

.result-sub-tabs :deep(.q-tab__icon) {
    font-size: 13px;
    line-height: 1;

    margin: 0;
}

.result-sub-tabs :deep(.q-tab__label) {
    font-size: 10px;
    line-height: 1;

    margin: 0;
}

.result-sub-tabs :deep(.q-tab__indicator) {
    height: 2px;
}
</style>
