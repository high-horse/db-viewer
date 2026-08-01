<template>
    <div
        class="h-screen w-screen bg-[#070a12] flex font-sans select-none overflow-hidden text-[#94a3b8]"
    >
        <!-- Left Hand Grid Panel: Historical Saved Server Workspace Targets -->
        <div
            class="w-80 h-full bg-[#0f172a]/30 border-r border-[#1f2937] flex flex-col"
        >
            <div
                class="p-4 border-b border-[#1f2937] flex items-center justify-between"
            >
                <span
                    class="text-xs font-bold uppercase tracking-wider text-slate-400"
                    >Saved Workspaces</span
                >
                <q-icon name="cloud_queue" size="16px" class="text-slate-500" />
            </div>

            <q-scroll-area class="flex-grow p-3">
                <div class="space-y-2">
                    <!-- Sample list row card profile -->
                    <div
                        v-for="idx in 3"
                        :key="idx"
                        class="p-3 rounded-lg border border-[#1f2937] bg-[#111827]/40 hover:bg-[#1e293b]/40 cursor-pointer transition-all duration-200 group"
                    >
                        <div class="flex items-center gap-2.5">
                            <q-icon
                                name="lan"
                                size="14px"
                                class="text-indigo-400 group-hover:scale-110 transition-transform"
                            />
                            <div>
                                <div class="text-xs font-semibold text-white">
                                    Local_Postgres_{{ idx }}
                                </div>
                                <div
                                    class="text-[10px] font-mono text-[#6b7280]"
                                >
                                    localhost:5432 • dvdrental
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </q-scroll-area>
        </div>

        <!-- Right Hand Grid Panel: Interactive Core Form Setup Panel -->
        <div
            class="flex-grow h-full flex items-center justify-center p-8 bg-radial-gradient"
        >
            <div
                class="w-full max-w-md bg-[#111827]/80 backdrop-blur-md p-6 rounded-xl border border-[#1f2937] shadow-2xl"
            >
                <div class="flex items-center gap-3 mb-6">
                    <div
                        class="p-2 rounded-lg bg-indigo-500/10 border border-indigo-500/20"
                    >
                        <q-icon
                            name="add_link"
                            size="24px"
                            class="text-indigo-400"
                        />
                    </div>
                    <div>
                        <h4 class="text-sm font-bold text-white tracking-wide">
                            Setup Database Session
                        </h4>
                        <p class="text-[11px] text-[#6b7280] m-0">
                            Initialize single engine configuration links
                            dynamically
                        </p>
                    </div>
                </div>

                <div class="space-y-4">
                    <!-- Selection Picker Flag Row -->
                    <div>
                        <label
                            class="block text-[11px] font-bold uppercase tracking-wider text-slate-400 mb-1.5"
                            >Engine Type</label
                        >
                        <div class="grid grid-cols-3 gap-2">
                            <div
                                v-for="engine in [
                                    'PostgreSQL',
                                    'MySQL',
                                    'SQLite',
                                ]"
                                :key="engine"
                                :class="
                                    form.type === engine
                                        ? 'border-indigo-500 bg-indigo-500/5 text-white'
                                        : 'border-[#1f2937] hover:border-slate-700'
                                "
                                class="p-2 border rounded-md text-center text-xs font-semibold cursor-pointer transition-all"
                                @click="form.type = engine"
                            >
                                {{ engine }}
                            </div>
                        </div>
                    </div>

                    <!-- Network Configurations Fields Grid Layout -->
                    <div
                        v-if="form.type !== 'SQLite'"
                        class="grid grid-cols-4 gap-3"
                    >
                        <div class="col-span-3">
                            <label
                                class="block text-[11px] font-bold text-slate-400 mb-1"
                                >Host address</label
                            >
                            <input
                                v-model="form.host"
                                type="text"
                                class="w-full h-8 bg-[#070a12] border border-[#1f2937] rounded px-3 text-xs text-slate-200 focus:outline-none focus:border-indigo-500 font-mono"
                            />
                        </div>
                        <div class="col-span-1">
                            <label
                                class="block text-[11px] font-bold text-slate-400 mb-1"
                                >Port</label
                            >
                            <input
                                v-model="form.port"
                                type="number"
                                class="w-full h-8 bg-[#070a12] border border-[#1f2937] rounded px-2 text-xs text-slate-200 focus:outline-none focus:border-indigo-500 font-mono"
                            />
                        </div>
                    </div>

                    <!-- Credential Authentication Parameter Layout Fields -->
                    <div
                        v-if="form.type !== 'SQLite'"
                        class="grid grid-cols-2 gap-3"
                    >
                        <div>
                            <label
                                class="block text-[11px] font-bold text-slate-400 mb-1"
                                >Username User</label
                            >
                            <input
                                v-model="form.user"
                                type="text"
                                class="w-full h-8 bg-[#070a12] border border-[#1f2937] rounded px-3 text-xs text-slate-200 focus:outline-none focus:border-indigo-500"
                            />
                        </div>
                        <div>
                            <label
                                class="block text-[11px] font-bold text-slate-400 mb-1"
                                >Password</label
                            >
                            <input
                                v-model="form.password"
                                type="password"
                                class="w-full h-8 bg-[#070a12] border border-[#1f2937] rounded px-3 text-xs text-slate-200 focus:outline-none focus:border-indigo-500"
                            />
                        </div>
                    </div>

                    <!-- Database Mapping Scope Field Target -->
                    <div>
                        <label
                            class="block text-[11px] font-bold text-slate-400 mb-1"
                        >
                            {{
                                form.type === "SQLite"
                                    ? "SQLite Database File Path"
                                    : "Target Database Catalog Name"
                            }}
                        </label>
                        <input
                            v-model="form.database"
                            type="text"
                            :placeholder="
                                form.type === 'SQLite'
                                    ? '/home/user/app.db'
                                    : 'dvdrental'
                            "
                            class="w-full h-8 bg-[#070a12] border border-[#1f2937] rounded px-3 text-xs text-slate-200 focus:outline-none focus:border-indigo-500"
                        />
                    </div>

                    <!-- Interactive Command Form Actions Submit Buttons -->
                    <div class="pt-2">
                        <q-btn
                            unelevated
                            color="indigo"
                            class="w-full font-bold text-white h-9 tracking-wide"
                            @click="$emit('select-db', form)"
                        >
                            ESTABLISH WORKSPACE HANDSHAKE
                        </q-btn>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const form = ref({
    type: "PostgreSQL",
    host: "localhost",
    port: 5432,
    user: "postgres",
    password: "",
    database: "dvdrental",
});
</script>

<style scoped>
.bg-radial-gradient {
    background: radial-gradient(circle at center, #0f172a 0%, #070a12 100%);
}
</style>
