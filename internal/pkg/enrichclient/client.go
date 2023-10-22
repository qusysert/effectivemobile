package enrichclient

import (
	"effectivemobile/internal/app/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EnrichClient struct {
	AgifyApi       string
	GenderizeApi   string
	NationalizeApi string
}

func New(agifyApi, genderizeApi, nationalizeApi string) *EnrichClient {
	return &EnrichClient{AgifyApi: agifyApi, GenderizeApi: genderizeApi, NationalizeApi: nationalizeApi}
}

func (ec *EnrichClient) Enrich(user model.UserInfo) (model.UserInfo, error) {
	var err error

	user, err = ec.Agify(user, ec.AgifyApi)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("cannot agify user %v: %w", user, err)
	}
	user, err = ec.Genderize(user, ec.GenderizeApi)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("cannot genderize user %v: %w", user, err)
	}
	user, err = ec.Nationalize(user, ec.NationalizeApi)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("cannot nationalyze user %v: %w", user, err)
	}
	return user, err
}

func (ec *EnrichClient) Agify(user model.UserInfo, apiUrl string) (model.UserInfo, error) {
	var data agifyResponse
	user, err := getInfo(user, apiUrl, &data)
	if err != nil {
		return model.UserInfo{}, err
	}

	user.Age = data.Age
	return user, nil
}

func (ec *EnrichClient) Genderize(user model.UserInfo, apiUrl string) (model.UserInfo, error) {
	var data genderizeResponse
	user, err := getInfo(user, apiUrl, &data)
	if err != nil {
		return model.UserInfo{}, err
	}

	user.Gender = data.Gender
	return user, nil
}

func (ec *EnrichClient) Nationalize(user model.UserInfo, apiUrl string) (model.UserInfo, error) {
	var data nationalizeResponse
	user, err := getInfo(user, apiUrl, &data)
	if err != nil {
		return model.UserInfo{}, err
	}

	user.Nation = data.Country[0].CountryId
	return user, nil
}

func getInfo(user model.UserInfo, apiUrl string, data interface{}) (model.UserInfo, error) {
	response, err := http.Get(apiUrl + user.Name)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("cannot get data from URL: %w", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("cannot read JSON from body: %w", err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("cannot unmarshal JSON to struct: %w", err)
	}

	return user, nil
}
