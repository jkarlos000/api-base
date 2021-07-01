package nurse

import (
	"context"
	"backend/internal/entity"
	"backend/pkg/log"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// Service encapsulates usecase logic for nurses.
type Service interface {
	Get(ctx context.Context, id string) (Nurse, error)
	Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]Nurse, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateNurseRequest) (Nurse, error)
	Update(ctx context.Context, id string, input UpdateNurseRequest) (Nurse, error)
	Delete(ctx context.Context, id string) (Nurse, error)
}

// Nurse represents the data about an nurse.
type Nurse struct {
	entity.Nurse
}

// CreateNurseRequest represents an user creation request.
type CreateNurseRequest struct {
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Username             string  `json:"username"`
	Password             string  `json:"password"`
	Latitude 			*float32 `json:"latitude"`
	Longitude 			*float32 `json:"longitude"`
}

// Validate validates the CreateNurseRequest fields.
func (m CreateNurseRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.LastName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Username, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Latitude, validation.Required),
		validation.Field(&m.Longitude, validation.Required),
	)
}
// UpdateNurseRequest represents an user update request.
type UpdateNurseRequest struct {
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Username             string  `json:"username"`
	Password             *string `json:"password"`
	Latitude 			*float32 `json:"latitude"`
	Longitude 			*float32 `json:"longitude"`
	Rol					entity.Role `json:"rol"`
	IsActive			*bool `json:"is_active"`
}

// Validate validates the CreateUserRequest fields.
func (m UpdateNurseRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.LastName, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Username, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Latitude, validation.Required),
		validation.Field(&m.Longitude, validation.Required),
		validation.Field(&m.Rol, validation.Required),
		validation.Field(&m.IsActive, validation.Required),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new nurse service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the nurse with the specified the nurse ID.
func (s service) Get(ctx context.Context, id string) (Nurse, error) {
	nurse, err := s.repo.Get(ctx, id)
	if err != nil {
		return Nurse{}, err
	}
	return Nurse{nurse}, nil
}

// Create creates a new nurse.
func (s service) Create(ctx context.Context, req CreateNurseRequest) (Nurse, error) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Username = strings.TrimSpace(req.Username)

	nurseID := entity.GenerateID()
	userID := entity.GenerateID()

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return Nurse{}, err
	}

	if err := req.Validate(); err != nil {
		return Nurse{}, err
	}
	now := time.Now()
	err = s.repo.Create(ctx, entity.User{
		ID:        userID,
		Username:  req.Username,
		Password:  string(hash),
		IsActive:  true,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: now,
		UpdatedAt: &now,
	}, entity.Nurse{
		ID:             nurseID,
		UserID:         userID,
		HaveFullAccess: false,
		IsActive:       true,
		IsWorking:      false,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
	})
	if err != nil {
		return Nurse{}, err
	}
	return s.Get(ctx, nurseID)
}

// Update updates the nurse with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateNurseRequest) (Nurse, error) {
	var user entity.User
	excludes := []string{"Roles", "Permissions", "CreatedAt"}

	if err := req.Validate(); err != nil {
		return Nurse{}, err
	}

	now := time.Now()

	res, err := s.Get(ctx, id)
	if err != nil {
		return res, err
	}
	nurse := res.Nurse
	nurse.Latitude = req.Latitude
	nurse.Longitude = req.Longitude

	user.ID = nurse.UserID
	user.Username = req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.UpdatedAt = &now

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.MinCost)
		if err != nil {
			fmt.Println(err)
		}

		user.Password = string(hash)
	} else {
		excludes = append(excludes, "Password")
	}

	if req.IsActive != nil && *req.IsActive == false {
		user.IsActive = false
		res.IsActive = false
	} else {
		user.IsActive = true
		res.IsActive = true
	}
	// Fixme refactor repo.Update - remove user and excludes
	if err := s.repo.Update(ctx, user, nurse, excludes); err != nil {
		return res, err
	}

	res.Username = req.Username
	res.FirstName = req.FirstName
	res.LastName = req.LastName
	res.Latitude = req.Latitude
	res.Longitude = req.Longitude

	return res, nil
}

// Delete deletes the nurse with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Nurse, error) {
	nurse, err := s.Get(ctx, id)
	if err != nil {
		return Nurse{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Nurse{}, err
	}
	return nurse, nil
}

// Count returns the number of nurses.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the nurses with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int, term string, filters map[string]interface{}) ([]Nurse, error) {
	items, err := s.repo.Query(ctx, offset, limit, term, filters)
	if err != nil {
		return nil, err
	}
	result := []Nurse{}
	for _, item := range items {
		result = append(result, Nurse{item})
	}
	return result, nil
}
