import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { DbService } from "@bindings/db-viewer/internal/app";
import { Notify } from "quasar";

export const useActiveConnection = defineStore("activeConnection", () => {
  const activeConnectionId = ref<string | null>(null);
  const activeConnection = ref(null);

  async function setActiveConnection() {
    try {
      const activeConnectionId = await DbService.GetActiveConnection();
    } catch (error) {
      Notify.create({ message: "Failed to set active connection", color: "negative" });
    }
  }

  return {
    activeConnectionId,
    activeConnection,
    setActiveConnection,
  }
})