package db

type DBService interface {
	Close() error

	// User methods
	AddUser(User) (string, error)
	GetUser(string) (*PublicUser, error)
	UpdateUserUserName(string, string) error
	UpdateUserEmail(string, string) error
	UpdateUserAvatar(string, string) error
	UpdateUserSex(int, string) error
	UpdateUserBdate(string, string) error

	GetListReminds(string) ([]Remind, error)
	UpdateRemindTo(string, string) error
	UpdateRemindMessage(string, string) error
	UpdateRemindDate(string, string) error
}
