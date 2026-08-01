<script setup lang="ts">
import { useDarkMode } from "../composables/useDarkMode";

const { dark, toggleDark } = useDarkMode();
</script>

<template>
  <!-- 
    FIX: Removed plain Tailwind bg utility from here. 
    Quasar naturally manages body backgrounds when Dark.set(true) is called.
    If you want your exact slate color, target the body element in your style.css instead.
  -->
  <div class="min-h-screen flex items-center justify-center p-8">
    
    <!-- 
      FIX: Removed hardcoded text-gray-900 and dark:text-white.
      Quasar natively targets q-card text scaling when Dark mode is active.
    -->
    <q-card class="w-full max-w-md rounded-xl shadow-xl">

      <q-card-section>
        <h1 class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">
          Wails 3 + Quasar + Tailwind
        </h1>

        <p class="mt-3 text-gray-600 dark:text-gray-300">
          Dark mode is working.
        </p>

        <div class="mt-6 flex gap-3">
          <q-btn color="primary" label="Quasar Button" />

          <button @click="toggleDark" class="
              px-4
              py-2
              rounded-lg
              bg-emerald-500
              text-white
              hover:bg-emerald-600
              dark:bg-emerald-600
              dark:hover:bg-emerald-700
            ">
            {{ dark ? "Light Mode" : "Dark Mode" }}
          </button>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>
