<template>
    <div
        class="h-full w-full min-w-0 min-h-0 flex flex-col bg-[#0c0b09] overflow-hidden"
    >

        <!-- Query info -->
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
                    <q-icon
                        name="play_arrow"
                        size="16px"
                        class="mr-1"
                    />

                    RUN
                </q-btn>
            </div>
        </div>

        <!-- SQL Editor -->
        <div
            v-if="activeTab"
            class="flex-grow min-h-0 min-w-0 relative overflow-hidden"
        >
            <SqlEditor
                ref="sqlEditorRef"
                :model-value="activeTab.sql"
                :db-driver="'pgx'"
                @update:model-value="handleSqlUpdate"
                @execute="executeActiveTab"
            />
        </div>

        <!-- No active query -->
        <div
            v-else
            class="flex-grow flex flex-col items-center justify-center gap-3 bg-[#0c0b09]"
        >
            <q-icon
                name="code"
                size="36px"
                class="text-[#4b4540]"
            />

            <div class="text-xs text-[#6b7280]">
                No query tab open
            </div>

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
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { QueryTab } from "@/types/queryTab";
import SqlEditor from "./SqlEditor.vue";

const sqlEditorRef =
    ref<InstanceType<typeof SqlEditor> | null>(null);

const props = defineProps<{
    activeTab: QueryTab | null;
    onExecute: (id: string, sql: string) => void;
}>();

const emit = defineEmits<{
    createTab: [];
    updateSql: [id: string, sql: string];
}>();

function createTab() {
    emit("createTab");
}

function handleSqlUpdate(sql: string) {
    if (!props.activeTab) {
        return;
    }

    emit(
        "updateSql",
        props.activeTab.id,
        sql,
    );
}

function executeActiveTab(sql: string) {
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
</script>