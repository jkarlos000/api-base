package user

import (
	"context"
	"backend/internal/entity"
	"backend/internal/test"
	"backend/pkg/log"
	"testing"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "users")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	now := time.Now()
	roleID := "f5461314-7287-47f1-8c62-164e92353873"
	userID := "b464df2a-ea63-4369-adfc-f4c82e6eff40"
	err = repo.Create(ctx, entity.User{
		ID:        userID,
		FirstName: "Ilmar",
		LastName:  "López",
		Username:  "ilmarlopez",
		Password:  "$2a$04$bRTPCB6nl7ddsDoGDMdmxuMzmcd2NZhIjuusbj2JN1mBS4dKIZXem",
		RoleID:    roleID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: &now,
	})
	assert.Nil(t, err)

	// insert role of user
	_, err = db.With(ctx).Insert("role_user", dbx.Params{
		"role_id": roleID,
		"user_id": userID,
	}).Execute()
	assert.Nil(t, err)

	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

}
