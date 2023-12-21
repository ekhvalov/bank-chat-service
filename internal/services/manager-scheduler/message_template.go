package managerscheduler

import (
	"fmt"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const managerAssignedMsgTemplate = "Manager %s will answer you"

func ManagerAssignedMessageText(managerID types.UserID) string {
	return fmt.Sprintf(managerAssignedMsgTemplate, managerID)
}
