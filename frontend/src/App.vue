<script setup lang="ts">
import { ref } from "vue";

import HelloWorld from "./components/HelloWorld.vue";
import Test from "./components/Test.vue";
import StateDatabaseNotSelected from "./layouts/StateDatabaseNotSelected.vue";
import StateDatabaseSelected from "./layouts/StateDatabaseSelected.vue";

const isConnected = ref(true);
</script>

<template>
    <!-- <Test /> -->
    <!-- <HelloWorld /> -->
    <!-- <StateDatabaseSelected
        v-if="isConnected"
        @disconnect="isConnected = false"
    /> -->
    <StateDatabaseNotSelected
    />
        <!-- v-else -->
        <!-- @select-db="handleEngineConnection" -->
</template>
