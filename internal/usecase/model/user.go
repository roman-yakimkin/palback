package model

import (
	"palback/internal/domain/model"
	"palback/internal/pkg/helpers"
)

type UserDetail struct {
	model.User
	Role model.Role
}

func CreateUserDetail(user model.User, role model.Role) UserDetail {
	return UserDetail{
		User: user,
		Role: role,
	}
}

type UserList struct {
	Items []UserDetail
}

func CreateUserList(users []model.User, roles map[model.RoleID]*model.Role) (result UserList) {
	result.Items = make([]UserDetail, 0, len(users))
	for _, user := range users {
		result.Items = append(result.Items, CreateUserDetail(user, helpers.FromPtr(roles[user.RoleID])))
	}

	return
}
