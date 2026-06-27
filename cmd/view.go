package main

import (
	// "db-viewer/internal/helpers"
	"fmt"

	tea "charm.land/bubbletea/v2"

	lipgloss "charm.land/lipgloss/v2"
)

func (m AppStateModel) View() tea.View {

	switch m.state {

	case Auth:
		return tea.NewView(renderAuth(m))

	case ConnList:
		return tea.NewView(renderConnections(m))

	case ConnForm:
		return tea.NewView(renderConnectionForm(m))

	case SchemaBrowser:
		return tea.NewView(renderSchema(m))

	case TableView:
		return tea.NewView(renderQuery(m))

	case ErrorState:
		return tea.NewView(renderError(m))

	default:
		return tea.NewView("loading...")
	}
}

func renderAuth(m AppStateModel) string {
	fullWidthTitle := sectionTitle.Width(m.width - 10).Render("Login")
	s :=  fullWidthTitle + "\n\n" +
		m.email.View() + "\n" +
		m.password.View() + "\n\n" +
		textDanger.Render(m.authErrMsg) + "\n\n" +
		"Tab to switch • Enter to login"

	return box.Width(m.width - 4).Height(m.height - 4).Render(s)
}

func renderConnections(m AppStateModel) string {
	s := "Connections\n\n"

	for i, c := range m.connections {

		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		line := fmt.Sprintf("%s %s (%s)", cursor, c.Name, c.Driver)

		s += lipgloss.NewStyle().
			PaddingLeft(1).
			Render(line) + "\n"
	}

	s += "\n\n(press 'n' to add new connection)"

	return box.Width(m.width - 4).Height(m.height - 4).Render(s)
}

func renderConnectionForm(m AppStateModel) string {
	f := m.connForm

	s := "New Connection\n\n"

	// %-10s left-aligns each label in a 10-character field
	s += fmt.Sprintf("%-10s", "Alias Name") + f.name.View() + "\n"
	s += fmt.Sprintf("%-10s", "Driver") + f.driver.View() + "\n"
	s += fmt.Sprintf("%-10s", "Host") + f.host.View() + "\n"
	s += fmt.Sprintf("%-10s", "Port") + f.port.View() + "\n"
	s += fmt.Sprintf("%-10s", "UserName") + f.user.View() + "\n"
	s += fmt.Sprintf("%-10s", "Password") + f.password.View() + "\n"
	s += fmt.Sprintf("%-10s", "Database") + f.dbname.View() + "\n"

	s += "\nTab: next • Shift+Tab: previous • Enter: save"

	return box.Width(m.width - 4).Height(m.height - 4).Render(s)
}

func renderSchema(m AppStateModel) string {
	s := "Tables\n\n"

	for i, t := range m.tables {

		cursor := " "
		if i == m.tableCursor {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, t)
	}

	return box.Width(m.width - 4).Height(m.height - 4).Render(s)
}

func renderQuery(m AppStateModel) string {
	return box.Render(
		"SQL Editor\n\n" +
			m.queryInput +
			"\n\nResult:\n" +
			m.queryResult +
			"\n\nPress ENTER to run",
	)
}

func renderError(m AppStateModel) string {
	if m.err == nil {
		return "no error"
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true).
		Render(m.err.Error())
}
