package entities

type ConnectionConfig struct {
	ID   string
	Name string
	Type string

	Host string
	Port int

	User     string
	Password string
	Database string

	SSL bool

	InMemory bool // in case of in-memory database
	ReadOnly bool
	
}
