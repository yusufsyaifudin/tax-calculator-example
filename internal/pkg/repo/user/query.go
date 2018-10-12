package user

import (
	"context"

	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/conn"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/model"
)

// Create will insert a new record in database.
func Create(parent context.Context, username, password string) (User *model.User, err error) {
	User = &model.User{}
	err = conn.GetDBConnection().Writer().Query(parent, User, sqlInsertUser, username, password)
	return
}

// FindByID will looking for user by primary id
func FindByID(parent context.Context, id int64) (User *model.User, err error) {
	User = &model.User{}
	err = conn.GetDBConnection().Reader().Query(parent, User, sqlFindUserByID, id)
	return
}

// FindByUsername will looking for user by username.
func FindByUsername(parent context.Context, username string) (User *model.User, err error) {
	User = &model.User{}
	err = conn.GetDBConnection().Reader().Query(parent, User, sqlFindUserByUsername, username)
	return
}
