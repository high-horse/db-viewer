package main

import (
	"database/sql"
	"db-viewer/internal/types"
	"strconv"

	tea "charm.land/bubbletea/v2"
)

func (m AppStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case LoadedConnectionsMsg:
		if msg.Error != nil {
			m.err = msg.Error
			m.state = ErrorState
			return m, nil
		}
		m.state = ConnList
		m.connections = msg.Connections
		m.cursor = 0
		return m, nil

	case ConnectionFormSubmitMsg:
		portInt, _ := strconv.Atoi(msg.Port)
		m.connections = append(m.connections, types.Connection{
			Name:   msg.Name,
			Driver: msg.Driver,
			Host:   msg.Host,
			Port: sql.NullInt64{
				Int64: int64(portInt),
				Valid: true,
			},
			User:     msg.User,
			Password: msg.Password,
			DBName:   msg.DBName,
		})

		m.state = ConnList
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		case "esc":
			if(m.state > ConnList) {
				m.state = ConnList
			}
		}
		

	}

	switch m.state {
	case Auth:
		{
			// var cmd tea.Cmd
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "tab":
					if m.email.Focused() {
						m.email.Blur()
						m.password.Focus()
					} else {
						m.password.Blur()
						m.email.Focus()
					}

				case "enter":
					ok, msg := validatePassword(m.email.Value(), m.password.Value())
					if !ok {
						m.authErrMsg = msg
						return m, nil
					} else {
						m.state = ConnList
					}

					m.state = ConnList
					return m, loadConnectionsCmd()
				}

				m.email, _ = m.email.Update(msg)
				m.password, _ = m.password.Update(msg)

				return m, nil

			}

		}

	case ConnList:
		{
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {

				case "up":
					if m.cursor > 0 {
						m.cursor--
					}

				case "down":
					if m.cursor < len(m.connections)-1 {
						m.cursor++
					}

				case "n", "N":
					m.state = ConnForm

				case "enter":
					m.currentConn = &m.connections[m.cursor]
					m.state = SchemaBrowser
				}
			}
		}
	case ConnForm:
		var cmd tea.Cmd
		m.connForm, cmd = m.connForm.Update(msg)
		return m, cmd

	}

	return m, nil
}

func validatePassword(uname, passwd string) (bool, string) {
	if len(uname) == 0 || len(passwd) == 0 {
		return false, "username and password cannot be empty"
	}

	if uname != "root" || passwd != "root" {
		return false, "invalid credetials"
	}

	return true, ""
}
