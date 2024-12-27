package interfaces

import (
	"github.com/gocql/gocql"
)

type QueryBuilder interface {
    InsertQuery(table string, data map[string]interface{}) *gocql.Query
    DeleteQuery(table string, id gocql.UUID) *gocql.Query
    SelectQuery(table string, id gocql.UUID) *gocql.Query
    SelectAllQuery(table string) *gocql.Query
	SelectConditionQuery(table string, column string, value string) *gocql.Query
    UpdateQuery(table string, column string, id gocql.UUID, data map[string]interface{}) *gocql.Query
}