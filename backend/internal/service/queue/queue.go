package queue

import (
	"backend/internal/service/db"
)

type QueueService interface {
	//SetLinkStatus(int, string) error
	SetRemind(remind db.Remind) error
}
