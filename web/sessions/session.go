package sessions

import (
	"bytes"
	"encoding/gob"
	"net/http"

	"github.com/google/uuid"
	"github.com/moroz/securecookie"
)

type Store struct {
	securecookie.Store
}

func NewStore(key []byte) (*Store, error) {
	s, err := securecookie.NewStore(key)
	if err != nil {
		return nil, err
	}
	return &Store{s}, nil
}

type Payload map[string]any

func init() {
	gob.Register(&Payload{})
	gob.Register(&uuid.UUID{})
}

func (s Store) DecodeSession(cookie *http.Cookie) (Payload, error) {
	result := make(Payload)

	if cookie == nil {
		return result, nil
	}

	binary, err := s.DecryptCookie(cookie.Value)
	if err != nil {
		return result, err
	}
	err = gob.NewDecoder(bytes.NewBuffer(binary)).Decode(&result)
	return result, err
}

func (s Store) EncodeSession(v any) (string, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		return "", err
	}

	return s.EncryptCookie(buf.Bytes())
}
