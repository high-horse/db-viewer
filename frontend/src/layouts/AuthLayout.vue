<template>
    <div class="h-screen flex flex-col bg-[#0c0b09]">
        <!-- Top Navbar -->
        <header
            class="h-14 flex items-center justify-between px-5 border-b border-[#292521] bg-[#161310]"
        >
            <div class="flex items-center gap-3">
                <q-icon name="storage" size="22px" class="text-amber-400" />

                <div>
                    <div class="text-sm font-bold text-white">DB Viewer</div>
                    <div class="text-[10px] text-gray-500">
                        Database Workspace
                    </div>
                </div>
            </div>

            <q-btn
                unelevated
                rounded
                color="amber"
                text-color="black"
                icon="add"
                label="New Connection"
                class="text-xs font-bold text-capitalize"
                @click="$router.push({ name: 'Welcome', query: {} })"
            />
        </header>

        <!-- Main Area -->
        <div class="flex flex-1 overflow-hidden">
            <!-- Left Connections Panel -->
            <aside class="w-72 border-r border-[#292521] bg-[#100e0c]">
                <q-scroll-area class="h-full p-3">
                    <div class="space-y-2">
                        <div
                            v-for="connection in connections"
                            :key="connection.id"
                            class="p-3 rounded-lg border border-[#292521] bg-[#161310]/60 hover:bg-[#231f1a] cursor-pointer transition select-none"
                            @dblclick="selectConnection(connection)"
                        >
                            <div class="flex justify-between">
                                <div class="flex items-center gap-2">
                                    <q-icon
                                        name="lan"
                                        size="16px"
                                        class="text-amber-400"
                                    />

                                    <div class="min-w-0">
                                        <div
                                            class="text-sm font-semibold text-white truncate"
                                        >
                                            {{ connection.name }}
                                        </div>

                                        <div class="text-xs text-gray-400">
                                            {{ connection.driver }}

                                            <span
                                                v-if="connection.host"
                                                class="text-[11px] text-gray-500 font-mono"
                                            >
                                                {{ connection.host }}:{{
                                                    connection.port.Int64
                                                }}
                                            </span>

                                            <span
                                                v-else
                                                class="text-[11px] text-gray-500 font-mono truncate"
                                            >
                                                {{
                                                    connection.dbname
                                                        ? connection.dbname
                                                              .split("/")
                                                              .pop()
                                                        : ""
                                                }}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                                <q-btn flat round dense icon="more_vert" color="amber" :aria-label="`Actions for ${connection.name}`" @click.stop @dblclick.stop>
                                    <q-menu class="bg-[#161310] text-gray-300 border border-[#292521]">
                                        <q-list dense>
                                            <q-item clickable v-close-popup @click="editConnection(connection.id)">
                                                <q-item-section avatar><q-icon name="edit" size="16px" /></q-item-section>
                                                <q-item-section>Edit connection</q-item-section>
                                            </q-item>
                                            <q-item clickable v-close-popup class="text-red-400" @click="connectionToDelete = connection">
                                                <q-item-section avatar><q-icon name="delete_outline" size="16px" /></q-item-section>
                                                <q-item-section>Delete connection</q-item-section>
                                            </q-item>
                                        </q-list>
                                    </q-menu>
                                </q-btn>
                            </div>
                        </div>
                    </div>
                </q-scroll-area>
            </aside>

            <!-- Router Content -->
            <main
                class="flex-1 flex items-center justify-center p-8 bg-radial-gradient"
            >
                <q-layout>
                    <q-page-container>
                        <router-view />
                    </q-page-container>
                </q-layout>
            </main>
        </div>
        <q-dialog :model-value="!!connectionToDelete" :persistent="deletingConnection" @update:model-value="value => { if (!value) connectionToDelete = null; }">
            <q-card class="bg-[#161310] text-gray-300" style="width: 400px; max-width: 90vw">
                <q-card-section class="text-base font-semibold">Delete connection</q-card-section>
                <q-card-section class="pt-0">
                    Delete “{{ connectionToDelete?.name }}” from saved connections? Your database and its data will remain unchanged.
                </q-card-section>
                <q-card-actions align="right">
                    <q-btn flat no-caps label="Cancel" :disable="deletingConnection" @click="connectionToDelete = null" />
                    <q-btn flat no-caps color="red" label="Delete" :loading="deletingConnection" @click="deleteConnection" />
                </q-card-actions>
            </q-card>
        </q-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useConnectionStore } from "@/stores/connectionStore";
import { useActiveConnection } from "@/stores/activeConnection";
import type { Connection } from "@bindings/db-viewer/internal/types";

const $router = useRouter();
const activeConnectionStore = useActiveConnection();
const store = useConnectionStore();
const connectionToDelete = ref<Connection | null>(null);
const deletingConnection = ref(false);

async function deleteConnection() {
    const connection = connectionToDelete.value;
    if (!connection || deletingConnection.value) return;
    deletingConnection.value = true;
    try {
        if (await store.deleteConnection(connection.id)) {
            connectionToDelete.value = null;
            if ($router.currentRoute.value.query.edit === String(connection.id)) {
                await $router.replace({ name: "Welcome", query: {} });
            }
        }
    } finally {
        deletingConnection.value = false;
    }
}

const { connections, showNewConnectionDialog, activeConnection } =
    storeToRefs(store);

function editConnection(id: number) {
    $router.push({ name: "Welcome", query: { edit: String(id) } });
}

async function selectConnection(connection: any) {
    store.setSelectedSession(connection);

    const connected = await store.connectToSession(connection);
    if (connected) {
        await activeConnectionStore.setActiveConnection();
        $router.push({ name: "WorkSpace" });
    }
}

watch(activeConnection, (newConnection) => {
    if (newConnection) {
        console.log("connection dd", newConnection);
    }
});

</script>
