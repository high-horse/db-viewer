import { defineStore } from "pinia";
import { ref } from "vue";
import { DbService } from "@bindings/db-viewer/internal/app";
import { Notify } from "quasar";

export const useActiveConnection = defineStore("activeConnection", () => {
  const activeConnectionId = ref<string | null>(null);

  async function setActiveConnection() {
    try {
      activeConnectionId.value = await DbService.GetActiveConnection();
    } catch (error) {
      activeConnectionId.value = null;
      Notify.create({ message: "Failed to set active connection", color: "negative" });
    }
  }

  function clear() {
    activeConnectionId.value = null;
  }

  return {
    activeConnectionId,
    setActiveConnection,
    clear,
  };
});