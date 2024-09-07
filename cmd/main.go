package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"sync"
	"time"

	userpb "awesomeProject4/user"
	"google.golang.org/grpc"
)

var jwtKey = []byte("your_secret_key")

type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

type server struct {
	userpb.UnimplementedUserServiceServer
	mu     sync.Mutex
	users  map[string]*User
	nextID int32
}
type Claims struct {
	jwt.StandardClaims
	ID int32 `json:"ID"`
}

func (s *server) GetAllUsers(ctx context.Context, _ *userpb.Empty) (*userpb.GetAllUsersResponse, error) {
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

func (s *server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
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

func generateJWT(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &server{
		users: make(map[string]*User), // Инициализация карты
	})

	fmt.Println("gRPC server is running on port :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
