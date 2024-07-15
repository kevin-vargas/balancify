package local

import (
	"bff-server/internal/entity"
	"context"
)

type repository struct {
}

func (s *repository) GetUser(context.Context) (entity.User, error) {
	return user, nil
}
