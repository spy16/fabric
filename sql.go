package fabric

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var _ Store = &SQLStore{}
var _ ReWeighter = &SQLStore{}
var _ Counter = &SQLStore{}

// SQLStore implements Store interface using the Go standard library
// sql package.
type SQLStore struct {
	DB *sql.DB
}

// Count returns the number of triples that match the given query.
func (ss *SQLStore) Count(ctx context.Context, query Query) (int, error) {
	sq := `SELECT count(*) FROM triples`

	where, args, err := getWhereClause(query)
	if err != nil {
		return 0, err
	}
	if where != "" {
		sq += fmt.Sprintf("WHERE %s", where)
	}

	if query.Limit > 0 {
		sq = fmt.Sprintf("%s LIMIT %d", sq, query.Limit)
	}

	var count int64
	row := ss.DB.QueryRowContext(ctx, sq, args...)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return int(count), nil
}

// Insert persists the given triple into the triples table.
func (ss *SQLStore) Insert(ctx context.Context, tri Triple) error {
	query := `INSERT INTO triples (source, predicate, target, weight) VALUES (?, ?, ?, ?)`

	_, err := ss.DB.ExecContext(ctx, query, tri.Source, tri.Predicate, tri.Target, tri.Weight)
	return err
}

// Query converts the given query object into SQL SELECT and fetches all the triples.
func (ss *SQLStore) Query(ctx context.Context, query Query) ([]Triple, error) {
	sq := `SELECT * FROM triples`

	where, args, err := getWhereClause(query)
	if err != nil {
		return nil, err
	}
	if where != "" {
		sq += fmt.Sprintf(" WHERE %s", where)
	}

	if query.Limit > 0 {
		sq = fmt.Sprintf("%s LIMIT %d", sq, query.Limit)
	}

	rows, err := ss.DB.QueryContext(ctx, sq, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	triples := []Triple{}
	for rows.Next() {
		var tri Triple
		if err := rows.Scan(&tri.Source, &tri.Predicate, &tri.Target, &tri.Weight); err != nil {
			return nil, err
		}

		triples = append(triples, tri)
	}

	return triples, nil
}

// Delete removes all the triples from the database that match the query.
func (ss *SQLStore) Delete(ctx context.Context, query Query) (int, error) {
	sq := `DELETE FROM triples WHERE %s`

	where, args, err := getWhereClause(query)
	if err != nil {
		return 0, err
	}
	if where == "" {
		return 0, errors.New("no query clause specified")
	}

	q := fmt.Sprintf(sq, where)

	res, err := ss.DB.ExecContext(ctx, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// ReWeight updates the weight of all the triples matching the given query.
func (ss *SQLStore) ReWeight(ctx context.Context, query Query, delta float64, replace bool) (int, error) {
	args := []interface{}{delta}
	sq := "UPDATE triples "
	if replace {
		sq += "SET weight=?"
	} else {
		sq += "SET weight=weight + ?"
	}

	where, tmp, err := getWhereClause(query)
	if err != nil {
		return 0, err
	}
	if where != "" {
		sq += fmt.Sprintf(" WHERE %s", where)
		args = append(args, tmp...)
	}

	res, err := ss.DB.ExecContext(ctx, sq, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Setup runs appropriate queries to setup all the required tables.
func (ss *SQLStore) Setup(ctx context.Context) error {
	_, err := ss.DB.ExecContext(ctx, sqlMigration)
	return err
}

func getWhereClause(query Query) (string, []interface{}, error) {
	where := []string{}
	args := []interface{}{}
	for col, clause := range query.Map() {
		sqlOp, value, err := toSQL(clause)
		if err != nil {
			return "", nil, err
		}

		where = append(where, fmt.Sprintf("%s %s ?", col, sqlOp))
		args = append(args, value)
	}
	return strings.TrimSpace(strings.Join(where, " AND ")), args, nil
}

func toSQL(clause Clause) (string, string, error) {
	switch clause.Type {
	case "=", "==", "equal":
		return "=", clause.Value, nil

	case "~", "~=", "like":
		return " LIKE ", strings.Replace(clause.Value, "*", "%", -1), nil

	case ">", "gt":
		return ">", clause.Value, nil

	case "<", "lt":
		return "<", clause.Value, nil

	case "<=", "lte":
		return "<=", clause.Value, nil

	case ">=", "gte":
		return ">=", clause.Value, nil
	}

	return "", "", fmt.Errorf("clause type '%s' not supported", clause.Type)
}

const sqlMigration = `
create table if not exists triples (
	source text not null,
	predicate text not null,
	target text not null,
	weight decimal not null default 0
);
create unique index if not exists triple_idx on triples (source, predicate, target);
`
