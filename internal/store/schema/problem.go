package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"time"
)

// Problem holds the schema definition for the Problem entity.
type Problem struct {
	ent.Schema
}

// Fields of the Problem.
func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.ProblemID{}).Default(types.NewProblemID),
		field.UUID("manager_id", types.UserID{}),
		field.Time("resolved_at").Default(time.Now),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).Ref("problems").Unique().Required(),
		edge.To("messages", Message.Type),
	}
}
