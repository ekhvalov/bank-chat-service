package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	"time"
)

// bodySizeLimit 3000 unicode symbols of 4 bytes each, and rounded up to kilobyte.
const bodySizeLimit = 1024 * 12

var timeNil = func() time.Time { return time.Unix(0, 0) }

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.NewMessageID()).Default(types.NewMessageID),
		field.UUID("author_id", types.UserID{}).Optional(),
		field.Bool("is_visible_for_client").Default(false),
		field.Bool("is_visible_for_manager").Default(false),
		field.Text("body").MaxLen(bodySizeLimit).Immutable(),
		field.Time("checked_at").Default(timeNil),
		field.Bool("is_blocked"),
		field.Bool("is_service").Immutable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem", Problem.Type).
			Ref("messages").
			Unique().
			Required(),
		edge.From("chat", Chat.Type).
			Ref("messages").
			Unique().
			Required(),
	}
}
