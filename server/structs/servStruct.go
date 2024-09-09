package structs

import (
	userpb "awesomeProject4/user"
	"sync"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	mu     sync.Mutex
	users  map[string]*User
	nextID int32
}
