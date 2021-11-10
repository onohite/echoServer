package db

type DBService interface {
	Close() error

	AddUser(User) (string, error)
}
