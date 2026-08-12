<template>
    <div
        class="h-full bg-[#100e0c] border-r border-[#292521] flex flex-col font-sans"
    >
        <!-- Sidebar Header -->
        <div
            class="p-3 border-b border-[#292521] flex justify-between items-center bg-[#161310]"
        >
            <div class="flex items-center gap-2">
                <q-icon
                    name="account_tree"
                    size="15px"
                    class="text-amber-400"
                />

                <span
                    class="text-[11px] font-bold uppercase tracking-wider text-[#6b7280]"
                >
                    Schema Explorer
                </span>
            </div>

            <q-btn
                flat
                dense
                round
                icon="refresh"
                size="xs"
                class="text-[#6b7280] hover:text-amber-400"
                :loading="loading"
                @click="refreshSchema"
            />
        </div>

        <!-- Connection Information -->
        <div class="px-3 py-2 border-b border-[#292521] bg-[#100e0c]">
            <div class="flex items-center gap-2">
                <q-icon name="storage" size="15px" class="text-amber-400" />

                <div class="min-w-0">
                    <div class="text-xs font-semibold text-white truncate">
                        dvdrental
                    </div>

                    <div class="text-[10px] text-[#6b7280] truncate">
                        PostgreSQL
                    </div>
                </div>
            </div>
        </div>

        <!-- Schema Tree -->
        <q-scroll-area class="flex-grow p-2">
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
                        @dblclick="handleNodeDoubleClick(prop.node)"
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

        <!-- Sidebar Footer -->
        <div
            class="h-7 px-3 flex items-center border-t border-[#292521] bg-[#161310]"
        >
            <span class="text-[10px] text-[#4b4540] font-mono">
                {{ tableCount }} tables
            </span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";

interface SchemaNode {
    id: string;
    label: string;
    icon: string;
    iconColor: string;
    type?: string;
    children?: SchemaNode[];
}

const emit = defineEmits<{
    refresh: [];
    selectTable: [node: SchemaNode];
}>();

const loading = ref(false);

const schemaNodes = ref<SchemaNode[]>([
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

const tableCount = computed(() => {
    return schemaNodes.value.reduce((count, node) => {
        return (
            count +
            (node.children?.filter((child) => child.type === "table").length ??
                0)
        );
    }, 0);
});

async function refreshSchema() {
    loading.value = true;

    try {
        // TODO:
        // Load schema from DbService here.

        await new Promise((resolve) => setTimeout(resolve, 300));

        emit("refresh");
    } finally {
        loading.value = false;
    }
}

function handleNodeDoubleClick(node: SchemaNode) {
    if (node.type === "table") {
        emit("selectTable", node);
    }
}
</script>
