To compile and run the example:

# When the remote interface example.Demo changes
- go run demo-generate-json.go
	updates ../demo.json
- edit ../demo.json by hand, if you wish to change the parameter names
- go run demo-generate-client.go
	generates ../generated/...

# To compile the server and client
	go build demo-client.go
	go build demo-client.go

# To run the server
	./demo-server
	
# To run the client
	./demo-client -h
	./demo-client repeat a 3
	./demo-client 
	
# To see the JSON request/response
	add "-trace" to the client or server command line