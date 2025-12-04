package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

type Cache struct {
	config    *oauth2.Config
	token     *oauth2.Token
	tokenPath string
	mu        sync.Mutex
}

func New(clientID, clientSecret, redirectUrl, authUrl, tokenUrl, tokenFilePath string) *Cache {
	return &Cache{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectUrl,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authUrl,
				TokenURL: tokenUrl,
			},
		},
		tokenPath: tokenFilePath,
	}
}

// Config returns the oauth config
func (c *Cache) Config() *oauth2.Config {
	return c.config
}

// Client returns an http client with auto-refresh capability
func (c *Cache) Client(ctx context.Context) (*http.Client, error) {
	tok, err := c.Token(ctx)
	if err != nil {
		return nil, err
	}
	return c.config.Client(ctx, tok), nil
}

// Token returns a valid token, refreshing or loading from disk as needed
func (c *Cache) Token(ctx context.Context) (*oauth2.Token, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Used cached token if valid
	if c.token != nil && c.token.Valid() {
		return c.token, nil
	}

	// Try loading from disk
	if c.token == nil {
		if err := c.loadToken(); err != nil {
			return nil, err
		}
		if c.token != nil && c.token.Valid() {
			return c.token, nil
		}
	}

	// Refresh token if needed and possible
	if c.token != nil && !c.token.Valid() && c.token.RefreshToken != "" {
		newTok, err := c.refresh(ctx)
		if err != nil {
			return nil, err
		}
		return newTok, nil
	}

	return nil, errors.New("no valid token available — user must authenticate")
}

// ExchangeCode exchanges the oauth "code" for an access and refresh token
func (c *Cache) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	tok, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	c.token = tok
	if err := c.saveToken(tok); err != nil {
		return nil, err
	}
	return tok, nil
}

func (c *Cache) refresh(ctx context.Context) (*oauth2.Token, error) {
	ts := c.config.TokenSource(ctx, c.token)
	newTok, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("refresh token failed: %w", err)
	}
	c.token = newTok
	return newTok, c.saveToken(newTok)
}

func (c *Cache) saveToken(tok *oauth2.Token) error {
	data, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.tokenPath, data, 0600)
}

func (c *Cache) loadToken() error {
	data, err := os.ReadFile(c.tokenPath)
	if err != nil {
		return err
	}

	var tok oauth2.Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return err
	}

	if tok.Expiry.IsZero() {
		tok.Expiry = time.Now().Add(-time.Hour)
	}

	c.token = &tok
	return nil
}
