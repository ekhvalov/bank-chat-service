// Code generated by ent, DO NOT EDIT.

package message

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/ekhvalov/bank-chat-service/internal/store/predicate"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

// ID filters vertices based on their ID field.
func ID(id types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.MessageID) predicate.Message {
	return predicate.Message(sql.FieldLTE(FieldID, id))
}

// AuthorID applies equality check predicate on the "author_id" field. It's identical to AuthorIDEQ.
func AuthorID(v types.UserID) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldAuthorID, v))
}

// IsVisibleForClient applies equality check predicate on the "is_visible_for_client" field. It's identical to IsVisibleForClientEQ.
func IsVisibleForClient(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsVisibleForClient, v))
}

// IsVisibleForManager applies equality check predicate on the "is_visible_for_manager" field. It's identical to IsVisibleForManagerEQ.
func IsVisibleForManager(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsVisibleForManager, v))
}

// Body applies equality check predicate on the "body" field. It's identical to BodyEQ.
func Body(v string) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldBody, v))
}

// CheckedAt applies equality check predicate on the "checked_at" field. It's identical to CheckedAtEQ.
func CheckedAt(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldCheckedAt, v))
}

// IsBlocked applies equality check predicate on the "is_blocked" field. It's identical to IsBlockedEQ.
func IsBlocked(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsBlocked, v))
}

// IsService applies equality check predicate on the "is_service" field. It's identical to IsServiceEQ.
func IsService(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsService, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldCreatedAt, v))
}

// AuthorIDEQ applies the EQ predicate on the "author_id" field.
func AuthorIDEQ(v types.UserID) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldAuthorID, v))
}

// AuthorIDNEQ applies the NEQ predicate on the "author_id" field.
func AuthorIDNEQ(v types.UserID) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldAuthorID, v))
}

// AuthorIDIn applies the In predicate on the "author_id" field.
func AuthorIDIn(vs ...types.UserID) predicate.Message {
	return predicate.Message(sql.FieldIn(FieldAuthorID, vs...))
}

// AuthorIDNotIn applies the NotIn predicate on the "author_id" field.
func AuthorIDNotIn(vs ...types.UserID) predicate.Message {
	return predicate.Message(sql.FieldNotIn(FieldAuthorID, vs...))
}

// AuthorIDGT applies the GT predicate on the "author_id" field.
func AuthorIDGT(v types.UserID) predicate.Message {
	return predicate.Message(sql.FieldGT(FieldAuthorID, v))
}

// AuthorIDGTE applies the GTE predicate on the "author_id" field.
func AuthorIDGTE(v types.UserID) predicate.Message {
	return predicate.Message(sql.FieldGTE(FieldAuthorID, v))
}

// AuthorIDLT applies the LT predicate on the "author_id" field.
func AuthorIDLT(v types.UserID) predicate.Message {
	return predicate.Message(sql.FieldLT(FieldAuthorID, v))
}

// AuthorIDLTE applies the LTE predicate on the "author_id" field.
func AuthorIDLTE(v types.UserID) predicate.Message {
	return predicate.Message(sql.FieldLTE(FieldAuthorID, v))
}

// AuthorIDIsNil applies the IsNil predicate on the "author_id" field.
func AuthorIDIsNil() predicate.Message {
	return predicate.Message(sql.FieldIsNull(FieldAuthorID))
}

// AuthorIDNotNil applies the NotNil predicate on the "author_id" field.
func AuthorIDNotNil() predicate.Message {
	return predicate.Message(sql.FieldNotNull(FieldAuthorID))
}

// IsVisibleForClientEQ applies the EQ predicate on the "is_visible_for_client" field.
func IsVisibleForClientEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsVisibleForClient, v))
}

// IsVisibleForClientNEQ applies the NEQ predicate on the "is_visible_for_client" field.
func IsVisibleForClientNEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldIsVisibleForClient, v))
}

// IsVisibleForManagerEQ applies the EQ predicate on the "is_visible_for_manager" field.
func IsVisibleForManagerEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsVisibleForManager, v))
}

// IsVisibleForManagerNEQ applies the NEQ predicate on the "is_visible_for_manager" field.
func IsVisibleForManagerNEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldIsVisibleForManager, v))
}

// BodyEQ applies the EQ predicate on the "body" field.
func BodyEQ(v string) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldBody, v))
}

// BodyNEQ applies the NEQ predicate on the "body" field.
func BodyNEQ(v string) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldBody, v))
}

// BodyIn applies the In predicate on the "body" field.
func BodyIn(vs ...string) predicate.Message {
	return predicate.Message(sql.FieldIn(FieldBody, vs...))
}

// BodyNotIn applies the NotIn predicate on the "body" field.
func BodyNotIn(vs ...string) predicate.Message {
	return predicate.Message(sql.FieldNotIn(FieldBody, vs...))
}

// BodyGT applies the GT predicate on the "body" field.
func BodyGT(v string) predicate.Message {
	return predicate.Message(sql.FieldGT(FieldBody, v))
}

// BodyGTE applies the GTE predicate on the "body" field.
func BodyGTE(v string) predicate.Message {
	return predicate.Message(sql.FieldGTE(FieldBody, v))
}

// BodyLT applies the LT predicate on the "body" field.
func BodyLT(v string) predicate.Message {
	return predicate.Message(sql.FieldLT(FieldBody, v))
}

// BodyLTE applies the LTE predicate on the "body" field.
func BodyLTE(v string) predicate.Message {
	return predicate.Message(sql.FieldLTE(FieldBody, v))
}

// BodyContains applies the Contains predicate on the "body" field.
func BodyContains(v string) predicate.Message {
	return predicate.Message(sql.FieldContains(FieldBody, v))
}

// BodyHasPrefix applies the HasPrefix predicate on the "body" field.
func BodyHasPrefix(v string) predicate.Message {
	return predicate.Message(sql.FieldHasPrefix(FieldBody, v))
}

// BodyHasSuffix applies the HasSuffix predicate on the "body" field.
func BodyHasSuffix(v string) predicate.Message {
	return predicate.Message(sql.FieldHasSuffix(FieldBody, v))
}

// BodyEqualFold applies the EqualFold predicate on the "body" field.
func BodyEqualFold(v string) predicate.Message {
	return predicate.Message(sql.FieldEqualFold(FieldBody, v))
}

// BodyContainsFold applies the ContainsFold predicate on the "body" field.
func BodyContainsFold(v string) predicate.Message {
	return predicate.Message(sql.FieldContainsFold(FieldBody, v))
}

// CheckedAtEQ applies the EQ predicate on the "checked_at" field.
func CheckedAtEQ(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldCheckedAt, v))
}

// CheckedAtNEQ applies the NEQ predicate on the "checked_at" field.
func CheckedAtNEQ(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldCheckedAt, v))
}

// CheckedAtIn applies the In predicate on the "checked_at" field.
func CheckedAtIn(vs ...time.Time) predicate.Message {
	return predicate.Message(sql.FieldIn(FieldCheckedAt, vs...))
}

// CheckedAtNotIn applies the NotIn predicate on the "checked_at" field.
func CheckedAtNotIn(vs ...time.Time) predicate.Message {
	return predicate.Message(sql.FieldNotIn(FieldCheckedAt, vs...))
}

// CheckedAtGT applies the GT predicate on the "checked_at" field.
func CheckedAtGT(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldGT(FieldCheckedAt, v))
}

// CheckedAtGTE applies the GTE predicate on the "checked_at" field.
func CheckedAtGTE(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldGTE(FieldCheckedAt, v))
}

// CheckedAtLT applies the LT predicate on the "checked_at" field.
func CheckedAtLT(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldLT(FieldCheckedAt, v))
}

// CheckedAtLTE applies the LTE predicate on the "checked_at" field.
func CheckedAtLTE(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldLTE(FieldCheckedAt, v))
}

// IsBlockedEQ applies the EQ predicate on the "is_blocked" field.
func IsBlockedEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsBlocked, v))
}

// IsBlockedNEQ applies the NEQ predicate on the "is_blocked" field.
func IsBlockedNEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldIsBlocked, v))
}

// IsServiceEQ applies the EQ predicate on the "is_service" field.
func IsServiceEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldIsService, v))
}

// IsServiceNEQ applies the NEQ predicate on the "is_service" field.
func IsServiceNEQ(v bool) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldIsService, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Message {
	return predicate.Message(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Message {
	return predicate.Message(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Message {
	return predicate.Message(sql.FieldLTE(FieldCreatedAt, v))
}

// HasProblem applies the HasEdge predicate on the "problem" edge.
func HasProblem() predicate.Message {
	return predicate.Message(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProblemTable, ProblemColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProblemWith applies the HasEdge predicate on the "problem" edge with a given conditions (other predicates).
func HasProblemWith(preds ...predicate.Problem) predicate.Message {
	return predicate.Message(func(s *sql.Selector) {
		step := newProblemStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasChat applies the HasEdge predicate on the "chat" edge.
func HasChat() predicate.Message {
	return predicate.Message(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ChatTable, ChatColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasChatWith applies the HasEdge predicate on the "chat" edge with a given conditions (other predicates).
func HasChatWith(preds ...predicate.Chat) predicate.Message {
	return predicate.Message(func(s *sql.Selector) {
		step := newChatStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Message) predicate.Message {
	return predicate.Message(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Message) predicate.Message {
	return predicate.Message(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Message) predicate.Message {
	return predicate.Message(sql.NotPredicates(p))
}
