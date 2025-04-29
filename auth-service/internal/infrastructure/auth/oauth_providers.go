package auth

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type OAuthProvider interface {
	AuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (map[string]interface{}, error)
}

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (g *GoogleProvider) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GoogleProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

// Implementar GetUserInfo...

type FacebookProvider struct {
	config *oauth2.Config
}

func NewFacebookProvider(clientID, clientSecret, redirectURL string) *FacebookProvider {
	return &FacebookProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email"},
			Endpoint:     facebook.Endpoint,
		},
	}
}

// Implementar métodos para Facebook...
