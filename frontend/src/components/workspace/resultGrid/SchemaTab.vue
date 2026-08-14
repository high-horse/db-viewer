<template>
    <div class="h-full overflow-auto bg-[#100e0c]">
        <table class="w-full border-collapse font-mono text-[11px]">
            <thead>
                <tr
                    class="bg-[#161310] text-amber-400 border-b-2 border-[#292521]"
                >
                    <th class="schema-cell text-left w-10">
                        #
                    </th>

                    <th class="schema-cell text-left">
                        Column
                    </th>

                    <th class="schema-cell text-left">
                        Type
                    </th>

                    <th class="schema-cell text-left">
                        Nullable
                    </th>

                    <th class="schema-cell text-left">
                        Default
                    </th>
                </tr>
            </thead>

            <tbody>
                <tr
                    v-for="(column, index) in result.Columns"
                    :key="column.Name"
                    class="schema-row"
                >
                    <td class="schema-cell text-[#4b5563]">
                        {{ index + 1 }}
                    </td>

                    <td
                        class="schema-cell text-white font-medium"
                    >
                        <div class="flex items-center gap-2">
                            <q-icon
                                name="view_column"
                                size="14px"
                                class="text-amber-500"
                            />

                            {{ column.Name }}
                        </div>
                    </td>

                    <td class="schema-cell text-blue-300">
                        {{ column.Type }}
                    </td>

                    <td class="schema-cell">
                        <span
                            :class="
                                column.Nullable
                                    ? 'text-teal-400'
                                    : 'text-amber-400'
                            "
                        >
                            {{
                                column.Nullable
                                    ? "YES"
                                    : "NO"
                            }}
                        </span>
                    </td>

                    <td class="schema-cell text-[#94a3b8]">
                        <span
                            v-if="column.DefaultValue"
                            class="text-purple-300"
                        >
                            {{ column.DefaultValue }}
                        </span>

                        <span
                            v-else
                            class="text-[#4b4540]"
                        >
                            —
                        </span>
                    </td>
                </tr>
            </tbody>
        </table>

        <!-- Empty schema -->
        <div
            v-if="!result.Columns.length"
            class="h-full flex items-center justify-center text-[#4b4540] text-xs"
        >
            No column information available
        </div>
    </div>
</template>

<script setup lang="ts">
import type { QueryResult } from "@/types/queryTab";

defineProps<{
    result: QueryResult;
}>();
</script>

<style scoped>
.schema-row {
    background-color: #100e0c;
}

.schema-row:nth-child(even) {
    background-color: #13110f;
}

.schema-row:hover {
    background-color: #231f1a;
}

.schema-cell {
    padding: 8px 10px;
    border-bottom: 1px solid #292521;
    white-space: nowrap;
}

thead th.schema-cell {
    position: sticky;
    top: 0;
    z-index: 2;

    padding-top: 7px;
    padding-bottom: 7px;

    font-weight: bold;
    font-size: 10px;
}

table {
    min-width: 100%;
}

tbody td {
    font-size: 11px;
}
</style>