// Code generated by ent, DO NOT EDIT.

package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ekhvalov/bank-chat-service/internal/store/chat"
	"github.com/ekhvalov/bank-chat-service/internal/store/message"
	"github.com/ekhvalov/bank-chat-service/internal/store/problem"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

// ChatCreate is the builder for creating a Chat entity.
type ChatCreate struct {
	config
	mutation *ChatMutation
	hooks    []Hook
}

// SetClientID sets the "client_id" field.
func (cc *ChatCreate) SetClientID(ti types.UserID) *ChatCreate {
	cc.mutation.SetClientID(ti)
	return cc
}

// SetCreatedAt sets the "created_at" field.
func (cc *ChatCreate) SetCreatedAt(t time.Time) *ChatCreate {
	cc.mutation.SetCreatedAt(t)
	return cc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cc *ChatCreate) SetNillableCreatedAt(t *time.Time) *ChatCreate {
	if t != nil {
		cc.SetCreatedAt(*t)
	}
	return cc
}

// SetID sets the "id" field.
func (cc *ChatCreate) SetID(ti types.ChatID) *ChatCreate {
	cc.mutation.SetID(ti)
	return cc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (cc *ChatCreate) SetNillableID(ti *types.ChatID) *ChatCreate {
	if ti != nil {
		cc.SetID(*ti)
	}
	return cc
}

// AddMessageIDs adds the "messages" edge to the Message entity by IDs.
func (cc *ChatCreate) AddMessageIDs(ids ...types.MessageID) *ChatCreate {
	cc.mutation.AddMessageIDs(ids...)
	return cc
}

// AddMessages adds the "messages" edges to the Message entity.
func (cc *ChatCreate) AddMessages(m ...*Message) *ChatCreate {
	ids := make([]types.MessageID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return cc.AddMessageIDs(ids...)
}

// AddProblemIDs adds the "problems" edge to the Problem entity by IDs.
func (cc *ChatCreate) AddProblemIDs(ids ...types.ProblemID) *ChatCreate {
	cc.mutation.AddProblemIDs(ids...)
	return cc
}

// AddProblems adds the "problems" edges to the Problem entity.
func (cc *ChatCreate) AddProblems(p ...*Problem) *ChatCreate {
	ids := make([]types.ProblemID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cc.AddProblemIDs(ids...)
}

// Mutation returns the ChatMutation object of the builder.
func (cc *ChatCreate) Mutation() *ChatMutation {
	return cc.mutation
}

// Save creates the Chat in the database.
func (cc *ChatCreate) Save(ctx context.Context) (*Chat, error) {
	cc.defaults()
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *ChatCreate) SaveX(ctx context.Context) *Chat {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *ChatCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *ChatCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *ChatCreate) defaults() {
	if _, ok := cc.mutation.CreatedAt(); !ok {
		v := chat.DefaultCreatedAt()
		cc.mutation.SetCreatedAt(v)
	}
	if _, ok := cc.mutation.ID(); !ok {
		v := chat.DefaultID()
		cc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *ChatCreate) check() error {
	if _, ok := cc.mutation.ClientID(); !ok {
		return &ValidationError{Name: "client_id", err: errors.New(`store: missing required field "Chat.client_id"`)}
	}
	if v, ok := cc.mutation.ClientID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "client_id", err: fmt.Errorf(`store: validator failed for field "Chat.client_id": %w`, err)}
		}
	}
	if _, ok := cc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`store: missing required field "Chat.created_at"`)}
	}
	if v, ok := cc.mutation.ID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`store: validator failed for field "Chat.id": %w`, err)}
		}
	}
	return nil
}

func (cc *ChatCreate) sqlSave(ctx context.Context) (*Chat, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*types.ChatID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *ChatCreate) createSpec() (*Chat, *sqlgraph.CreateSpec) {
	var (
		_node = &Chat{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(chat.Table, sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID))
	)
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cc.mutation.ClientID(); ok {
		_spec.SetField(chat.FieldClientID, field.TypeUUID, value)
		_node.ClientID = value
	}
	if value, ok := cc.mutation.CreatedAt(); ok {
		_spec.SetField(chat.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := cc.mutation.MessagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.MessagesTable,
			Columns: []string{chat.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.ProblemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.ProblemsTable,
			Columns: []string{chat.ProblemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ChatCreateBulk is the builder for creating many Chat entities in bulk.
type ChatCreateBulk struct {
	config
	err      error
	builders []*ChatCreate
}

// Save creates the Chat entities in the database.
func (ccb *ChatCreateBulk) Save(ctx context.Context) ([]*Chat, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Chat, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ChatMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *ChatCreateBulk) SaveX(ctx context.Context) []*Chat {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *ChatCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *ChatCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
