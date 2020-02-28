package session

import (
	"io/ioutil"
	"os"
)

// TokenLoader can load an access or refresh token.
type TokenLoader interface {
	// Load loads the token and returns it.
	Load() (string, error)
}

// TokenStorer can store a refresh or access token for later usage.
type TokenStorer interface {
	// Store stores the token.
	Store(token string) error
}

// TokenStore is responsible for storing either refresh or access tokens.
type TokenStore interface {
	TokenStorer
	TokenLoader
}

// FileTokenStore implements a TokenStore that persists the token
// in the local file system.
type FileTokenStore struct {
	path string
}

// Store implements TokenStoreer and persists the token in the
// local filesystem.
func (fs *FileTokenStore) Store(token string) error {
	return ioutil.WriteFile(fs.path, []byte(token), 0600)
}

// Load implements TokenLoader and loads the token from the
// local filesystem.
func (fs *FileTokenStore) Load() (string, error) {
	content, err := ioutil.ReadFile(fs.path)
	return string(content), err
}

// EnvStore is a TokenLoader and TokenStorer that uses an environment
// variable to retrieve and store an access or refresh token.
type EnvStore struct {
	key string
}

// NewEnvLoader returns a new token loader that uses
// the given environment key.
func NewEnvLoader(key string) *EnvStore {
	return &EnvStore{key}
}

// Load implements the TokenLoader interface and returns the value
// of the environment variable configures during NewEnvLoader()
func (env *EnvStore) Load() (string, error) {
	value := os.Getenv(env.key)
	if value == "" {
		return "", os.ErrNotExist
	}

	return value, nil
}

// Store implements the TokenStorer interface and stores the given
// token in an environment variable.
func (env *EnvStore) Store(token string) error {
	os.Setenv(env.key, token)
	return nil
}
