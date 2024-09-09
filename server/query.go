package server

import (
	userpb "awesomeProject4/user"
	"context"
)

func (s *structs.server) GetAllUsers(ctx context.Context, _ *userpb.Empty) (*userpb.GetAllUsersResponse, error) {
	users := []*userpb.User{}
	for _, user := range s.users {
		users = append(users, &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Name:     "",
		})
	}

	return &userpb.GetAllUsersResponse{
		Users: users,
	}, nil
}
