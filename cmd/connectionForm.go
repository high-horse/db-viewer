package main

import (
	textinput "charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

var drivers = []string{
	"postgres",
	"mysql",
	"sqlite",
	"mariadb",
	"sqlserver",
	"oracle",
}
var width = 100
type ConnectionForm struct {
	focus int

	name     textinput.Model
	driver   textinput.Model
	host     textinput.Model
	port     textinput.Model
	user     textinput.Model
	password textinput.Model
	dbname   textinput.Model
}

func newConnectionForm() ConnectionForm {
	name := textinput.New()
	name.Placeholder = "Alias Name"
	name.SetWidth(width)

	driver := textinput.New()
	driver.Placeholder = "Driver (postgres/mysql/sqlite)"
	driver.SetWidth(width)

	host := textinput.New()
	host.Placeholder = "Host"
	host.SetWidth(width)

	port := textinput.New()
	port.Placeholder = "Port"
	port.SetWidth(width)

	user := textinput.New()
	user.Placeholder = "User"
	user.SetWidth(width)

	password := textinput.New()
	password.Placeholder = "Password"
	password.EchoMode = textinput.EchoPassword
	password.SetWidth(width)

	dbname := textinput.New()
	dbname.Placeholder = "Database Name"
	dbname.SetWidth(width)

	form := ConnectionForm{
		name:     name,
		driver:   driver,
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}

	form.focus = 0
	form.applyFocus()

	return form
}

func (f ConnectionForm) Update(msg tea.Msg) (ConnectionForm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			f.focus = (f.focus + 1) % 7
			f.applyFocus()

		case "shift+tab":
			f.focus--
			if f.focus < 0 {
				f.focus = 6
			}
			f.applyFocus()

		case "enter":
			return f, func() tea.Msg {
				return ConnectionFormSubmitMsg{
					Name:     f.name.Value(),
					Driver:   f.driver.Value(),
					Host:     f.host.Value(),
					Port:     f.port.Value(),
					User:     f.user.Value(),
					Password: f.password.Value(),
					DBName:   f.dbname.Value(),
				}
			}
		}

	}

	var cmd tea.Cmd

	f.name, cmd = f.name.Update(msg)
	f.driver, _ = f.driver.Update(msg)
	f.host, _ = f.host.Update(msg)
	f.port, _ = f.port.Update(msg)
	f.user, _ = f.user.Update(msg)
	f.password, _ = f.password.Update(msg)
	f.dbname, _ = f.dbname.Update(msg)

	return f, cmd

}

func (f *ConnectionForm) applyFocus() {
	f.name.Blur()
	f.driver.Blur()
	f.host.Blur()
	f.port.Blur()
	f.user.Blur()
	f.password.Blur()
	f.dbname.Blur()

	switch f.focus {
	case 0:
		f.name.Focus()
	case 1:
		f.driver.Focus()
	case 2:
		f.host.Focus()
	case 3:
		f.port.Focus()
	case 4:
		f.user.Focus()
	case 5:
		f.password.Focus()
	case 6:
		f.dbname.Focus()
	}
}
