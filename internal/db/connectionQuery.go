package db

import "db-viewer/internal/types"

func GetConnectionList() ([]types.Connection, error) {

	query := "SELECT * from connections"
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []types.Connection
	for rows.Next() {
		var connection types.Connection
		if err := rows.Scan(&connection.Id, &connection.Name, &connection.Driver, &connection.Host, &connection.Port, &connection.User, &connection.Password, &connection.DBName); err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}
	return connections, nil
}

func StoreConnection(conn types.Connection) error {
	return nil
}