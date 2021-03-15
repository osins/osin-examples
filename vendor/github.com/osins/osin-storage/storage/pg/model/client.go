package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Client define
type Client struct {
	Id          uuid.UUID `gorm:"primaryKey;->;<-:create;"`
	Secret      string
	RedirectUri string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	NeedLogin   bool
	NeedRefresh bool
}

// GetId method define
func (d *Client) GetId() string {
	return d.Id.String()
}

func (d *Client) GetNeedLogin() bool {
	return d.NeedLogin
}

// GetSecret method define
func (d *Client) GetSecret() string {

	return d.Secret
}

// GetRedirectUri method define
func (d *Client) GetRedirectUri() string {

	return d.RedirectUri
}

func (d *Client) GetNeedRefresh() bool {
	return d.NeedRefresh
}

// Implement the ClientSecretMatcher interface
// ClientSecretMatches method define
func (d *Client) ClientSecretMatches(secret string) bool {

	return d.Secret == secret
}
