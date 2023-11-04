package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

// jobMaxAttempts is some limit as protection from endless retries of outbox jobs.
const jobMaxAttempts = 30

type Job struct {
	ent.Schema
}

func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.JobID{}).Default(types.NewJobID).Unique().Immutable(),

		// FIXME: Имя джобы (name), не может быть пустым, заполняется единожды.
		field.String("name").NotEmpty().Immutable(),

		// FIXME: Данные для выполнения (payload), не могут быть пустыми, заполняется единожды.
		field.String("payload").NotEmpty().Immutable(),

		// FIXME: Количество попыток выполнения (attempts), по умолчанию 0, входит в [0, N], где N – какой-то лимит.
		field.Int("attempts").Default(0).Max(jobMaxAttempts),

		// FIXME: Время, когда можно выполнять джобу (available_at). Полезно для отложенных задач, заполняется единожды.
		field.Time("available_at").Immutable(),

		// FIXME: Время, которое даётся на выполнение задачи (reserved_until). Пока не наступило это время,
		// FIXME: другие горутины не могут взять данную задачу. Когда одна из горутин берёт задачу, она выставляет
		// FIXME: значение этого поля в <time.Now() + some timeout>, как бы резервируя её под себя.
		field.Time("reserved_until").Default(time.Time{}),

		// FIXME: Время создания задачи (created_at), заполняется единожды, по умолчанию time.Now().
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Job) Indexes() []ent.Index {
	// FIXME: Расставь индексы на основе запросов в сервисе Outbox.
	return nil
}

type FailedJob struct {
	ent.Schema
}

func (FailedJob) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.FailedJobID{}).Default(types.NewFailedJobID).Unique().Immutable(),

		// FIXME: Имя задачи (name) – строка, заполняется единожды, не может быть пустым.
		field.String("name").NotEmpty().Immutable(),

		// FIXME: Данные для выполнения (payload) – строка, заполняется единожды, не может быть пустым.
		field.String("payload").NotEmpty().Immutable(),

		// FIXME: Причина неудачи (reason) – строка, заполняется единожды, не может быть пустым.
		field.String("reason").NotEmpty().Immutable(),

		// FIXME: Время создания задачи (created_at), заполняется единожды, по умолчанию time.Now().
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
