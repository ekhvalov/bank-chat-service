// Code generated by ent, DO NOT EDIT.

package problem

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/ekhvalov/bank-chat-service/internal/store/predicate"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

// ID filters vertices based on their ID field.
func ID(id types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.ProblemID) predicate.Problem {
	return predicate.Problem(sql.FieldLTE(FieldID, id))
}

// ManagerID applies equality check predicate on the "manager_id" field. It's identical to ManagerIDEQ.
func ManagerID(v types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldManagerID, v))
}

// ResolvedAt applies equality check predicate on the "resolved_at" field. It's identical to ResolvedAtEQ.
func ResolvedAt(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldResolvedAt, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldCreatedAt, v))
}

// ManagerIDEQ applies the EQ predicate on the "manager_id" field.
func ManagerIDEQ(v types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldManagerID, v))
}

// ManagerIDNEQ applies the NEQ predicate on the "manager_id" field.
func ManagerIDNEQ(v types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldNEQ(FieldManagerID, v))
}

// ManagerIDIn applies the In predicate on the "manager_id" field.
func ManagerIDIn(vs ...types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldIn(FieldManagerID, vs...))
}

// ManagerIDNotIn applies the NotIn predicate on the "manager_id" field.
func ManagerIDNotIn(vs ...types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldNotIn(FieldManagerID, vs...))
}

// ManagerIDGT applies the GT predicate on the "manager_id" field.
func ManagerIDGT(v types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldGT(FieldManagerID, v))
}

// ManagerIDGTE applies the GTE predicate on the "manager_id" field.
func ManagerIDGTE(v types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldGTE(FieldManagerID, v))
}

// ManagerIDLT applies the LT predicate on the "manager_id" field.
func ManagerIDLT(v types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldLT(FieldManagerID, v))
}

// ManagerIDLTE applies the LTE predicate on the "manager_id" field.
func ManagerIDLTE(v types.UserID) predicate.Problem {
	return predicate.Problem(sql.FieldLTE(FieldManagerID, v))
}

// ResolvedAtEQ applies the EQ predicate on the "resolved_at" field.
func ResolvedAtEQ(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldResolvedAt, v))
}

// ResolvedAtNEQ applies the NEQ predicate on the "resolved_at" field.
func ResolvedAtNEQ(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldNEQ(FieldResolvedAt, v))
}

// ResolvedAtIn applies the In predicate on the "resolved_at" field.
func ResolvedAtIn(vs ...time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldIn(FieldResolvedAt, vs...))
}

// ResolvedAtNotIn applies the NotIn predicate on the "resolved_at" field.
func ResolvedAtNotIn(vs ...time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldNotIn(FieldResolvedAt, vs...))
}

// ResolvedAtGT applies the GT predicate on the "resolved_at" field.
func ResolvedAtGT(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldGT(FieldResolvedAt, v))
}

// ResolvedAtGTE applies the GTE predicate on the "resolved_at" field.
func ResolvedAtGTE(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldGTE(FieldResolvedAt, v))
}

// ResolvedAtLT applies the LT predicate on the "resolved_at" field.
func ResolvedAtLT(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldLT(FieldResolvedAt, v))
}

// ResolvedAtLTE applies the LTE predicate on the "resolved_at" field.
func ResolvedAtLTE(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldLTE(FieldResolvedAt, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Problem {
	return predicate.Problem(sql.FieldLTE(FieldCreatedAt, v))
}

// HasChat applies the HasEdge predicate on the "chat" edge.
func HasChat() predicate.Problem {
	return predicate.Problem(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ChatTable, ChatColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasChatWith applies the HasEdge predicate on the "chat" edge with a given conditions (other predicates).
func HasChatWith(preds ...predicate.Chat) predicate.Problem {
	return predicate.Problem(func(s *sql.Selector) {
		step := newChatStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMessages applies the HasEdge predicate on the "messages" edge.
func HasMessages() predicate.Problem {
	return predicate.Problem(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MessagesTable, MessagesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMessagesWith applies the HasEdge predicate on the "messages" edge with a given conditions (other predicates).
func HasMessagesWith(preds ...predicate.Message) predicate.Problem {
	return predicate.Problem(func(s *sql.Selector) {
		step := newMessagesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Problem) predicate.Problem {
	return predicate.Problem(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Problem) predicate.Problem {
	return predicate.Problem(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Problem) predicate.Problem {
	return predicate.Problem(sql.NotPredicates(p))
}
