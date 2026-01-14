package domain

import (
	"time"
)

type User struct {
	id        UserId
	email     Email
	username  UserName
	pwdHash   PwdHash
	createdAt time.Time
}

func NewUser(email Email, userName UserName, pwdHash PwdHash) *User {
	now := time.Now().UTC()
	u := &User{
		id:        NewUserId(),
		email:     email,
		username:  userName,
		pwdHash:   pwdHash,
		createdAt: now,
	}

	return u
}

func RehydrateUser(
	id UserId,
	email Email,
	userName UserName,
	pwdHash PwdHash,
	createdAt time.Time,
) *User {
	return &User{
		id:        id,
		email:     email,
		username:  userName,
		pwdHash:   pwdHash,
		createdAt: createdAt,
	}
}

func (u *User) Id() UserId           { return u.id }
func (u *User) Email() Email         { return u.email }
func (u *User) Username() UserName   { return u.username }
func (u *User) PwdHash() PwdHash     { return u.pwdHash }
func (u *User) CreatedAt() time.Time { return u.createdAt }
