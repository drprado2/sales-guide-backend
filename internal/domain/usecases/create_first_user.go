package usecases

import (
	"context"
	"errors"
	"github.com/drprado2/react-redux-typescript/configs"
	"github.com/drprado2/react-redux-typescript/internal/domain"
	"github.com/drprado2/react-redux-typescript/internal/domain/entities"
	"github.com/drprado2/react-redux-typescript/internal/domain/valueobjects"
	"github.com/drprado2/react-redux-typescript/pkg/ctxvals"
	"github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/drprado2/react-redux-typescript/pkg/pointers"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gopkg.in/auth0.v5/management"
	"strings"
	"time"
)

type (
	CreateFirstUserInput struct {
		CompanyId   string     `json:"companyId" validate:"required,min=1" name:"empresa"`
		Name        string     `json:"name" validate:"required,min=1" name:"nome"`
		Email       string     `json:"email" validate:"required,email" name:"e-mail"`
		Phone       string     `json:"phone" validate:"numeric" name:"telefone"`
		BirthDate   *time.Time `json:"birthDate" name:"data de nascimento"`
		Password    string     `json:"password" validate:"required,min=8" name:"senha"`
		AvatarImage string     `json:"avatarImage" name:"foto"`
	}

	CreateFirstUserOutput struct {
		ID                  string     `json:"id"`
		CompanyId           string     `json:"companyId"`
		Name                string     `json:"name"`
		Email               string     `json:"email"`
		Phone               string     `json:"phone"`
		BirthDate           *time.Time `json:"birthDate"`
		Password            string     `json:"password"`
		AvatarImage         string     `json:"avatarImage"`
		RecordCreationCount int        `json:"recordCreationCount"`
		RecordEditingCount  int        `json:"recordEditingCount"`
		RecordDeletionCount int        `json:"recordDeletionCount"`
		LastAccess          *time.Time `json:"lastAccess"`
		CreatedAt           time.Time  `json:"createdAt"`
		UpdatedAt           time.Time  `json:"updatedAt"`
		RowVersion          uint32     `json:"rowVersion"`
	}
)

func CreateFirstUser(ctx context.Context, input *CreateFirstUserInput) (*CreateFirstUserOutput, error) {
	span, ctx := tracer.SpanFromContext(ctx)
	defer span.Finish()

	location := ctxvals.LocationOrDefault(ctx)
	timeoffset := ctxvals.TimeOffsetOrDefault(ctx)

	if err := payloadValidator.Struct(input); err != nil {
		return nil, domain.PayloadErrorFromValidator(err, validatorTranslates)
	}

	errGroup := &domain.ErrorGroup{}

	user := entities.NewUser(uuid.NewString(), input.CompanyId, input.Name, input.Email)
	if input.AvatarImage != "" {
		avatarUri, err := valueobjects.NewUri(input.AvatarImage)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating user with invalid avatar URI")
			errGroup.Append(domain.NewConstraintError(errors.New("avatar inválido")))
		}
		user.AvatarImage = avatarUri
	}
	if input.BirthDate != nil {
		d := input.BirthDate.Add(time.Hour * time.Duration(timeoffset))
		user.BirthDate = &d
	}
	user.Phone = input.Phone

	if err := user.ValidToSave(); err != nil {
		logs.Logger(ctx).WithError(err).Warn("user invalid to save")
		errGroup.Append(domain.NewConstraintError(err))
	}

	if errGroup.HasError() {
		return nil, errGroup.AsError()
	}

	tx, err := userRepository.BeginTx(ctx)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	rowversion, err := userRepository.CreateTx(ctx, tx, user)
	if err != nil {
		tx.Rollback(ctx)
		if err, ok := err.(*pq.Error); ok {
			if err.Code == domain.PgUniqueConstraintErrorCode {
				logs.Logger(ctx).WithError(err).Error("fail creating user with unique constraint error")
				return nil, domain.NewConstraintError(errors.New("já existe uma usuário com esse e-mail"))
			}
			if err.Code == domain.PgForeignErrorCode {
				logs.Logger(ctx).WithError(err).Error("fail creating user with foreign key error")
				return nil, domain.NewConstraintError(errors.New("empresa inválida"))
			}
		}
		logs.Logger(ctx).WithError(err).Error("fail creating company with error in repository")
		return nil, domain.NewInternalError(err)
	}

	nameParts := strings.Split(user.Name, " ")
	emailSeparator := strings.IndexRune(user.Email, '@')
	nickname := user.Email[0:emailSeparator]
	authUser := &management.User{
		UserMetadata: map[string]interface{}{"company_id": user.CompanyID},
		ID:           pointers.SafeString(user.ID),
		Email:        &user.Email,
		Name:         &user.Name,
		GivenName:    &nameParts[0],
		FamilyName:   &nameParts[len(nameParts)-1],
		Nickname:     &nickname,
		Password:     &input.Password,
		VerifyEmail:  pointers.Bool(configs.Get().Auth0VerifyEmail),
		Picture:      pointers.String(user.AvatarImage.AsString()),
		Blocked:      pointers.Bool(false),
		Connection:   pointers.String("Username-Password-Authentication"),
	}
	if err := auth0manager.User.Create(authUser); err != nil {
		tx.Rollback(ctx)
		logs.Logger(ctx).WithError(err).Error("fail creating auth0 user")
		return nil, domain.NewInternalError(err)
	}

	tx.Commit(ctx)

	return &CreateFirstUserOutput{
		ID:                  user.ID,
		CompanyId:           user.CompanyID,
		Name:                user.Name,
		Email:               user.Email,
		Phone:               user.Phone,
		BirthDate:           input.BirthDate,
		Password:            entities.UserHiddenPassword,
		AvatarImage:         user.AvatarImage.AsString(),
		RecordCreationCount: user.RecordCreationCount,
		RecordEditingCount:  user.RecordEditingCount,
		RecordDeletionCount: user.RecordDeletionCount,
		LastAccess:          user.LastAccess,
		CreatedAt:           user.CreatedAt.In(location),
		UpdatedAt:           user.UpdatedAt.In(location),
		RowVersion:          rowversion,
	}, nil
}
