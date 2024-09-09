package server

import (
	userpb "awesomeProject4/user"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *structs.Server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("CreateUser gRPC method invoked")

	if _, exists := s.users[req.GetEmail()]; exists {
		return nil, errors.New("Email already exists")
	}

	userID := uuid.New().String()
	// Создание пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Ошибка шифрования пароля Blehh", err)
	}
	user := &User{
		ID:       userID,
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: string(hashedPassword),
	}
	s.users[req.GetEmail()] = user

	return &userpb.CreateUserResponse{
		Message: "User created successfully!",
		User: &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Name:     "", // Пока что не используется, можно оставить пустым
		},
	}, nil
}

func (s *server) LoginUser(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("LoginUser gRPC method invoked")

	user, exists := s.users[req.Email]
	if !exists {
		return nil, errors.New("Пользователь не найден")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("Неверный адрес почты или пароль")
	}

	token, err := generateJWT(*user)
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		Token:   token,
		Message: "Вход выполнен успешно !",
	}, nil
}
