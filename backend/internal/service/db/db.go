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

	// Game methods
	GetAllGames() (*Games, error)
	GetGameRanks(int) (*Ranks, error)
	FindGameProfile(int) (string, error)
	CreateGameProfile(string) (int, error)
	UpdateGameProfileDescription(string, string, int) error
	UpdateGameProfileContact(string, string, int) error
}
