// Code generated by ent, DO NOT EDIT.

package gen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/job"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/predicate"
)

// JobUpdate is the builder for updating Job entities.
type JobUpdate struct {
	config
	hooks    []Hook
	mutation *JobMutation
}

// Where appends a list predicates to the JobUpdate builder.
func (ju *JobUpdate) Where(ps ...predicate.Job) *JobUpdate {
	ju.mutation.Where(ps...)
	return ju
}

// SetAttempts sets the "attempts" field.
func (ju *JobUpdate) SetAttempts(i int) *JobUpdate {
	ju.mutation.ResetAttempts()
	ju.mutation.SetAttempts(i)
	return ju
}

// SetNillableAttempts sets the "attempts" field if the given value is not nil.
func (ju *JobUpdate) SetNillableAttempts(i *int) *JobUpdate {
	if i != nil {
		ju.SetAttempts(*i)
	}
	return ju
}

// AddAttempts adds i to the "attempts" field.
func (ju *JobUpdate) AddAttempts(i int) *JobUpdate {
	ju.mutation.AddAttempts(i)
	return ju
}

// SetReservedUntil sets the "reserved_until" field.
func (ju *JobUpdate) SetReservedUntil(t time.Time) *JobUpdate {
	ju.mutation.SetReservedUntil(t)
	return ju
}

// SetNillableReservedUntil sets the "reserved_until" field if the given value is not nil.
func (ju *JobUpdate) SetNillableReservedUntil(t *time.Time) *JobUpdate {
	if t != nil {
		ju.SetReservedUntil(*t)
	}
	return ju
}

// Mutation returns the JobMutation object of the builder.
func (ju *JobUpdate) Mutation() *JobMutation {
	return ju.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ju *JobUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ju.sqlSave, ju.mutation, ju.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ju *JobUpdate) SaveX(ctx context.Context) int {
	affected, err := ju.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ju *JobUpdate) Exec(ctx context.Context) error {
	_, err := ju.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ju *JobUpdate) ExecX(ctx context.Context) {
	if err := ju.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ju *JobUpdate) check() error {
	if v, ok := ju.mutation.Attempts(); ok {
		if err := job.AttemptsValidator(v); err != nil {
			return &ValidationError{Name: "attempts", err: fmt.Errorf(`gen: validator failed for field "Job.attempts": %w`, err)}
		}
	}
	return nil
}

func (ju *JobUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ju.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(job.Table, job.Columns, sqlgraph.NewFieldSpec(job.FieldID, field.TypeUUID))
	if ps := ju.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ju.mutation.Attempts(); ok {
		_spec.SetField(job.FieldAttempts, field.TypeInt, value)
	}
	if value, ok := ju.mutation.AddedAttempts(); ok {
		_spec.AddField(job.FieldAttempts, field.TypeInt, value)
	}
	if value, ok := ju.mutation.ReservedUntil(); ok {
		_spec.SetField(job.FieldReservedUntil, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ju.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{job.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ju.mutation.done = true
	return n, nil
}

// JobUpdateOne is the builder for updating a single Job entity.
type JobUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *JobMutation
}

// SetAttempts sets the "attempts" field.
func (juo *JobUpdateOne) SetAttempts(i int) *JobUpdateOne {
	juo.mutation.ResetAttempts()
	juo.mutation.SetAttempts(i)
	return juo
}

// SetNillableAttempts sets the "attempts" field if the given value is not nil.
func (juo *JobUpdateOne) SetNillableAttempts(i *int) *JobUpdateOne {
	if i != nil {
		juo.SetAttempts(*i)
	}
	return juo
}

// AddAttempts adds i to the "attempts" field.
func (juo *JobUpdateOne) AddAttempts(i int) *JobUpdateOne {
	juo.mutation.AddAttempts(i)
	return juo
}

// SetReservedUntil sets the "reserved_until" field.
func (juo *JobUpdateOne) SetReservedUntil(t time.Time) *JobUpdateOne {
	juo.mutation.SetReservedUntil(t)
	return juo
}

// SetNillableReservedUntil sets the "reserved_until" field if the given value is not nil.
func (juo *JobUpdateOne) SetNillableReservedUntil(t *time.Time) *JobUpdateOne {
	if t != nil {
		juo.SetReservedUntil(*t)
	}
	return juo
}

// Mutation returns the JobMutation object of the builder.
func (juo *JobUpdateOne) Mutation() *JobMutation {
	return juo.mutation
}

// Where appends a list predicates to the JobUpdate builder.
func (juo *JobUpdateOne) Where(ps ...predicate.Job) *JobUpdateOne {
	juo.mutation.Where(ps...)
	return juo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (juo *JobUpdateOne) Select(field string, fields ...string) *JobUpdateOne {
	juo.fields = append([]string{field}, fields...)
	return juo
}

// Save executes the query and returns the updated Job entity.
func (juo *JobUpdateOne) Save(ctx context.Context) (*Job, error) {
	return withHooks(ctx, juo.sqlSave, juo.mutation, juo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (juo *JobUpdateOne) SaveX(ctx context.Context) *Job {
	node, err := juo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (juo *JobUpdateOne) Exec(ctx context.Context) error {
	_, err := juo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (juo *JobUpdateOne) ExecX(ctx context.Context) {
	if err := juo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (juo *JobUpdateOne) check() error {
	if v, ok := juo.mutation.Attempts(); ok {
		if err := job.AttemptsValidator(v); err != nil {
			return &ValidationError{Name: "attempts", err: fmt.Errorf(`gen: validator failed for field "Job.attempts": %w`, err)}
		}
	}
	return nil
}

func (juo *JobUpdateOne) sqlSave(ctx context.Context) (_node *Job, err error) {
	if err := juo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(job.Table, job.Columns, sqlgraph.NewFieldSpec(job.FieldID, field.TypeUUID))
	id, ok := juo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`gen: missing "Job.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := juo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, job.FieldID)
		for _, f := range fields {
			if !job.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("gen: invalid field %q for query", f)}
			}
			if f != job.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := juo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := juo.mutation.Attempts(); ok {
		_spec.SetField(job.FieldAttempts, field.TypeInt, value)
	}
	if value, ok := juo.mutation.AddedAttempts(); ok {
		_spec.AddField(job.FieldAttempts, field.TypeInt, value)
	}
	if value, ok := juo.mutation.ReservedUntil(); ok {
		_spec.SetField(job.FieldReservedUntil, field.TypeTime, value)
	}
	_node = &Job{config: juo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, juo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{job.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	juo.mutation.done = true
	return _node, nil
}
