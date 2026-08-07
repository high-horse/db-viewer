<template>
    <div
        class="w-full max-w-md bg-[#111827]/80 backdrop-blur-md p-6 rounded-xl border border-[#1f2937] shadow-2xl"
    >
       
        <!-- Header Section (Flex Container) -->
        <div class="flex items-center gap-3 mb-6">
            <div class="p-2 rounded-lg bg-indigo-500/10 border border-indigo-500/20">
                <q-icon name="add_link" size="24px" class="text-indigo-400" />
            </div>
            <div>
                <h4 class="text-sm font-bold text-white">
                    Setup Database Session
                </h4>
                <p class="text-[11px] text-[#6b7280]">
                    Initialize single engine configuration links dynamically
                </p>
            </div>
        </div>

        <!-- Tabs (Full Width) -->
        <q-tabs
            v-model="form.type"
            dense
            active-color="indigo"
            indicator-color="indigo"
            align="left"
            class="w-full"
        >
            <q-tab name="postgres" label="PostgreSQL" />
            <q-tab name="mysql" label="MySQL" />
            <q-tab name="sqlite" label="SQLite" />
        </q-tabs>

        <!-- Form (Full Width) -->
        <q-form ref="formRef" @submit="connect" class="q-px-none q-pt-md">
            <!-- PostgreSQL / MySQL fields -->
            <template v-if="form.type !== 'sqlite'">
                <q-input
                    v-model="form.host"
                    label="Host"
                    class="w-full"
                    :rules="[required]"
                    outlined
                    dense
                />

                <q-input
                    v-model="form.port"
                    label="Port"
                    type="number"
                    class="w-full"
                    :rules="[required]"
                    outlined
                    dense
                />

                <q-input
                    v-model="form.user"
                    label="User"
                    class="w-full"
                    :rules="[required]"
                    outlined
                    dense
                />

                <q-input
                    v-model="form.password"
                    label="Password"
                    type="password"
                    class="w-full"
                    :rules="[required]"
                    outlined
                    dense
                />

                <q-input
                    v-model="form.database"
                    label="Database"
                    class="w-full"
                    :rules="[required]"
                    outlined
                    dense
                />
            </template>

            <!-- SQLite fields -->
            <template v-else>
                <q-input
                    v-model="form.database"
                    label="Database File"
                    hint="/path/to/database.db"
                    class="w-full"
                    :rules="[required]"
                    outlined
                    dense
                >
                    <template #append>
                        <q-icon name="folder" />
                    </template>
                </q-input>
            </template>

            <q-btn unelevated color="indigo" class="w-full text-capitalize q-mt-md" type="submit" rounded >
                Ping Connection
            </q-btn>
        </q-form>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { QForm } from "quasar";
import { useConnectionStore } from "@/stores/connectionStore";

const store = useConnectionStore();

const formRef = ref<QForm | null>(null);

const form = ref({
  type: "postgres",
    host: "",
    port: "",
    user: "",
    password: "",
    database: "",
});

const required = (value: string | number | null) =>
    !!value || "This field is required";

async function connect() {
    const valid = await formRef.value?.validate();

    if (!valid) {
        return;
    }

    console.log("connecting", form.value);

    // call your DbService.Connect here
}

onMounted(() => {
    store.getConnections();
});
</script>
