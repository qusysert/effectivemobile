package model

type UserInfo struct {
	Id         int    `json:"id" example:"1"`
	Name       string `json:"name" example:"Peter"`
	Surname    string `json:"surname" example:"Ivanovich"`
	Patronymic string `json:"patronymic" example:""`
	Age        int    `json:"age" example:"68"`
	Gender     string `json:"gender" example:"male"`
	Nation     string `json:"nation" example:"RU"`
}

type UserFilter struct {
	Id       int    `json:"id"`
	NameLike string `json:"name_like"`
	AgeFrom  int    `json:"age_from"`
	AgeTo    int    `json:"age_to"`
	Gender   string `json:"gender"`
	Nation   string `json:"nation"`
	PageSize int    `json:"page_size"`
	PageNum  int    `json:"page_num"`
}
