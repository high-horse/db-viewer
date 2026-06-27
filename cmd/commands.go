package main

import (
	"db-viewer/internal/db"
	"db-viewer/internal/types"

	tea "charm.land/bubbletea/v2"
)

type LoadedConnectionsMsg struct {
	Connections []types.Connection
	Error       error
}



func loadConnectionsCmd() tea.Cmd {
	return func() tea.Msg {
		cons, err := db.GetConnectionList()
		return LoadedConnectionsMsg {
			Connections: cons,
			Error: err,
		}
	}
}

func saveConnectionCmd(conn ConnectionFormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		connection := types.Connection{
			Name: conn.Name,
			Driver: conn.Driver,
			Host: conn.Host,
		}
		db.StoreConnection(connection)
		return loadConnectionsCmd()
	}
}