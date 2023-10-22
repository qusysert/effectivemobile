package handler

import "effectivemobile/internal/app/model"

type User struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Nation     string `json:"nation"`
}

func toModelUser(user interface{}) model.UserInfo {
	switch u := user.(type) {
	case AddUserRequest:
		return model.UserInfo{
			Name:       u.Name,
			Surname:    u.Surname,
			Patronymic: u.Patronymic,
		}
	case UpdateUserRequest:
		return model.UserInfo{
			Name:       u.Name,
			Surname:    u.Surname,
			Patronymic: u.Patronymic,
		}
	default:

		panic("unsupported user type")
	}
}

func fromModelUser(userInfo model.UserInfo) User {
	return User{
		Name:       userInfo.Name,
		Surname:    userInfo.Surname,
		Patronymic: userInfo.Patronymic,
		Age:        userInfo.Age,
		Gender:     userInfo.Gender,
		Nation:     userInfo.Nation,
	}
}
