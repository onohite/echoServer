package db

type DBService interface {
	Close() error
	GetAllUsers() (*[]User, error)
}
