make:
	go mod init leisurely && go mod tidy 
	
run:
	go run server.go

clean:
	rm ./go.mod
	rm ./go.sum