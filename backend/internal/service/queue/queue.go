package queue

type QueueService interface {
	SetLinkStatus(int, string) error
}
