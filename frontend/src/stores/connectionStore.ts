import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { DatabaseService } from "@bindings/db-viewer/internal/app";
import type { Connection } from "@bindings/db-viewer/internal/types";
import { DbService } from "@bindings/db-viewer/internal/app";
import type { ConnectionConfig } from "@bindings/db-viewer/internal/engine/entities";
import { Notify } from "quasar";

export const useConnectionStore = defineStore("connection", () => {
  const selectedConnection = ref<Connection | null>(null);
  const activeConnection = ref<Connection | null>(null);
  const connections = ref<Connection[]>([]);
  const showNewConnectionDialog = ref(false);
  const loadingStates = ref<Record<string, boolean>>({
    getConnections: false,
    saveConnection: false,
    connecting: false,
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

  function getActiveSession() {
    return activeConnection.value;
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

  // Shared mapper — used by both ping and connect so they can never drift apart.
  function toConfig(connection: Connection): ConnectionConfig {
    return {
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
            PrivateKey: connection.ssh_config.private_key.Valid
              ? connection.ssh_config.private_key.String
              : "",
            Passphrase: connection.ssh_config.passphrase.Valid
              ? connection.ssh_config.passphrase.String
              : "",
            Password: connection.ssh_config.password.Valid
              ? connection.ssh_config.password.String
              : "",
          }
        : null,
      InMemory: false,
      ReadOnly: false,
      Color: connection.color.Valid ? connection.color.String : "",
    };
  }

  async function pingConnection(connection: Connection): Promise<boolean> {
    try {
      return await DbService.PingConfig(toConfig(connection));
    } catch (error) {
      Notify.create({
        type: "negative",
        message: error instanceof Error ? error.message : String(error),
      });
      console.error(error);
      return false;
    }
  }

  // Ping first, then actually register/activate the connection on the backend.
  async function connectToSession(connection: Connection): Promise<boolean> {
    loadingStates.value.connecting = true;
    try {
      const reachable = await pingConnection(connection);
      if (!reachable) {
        Notify.create({ type: "negative", message: "Failed to connect" });
        return false;
      }

      // Connection already has a saved row (has an id) — use Connect, not SaveAndConnect,
      // or you'll insert a duplicate row every time someone double-clicks it.
      const connected = await DbService.Connect(toConfig(connection));
      if (!connected) {
        Notify.create({ type: "negative", message: "Failed to connect" });
        return false;
      }

      setActiveSession(connection);
      Notify.create({ type: "positive", message: "Connection established" });
      return true;
    } catch (error) {
      Notify.create({
        type: "negative",
        message: error instanceof Error ? error.message : String(error),
      });
      console.error(error);
      return false;
    } finally {
      loadingStates.value.connecting = false;
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
    getActiveSession,
    clearActiveSession,
    getConnections,
    pingConnection,
    connectToSession,
  };
});