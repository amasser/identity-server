package authn

import (
	"strconv"

	"github.com/keratin/authn-go/authn"
)

// Account is a user account managed by authn-server
type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Locked   bool   `json:"locked"`
	Deleted  bool   `json:"deleted"`
}

// Service provides access to authn-server
type Service interface {
	// ImportAccount creates a new account at authn-server
	ImportAccount(username, password string, locked bool) (int, error)

	// GetAccount returns an authn account
	GetAccount(id int) (Account, error)

	// LockAccount locks a user account
	LockAccount(accountID int) error

	// UnlockAccount unlocks a user account
	UnlockAccount(accountID int) error

	// ArchiveAccount archives a user account
	ArchiveAccount(id int) error
}

type service struct {
	cli *authn.Client
}

// NewService returns a new authn-service
func NewService(cfg Config) (Service, error) {
	cli, err := authn.NewClient(authn.Config{
		Issuer:         cfg.Issuer,
		PrivateBaseURL: cfg.PrivateBaseAddress,
		Audience:       cfg.Audience,
		Username:       cfg.Username,
		Password:       cfg.Password,
	})
	if err != nil {
		return nil, err
	}

	return &service{
		cli: cli,
	}, nil
}

func (s *service) ImportAccount(username, password string, locked bool) (int, error) {
	return s.cli.ImportAccount(username, password, locked)
}

func (s *service) LockAccount(accountID int) error {
	return s.cli.LockAccount(strconv.Itoa(accountID))
}

func (s *service) UnlockAccount(accountID int) error {
	return s.cli.UnlockAccount(strconv.Itoa(accountID))
}

func (s *service) ArchiveAccount(id int) error {
	return s.cli.ArchiveAccount(strconv.Itoa(id))
}

func (s *service) GetAccount(id int) (Account, error) {
	ac, err := s.cli.GetAccount(strconv.Itoa(id))
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:       ac.ID,
		Username: ac.Username,
		Locked:   ac.Locked,
		Deleted:  ac.Deleted,
	}, nil
}
