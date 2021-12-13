package db

type DBService interface {
	Close() error

	GetReminds() ([]Remind, error)
	UpdateStatusReminds(string) error
}
