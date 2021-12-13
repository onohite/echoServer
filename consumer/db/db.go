package db

type DBService interface {
	Close() error

	AddRemind(Remind) error
}
