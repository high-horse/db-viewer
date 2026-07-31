package main

import (
	"context"
	"db-viewer/internal/db"
	manager "db-viewer/internal/engine/connectionManager"
	"db-viewer/internal/engine/drivers/mysql"
	"db-viewer/internal/engine/drivers/postgres"
	"db-viewer/internal/engine/drivers/sqlite"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/factory"
	"db-viewer/internal/engine/transports"
	"fmt"

	"log"

	tea "charm.land/bubbletea/v2"
)

var driverFactory *factory.Factory

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	driverFactory = factory.New()

	driverFactory.Register(mysql.NewDriver())
	log.Println("registring driver to factory pgx")

	driverFactory.Register(postgres.NewDriver())
	log.Println("drivers registered")

	driverFactory.Register(sqlite.NewDriver())
	log.Println("registring driver to factory sqlite")
}

func main() {

	log.Println("registring driver to factory mysql")
	runPostgres()
}

func runPostgres() {
	manager := manager.NewConnectionManager()
	config := entities.ConnectionConfig{
		ID:       "postgres_local",
		Name:     "dvdrental",
		Type:     "pgx",
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",      // Added default administrative user
		Password: "your_password", // Replace with your actual postgres user password
		Database: "dvdrental",
		SSL:      false, // pgx translates false to sslmode=disable
		InMemory: false,
		ReadOnly: false,
	}

	transport := transports.NewDirect(config.Host, config.Port)

	conn, err := driverFactory.Create(context.Background(), config, transport)
	if err != nil {
		log.Fatal("error: ", err)
	}
	manager.Add(conn)
	log.Println("conn status", conn.IsConnected(), conn.Name())

	if err := conn.Connect(context.Background()); err != nil {
		log.Fatal("connection failed:", err)
	}
	defer conn.Disconnect()

	log.Println("conn status", conn.IsConnected(), conn.Name())
	driver, err := driverFactory.Driver(conn.Type())
	if err != nil {
		log.Fatal("driver error ", err)
	}
	log.Println("driver", driver.Name())

	databases, err := driver.Inspector().ListDatabases(context.Background(), conn)
	if err != nil {
		log.Fatal("list databases err ", err)
	}
	log.Println("Databases ", databases)

	tables, err := driver.Inspector().ListTables(context.Background(), conn)
	if err != nil {
		log.Fatal("tables list err", err)
	}
	for _, table := range tables {
		log.Println("tables", table.Name)
	}

	columns, err := driver.Inspector().ListColumns(context.Background(), conn, tables[0])
	if err != nil {
		log.Fatal("list columns error:", err)
	}

	log.Println("columns count for table", len(columns), tables[0].Name)
	for _, column := range columns {
		log.Println("columns ", column.Name)
	}

	q := `
		select * from actor where first_name like "%Joe%;
	`

	result, err := driver.Executor().Execute(context.Background(), conn, q)
	if err != nil {
		log.Fatal("err executing", err)
	}
	log.Println("Query executed successfully")
	log.Println("Duration:", result.Duration)
	log.Println("Rows affected:", result.RowsAffected)
}

func runSQLite() {
	manager := manager.NewConnectionManager()
	config := entities.ConnectionConfig{
		ID:       "local",
		Name:     "Local SQLite",
		Type:     "sqlite",
		Host:     "",
		Port:     0,
		User:     "",
		Password: "",
		Database: "/home/camel/Desktop/go/chaarm/db-viewer/static/app.db",
		SSL:      false,
		InMemory: false,
		ReadOnly: false,
	}

	transport := transports.NewDirect(config.Host, config.Port)

	conn, err := driverFactory.Create(context.Background(), config, transport)
	if err != nil {
		log.Fatal("error ", err)
	}
	manager.Add(conn)
	log.Println("conn status", conn.IsConnected(), conn.Name())

	if err := conn.Connect(context.Background()); err != nil {
		log.Fatal("connection failed:", err)
	}
	defer conn.Disconnect()

	log.Println("connection status", conn.IsConnected(), conn.Name())
	driver, err := driverFactory.Driver(conn.Type())
	if err != nil {
		log.Fatal("driver error:", err)
	}
	log.Println("driver", driver.Name())

	databases, err := driver.Inspector().ListDatabases(context.Background(), conn)
	if err != nil {
		log.Fatal("list databases error:", err)
	}
	log.Println("databases", databases)

	tables, err := driver.Inspector().ListTables(context.Background(), conn)
	if err != nil {
		log.Fatal("list tables error:", err)
	}
	for _, table := range tables {
		log.Println("tables", table.Name)
	}

	columns, err := driver.Inspector().ListColumns(context.Background(), conn, tables[0])
	if err != nil {
		log.Fatal("list columns error:", err)
	}
	log.Println("columns count for table", len(columns), tables[0].Name)
	for _, column := range columns {
		log.Println("columns ", column.Name)
	}
}

func runMysql() {

	manager := manager.NewConnectionManager()

	config := entities.ConnectionConfig{
		ID:       "local",
		Name:     "Local MySQL",
		Type:     "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "app_user",
		Password: "strong_password",
		Database: "app_database",
	}
	transport := transports.NewDirect(config.Host, config.Port)

	conn, err := driverFactory.Create(context.TODO(), config, transport)
	if err != nil {
		log.Fatal("error ", err)
	}
	manager.Add(conn)

	if err := conn.Connect(context.Background()); err != nil {
		log.Fatal("connection failed:", err)
	}
	log.Println("conn status", conn.IsConnected(), conn.Name())

	driver, err := driverFactory.Driver(conn.Type())
	if err != nil {
		log.Fatal("executor selection failed:", err)
	}

	exec := driver.Executor()

	query := `
			SELECT
				e.employeeNumber,
				CONCAT(e.firstName, ' ', e.lastName) AS employee_name,
				o.city,
				COALESCE(o.phone, 'No phone') AS office_phone
			FROM employees e
			LEFT JOIN offices o
				ON e.officeCode = o.officeCode
			ORDER BY employee_name;
		`
	result, err := exec.Execute(
		context.Background(),
		conn,
		query,
	)

	if err != nil {
		log.Fatal("query execution failed:", err)
	}

	log.Println("Query executed successfully")
	log.Println("Duration:", result.Duration)
	log.Println("Rows affected:", result.RowsAffected)
	// fmt.Println("Columns:")
	// for _, col := range result.Columns {
	// 	fmt.Println("Name:", col.Name)
	// 	fmt.Println("Type:", col.DatabaseType)
	// }

	// fmt.Println("Rows:")
	// for _, row := range result.Rows {
	// 	for i, value := range row {
	// 		fmt.Printf("%s = %v\n", result.Columns[i].Name, value)
	// 	}
	// }

	inspect := driver.Inspector()

	dbs, err := inspect.ListDatabases(context.Background(), conn)
	if err != nil {
		fmt.Println("err while listing datase")
	}
	fmt.Println("datanases len and names", len(dbs), dbs)

	tables, err := inspect.ListTables(context.Background(), conn)
	if err != nil {
		fmt.Println("error in insptect tables", err)
	}

	fmt.Println("tables count ", len(tables))

	cols, err := inspect.ListColumns(context.Background(), conn, tables[0])
	if err != nil {
		fmt.Println("error in inspect columns", err)
	}
	fmt.Println("columns count", len(cols))
	fmt.Println("columns:", cols)

}

func main_old() {
	_, err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(initAppStateModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
