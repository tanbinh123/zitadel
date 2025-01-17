package idp

import (
	"context"

	"golang.org/x/text/language"
)

// Provider is the minimal implementation for a 3rd party authentication provider
type Provider interface {
	Name() string
	BeginAuth(ctx context.Context, state string, params ...any) (Session, error)
	IsLinkingAllowed() bool
	IsCreationAllowed() bool
	IsAutoCreation() bool
	IsAutoUpdate() bool
}

// User contains the information of a federated user.
type User interface {
	GetID() string
	GetFirstName() string
	GetLastName() string
	GetDisplayName() string
	GetNickname() string
	GetPreferredUsername() string
	GetEmail() string
	IsEmailVerified() bool
	GetPhone() string
	IsPhoneVerified() bool
	GetPreferredLanguage() language.Tag
	GetAvatarURL() string
	GetProfile() string
}
