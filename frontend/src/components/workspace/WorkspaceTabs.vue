<template>
    <div
        class="h-9 shrink-0 bg-[#161310] border-b border-[#292521] flex items-center min-w-0"
    >
        <!-- Horizontally scrollable tabs -->
        <q-scroll-area
            horizontal
            class="h-full flex-1 min-w-0 tabs-scroll-area"
        >
            <div class="tabs-content h-full flex items-center">
                <!-- Tabs -->
                <div
                    v-for="tab in tabs"
                    :key="tab.id"
                    class="h-full shrink-0 flex items-center border-r border-[#292521]"
                    :class="
                        tab.id === activeTabId
                            ? 'bg-[#0c0b09]'
                            : 'bg-[#161310] hover:bg-[#231f1a]'
                    "
                >
                    <!-- Tab button -->
                    <q-btn
                        flat
                        dense
                        no-caps
                        unelevated
                        class="tab-button"
                        :class="
                            tab.id === activeTabId
                                ? 'text-amber-400 active-tab'
                                : 'text-[#6b7280] hover:text-[#d1d5db]'
                        "
                        @click="selectTab(tab.id)"
                    >
                        <div
                            class="flex items-center gap-2 min-w-0 whitespace-nowrap"
                        >
                            <!-- Query / Result icon -->
                            <q-icon
                                :name="
                                    tab.type === 'result'
                                        ? 'table_view'
                                        : 'code'
                                "
                                size="13px"
                            />

                            <!-- Tab title -->
                            <span class="max-w-32 truncate">
                                {{ tab.title }}
                            </span>

                            <!-- Dirty indicator -->
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
                        </div>
                    </q-btn>

                    <!-- Close -->
                    <q-btn
                        flat
                        dense
                        round
                        size="xs"
                        icon="close"
                        class="tab-close"
                        :class="
                            tab.id === activeTabId
                                ? 'text-[#6b7280]'
                                : 'text-[#4b4540]'
                        "
                        @click.stop="closeTab(tab.id)"
                    >
                        <q-tooltip>
                            Close {{ tab.title }}
                        </q-tooltip>
                    </q-btn>
                </div>

                <!-- New query tab -->
                <q-btn
                    flat
                    dense
                    square
                    icon="add"
                    class="new-tab-button shrink-0"
                    color="grey-6"
                    title="New Query"
                    @click="createTab"
                >
                    <q-tooltip>New Query</q-tooltip>
                </q-btn>
            </div>
        </q-scroll-area>
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

<style scoped>
.tabs-scroll-area {
    min-width: 0;
    min-height: 0;
}

/*
 * Keep horizontal scrolling enabled.
 * Disable vertical scrolling.
 */
.tabs-scroll-area :deep(.q-scrollarea__container) {
    overflow-x: auto !important;
    overflow-y: hidden !important;

    /* Firefox */
    scrollbar-width: none;

    /* IE / old Edge */
    -ms-overflow-style: none;
}

/*
 * Chrome / Safari / Edge
 */
.tabs-scroll-area :deep(.q-scrollarea__container::-webkit-scrollbar) {
    width: 0 !important;
    height: 0 !important;
    display: none !important;
}

/*
 * Hide Quasar's custom scrollbar.
 *
 * Scrolling still works because the container remains
 * horizontally scrollable.
 */
.tabs-scroll-area :deep(.q-scrollarea__bar),
.tabs-scroll-area :deep(.q-scrollarea__thumb) {
    display: none !important;
    opacity: 0 !important;
    visibility: hidden !important;
}


.tabs-content {
    width: max-content;
    min-width: 100%;

    /*
     * Very important:
     * Never allow tabs to wrap onto another line.
     */
    flex-wrap: nowrap !important;
}


.tab-button {
    height: 36px !important;
    min-height: 36px !important;

    padding: 0 10px !important;

    border-radius: 0 !important;

    font-size: 11px;
    font-family: monospace;

    position: relative;
}

/*
 * Active tab bottom border
 */
.tab-button.active-tab::after {
    content: "";

    position: absolute;

    left: 0;
    right: 0;
    bottom: 0;

    height: 2px;

    background: #f59e0b;
}

/*
 * Remove Quasar focus helper.
 */
.tab-button :deep(.q-focus-helper) {
    display: none;
}

.tab-close {
    width: 22px;
    height: 28px;
    min-height: 28px;

    margin-right: 2px;

    border-radius: 3px !important;
}

.tab-close:hover {
    color: #ffffff !important;
    background: #292521 !important;
}


.new-tab-button {
    width: 36px;
    height: 36px;
    min-height: 36px;

    border-radius: 0 !important;
}

.new-tab-button:hover {
    color: #f59e0b !important;
    background: #231f1a !important;
}
</style>