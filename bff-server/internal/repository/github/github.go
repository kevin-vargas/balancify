package github

import (
	"bff-server/internal/entity"
	"bff-server/internal/middleware"
	"bff-server/internal/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type userEmailResponse struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

type repository struct {
	base string
}

func (s *repository) GetUser(ctx context.Context) (*entity.User, error) {
	t, ok := middleware.GetToken(ctx)
	if !ok {
		return nil, errors.New("invalid token")
	}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user", s.base), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d", res.StatusCode)
	}

	type user struct {
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
	}
	var u user
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		return nil, err
	}
	eu := entity.User(u)

	if eu.Email == "" {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user/emails", s.base), nil)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("invalid status code %d", res.StatusCode)
		}
		emails := []userEmailResponse{}
		if err := json.NewDecoder(res.Body).Decode(&emails); err != nil {
			return nil, err
		}
		for _, email := range emails {
			if email.Primary {
				eu.Email = email.Email
				return &eu, nil
			}
		}
		return nil, errors.New("not primary email on github")
	}
	return &eu, nil

}

func New(base string) service.Repository {
	return &repository{
		base: base,
	}
}
