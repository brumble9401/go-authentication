package querybuilder

import (
	"strings"

	"github.com/gocql/gocql"
)

type ScyllaQueryBuilder struct {
    session *gocql.Session
}

func NewScyllaQueryBuilder(session *gocql.Session) *ScyllaQueryBuilder {
    return &ScyllaQueryBuilder{session: session}
}

func (qb *ScyllaQueryBuilder) InsertQuery(table string, data map[string]interface{}) *gocql.Query {
    columns := make([]string, 0, len(data))
    values := make([]interface{}, 0, len(data))
    placeholders := make([]string, 0, len(data))

    for col, val := range data {
        columns = append(columns, col)
        values = append(values, val)
        placeholders = append(placeholders, "?")
    }

    query := qb.session.Query(
        "INSERT INTO "+table+" ("+strings.Join(columns, ", ")+") VALUES ("+strings.Join(placeholders, ", ")+")",
        values...,
    )
    return query
}

func (qb *ScyllaQueryBuilder) DeleteQuery(table string, id gocql.UUID) *gocql.Query {
    query := qb.session.Query(
        "DELETE FROM "+table+" WHERE id = ?",
        id,
    )
    return query
}

func (qb *ScyllaQueryBuilder) SelectQuery(table string, id gocql.UUID) *gocql.Query {
    query := qb.session.Query(
        "SELECT * FROM "+ table +" WHERE user_id = ?",
        id,
    )
    return query
}

func (qb *ScyllaQueryBuilder) SelectAllQuery(table string) *gocql.Query {
    query := qb.session.Query(
        "SELECT * FROM " + table,
    )
    return query
}

func (qb *ScyllaQueryBuilder) SelectConditionQuery(table string, column string, value string) *gocql.Query {
	query := qb.session.Query(
        "SELECT * FROM "+table+" WHERE "+column+" = ? ALLOW FILTERING",
        value,
    )
    return query
}

func (qb *ScyllaQueryBuilder) UpdateQuery(table string, column string, id gocql.UUID, data map[string]interface{}) *gocql.Query {
    setClauses := make([]string, 0, len(data))
    values := make([]interface{}, 0, len(data)+1)

    for col, val := range data {
        setClauses = append(setClauses, col+" = ?")
        values = append(values, val)
    }
    values = append(values, id)

    query := qb.session.Query(
        "UPDATE "+table+" SET "+strings.Join(setClauses, ", ")+" WHERE " + column + " = ?",
        values...,
    )
    return query
}