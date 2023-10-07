// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package sqlite

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"ariga.io/atlas/sql/internal/sqlx"
	"ariga.io/atlas/sql/schema"
)

// DefaultDiff provides basic diffing capabilities for MySQL dialects.
// Note, it is recommended to call Open, create a new Driver and use its
// Differ when a database connection is available.
var DefaultDiff schema.Differ = &sqlx.Diff{DiffDriver: &diff{}}

// A diff provides a SQLite implementation for sqlx.DiffDriver.
type diff struct{}

// SchemaAttrDiff returns a changeset for migrating schema attributes from one state to the other.
func (*diff) SchemaAttrDiff(_, _ *schema.Schema) []schema.Change {
	// No special schema attribute diffing for SQLite.
	return nil
}

// SchemaObjectDiff returns a changeset for migrating schema objects from
// one state to the other.
func (*diff) SchemaObjectDiff(_, _ *schema.Schema) ([]schema.Change, error) {
	return nil, nil
}

// TableAttrDiff returns a changeset for migrating table attributes from one state to the other.
func (d *diff) TableAttrDiff(from, to *schema.Table) ([]schema.Change, error) {
	var changes []schema.Change
	for _, a := range []schema.Attr{&WithoutRowID{}, &Strict{}} {
		switch {
		case sqlx.Has(from.Attrs, a) && !sqlx.Has(to.Attrs, a):
			changes = append(changes, &schema.DropAttr{
				A: a,
			})
		case !sqlx.Has(from.Attrs, a) && sqlx.Has(to.Attrs, a):
			changes = append(changes, &schema.AddAttr{
				A: a,
			})
		}
	}
	return append(changes, sqlx.CheckDiff(from, to)...), nil
}

func (d *diff) ViewAttrChanged(_, _ *schema.View) bool {
	return false // Not implemented.
}

// ColumnChange returns the schema changes (if any) for migrating one column to the other.
func (d *diff) ColumnChange(_ *schema.Table, from, to *schema.Column) (schema.ChangeKind, error) {
	change := sqlx.CommentChange(from.Attrs, to.Attrs)
	if from.Type.Null != to.Type.Null {
		change |= schema.ChangeNull
	}
	changed, err := d.typeChanged(from, to)
	if err != nil {
		return schema.NoChange, err
	}
	if changed {
		change |= schema.ChangeType
	}
	if changed := d.defaultChanged(from, to); changed {
		change |= schema.ChangeDefault
	}
	if d.generatedChanged(from, to) {
		change |= schema.ChangeGenerated
	}
	return change, nil
}

// typeChanged reports if the column type was changed.
func (d *diff) typeChanged(from, to *schema.Column) (bool, error) {
	fromT, toT := from.Type.Type, to.Type.Type
	if fromT == nil || toT == nil {
		return false, fmt.Errorf("sqlite: missing type information for column %q", from.Name)
	}
	// Types are mismatched if they do not have the same "type affinity".
	return reflect.TypeOf(fromT) != reflect.TypeOf(toT), nil
}

// defaultChanged reports if the default value of a column was changed.
func (d *diff) defaultChanged(from, to *schema.Column) bool {
	d1, ok1 := sqlx.DefaultValue(from)
	d2, ok2 := sqlx.DefaultValue(to)
	if ok1 != ok2 {
		return true
	}
	if d1 == d2 {
		return false
	}
	x1, err1 := sqlx.Unquote(d1)
	x2, err2 := sqlx.Unquote(d2)
	return err1 != nil || err2 != nil || x1 != x2
}

// generatedChanged reports if the generated expression of a column was changed.
func (*diff) generatedChanged(from, to *schema.Column) bool {
	var (
		fromX, toX     schema.GeneratedExpr
		fromHas, toHas = sqlx.Has(from.Attrs, &fromX), sqlx.Has(to.Attrs, &toX)
	)
	return fromHas != toHas || fromHas && (sqlx.MayWrap(fromX.Expr) != sqlx.MayWrap(toX.Expr) || storedOrVirtual(fromX.Type) != storedOrVirtual(toX.Type))
}

// IsGeneratedIndexName reports if the index name was generated by the database.
// See: https://github.com/sqlite/sqlite/blob/e937df8/src/build.c#L3583.
func (d *diff) IsGeneratedIndexName(t *schema.Table, idx *schema.Index) bool {
	p := fmt.Sprintf("sqlite_autoindex_%s_", t.Name)
	if !strings.HasPrefix(idx.Name, p) {
		return false
	}
	i, err := strconv.ParseInt(strings.TrimPrefix(idx.Name, p), 10, 64)
	return err == nil && i > 0
}

// FindGeneratedIndex finds the table index that represents the generated index.
// This is useful because unlike MySQL/PostgreSQL, SQLite does not allow creating
// the generated indexes with their internal names. Therefore, they are renamed in
// normalization phase. See migrate.go#normalizeIdxName for more details.
func (d *diff) FindGeneratedIndex(t *schema.Table, idx *schema.Index) (*schema.Index, bool) {
	nr := schema.NewIndex(idx.Name)
	if normalizeIdxName(nr, t) != nil {
		return nil, false
	}
	return t.Index(nr.Name)
}

// IndexAttrChanged reports if the index attributes were changed.
func (*diff) IndexAttrChanged(from, to []schema.Attr) bool {
	var p1, p2 IndexPredicate
	return sqlx.Has(from, &p1) != sqlx.Has(to, &p2) || (p1.P != p2.P && p1.P != sqlx.MayWrap(p2.P))
}

// IndexPartAttrChanged reports if the index-part attributes were changed.
func (*diff) IndexPartAttrChanged(_, _ *schema.Index, _ int) bool {
	return false
}

// ReferenceChanged reports if the foreign key referential action was changed.
func (*diff) ReferenceChanged(from, to schema.ReferenceOption) bool {
	// According to SQLite, if an action is not explicitly
	// specified, it defaults to "NO ACTION".
	if from == "" {
		from = schema.NoAction
	}
	if to == "" {
		to = schema.NoAction
	}
	return from != to
}

// Normalize implements the sqlx.Normalizer interface.
func (d *diff) Normalize(from, to *schema.Table) error {
	used := make([]bool, len(to.ForeignKeys))
	// In SQLite, there is no easy way to get the foreign-key constraint
	// name, except for parsing the CREATE statement. Therefore, we check
	// if there is a foreign-key with identical properties.
	for _, fk1 := range from.ForeignKeys {
		for i, fk2 := range to.ForeignKeys {
			if used[i] {
				continue
			}
			if fk2.Symbol == fk1.Symbol && !sqlx.IsUint(fk1.Symbol) || sameFK(fk1, fk2) {
				fk1.Symbol = fk2.Symbol
				used[i] = true
			}
		}
	}
	// Normalize names of indexes generated by UNIQUE constraints before
	// comparing. See the normalizeIdxName function for details.
	for _, idx := range to.Indexes {
		if err := normalizeIdxName(idx, to); err != nil {
			return err
		}
	}
	return nil
}

func sameFK(fk1, fk2 *schema.ForeignKey) bool {
	if fk1.Table.Name != fk2.Table.Name || fk1.RefTable.Name != fk2.RefTable.Name ||
		len(fk1.Columns) != len(fk2.Columns) || len(fk1.RefColumns) != len(fk2.RefColumns) {
		return false
	}
	for i, c1 := range fk1.Columns {
		if c1.Name != fk2.Columns[i].Name {
			return false
		}
	}
	for i, c1 := range fk1.RefColumns {
		if c1.Name != fk2.RefColumns[i].Name {
			return false
		}
	}
	return true
}
