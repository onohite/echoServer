package graph

type GraphService interface {
	Close() error
	SetProfile(GameProfile) (string, error)
	GetProfile(string) (*GameProfile, error)
}
