package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ThomasITU/DISYSPrep/Proto"
	h "github.com/ThomasITU/DISYSPrep/helpermethod"
	"google.golang.org/grpc"
)

const (
	SERVER_ADDRESS = "localhost:5000"
)

type User struct {
	userId           int64
	lamportTimeStamp int64
	arbiter          sync.Mutex
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
	user := User{userId: int64(id), lamportTimeStamp: 0, arbiter: sync.Mutex{}}

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
				u.arbiter.Lock()
				response, err := client.JoinService(ctx, &Proto.JoinRequest{UserId: u.userId, Timestamp: (u.lamportTimeStamp + 1)})
				h.CheckError(err, "listenForInput joinchat")
				u.IncrementLamportTimestamp(response.GetTimestamp())
				u.arbiter.Unlock()

				fmt.Println(response.GetMsg())
			case "getvalue":
				u.arbiter.Lock()

				response, err := client.GetValue(ctx, &Proto.GetRequest{Timestamp: (u.lamportTimeStamp + 1)})
				h.CheckError(err, "ListenForInput getvalue")
				u.IncrementLamportTimestamp(response.GetTimestamp())
				u.arbiter.Unlock()

				fmt.Println(response)
			case "setvalue":
				var value int64
				for value == 0 {
					fmt.Println("Choose an none 0 integer you want the value set to")
					fmt.Scanln(&value)
				}

				u.arbiter.Lock()
				response, err := client.SetValue(ctx, &Proto.SetRequest{UserId: u.userId, RequestedValue: value, Timestamp: (u.lamportTimeStamp + 1)})
				h.CheckError(err, "ListenForInput setvalue")
				u.IncrementLamportTimestamp(response.GetTimestamp())
				u.arbiter.Unlock()

				fmt.Println(response.GetMsg())
				value = 0 // reset the value
			}
		}
	}
}

func (u *User) IncrementLamportTimestamp(serverTimeStamp int64) {
	u.lamportTimeStamp = (h.Max(u.lamportTimeStamp, serverTimeStamp) + 1)
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
