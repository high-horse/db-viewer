<template>
    <div
        class="w-full max-w-xl bg-[#161310] backdrop-blur-md p-6 rounded-xl border shadow-2xl"
    >
        <!-- Header Section (Flex Container) -->
        <!-- <div class="flex items-center gap-3 mb-6">
            <div>
                <h4 class="text-sm font-bold text-white">
                    Setup Database Session
                </h4>
                <p class="text-[11px] text-[#8a8478]">
                    Initialize single engine configuration links dynamically
                </p>
            </div>
        </div> -->

        <!-- Tabs (Full Width) -->
        <q-tabs
            v-model="form.type"
            dense
            active-color="amber"
            indicator-color="amber"
            align="left"
            class="w-full"
        >
            <q-tab name="pgx" label="PostgreSQL" />
            <q-tab name="mysql" label="MySQL" />
            <q-tab name="sqlite" label="SQLite" />
        </q-tabs>

        <!-- Form (Full Width) -->
        <q-form
            ref="formRef"
            @submit="connect"
            class="flex-1 flex flex-col min-h-0 px-6 pb-6 q-pt-md"
        >
            <div class="flex-1 min-h-0">
                <!-- <q-scroll-area style="height: 100%; min-width: 300px"> -->
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
                            class="w-full q-pb-md"
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

                        <q-input
                            v-model="form.connection_name"
                            label="Connection Name"
                            class="w-full"
                            :rules="[required]"
                            outlined
                            dense
                        />

                        <div class="q-mt-md">
                            <div class="text-caption text-grey-5 q-mb-sm">
                                Connection Color
                            </div>

                            <div class="flex items-center gap-2">
                                <button
                                    v-for="color in colors"
                                    :key="color"
                                    type="button"
                                    class="w-7 h-7 rounded-md flex items-center justify-center border-2 border-transparent cursor-pointer transition-transform duration-150 hover:scale-110"
                                    :class="
                                        form.color === color
                                            ? 'border-white ring-2 ring-white/40'
                                            : ''
                                    "
                                    :style="{ backgroundColor: color }"
                                    @click="form.color = color"
                                >
                                    <q-icon
                                        v-if="form.color === color"
                                        name="check"
                                        size="16px"
                                        color="white"
                                    />
                                </button>

                                <span class="text-xs text-grey-5 ml-2">
                                    {{ form.color }}
                                </span>
                            </div>
                        </div>

                        <q-checkbox
                            v-model="form.save_connection"
                            label="Save connection"
                            class="w-1/2 q-pt-md"
                            :rules="[required]"
                            outlined
                            dense
                        />
                        <q-checkbox
                            v-model="form.readonly"
                            label="Readonly"
                            class="w-1/2 q-pt-md"
                            :rules="[required]"
                            outlined
                            dense
                        />
                        <q-checkbox
                            v-model="form.use_ssh"
                            label="Connect through SSH tunnel"
                            class="w-full q-pt-md"
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
                <!-- </q-scroll-area> -->
            </div>
            <q-btn
                unelevated
                color="amber"
                class="w-full text-capitalize q-mt-md text-black"
                type="submit"
                rounded
            >
                Test Connection
            </q-btn>
        </q-form>

        <q-dialog v-model="sshDialog" persistent>
            <q-card class="w-full max-w-lg bg-[#161310] text-white">
                <q-card-section class="flex items-center justify-between">
                    <div>
                        <div class="text-base font-bold">SSH Configuration</div>

                        <div class="text-xs text-grey-5 q-mt-xs">
                            Configure the SSH tunnel for this connection
                        </div>
                    </div>

                    <q-btn
                        flat
                        round
                        dense
                        icon="close"
                        color="grey-5"
                        @click="cancelSsh"
                    />
                </q-card-section>

                <q-separator dark />

                <q-card-section class="q-gutter-md">
                    <q-input
                        v-model="sshForm.name"
                        label="Configuration Name"
                        outlined
                        dense
                        color="amber"
                        :rules="[required]"
                    />

                    <div class="row q-col-gutter-md">
                        <div class="col-8">
                            <q-input
                                v-model="sshForm.host"
                                label="SSH Host"
                                outlined
                                dense
                                color="amber"
                                :rules="[required]"
                            />
                        </div>

                        <div class="col-4">
                            <q-input
                                v-model.number="sshForm.port"
                                label="Port"
                                type="number"
                                outlined
                                dense
                                color="amber"
                                :rules="[required]"
                            />
                        </div>
                    </div>

                    <q-input
                        v-model="sshForm.username"
                        label="Username"
                        outlined
                        dense
                        color="amber"
                        :rules="[required]"
                    />

                    <q-select
                        v-model="sshForm.auth_method"
                        label="Authentication Method"
                        outlined
                        dense
                        color="amber"
                        emit-value
                        map-options
                        :options="[
                            {
                                label: 'Password',
                                value: 'password',
                            },
                            {
                                label: 'Private Key',
                                value: 'private_key',
                            },
                        ]"
                    />

                    <!-- Password authentication -->
                    <q-input
                        v-if="sshForm.auth_method === 'password'"
                        v-model="sshForm.password"
                        label="SSH Password"
                        type="password"
                        outlined
                        dense
                        color="amber"
                    />

                    <!-- Private key authentication -->
                    <template v-if="sshForm.auth_method === 'private_key'">
                        <q-input
                            v-model="sshForm.private_key"
                            label="Private Key"
                            type="textarea"
                            outlined
                            dense
                            color="amber"
                            autogrow
                            hint="Paste your SSH private key"
                        />

                        <q-input
                            v-model="sshForm.passphrase"
                            label="Key Passphrase"
                            type="password"
                            outlined
                            dense
                            color="amber"
                        />
                    </template>
                </q-card-section>

                <q-separator dark />

                <q-card-actions align="right" class="q-pa-md">
                    <q-btn
                        flat
                        label="Cancel"
                        color="grey-5"
                        @click="cancelSsh"
                    />

                    <q-btn
                        unelevated
                        label="Save SSH Configuration"
                        color="amber"
                        text-color="black"
                        @click="saveSshConfig"
                    />
                </q-card-actions>
            </q-card>
        </q-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, hydrateOnIdle } from "vue";
import type { QForm } from "quasar";
import { useConnectionStore } from "@/stores/connectionStore";
import { DbService } from "@bindings/db-viewer/internal/app";
import type { ConnectionConfig } from "@bindings/db-viewer/internal/engine/entities";
import { useQuasar, Dialog, Notify } from "quasar";
import { useActiveConnection } from "@/stores/activeConnection";
import { useRouter } from "vue-router";

const $q = useQuasar();
const $router = useRouter();
const store = useConnectionStore();
const activeConnectionStore = useActiveConnection();

const formRef = ref<QForm | null>(null);

const form = ref({
    type: "pgx",
    host: "",
    port: 5432,
    user: "",
    password: "",
    database: "",
    connection_name: "",
    color: "",
    save_connection: true,
    readonly: false,
    use_ssh: false,
});
const sshForm = ref({
    name: "",
    host: "",
    port: 22,
    username: "",
    auth_method: "password",
    private_key: "",
    passphrase: "",
    password: "",
});

const sshDialog = ref(false);

const colors = [
    "#EF4444", // Red
    "#F97316", // Orange
    "#EAB308", // Yellow
    "#84CC16", // Lime
    "#22C55E", // Green
    "#14B8A6", // Teal
    "#06B6D4", // Cyan
    "#3B82F6", // Blue
    "#6366F1", // Indigo
    "#8B5CF6", // Violet
    "#A855F7", // Purple
    "#EC4899", // Pink
];

const required = (value: string | number | null) =>
    !!value || "This field is required";

async function connect() {
    const valid = await formRef.value?.validate();

    if (!valid) {
        return;
    }

    try {
        const response = await DbService.PingConfig(parseToConfig());
        if (response) {
            Dialog.create({
                title: "Ping Successful",
                message: "Connect to the database?",
                ok: "OK",
                cancel: "Cancel",
            }).onOk(async () => {
                await DbService.SaveAndConnect(parseToConfig());
                activeConnectionStore.setActiveConnection();
              // $router.push({ name: "workspace" });
                $router.push({ name: "WorkSpace" });
            });
        }
    } catch (error: any) {
        Dialog.create({
            message: error.message || "Failed to ping config",
            color: "negative",
        });
        console.error(error);
    }
}

function parseToConfig() {
    return <ConnectionConfig>{
        Name: form.value.connection_name,
        Type: form.value.type,

        Host: form.value.host,
        Port: form.value.port,

        User: form.value.user,
        Password: form.value.password,
        Database: form.value.database,

        SSL: false,

        SSHConfigID: null,

        SSHConfig: form.value.use_ssh
            ? {
                  Name: sshForm.value.name,
                  Host: sshForm.value.host,
                  Port: sshForm.value.port,
                  Username: sshForm.value.username,
                  AuthMethod: sshForm.value.auth_method,
                  PrivateKey: sshForm.value.private_key,
                  Passphrase: sshForm.value.passphrase,
                  Password: sshForm.value.password,
              }
            : null,

        InMemory: false,
        ReadOnly: form.value.readonly,
        Color: form.value.color,
    };
}

function saveSshConfig() {
    sshDialog.value = false;
}

function cancelSsh() {
    sshDialog.value = false;
    form.value.use_ssh = false;
    resetSshForm();
}

function resetSshForm() {
    sshForm.value = {
        name: "",
        host: "",
        port: 22,
        username: "",
        auth_method: "password",
        private_key: "",
        passphrase: "",
        password: "",
    };
}
watch(
    () => form.value.type,
    (newVal) => {
        if (newVal === "pgx") {
            form.value.port = 5432;
        } else if (newVal === "mysql") {
            form.value.port = 3306;
        } else if (newVal === "sqlite") {
            form.value.port = 0;
        }
    },
);

watch(
    () => form.value.use_ssh,
    (enabled) => {
        if (enabled) {
            sshDialog.value = true;
        } else {
            resetSshForm();
        }
    },
);

onMounted(() => {
    store.getConnections();
});
</script>
