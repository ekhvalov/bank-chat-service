package jobpayload

import (
	"errors"
	"fmt"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func Marshal(messageID types.MessageID) (string, error) {
	if messageID.IsZero() {
		return "", errors.New("empty message id")
	}
	return messageID.String(), nil
}

func Unmarshal(payload string) (types.MessageID, error) {
	msgID, err := types.Parse[types.MessageID](payload)
	if err != nil {
		return types.MessageIDNil, fmt.Errorf("parse message id: %v", err)
	}
	if msgID.IsZero() {
		return types.MessageIDNil, fmt.Errorf("empty id")
	}
	return msgID, nil
}
