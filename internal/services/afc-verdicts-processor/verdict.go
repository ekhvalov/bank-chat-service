package afcverdictsprocessor

import (
	"github.com/golang-jwt/jwt"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

type Status string

const (
	StatusOK         = "ok"
	StatusSuspicious = "suspicious"
)

type verdict struct {
	jwt.MapClaims
	ChatID    types.ChatID    `json:"chatId" validate:"required"`
	MessageID types.MessageID `json:"messageId" validate:"required"`
	Status    Status          `json:"status" validate:"required"`
}
