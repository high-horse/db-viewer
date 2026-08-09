<template>
    <header
        class="h-12 bg-[#161310] border-b border-[#292521] flex items-center justify-between px-4 z-10"
    >
        <div class="flex items-center gap-3">
            <q-icon name="dns" size="20px" class="text-amber-400" />

            <span class="font-bold text-white tracking-wide text-sm">
                SQL Client Pro
            </span>

            <q-badge
                color="amber-10"
                text-color="amber-4"
                label="Connected"
                class="text-xs font-medium border border-amber-9"
            />
        </div>

        <div class="flex items-center gap-2">
            <q-btn
                flat
                dense
                round
                icon="tune"
                size="sm"
                class="text-[#6b7280] hover:text-white"
            />

            <q-btn
                flat
                dense
                round
                icon="logout"
                size="sm"
                class="text-amber-400 hover:bg-amber-500/10"
                @click="$emit('disconnect')"
            />
        </div>
    </header>
</template>

<script setup lang="ts">
defineEmits<{
    disconnect: [];
}>();
</script>
