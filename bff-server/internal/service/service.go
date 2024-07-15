package service

import (
	"bff-server/internal/entity"
	"bff-server/internal/server"
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Repository interface {
	GetUser(context.Context) (*entity.User, error)
}

type service struct {
	c         *http.Client
	UploadUri string
	r         Repository
}

func (s *service) GetUser(ctx context.Context) (*entity.Data[entity.User], error) {
	u, err := s.r.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if errs, ok := entity.Validate(u); !ok {
		return nil, errors.New(strings.Join(errs, ","))
	}
	return &entity.Data[entity.User]{
		Data: *u,
	}, nil
}

func (s *service) UploadFile(ctx context.Context, r io.Reader) error {
	u, err := s.r.GetUser(ctx)
	if err != nil {
		return err
	}
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multipartWriter.CreateFormFile("file", "file.csv")
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return err
	}
	err = multipartWriter.WriteField("email", u.Email)
	if err != nil {
		return err
	}
	err = multipartWriter.Close()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.UploadUri, &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	resp, err := s.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(raw))
	}
	return nil
}

func New(r Repository, u string, c *http.Client) server.Service {
	return &service{
		UploadUri: u,
		c:         c,
		r:         r,
	}
}
