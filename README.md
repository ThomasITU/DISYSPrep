# DISYSPrep
Some hopefully useful code snippets for DISYS exam

## Protoc commands from cmd in parent folder - ./DISYSPrep
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Proto/renameService.proto

# How to run
- Start one server, cd into /server 

`go run .`

- Start multiple clients, cd into /client 

`go run .`

## 3 commands from client further instructions will be displayed in the terminal

`joinchat`

`getvalue`

`setvalue` 


# Details
Server log is written everytime setvalue is invoked
