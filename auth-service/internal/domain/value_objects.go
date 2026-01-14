package domain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 32
	MinUserNameLength = 3
	MaxUserNameLength = 32
)

type UserId struct {
	value uuid.UUID
}

func NewUserId() UserId {
	return UserId{uuid.New()}
}

func UserIdFromUUID(id uuid.UUID) (UserId, error) {
	if id == uuid.Nil {
		return UserId{}, ErrUserIdRequired
	}
	return UserId{value: id}, nil
}

func UserIdFromString(s string) (UserId, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return UserId{}, fmt.Errorf("%w: %v", ErrUserIdRequired, err)
	}
	return UserIdFromUUID(id)
}

func (id UserId) UUID() uuid.UUID { return id.value }
func (id UserId) String() string  { return id.value.String() }

type UserName string

func NewUserName(raw string) (UserName, error) {
	v := strings.TrimSpace(raw)
	if len(v) > MaxUserNameLength {
		return "", ErrUserNameIsToLong
	}
	if len(v) < MinUserNameLength {
		return "", ErrUserNameIsToShort
	}
	return UserName(v), nil
}

func (u UserName) String() string { return string(u) }

type Email string

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`)

func NewEmail(raw string) (Email, error) {
	v := strings.TrimSpace(raw)
	if !emailRe.MatchString(v) {
		return "", ErrInvalidEmail
	}
	return Email(v), nil
}

func (e Email) String() string { return string(e) }

type Pwd string

func NewPwd(raw string) (Pwd, error) {
	v := strings.TrimSpace(raw)
	if len(v) > MaxPasswordLength {
		return "", ErrPwdIsToLong
	}
	if len(v) < MinPasswordLength {
		return "", ErrPwdIsToShort
	}
	return Pwd(v), nil
}

func (p Pwd) String() string { return string(p) }

type PwdHash string

func NewPwdHash(value string) (PwdHash, error) {
	return PwdHash(value), nil
}

func (p PwdHash) String() string { return string(p) }

type AccessToken string

func NewAccessToken(value string) (AccessToken, error) {
	return AccessToken(value), nil
}

func (p AccessToken) String() string { return string(p) }
