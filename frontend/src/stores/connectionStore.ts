import { defineStore } from "pinia";
import { ref, computed } from "vue";

import { DatabaseService } from "@bindings/db-viewer/internal/app";
import type { Connection } from "@bindings/db-viewer/internal/types";
import { DbService } from "@bindings/db-viewer/internal/app";
import { Notify } from "quasar";

export const useConnectionStore = defineStore("connection", () => {
  const selectedConnection = ref<Connection | null>(null);
  const activeConnection = ref<Connection | null>(null);
  const connections = ref<Connection[]>([]);
  const showNewConnectionDialog = ref(false);

  const loadingStates = ref<Record<string, boolean>>({
    getConnections: false,
    saveConnection: false,
  });

  const isConnected = computed(() => activeConnection.value !== null);

  function setSelectedSession(session: Connection) {
    selectedConnection.value = session;
  }

  function clearSelectedSession() {
    selectedConnection.value = null;
  }

  function setActiveSession(session: Connection) {
    activeConnection.value = session;
  }

  function clearActiveSession() {
    activeConnection.value = null;
  }

  async function getConnections() {
    try {
      loadingStates.value.getConnections = true;
      connections.value = (await DatabaseService.GetConnections()) || [];
    } catch (error) {
      console.error(error);
    } finally {
      loadingStates.value.getConnections = false;
    }
  }

  async function pingConnection(): Promise<boolean> {
    if (!activeConnection.value) return false;
  
    try {
      const result = await DbService.PingConfig({
        ID: String(activeConnection.value.id),
        Name: activeConnection.value.name,
        Host: activeConnection.value.host,
        Port: Number(activeConnection.value.port.Int64),
        User: activeConnection.value.user,
        Password: activeConnection.value.password,
        Database: activeConnection.value.dbname,
        Type: activeConnection.value.driver,
  
        SSL: false,
        InMemory: false,
        ReadOnly: false,
      });
  
      console.log("ping result:", result);
  
      return result;
    } catch (error) {
      Notify.create({
        type: "negative",
        message: error instanceof Error ? error.message : String(error),
      });
  
      console.error(error);
      return false;
    }
  }
  return {
    selectedConnection,
    activeConnection,
    isConnected,
    connections,
    loadingStates,
    showNewConnectionDialog,

    setSelectedSession,
    clearSelectedSession,
    setActiveSession,
    clearActiveSession,
    getConnections,
    pingConnection,
  };
});
