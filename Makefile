.PHONY:client

client:
	cd client && yarn start

server:
	go run server.go