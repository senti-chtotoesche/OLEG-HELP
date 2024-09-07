// client.go
package main

import (
	userpb "awesomeProject4/user"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	createResp, err := client.CreateUser(context.Background(), &userpb.CreateUserRequest{
		Username: "Jane_Doe",
		Email:    "janeRat@.mail",
		Password: "ABOBA",
	})
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	fmt.Println("CreateUserResponse: %v\n", createResp)

	loginResp, err := client.LoginUser(context.Background(), &userpb.LoginRequest{
		Email:    "janedoe@exmpl.com",
		Password: "ABOBA",
	})
	if err != nil {
		log.Fatalf("Ошибка при входе пользователя : %v", err)
	}
	fmt.Printf("Login response: %v\n", loginResp)

	jwtToken := loginResp.GetToken()
	fmt.Printf("JWTTOKEN:%v", jwtToken)
	//loginResp, err := client.LoginUser(context.Background(), &userpb.LoginRequest{
	//	Email:    "janeRat@.mail",
	//	Password: "ABOBA",
	//})
	//if err != nil {
	//	log.Fatalf("Error logging in: %v", err)
	//}
	//fmt.Printf("LoginResponse: %v\n", loginResp)
}
