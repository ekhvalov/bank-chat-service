package clientv1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

var stub = MessagesPage{Messages: []Message{
	{
		AuthorId:  types.NewUserID(),
		Body:      "Здравствуйте! Разберёмся.",
		CreatedAt: time.Now(),
		Id:        types.NewMessageID(),
	},
	{
		AuthorId:  types.MustParse[types.UserID]("fec01fe8-483b-4cad-a0f6-ad0d431b433f"),
		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
		CreatedAt: time.Now().Add(-time.Minute),
		Id:        types.NewMessageID(),
	},
}}

func (h Handlers) PostGetHistory(eCtx echo.Context, _ PostGetHistoryParams) error {
	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"data": stub,
	})
}
