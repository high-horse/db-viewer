<template>
    <div
        class="h-full bg-[#100e0c] border-r border-[#292521] flex flex-col font-sans"
    >
        <!-- Sidebar Header -->
        <div
            class="p-3 border-b border-[#292521] flex justify-between items-center bg-[#161310]"
        >
            <span
                class="text-[11px] font-bold uppercase tracking-wider text-[#6b7280]"
            >
                Schema Explorer
            </span>

            <q-btn
                flat
                dense
                round
                icon="refresh"
                size="xs"
                class="text-[#6b7280] hover:text-amber-400"
                @click="refreshSchema"
            />
        </div>

        <!-- Schema Tree -->
        <q-scroll-area class="grow p-2">
            <q-tree
                :nodes="schemaNodes"
                node-key="id"
                dark
                no-connectors
                class="text-xs"
            >
                <template #default-header="prop">
                    <div
                        class="flex items-center gap-2 py-1 px-1 rounded cursor-pointer group w-full hover:bg-[#231f1a]"
                    >
                        <q-icon
                            :name="prop.node.icon"
                            :color="prop.node.iconColor"
                            size="14px"
                        />

                        <span
                            :class="
                                prop.node.type === 'table'
                                    ? 'text-white font-medium'
                                    : 'text-[#94a3b8]'
                            "
                        >
                            {{ prop.node.label }}
                        </span>
                    </div>
                </template>
            </q-tree>
        </q-scroll-area>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const schemaNodes = ref([
    {
        id: "db-root",
        label: "dvdrental",
        icon: "storage",
        iconColor: "amber-5",
        children: [
            {
                id: "t1",
                label: "actor",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
            {
                id: "t2",
                label: "customer",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
            {
                id: "t3",
                label: "film",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
            {
                id: "t4",
                label: "inventory",
                icon: "table_chart",
                iconColor: "amber-4",
                type: "table",
            },
        ],
    },
]);

function refreshSchema() {
    // TODO:
    // Load schema from backend here.
    console.log("Refreshing schema...");
}
</script>
