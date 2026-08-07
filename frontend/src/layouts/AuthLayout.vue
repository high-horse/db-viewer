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
                @click="() => showNewConnectionDialog = true"
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
                                    </div>

                                    <div
                                        v-if="connection.host"
                                        class="text-[11px] text-gray-500 font-mono"
                                    >
                                        {{ connection.host }}:{{
                                            connection.port.Int64
                                        }}
                                    </div>

                                    <div
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
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </q-scroll-area>
            </aside>

            <!-- Router Content -->
            <main
                class="flex-1 flex items-center justify-center p-8 bg-radial-gradient"
            >
                <router-view />
            </main>
        </div>
    </div>
</template>

<script setup lang="ts">
import { watch } from "vue";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useConnectionStore } from "@/stores/connectionStore";

const router = useRouter();

const store = useConnectionStore();

const { connections, showNewConnectionDialog, activeConnection } = storeToRefs(store);

function selectConnection(connection: any) {
  store.setSelectedSession(connection);
  store.pingConnection();
}

watch(activeConnection, (newConnection) => {
    if (newConnection) {
        console.log("connection dd",newConnection);
    }
});
</script>