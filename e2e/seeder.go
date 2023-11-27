package e2e

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
)

func CreateUserModels(size int) []*model.User {
	users := make([]*model.User, size)

	for i := 0; i < size; i++ {
		users[i] = CreateUserModel(i, types.Registered)
	}

	return users
}

func CreateUserModel(i int, role types.Role) *model.User {
	userID := i

	date := time.Now().Add(-6 * time.Hour).Add(time.Duration(i) * time.Minute)

	ep, err := apputils.EncryptPassword("12341234")
	if err != nil {
		return nil
	}

	return &model.User{
		BaseModel: model.BaseModel{
			ID:        uint64(userID),
			CreatedAt: date,
			UpdatedAt: date,
			DeletedAt: nil,
			CreatedBy: strconv.Itoa(userID),
			UpdatedBy: strconv.Itoa(userID),
			Status:    types.Active,
		},
		Email:             fmt.Sprintf("%v@example.com", i),
		Role:              role,
		Username:          fmt.Sprintf("username-%v", i),
		EncryptedPassword: ep,
	}
}
