package dto

import "control-backend/login-kit/model"

type UserDto struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

func ToUserDto(user *model.User) UserDto {
	return UserDto{
		Name: user.Name,
		Uid:  user.Uid,
	}
}
