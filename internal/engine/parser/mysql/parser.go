package mysqlQueryParser

import "db-viewer/internal/engine/entities"

type Parser struct {}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) Parse(sql string) (*entities.SQLQueryEntity, error) {
    // stmt, err := mysqlParser.Parse(sql)
    // if err != nil {
    //     return nil, err
    // }

    return &entities.SQLQueryEntity{
        RawSQL:        sql,
        Dialect:       entities.DialectMySQL,
        // StatementType: determineStatementType(stmt),
    }, nil
}