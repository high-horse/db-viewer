<template>
    <div
        class="h-full w-full min-w-0 min-h-0 flex flex-col bg-[#0c0b09] overflow-hidden"
    >
        <div
            class="h-9 shrink-0 bg-[#161310] border-b border-[#292521] flex items-center"
        >
            <!-- Tabs Container -->
            <div
                class="flex items-center h-full flex-1 min-w-0 overflow-x-auto"
            >
                <!-- Existing Tabs -->
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
                    <!-- Tab Select -->
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
                        <q-icon name="code" size="12px" />

                        <span class="max-w-32 truncate">
                            {{ tab.title }}
                        </span>

                        <!-- Dirty Indicator -->
                        <span
                            v-if="tab.dirty"
                            class="text-amber-500 text-[9px]"
                        >
                            ●
                        </span>

                        <!-- Loading Indicator -->
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
                        <q-icon name="close" size="13px" />
                    </button>
                </div>

                <!-- New Tab -->
                <button
                    type="button"
                    class="h-full w-9 shrink-0 flex items-center justify-center text-[#6b7280] hover:text-amber-400 hover:bg-[#231f1a]"
                    title="New Query"
                    @click="createTab"
                >
                    <q-icon name="add" size="16px" />
                </button>
            </div>

            <!-- Run Button -->
            <div
                class="h-full flex items-center px-2 border-l border-[#292521]"
            >
                <q-btn
                    unelevated
                    color="amber"
                    text-color="black"
                    size="sm"
                    class="font-bold px-3 tracking-wide"
                    :disable="!activeTab || !activeTab.sql.trim()"
                    :loading="activeTab?.loading"
                    @click.prevent="sqlEditorRef?.execute()"
                >
                    <q-icon name="play_arrow" size="16px" class="mr-1" />

                    RUN
                </q-btn>
            </div>
        </div>

        <div
            class="h-7 shrink-0 bg-[#100e0c] border-b border-[#292521] flex items-center justify-between px-3"
        >
            <div class="flex items-center gap-3">
                <span
                    v-if="activeTab"
                    class="text-[10px] font-mono text-[#6b7280]"
                >
                    {{ activeTab.title }}
                </span>

                <span
                    v-if="activeTab?.dirty"
                    class="text-[10px] text-amber-500"
                >
                    Modified
                </span>
            </div>

            <div
                v-if="activeTab"
                class="flex items-center gap-2 text-[10px] text-[#4b4540]"
            >
                <span>Ctrl + Enter</span>
                <span>Run</span>
            </div>
        </div>

        <div
            v-if="activeTab"
            class="flex-grow min-h-0 min-w-0 relative overflow-hidden"
        >
            <SqlEditor
                v-if="activeTab"
                 ref="sqlEditorRef"
                :model-value="activeTab.sql"
                :db-driver="'pgx'"
                @update:model-value="handleSqlUpdate"
                @execute="executeActiveTab"
            />
        </div>

        <div
            v-else
            class="flex-grow flex flex-col items-center justify-center gap-3 bg-[#0c0b09]"
        >
            <q-icon name="code" size="36px" class="text-[#4b4540]" />

            <div class="text-xs text-[#6b7280]">No query tab open</div>

            <q-btn
                unelevated
                color="amber"
                text-color="black"
                size="sm"
                label="New Query"
                icon="add"
                @click="createTab"
            />
        </div>

        <div
            v-if="activeTab?.error && false"
            class="min-h-8 max-h-24 overflow-auto bg-red-950/20 border-t border-red-900/50 px-3 py-2"
        >
            <div class="flex items-start gap-2">
                <q-icon
                    name="error_outline"
                    size="15px"
                    class="text-red-400 mt-0.5"
                />

                <span class="text-[11px] text-red-300 font-mono">
                    {{ activeTab?.error }}
                </span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, onBeforeUnmount } from "vue";
import type { QueryTab } from "@/types/queryTab";
import SqlEditor from "./SqlEditor.vue";

const sqlEditorRef = ref<InstanceType<typeof SqlEditor> | null>(null);
const props = defineProps<{
    tabs: QueryTab[];
    activeTabId: string | null;
    activeTab: QueryTab | null;
  
    onExecute: (id: string, sql: string) => void;
}>();

const emit = defineEmits<{
    createTab: [];
    selectTab: [id: string];
    closeTab: [id: string];
    updateSql: [id: string, sql: string];
    toggleResultTab: [];
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

function handleSqlUpdate(sql: string) {
    if (!props.activeTab) {
        return;
    }

    emit("updateSql", props.activeTab.id, sql);
}

function executeActiveTab(sql: string) {
    console.log("SQL received from editor:", sql);

    if (!props.activeTab) {
        return;
    }

    if (!sql.trim()) {
        return;
    }

    props.onExecute(
        props.activeTab.id,
        sql,
    );
}

function handleToggleResultTabShortcut(event: KeyboardEvent) {
    if (event.ctrlKey && event.key.toLowerCase() === "j") {
        emit("toggleResultTab");
        event.preventDefault();
    }
}

function handleCreateTabShortcut(event: KeyboardEvent) {
    if (event.ctrlKey && event.key.toLowerCase() === "t") {
        event.preventDefault();
        createTab();
    }
}

function handleCloseTabShortcut(event: KeyboardEvent) {
    if (event.ctrlKey && event.key.toLowerCase() === "w") {
        event.preventDefault();
        closeTab(props.activeTabId ?? "");
    }
}

onBeforeMount(() => {
    window.addEventListener("keydown", handleCreateTabShortcut);
    window.addEventListener("keydown", handleCloseTabShortcut);
    window.addEventListener("keydown", handleToggleResultTabShortcut);
});

onBeforeUnmount(() => {
    window.removeEventListener("keydown", handleCreateTabShortcut);
    window.removeEventListener("keydown", handleCloseTabShortcut);
    window.removeEventListener("keydown", handleToggleResultTabShortcut);
});

</script>
