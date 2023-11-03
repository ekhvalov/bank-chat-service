// Code generated by ent, DO NOT EDIT.

package gen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/chat"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/message"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/problem"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

// MessageCreate is the builder for creating a Message entity.
type MessageCreate struct {
	config
	mutation *MessageMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetChatID sets the "chat_id" field.
func (mc *MessageCreate) SetChatID(ti types.ChatID) *MessageCreate {
	mc.mutation.SetChatID(ti)
	return mc
}

// SetProblemID sets the "problem_id" field.
func (mc *MessageCreate) SetProblemID(ti types.ProblemID) *MessageCreate {
	mc.mutation.SetProblemID(ti)
	return mc
}

// SetAuthorID sets the "author_id" field.
func (mc *MessageCreate) SetAuthorID(ti types.UserID) *MessageCreate {
	mc.mutation.SetAuthorID(ti)
	return mc
}

// SetNillableAuthorID sets the "author_id" field if the given value is not nil.
func (mc *MessageCreate) SetNillableAuthorID(ti *types.UserID) *MessageCreate {
	if ti != nil {
		mc.SetAuthorID(*ti)
	}
	return mc
}

// SetInitialRequestID sets the "initial_request_id" field.
func (mc *MessageCreate) SetInitialRequestID(ti types.RequestID) *MessageCreate {
	mc.mutation.SetInitialRequestID(ti)
	return mc
}

// SetNillableInitialRequestID sets the "initial_request_id" field if the given value is not nil.
func (mc *MessageCreate) SetNillableInitialRequestID(ti *types.RequestID) *MessageCreate {
	if ti != nil {
		mc.SetInitialRequestID(*ti)
	}
	return mc
}

// SetIsVisibleForClient sets the "is_visible_for_client" field.
func (mc *MessageCreate) SetIsVisibleForClient(b bool) *MessageCreate {
	mc.mutation.SetIsVisibleForClient(b)
	return mc
}

// SetNillableIsVisibleForClient sets the "is_visible_for_client" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsVisibleForClient(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsVisibleForClient(*b)
	}
	return mc
}

// SetIsVisibleForManager sets the "is_visible_for_manager" field.
func (mc *MessageCreate) SetIsVisibleForManager(b bool) *MessageCreate {
	mc.mutation.SetIsVisibleForManager(b)
	return mc
}

// SetNillableIsVisibleForManager sets the "is_visible_for_manager" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsVisibleForManager(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsVisibleForManager(*b)
	}
	return mc
}

// SetBody sets the "body" field.
func (mc *MessageCreate) SetBody(s string) *MessageCreate {
	mc.mutation.SetBody(s)
	return mc
}

// SetCheckedAt sets the "checked_at" field.
func (mc *MessageCreate) SetCheckedAt(t time.Time) *MessageCreate {
	mc.mutation.SetCheckedAt(t)
	return mc
}

// SetNillableCheckedAt sets the "checked_at" field if the given value is not nil.
func (mc *MessageCreate) SetNillableCheckedAt(t *time.Time) *MessageCreate {
	if t != nil {
		mc.SetCheckedAt(*t)
	}
	return mc
}

// SetIsBlocked sets the "is_blocked" field.
func (mc *MessageCreate) SetIsBlocked(b bool) *MessageCreate {
	mc.mutation.SetIsBlocked(b)
	return mc
}

// SetNillableIsBlocked sets the "is_blocked" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsBlocked(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsBlocked(*b)
	}
	return mc
}

// SetIsService sets the "is_service" field.
func (mc *MessageCreate) SetIsService(b bool) *MessageCreate {
	mc.mutation.SetIsService(b)
	return mc
}

// SetNillableIsService sets the "is_service" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsService(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsService(*b)
	}
	return mc
}

// SetCreatedAt sets the "created_at" field.
func (mc *MessageCreate) SetCreatedAt(t time.Time) *MessageCreate {
	mc.mutation.SetCreatedAt(t)
	return mc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (mc *MessageCreate) SetNillableCreatedAt(t *time.Time) *MessageCreate {
	if t != nil {
		mc.SetCreatedAt(*t)
	}
	return mc
}

// SetID sets the "id" field.
func (mc *MessageCreate) SetID(ti types.MessageID) *MessageCreate {
	mc.mutation.SetID(ti)
	return mc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (mc *MessageCreate) SetNillableID(ti *types.MessageID) *MessageCreate {
	if ti != nil {
		mc.SetID(*ti)
	}
	return mc
}

// SetProblem sets the "problem" edge to the Problem entity.
func (mc *MessageCreate) SetProblem(p *Problem) *MessageCreate {
	return mc.SetProblemID(p.ID)
}

// SetChat sets the "chat" edge to the Chat entity.
func (mc *MessageCreate) SetChat(c *Chat) *MessageCreate {
	return mc.SetChatID(c.ID)
}

// Mutation returns the MessageMutation object of the builder.
func (mc *MessageCreate) Mutation() *MessageMutation {
	return mc.mutation
}

// Save creates the Message in the database.
func (mc *MessageCreate) Save(ctx context.Context) (*Message, error) {
	mc.defaults()
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MessageCreate) SaveX(ctx context.Context) *Message {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MessageCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MessageCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mc *MessageCreate) defaults() {
	if _, ok := mc.mutation.InitialRequestID(); !ok {
		v := message.DefaultInitialRequestID()
		mc.mutation.SetInitialRequestID(v)
	}
	if _, ok := mc.mutation.IsVisibleForClient(); !ok {
		v := message.DefaultIsVisibleForClient
		mc.mutation.SetIsVisibleForClient(v)
	}
	if _, ok := mc.mutation.IsVisibleForManager(); !ok {
		v := message.DefaultIsVisibleForManager
		mc.mutation.SetIsVisibleForManager(v)
	}
	if _, ok := mc.mutation.IsBlocked(); !ok {
		v := message.DefaultIsBlocked
		mc.mutation.SetIsBlocked(v)
	}
	if _, ok := mc.mutation.IsService(); !ok {
		v := message.DefaultIsService
		mc.mutation.SetIsService(v)
	}
	if _, ok := mc.mutation.CreatedAt(); !ok {
		v := message.DefaultCreatedAt()
		mc.mutation.SetCreatedAt(v)
	}
	if _, ok := mc.mutation.ID(); !ok {
		v := message.DefaultID()
		mc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MessageCreate) check() error {
	if _, ok := mc.mutation.ChatID(); !ok {
		return &ValidationError{Name: "chat_id", err: errors.New(`gen: missing required field "Message.chat_id"`)}
	}
	if v, ok := mc.mutation.ChatID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "chat_id", err: fmt.Errorf(`gen: validator failed for field "Message.chat_id": %w`, err)}
		}
	}
	if _, ok := mc.mutation.ProblemID(); !ok {
		return &ValidationError{Name: "problem_id", err: errors.New(`gen: missing required field "Message.problem_id"`)}
	}
	if v, ok := mc.mutation.ProblemID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "problem_id", err: fmt.Errorf(`gen: validator failed for field "Message.problem_id": %w`, err)}
		}
	}
	if v, ok := mc.mutation.AuthorID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "author_id", err: fmt.Errorf(`gen: validator failed for field "Message.author_id": %w`, err)}
		}
	}
	if _, ok := mc.mutation.InitialRequestID(); !ok {
		return &ValidationError{Name: "initial_request_id", err: errors.New(`gen: missing required field "Message.initial_request_id"`)}
	}
	if v, ok := mc.mutation.InitialRequestID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "initial_request_id", err: fmt.Errorf(`gen: validator failed for field "Message.initial_request_id": %w`, err)}
		}
	}
	if _, ok := mc.mutation.IsVisibleForClient(); !ok {
		return &ValidationError{Name: "is_visible_for_client", err: errors.New(`gen: missing required field "Message.is_visible_for_client"`)}
	}
	if _, ok := mc.mutation.IsVisibleForManager(); !ok {
		return &ValidationError{Name: "is_visible_for_manager", err: errors.New(`gen: missing required field "Message.is_visible_for_manager"`)}
	}
	if _, ok := mc.mutation.Body(); !ok {
		return &ValidationError{Name: "body", err: errors.New(`gen: missing required field "Message.body"`)}
	}
	if v, ok := mc.mutation.Body(); ok {
		if err := message.BodyValidator(v); err != nil {
			return &ValidationError{Name: "body", err: fmt.Errorf(`gen: validator failed for field "Message.body": %w`, err)}
		}
	}
	if _, ok := mc.mutation.IsBlocked(); !ok {
		return &ValidationError{Name: "is_blocked", err: errors.New(`gen: missing required field "Message.is_blocked"`)}
	}
	if _, ok := mc.mutation.IsService(); !ok {
		return &ValidationError{Name: "is_service", err: errors.New(`gen: missing required field "Message.is_service"`)}
	}
	if _, ok := mc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`gen: missing required field "Message.created_at"`)}
	}
	if v, ok := mc.mutation.ID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`gen: validator failed for field "Message.id": %w`, err)}
		}
	}
	if _, ok := mc.mutation.ProblemID(); !ok {
		return &ValidationError{Name: "problem", err: errors.New(`gen: missing required edge "Message.problem"`)}
	}
	if _, ok := mc.mutation.ChatID(); !ok {
		return &ValidationError{Name: "chat", err: errors.New(`gen: missing required edge "Message.chat"`)}
	}
	return nil
}

func (mc *MessageCreate) sqlSave(ctx context.Context) (*Message, error) {
	if err := mc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*types.MessageID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	mc.mutation.id = &_node.ID
	mc.mutation.done = true
	return _node, nil
}

func (mc *MessageCreate) createSpec() (*Message, *sqlgraph.CreateSpec) {
	var (
		_node = &Message{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(message.Table, sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = mc.conflict
	if id, ok := mc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := mc.mutation.AuthorID(); ok {
		_spec.SetField(message.FieldAuthorID, field.TypeUUID, value)
		_node.AuthorID = value
	}
	if value, ok := mc.mutation.InitialRequestID(); ok {
		_spec.SetField(message.FieldInitialRequestID, field.TypeUUID, value)
		_node.InitialRequestID = value
	}
	if value, ok := mc.mutation.IsVisibleForClient(); ok {
		_spec.SetField(message.FieldIsVisibleForClient, field.TypeBool, value)
		_node.IsVisibleForClient = value
	}
	if value, ok := mc.mutation.IsVisibleForManager(); ok {
		_spec.SetField(message.FieldIsVisibleForManager, field.TypeBool, value)
		_node.IsVisibleForManager = value
	}
	if value, ok := mc.mutation.Body(); ok {
		_spec.SetField(message.FieldBody, field.TypeString, value)
		_node.Body = value
	}
	if value, ok := mc.mutation.CheckedAt(); ok {
		_spec.SetField(message.FieldCheckedAt, field.TypeTime, value)
		_node.CheckedAt = value
	}
	if value, ok := mc.mutation.IsBlocked(); ok {
		_spec.SetField(message.FieldIsBlocked, field.TypeBool, value)
		_node.IsBlocked = value
	}
	if value, ok := mc.mutation.IsService(); ok {
		_spec.SetField(message.FieldIsService, field.TypeBool, value)
		_node.IsService = value
	}
	if value, ok := mc.mutation.CreatedAt(); ok {
		_spec.SetField(message.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := mc.mutation.ProblemIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ProblemTable,
			Columns: []string{message.ProblemColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ProblemID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.ChatIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ChatTable,
			Columns: []string{message.ChatColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ChatID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Message.Create().
//		SetChatID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MessageUpsert) {
//			SetChatID(v+v).
//		}).
//		Exec(ctx)
func (mc *MessageCreate) OnConflict(opts ...sql.ConflictOption) *MessageUpsertOne {
	mc.conflict = opts
	return &MessageUpsertOne{
		create: mc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Message.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mc *MessageCreate) OnConflictColumns(columns ...string) *MessageUpsertOne {
	mc.conflict = append(mc.conflict, sql.ConflictColumns(columns...))
	return &MessageUpsertOne{
		create: mc,
	}
}

type (
	// MessageUpsertOne is the builder for "upsert"-ing
	//  one Message node.
	MessageUpsertOne struct {
		create *MessageCreate
	}

	// MessageUpsert is the "OnConflict" setter.
	MessageUpsert struct {
		*sql.UpdateSet
	}
)

// SetChatID sets the "chat_id" field.
func (u *MessageUpsert) SetChatID(v types.ChatID) *MessageUpsert {
	u.Set(message.FieldChatID, v)
	return u
}

// UpdateChatID sets the "chat_id" field to the value that was provided on create.
func (u *MessageUpsert) UpdateChatID() *MessageUpsert {
	u.SetExcluded(message.FieldChatID)
	return u
}

// SetProblemID sets the "problem_id" field.
func (u *MessageUpsert) SetProblemID(v types.ProblemID) *MessageUpsert {
	u.Set(message.FieldProblemID, v)
	return u
}

// UpdateProblemID sets the "problem_id" field to the value that was provided on create.
func (u *MessageUpsert) UpdateProblemID() *MessageUpsert {
	u.SetExcluded(message.FieldProblemID)
	return u
}

// SetIsVisibleForClient sets the "is_visible_for_client" field.
func (u *MessageUpsert) SetIsVisibleForClient(v bool) *MessageUpsert {
	u.Set(message.FieldIsVisibleForClient, v)
	return u
}

// UpdateIsVisibleForClient sets the "is_visible_for_client" field to the value that was provided on create.
func (u *MessageUpsert) UpdateIsVisibleForClient() *MessageUpsert {
	u.SetExcluded(message.FieldIsVisibleForClient)
	return u
}

// SetIsVisibleForManager sets the "is_visible_for_manager" field.
func (u *MessageUpsert) SetIsVisibleForManager(v bool) *MessageUpsert {
	u.Set(message.FieldIsVisibleForManager, v)
	return u
}

// UpdateIsVisibleForManager sets the "is_visible_for_manager" field to the value that was provided on create.
func (u *MessageUpsert) UpdateIsVisibleForManager() *MessageUpsert {
	u.SetExcluded(message.FieldIsVisibleForManager)
	return u
}

// SetCheckedAt sets the "checked_at" field.
func (u *MessageUpsert) SetCheckedAt(v time.Time) *MessageUpsert {
	u.Set(message.FieldCheckedAt, v)
	return u
}

// UpdateCheckedAt sets the "checked_at" field to the value that was provided on create.
func (u *MessageUpsert) UpdateCheckedAt() *MessageUpsert {
	u.SetExcluded(message.FieldCheckedAt)
	return u
}

// ClearCheckedAt clears the value of the "checked_at" field.
func (u *MessageUpsert) ClearCheckedAt() *MessageUpsert {
	u.SetNull(message.FieldCheckedAt)
	return u
}

// SetIsBlocked sets the "is_blocked" field.
func (u *MessageUpsert) SetIsBlocked(v bool) *MessageUpsert {
	u.Set(message.FieldIsBlocked, v)
	return u
}

// UpdateIsBlocked sets the "is_blocked" field to the value that was provided on create.
func (u *MessageUpsert) UpdateIsBlocked() *MessageUpsert {
	u.SetExcluded(message.FieldIsBlocked)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Message.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(message.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MessageUpsertOne) UpdateNewValues() *MessageUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(message.FieldID)
		}
		if _, exists := u.create.mutation.AuthorID(); exists {
			s.SetIgnore(message.FieldAuthorID)
		}
		if _, exists := u.create.mutation.InitialRequestID(); exists {
			s.SetIgnore(message.FieldInitialRequestID)
		}
		if _, exists := u.create.mutation.Body(); exists {
			s.SetIgnore(message.FieldBody)
		}
		if _, exists := u.create.mutation.IsService(); exists {
			s.SetIgnore(message.FieldIsService)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(message.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Message.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *MessageUpsertOne) Ignore() *MessageUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MessageUpsertOne) DoNothing() *MessageUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MessageCreate.OnConflict
// documentation for more info.
func (u *MessageUpsertOne) Update(set func(*MessageUpsert)) *MessageUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MessageUpsert{UpdateSet: update})
	}))
	return u
}

// SetChatID sets the "chat_id" field.
func (u *MessageUpsertOne) SetChatID(v types.ChatID) *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.SetChatID(v)
	})
}

// UpdateChatID sets the "chat_id" field to the value that was provided on create.
func (u *MessageUpsertOne) UpdateChatID() *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateChatID()
	})
}

// SetProblemID sets the "problem_id" field.
func (u *MessageUpsertOne) SetProblemID(v types.ProblemID) *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.SetProblemID(v)
	})
}

// UpdateProblemID sets the "problem_id" field to the value that was provided on create.
func (u *MessageUpsertOne) UpdateProblemID() *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateProblemID()
	})
}

// SetIsVisibleForClient sets the "is_visible_for_client" field.
func (u *MessageUpsertOne) SetIsVisibleForClient(v bool) *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.SetIsVisibleForClient(v)
	})
}

// UpdateIsVisibleForClient sets the "is_visible_for_client" field to the value that was provided on create.
func (u *MessageUpsertOne) UpdateIsVisibleForClient() *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateIsVisibleForClient()
	})
}

// SetIsVisibleForManager sets the "is_visible_for_manager" field.
func (u *MessageUpsertOne) SetIsVisibleForManager(v bool) *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.SetIsVisibleForManager(v)
	})
}

// UpdateIsVisibleForManager sets the "is_visible_for_manager" field to the value that was provided on create.
func (u *MessageUpsertOne) UpdateIsVisibleForManager() *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateIsVisibleForManager()
	})
}

// SetCheckedAt sets the "checked_at" field.
func (u *MessageUpsertOne) SetCheckedAt(v time.Time) *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.SetCheckedAt(v)
	})
}

// UpdateCheckedAt sets the "checked_at" field to the value that was provided on create.
func (u *MessageUpsertOne) UpdateCheckedAt() *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateCheckedAt()
	})
}

// ClearCheckedAt clears the value of the "checked_at" field.
func (u *MessageUpsertOne) ClearCheckedAt() *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.ClearCheckedAt()
	})
}

// SetIsBlocked sets the "is_blocked" field.
func (u *MessageUpsertOne) SetIsBlocked(v bool) *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.SetIsBlocked(v)
	})
}

// UpdateIsBlocked sets the "is_blocked" field to the value that was provided on create.
func (u *MessageUpsertOne) UpdateIsBlocked() *MessageUpsertOne {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateIsBlocked()
	})
}

// Exec executes the query.
func (u *MessageUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("gen: missing options for MessageCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MessageUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *MessageUpsertOne) ID(ctx context.Context) (id types.MessageID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("gen: MessageUpsertOne.ID is not supported by MySQL driver. Use MessageUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *MessageUpsertOne) IDX(ctx context.Context) types.MessageID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// MessageCreateBulk is the builder for creating many Message entities in bulk.
type MessageCreateBulk struct {
	config
	err      error
	builders []*MessageCreate
	conflict []sql.ConflictOption
}

// Save creates the Message entities in the database.
func (mcb *MessageCreateBulk) Save(ctx context.Context) ([]*Message, error) {
	if mcb.err != nil {
		return nil, mcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Message, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MessageMutation)
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
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MessageCreateBulk) SaveX(ctx context.Context) []*Message {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MessageCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MessageCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Message.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MessageUpsert) {
//			SetChatID(v+v).
//		}).
//		Exec(ctx)
func (mcb *MessageCreateBulk) OnConflict(opts ...sql.ConflictOption) *MessageUpsertBulk {
	mcb.conflict = opts
	return &MessageUpsertBulk{
		create: mcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Message.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mcb *MessageCreateBulk) OnConflictColumns(columns ...string) *MessageUpsertBulk {
	mcb.conflict = append(mcb.conflict, sql.ConflictColumns(columns...))
	return &MessageUpsertBulk{
		create: mcb,
	}
}

// MessageUpsertBulk is the builder for "upsert"-ing
// a bulk of Message nodes.
type MessageUpsertBulk struct {
	create *MessageCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Message.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(message.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MessageUpsertBulk) UpdateNewValues() *MessageUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(message.FieldID)
			}
			if _, exists := b.mutation.AuthorID(); exists {
				s.SetIgnore(message.FieldAuthorID)
			}
			if _, exists := b.mutation.InitialRequestID(); exists {
				s.SetIgnore(message.FieldInitialRequestID)
			}
			if _, exists := b.mutation.Body(); exists {
				s.SetIgnore(message.FieldBody)
			}
			if _, exists := b.mutation.IsService(); exists {
				s.SetIgnore(message.FieldIsService)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(message.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Message.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *MessageUpsertBulk) Ignore() *MessageUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MessageUpsertBulk) DoNothing() *MessageUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MessageCreateBulk.OnConflict
// documentation for more info.
func (u *MessageUpsertBulk) Update(set func(*MessageUpsert)) *MessageUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MessageUpsert{UpdateSet: update})
	}))
	return u
}

// SetChatID sets the "chat_id" field.
func (u *MessageUpsertBulk) SetChatID(v types.ChatID) *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.SetChatID(v)
	})
}

// UpdateChatID sets the "chat_id" field to the value that was provided on create.
func (u *MessageUpsertBulk) UpdateChatID() *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateChatID()
	})
}

// SetProblemID sets the "problem_id" field.
func (u *MessageUpsertBulk) SetProblemID(v types.ProblemID) *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.SetProblemID(v)
	})
}

// UpdateProblemID sets the "problem_id" field to the value that was provided on create.
func (u *MessageUpsertBulk) UpdateProblemID() *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateProblemID()
	})
}

// SetIsVisibleForClient sets the "is_visible_for_client" field.
func (u *MessageUpsertBulk) SetIsVisibleForClient(v bool) *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.SetIsVisibleForClient(v)
	})
}

// UpdateIsVisibleForClient sets the "is_visible_for_client" field to the value that was provided on create.
func (u *MessageUpsertBulk) UpdateIsVisibleForClient() *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateIsVisibleForClient()
	})
}

// SetIsVisibleForManager sets the "is_visible_for_manager" field.
func (u *MessageUpsertBulk) SetIsVisibleForManager(v bool) *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.SetIsVisibleForManager(v)
	})
}

// UpdateIsVisibleForManager sets the "is_visible_for_manager" field to the value that was provided on create.
func (u *MessageUpsertBulk) UpdateIsVisibleForManager() *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateIsVisibleForManager()
	})
}

// SetCheckedAt sets the "checked_at" field.
func (u *MessageUpsertBulk) SetCheckedAt(v time.Time) *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.SetCheckedAt(v)
	})
}

// UpdateCheckedAt sets the "checked_at" field to the value that was provided on create.
func (u *MessageUpsertBulk) UpdateCheckedAt() *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateCheckedAt()
	})
}

// ClearCheckedAt clears the value of the "checked_at" field.
func (u *MessageUpsertBulk) ClearCheckedAt() *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.ClearCheckedAt()
	})
}

// SetIsBlocked sets the "is_blocked" field.
func (u *MessageUpsertBulk) SetIsBlocked(v bool) *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.SetIsBlocked(v)
	})
}

// UpdateIsBlocked sets the "is_blocked" field to the value that was provided on create.
func (u *MessageUpsertBulk) UpdateIsBlocked() *MessageUpsertBulk {
	return u.Update(func(s *MessageUpsert) {
		s.UpdateIsBlocked()
	})
}

// Exec executes the query.
func (u *MessageUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("gen: OnConflict was set for builder %d. Set it on the MessageCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("gen: missing options for MessageCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MessageUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
