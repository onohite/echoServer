package db

type DBService interface {
	Close() error

	//Users
	GetAllUsers() (*[]User, error)

	//Links
	GetAllLinks() (*[]Link, error)
	GetLinkById(id string) (*Link, error)
	AddLink(link Link) error
}
