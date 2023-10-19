package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"time"
)

// bodySizeLimit 3000 unicode symbols of 4 bytes each, and rounded up to kilobyte.
const bodySizeLimit = 1024 * 12

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).Default(types.NewMessageID).Unique().Immutable(),
		field.UUID("chat_id", types.ChatID{}),
		field.UUID("problem_id", types.ProblemID{}),
		field.UUID("author_id", types.UserID{}).Optional().Immutable(),
		field.Bool("is_visible_for_client").Default(false),
		field.Bool("is_visible_for_manager").Default(false),
		field.Text("body").NotEmpty().MaxLen(bodySizeLimit).Immutable(),
		field.Time("checked_at").Optional(),
		field.Bool("is_blocked").Default(false),
		field.Bool("is_service").Default(false).Immutable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem", Problem.Type).
			Ref("messages").
			Field("problem_id").
			Unique().
			Required(),
		edge.From("chat", Chat.Type).
			Ref("messages").
			Field("chat_id").
			Unique().
			Required(),
	}
}

func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("chat_id"),
		index.Fields("author_id"),
		index.Fields("created_at"),
		index.Fields("is_visible_for_client"),
	}
}
