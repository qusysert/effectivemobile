package model

type UserInfo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Nation     string `json:"nation"`
}

type UserFilter struct {
	NameLike string `json:"name_like"`
	AgeFrom  int    `json:"age_from"`
	AgeTo    int    `json:"age_to"`
	Gender   string `json:"gender"`
	Nation   string `json:"nation"`
	PageSize int    `json:"page_size"`
	PageNum  int    `json:"page_num"`
}
