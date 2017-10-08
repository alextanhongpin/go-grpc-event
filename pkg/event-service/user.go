package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// UserInfo represents the schema from the auth0 userinfo endpoint
type UserInfo struct {
	Email    string `json:"email"`    // "test.account@userinfo.com"
	Name     string `json:"name"`     //  "test.account@userinfo.com"
	Picture  string `json:"picture"`  // "https://s.gravatar.com/avatar/dummy.png"
	UserID   string `json:"user_id"`  // "auth0|58454..."
	Nickname string `json:"nickname"` // "test.account"
	Sub      string `json:"sub"`      // "auth0|58454..."
	Admin    bool   `json:"-"`        // false
}

// IsAuthorized checks if the user is authorized
func (u UserInfo) IsAuthorized() bool {
	return u.UserID != ""
}

// IsAdmin checks if the user is admin
func (u UserInfo) IsAdmin() bool {
	return u.Admin
}

// Extract returns the user metadata from the context
func (u *UserInfo) Extract(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	meta := Metadata{md}
	u.Email = meta.Get("email")
	u.Name = meta.Get("name")
	u.Picture = meta.Get("picture")
	u.UserID = meta.Get("userid")
	u.Nickname = meta.Get("nickname")
	u.Sub = meta.Get("sub")
	u.Admin = meta.Get("admin") != ""
}
