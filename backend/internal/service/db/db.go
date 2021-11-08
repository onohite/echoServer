package db

type DBService interface {
	Close() error

	//Users
	GetAllUsers() (*[]User, error)

	//Links
	GetAllLinks() (*[]ResponseLink, error)
	GetLinkById(id string) (*ResponseLink, error)
	AddLink(link Link) (int, error)
	// UpdateStatusLink status - status of url , id - privary key
	UpdateStatusLink(int, string) error
}
