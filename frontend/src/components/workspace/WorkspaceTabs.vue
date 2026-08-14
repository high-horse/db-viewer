<template>
    <div
        class="h-9 shrink-0 bg-[#161310] border-b border-[#292521] flex items-center"
    >
        <!-- Tabs -->
        <div
            class="flex items-center h-full flex-1 min-w-0 overflow-x-auto"
        >
            <div
                v-for="tab in tabs"
                :key="tab.id"
                class="h-full flex items-center shrink-0 border-r border-[#292521]"
                :class="
                    tab.id === activeTabId
                        ? 'bg-[#0c0b09] border-b-2 border-b-amber-500'
                        : 'bg-[#161310] hover:bg-[#231f1a]'
                "
            >
                <!-- Select -->
                <button
                    type="button"
                    class="h-full px-3 flex items-center gap-2 text-xs font-mono min-w-0"
                    :class="
                        tab.id === activeTabId
                            ? 'text-amber-400'
                            : 'text-[#6b7280] hover:text-[#d1d5db]'
                    "
                    @click="selectTab(tab.id)"
                >
                    <!-- Different icon based on type -->
                    <q-icon
                        :name="
                            tab.type === 'result'
                                ? 'table_view'
                                : 'code'
                        "
                        size="12px"
                    />

                    <span class="max-w-32 truncate">
                        {{ tab.title }}
                    </span>

                    <!-- Query dirty indicator -->
                    <span
                        v-if="tab.type === 'query' && tab.dirty"
                        class="text-amber-500 text-[9px]"
                    >
                        ●
                    </span>

                    <!-- Loading -->
                    <q-spinner
                        v-if="tab.loading"
                        color="amber"
                        size="11px"
                    />
                </button>

                <!-- Close -->
                <button
                    type="button"
                    class="w-5 h-full flex items-center justify-center text-[#4b4540] hover:text-white hover:bg-[#231f1a]"
                    @click.stop="closeTab(tab.id)"
                >
                    <q-icon
                        name="close"
                        size="13px"
                    />
                </button>
            </div>

            <!-- New Query -->
            <button
                type="button"
                class="h-full w-9 shrink-0 flex items-center justify-center text-[#6b7280] hover:text-amber-400 hover:bg-[#231f1a]"
                title="New Query"
                @click="createTab"
            >
                <q-icon
                    name="add"
                    size="16px"
                />
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { QueryTab } from "@/types/queryTab";

defineProps<{
    tabs: QueryTab[];
    activeTabId: string | null;
}>();

const emit = defineEmits<{
    createTab: [];
    selectTab: [id: string];
    closeTab: [id: string];
}>();

function createTab() {
    emit("createTab");
}

function selectTab(id: string) {
    emit("selectTab", id);
}

function closeTab(id: string) {
    emit("closeTab", id);
}
</script>