package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ThomasITU/DISYSPrep/Proto"
	h "github.com/ThomasITU/DISYSPrep/helpermethod"
	"google.golang.org/grpc"
)

const (
	SERVER_ADDRESS = "localhost:5000"
)

type User struct {
	userId int64
}

func main() {

	//init
	//setup a connection, this is a blocking call
	conn, err := grpc.Dial(SERVER_ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//create user
	id := ChooseUserId()
	user := User{userId: int64(id)}

	//create client
	ctx := context.Background()
	client := Proto.NewProtoServiceClient(conn)

	//wait for input
	user.ListenForInput(client, ctx)
	fmt.Println("main has ended")
}

// listen for user input as a string and send request accordingly, implement logic in seperate methods as see fit
func (u *User) ListenForInput(client Proto.ProtoServiceClient, ctx context.Context) {
	for {
		var input string
		fmt.Scanln(&input)
		if len(input) > 0 {
			switch input {
			case "joinchat":
				response, err := client.JoinService(ctx, &Proto.JoinRequest{UserId: u.userId})
				h.CheckError(err, "listenForInput joinchat")
				fmt.Println(response.GetMsg())
			case "getvalue":
				response, err := client.GetValue(ctx, &Proto.GetRequest{})
				h.CheckError(err, "ListenForInput getvalue")
				fmt.Println(response)
			case "setvalue":
				fmt.Println("Choose an integer you want to set the value")
				var value int64
				fmt.Scanln(&value)
				response, err := client.SetValue(ctx, &Proto.SetRequest{UserId: u.userId, RequestedValue: value})
				h.CheckError(err, "ListenForInput setvalue")
				fmt.Println(response.GetMsg())
			}
		}
	}
}

// helper method wait for user to input a wanted userid
func ChooseUserId() int {
	for {
		var userId int
		fmt.Println("Choose a positive integer as id:")
		fmt.Scanln(&userId)
		if userId > 0 {
			return userId
		} else {
			fmt.Printf("\nSomething went wrong when processing your input: %v\nTry a positive integer like 1, 2, 3, ...", userId)
		}
	}
}
