package gormdb

import (
	"fmt"

	"gorm.io/gorm/clause"
)

type OrderByCase struct {
	Column clause.Column
	Values map[string]int
	Desc   bool
}

// Name where clause name
func (orderBy OrderByCase) Name() string {
	return "ORDER BY"
}

// Build build where clause
func (orderByCase OrderByCase) Build(builder clause.Builder) {
	builder.WriteString("CASE ")
	builder.WriteQuoted(orderByCase.Column)
	for field, weight := range orderByCase.Values {
		builder.WriteString(fmt.Sprintf(" WHEN '%s' THEN %d", field, weight)) // potential sql injection risk
	}
	builder.WriteString(" END")
	if orderByCase.Desc {
		builder.WriteString(" DESC")
	}
}

// MergeClause merge order by clauses
func (orderByCase OrderByCase) MergeClause(clause *clause.Clause) {
	clause.Expression = orderByCase
}
