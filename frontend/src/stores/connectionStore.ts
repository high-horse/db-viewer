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
  
    const connection = activeConnection.value;
  
    try {

      const result = await DbService.PingConfig({
          ID: String(connection.id),
          Name: connection.name,
          Host: connection.host,
          Port: Number(connection.port.Int64),
          User: connection.user,
          Password: connection.password,
          Database: connection.dbname,
          Type: connection.driver,
      
          SSL: false,
      
          SSHConfigID: connection.ssh_config_id?.Valid
              ? Number(connection.ssh_config_id.Int64)
              : null,
      
          SSHConfig: connection.ssh_config?.id
              ? {
                    ID: connection.ssh_config.id,
                    Name: connection.ssh_config.name,
                    Host: connection.ssh_config.host,
                    Port: Number(connection.ssh_config.port),
                    Username: connection.ssh_config.username,
                    AuthMethod: connection.ssh_config.auth_method,
      
                    PrivateKey:
                        connection.ssh_config.private_key.Valid
                            ? connection.ssh_config.private_key.String
                            : "",
      
                    Passphrase:
                        connection.ssh_config.passphrase.Valid
                            ? connection.ssh_config.passphrase.String
                            : "",
      
                    Password:
                        connection.ssh_config.password.Valid
                            ? connection.ssh_config.password.String
                            : "",
                }
              : null,
      
          InMemory: false,
          ReadOnly: false,
      
          // ADD THIS
          Color: connection.color.Valid
              ? connection.color.String
              : "",
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
