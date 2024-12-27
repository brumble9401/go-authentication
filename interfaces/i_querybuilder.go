package interfaces

import (
	"github.com/gocql/gocql"
)

type QueryBuilder interface {
    NewBatch(batchType gocql.BatchType) *gocql.Batch
	InsertQuery(table string, data map[string]interface{}) *gocql.Query
	InsertToBatch(batch *gocql.Batch, table string, data map[string]interface{})
	DeleteQuery(table string, id gocql.UUID) *gocql.Query
	AddDeleteToBatch(batch *gocql.Batch, table string, column string, id gocql.UUID)
	SelectQuery(table string, id gocql.UUID) *gocql.Query
	SelectAllQuery(table string) *gocql.Query
	SelectConditionQuery(table string, column string, value string) *gocql.Query
	UpdateQuery(table string, column string, id gocql.UUID, data map[string]interface{}) *gocql.Query
}
