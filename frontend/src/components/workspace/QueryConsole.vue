<template>
    <div class="h-full flex flex-col bg-[#0c0b09]">
        <!-- Tab Strip -->
        <div
            class="h-9 bg-[#161310] border-b border-[#292521] flex items-center justify-between px-2"
        >
            <div class="flex items-center h-full">
                <div
                    class="h-full px-4 border-r border-[#292521] bg-[#0c0b09] text-xs font-mono text-amber-400 flex items-center gap-2 border-b-2 border-b-amber-500"
                >
                    <q-icon name="code" size="12px" />

                    console_1.sql
                </div>
            </div>

            <q-btn
                unelevated
                color="amber"
                text-color="black"
                size="sm"
                class="font-bold px-3 tracking-wide"
                :loading="executing"
                @click="executeQuery"
            >
                <q-icon name="play_arrow" size="16px" class="mr-1" />

                RUN
            </q-btn>
        </div>

        <!-- SQL Editor -->
        <textarea
            v-model="model"
            class="grow w-full p-4 bg-[#0c0b09] text-amber-100/90 font-mono text-xs border-none resize-none focus:outline-none leading-relaxed tracking-wide"
            spellcheck="false"
            placeholder="Write your SQL query..."
        ></textarea>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";

const props = defineProps<{
    modelValue: string;
}>();

const emit = defineEmits<{
    "update:modelValue": [value: string];
    execute: [query: string];
}>();

const executing = ref(false);

const model = computed({
    get: () => props.modelValue,
    set: (value: string) => {
        emit("update:modelValue", value);
    },
});

async function executeQuery() {
    if (!model.value.trim()) {
        return;
    }

    executing.value = true;

    try {
        emit("execute", model.value);
    } finally {
        executing.value = false;
    }
}
</script>
