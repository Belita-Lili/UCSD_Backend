package auth

import (
	"context"
	"errors"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakAdapter struct {
	client       *gocloak.GoCloak
	realm        string
	clientID     string
	clientSecret string
}

func NewKeycloakAdapter(baseURL, realm, clientID, clientSecret string) *KeycloakAdapter {
	client := gocloak.NewClient(baseURL)
	return &KeycloakAdapter{
		client:       client,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (k *KeycloakAdapter) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	token, err := k.client.Login(ctx, k.clientID, k.clientSecret, k.realm, username, password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return token, nil
}

// Implementar métodos para OAuth (Google/Facebook), registro, recuperación, etc.
