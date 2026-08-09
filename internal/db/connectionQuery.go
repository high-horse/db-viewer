package db

import (
	"database/sql"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/types"
)


func GetConnectionList() ([]types.Connection, error) {
	query := `
		SELECT
			c.id, c.name, c.driver, c.host, c.port, c.user, c.password, c.dbname, c.pinned, c.color,
			s.id, s.name, s.host, s.port, s.username, s.auth_method, s.private_key, s.passphrase, s.password
		FROM connections c
		LEFT JOIN ssh_configs s ON s.id = c.ssh_config_id
	`
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []types.Connection
	for rows.Next() {
		var c types.Connection
		var sshID sql.NullInt64
		var sshName, sshHost, sshAuthMethod sql.NullString
		var sshPort sql.NullInt64
		var sshUsername sql.NullString

		if err := rows.Scan(
			&c.Id, &c.Name, &c.Driver, &c.Host, &c.Port, &c.User, &c.Password, &c.DBName, &c.Pinned, &c.Color,
			&sshID, &sshName, &sshHost, &sshPort, &sshUsername, &sshAuthMethod,
			&c.SSHConfig.PrivateKey, &c.SSHConfig.Passphrase, &c.SSHConfig.Password,
		); err != nil {
			return nil, err
		}

		if sshID.Valid {
			c.SSHConfigId = sshID
			c.SSHConfig.Id = int(sshID.Int64)
			c.SSHConfig.Name = sshName.String
			c.SSHConfig.Host = sshHost.String
			c.SSHConfig.Port = int(sshPort.Int64)
			c.SSHConfig.Username = sshUsername.String
			c.SSHConfig.AuthMethod = sshAuthMethod.String
		}

		connections = append(connections, c)
	}
	return connections, nil
}



func StoreConnection(conn entities.ConnectionConfig) (int64, error) {
	tx, err := Conn.Begin()
	if err != nil {
		return 0, err
	}

	var sshConfigID any // nil or int64, goes straight into the query as NULL/value

	if conn.SSHConfig != nil {
		res, err := tx.Exec(
			`INSERT INTO ssh_configs (name, host, port, username, auth_method, private_key, passphrase, password)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			conn.SSHConfig.Name,
			conn.SSHConfig.Host,
			conn.SSHConfig.Port,
			conn.SSHConfig.Username,
			conn.SSHConfig.AuthMethod,
			nullable(conn.SSHConfig.PrivateKey),
			nullable(conn.SSHConfig.Passphrase),
			nullable(conn.SSHConfig.Password),
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		sshConfigID = id
	}

	res, err := tx.Exec(
		`INSERT INTO connections (name, driver, host, port, user, password, dbname, pinned, color, ssh_config_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		conn.Name,
		conn.Type,
		conn.Host,
		conn.Port,
		conn.User,
		conn.Password,
		conn.Database,
		false,
		nullable(conn.Color),
		sshConfigID,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	newID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newID, nil
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}