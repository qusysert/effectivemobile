package test

import (
	"effectivemobile/internal/app/model"
	"math/rand"
	"strings"
	"testing"
)

func TestApp(t *testing.T) {
	env, rollback := NewEnv()
	defer rollback()

	rand.Seed(0)

	users := []model.UserInfo{
		{Name: "John", Surname: "Doe", Patronymic: "Smith"},
		{Name: "Emma", Surname: "Johnson", Patronymic: "Brown"},
		{Name: "Michael", Surname: "Williams", Patronymic: "Davis"},
		{Name: "Olivia", Surname: "Taylor", Patronymic: "Anderson"},
	}

	for _, user := range users {
		_, err := env.Srv.AddUser(env.Ctx, user)
		if err != nil {
			t.Errorf("cannot add user: %v", err)
		}
	}

	user0, err := env.Srv.GetUser(env.Ctx, model.UserFilter{NameLike: users[0].Name + " " + users[0].Surname})
	if err != nil {
		t.Errorf("cannot get user0: %v", err)
	}
	if user0[0].Name != users[0].Name {
		t.Errorf("wrong name while getting user0 with nameLike")
	}

	user0Age := user0[0].Age
	ageFiltered, err := env.Srv.GetUser(env.Ctx, model.UserFilter{AgeFrom: user0Age})
	if err != nil {
		t.Errorf("error getting users filtered by ageFrom from user0 age: %v", err)
	}
	if !(len(ageFiltered) >= 1) {
		t.Errorf("error getting users filtered by ageFrom from user0 age: got not enough users")
	}

	user1Feed, err := env.Srv.GetUser(env.Ctx, model.UserFilter{NameLike: users[1].Name + " " + users[1].Surname})
	user1NewName := getRandomName()

	err = env.Srv.UpdateUser(env.Ctx, model.UserInfo{Id: user1Feed[0].Id, Name: user1NewName})
	if err != nil {
		t.Errorf("error updating user1 name: %v", err)
	}

	user1UpdatedFeed, err := env.Srv.GetUser(env.Ctx, model.UserFilter{Id: user1Feed[0].Id})
	if err != nil {
		t.Errorf("error getting updated user 1: %v", err)
	}

	if user1UpdatedFeed[0].Name != user1NewName {
		t.Errorf("error updating user1 name: %v", user1UpdatedFeed)
	}

	paginationFeed, err := env.Srv.GetUser(env.Ctx, model.UserFilter{PageSize: 2, PageNum: 2})
	if err != nil {
		t.Errorf("error getting paginated feed: %v", err)
	}
	if len(paginationFeed) != 2 {
		t.Errorf("error getting paginated feed: wrong length")
	}

	_, err = env.Srv.GetUser(env.Ctx, model.UserFilter{PageSize: 0, PageNum: 2})
	if err == nil && !strings.Contains(err.Error(), "error getting user: wrong pagination options") {
		t.Errorf("error getting user: wrong pagination options worked")
	}
}

func getRandomName() string {
	names := []string{
		"John", "Emma", "Michael", "Olivia", "William",
		"Ava", "James", "Sophia", "Benjamin", "Isabella",
		"Jacob", "Mia", "Elijah", "Charlotte", "Alexander",
		"Amelia", "Daniel", "Harper", "Matthew", "Evelyn",
	}

	randomIndex := rand.Intn(len(names))
	return names[randomIndex]
}
