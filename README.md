# DISYSPrep
Some hopefully useful code snippets for DISYS exam

## Protoc commands from cmd in parent folder - ./DISYSPrep
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Proto/renameService.proto

# How to run
- Start one front end

`go run .\FrontEnd\FrontEnd.go`

- start multiple servers (up to MAX_REPLICAS, default = 5)

`go run .\Server\Server.go`

- Start multiple clients

`go run .\Client\Client.go`

## 3 commands from client further instructions will be displayed in the terminal

`joinchat`

`getvalue`

`setvalue` 


# Details
This is active replication, meaning **client** request are handle by the **front end** and executed on all the **servers**

A Front end log is written everytime setvalue is invoked

A Server log for each port is written everytime setvalue is invoked
