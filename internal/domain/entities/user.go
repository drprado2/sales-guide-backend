package entities

import (
	"github.com/drprado2/sales-guide/internal/domain/errors"
	"github.com/drprado2/sales-guide/internal/domain/valueobjects"
	"strings"
	"time"
)

const (
	DefaultAvatarImage = "https://w7.pngwing.com/pngs/340/946/png-transparent-avatar-user-computer-icons-software-developer-avatar-child-face-heroes.png"
	UserHiddenPassword = "********"
)

type (
	User struct {
		ID                  string
		CompanyID           string
		Name                string
		Email               string
		Phone               string
		BirthDate           *time.Time
		AvatarImage         *valueobjects.Uri
		RecordCreationCount int
		RecordEditingCount  int
		RecordDeletionCount int
		LastAccess          *time.Time
		CreatedAt           time.Time
		UpdatedAt           time.Time
		RowVersion          uint32
	}
)

func NewUser(id string, companyId string, name string, email string) *User {
	defaultAvatar, _ := valueobjects.NewUri(DefaultAvatarImage)

	return &User{
		ID:                  id,
		CompanyID:           strings.ToLower(companyId),
		Name:                strings.ToLower(name),
		Email:               strings.ToLower(email),
		AvatarImage:         defaultAvatar,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
		RecordCreationCount: 0,
		RecordEditingCount:  0,
		RecordDeletionCount: 0,
	}
}

func (u *User) Validate() error {
	if u.ID == "" ||
		u.Name == "" ||
		u.Email == "" ||
		u.AvatarImage == nil {
		return errors.CompanyInvalidToSaveError
	}
	return nil
}
