jsoncall is a Go module that facilitates creating HTTP web services from Go interfaces,
using reflection to marshal/unmarshal Go method inputs/outputs to/from a single JSON object.

# Goals/Features
- Use reflection to automatically marshal/unmarshal Go method inputs/outputs.
- The server uses a well-defined protocol that can be called from other languages.
- The optional Go client code uses a generated stub to implement a Go interface that
  is implemented by marshalling method parameters to HTTP calls 
  and unmarshalling the outputs from the response.
- Encapsulate all method parameters into a single JSON object
- Encapsulate all method outputs into a single JSON object
- Support any input parameter or output value that can be marshalled/unmarshalled to/from JSON.
- Marshal/unmarshal inputs/outputs, using encoding/json with a single call to json.Marshal, json.Unmarshal.

# steps
To create a Go client and a Go server that communicate via a HTTP,
the following steps need to be done:
- create an Api interface type
- create a struct type that implements the Api.  It will be run on the server.
- generate a JSON descriptor file, via the descriptor generator (optional)
- edit the descriptor file by hand to customize method, input, and output names (optional)
- generate a client stub, via a provided code generator.  (needed only by the client).
- compile the client and the server

# Example server (full)
This example uses an interface, and a custom api descriptor.
## common code
	package example
	import (
		_ "embed"
		"net/http"
		
		"melato.org/jsoncall"
	)

	type Demo interface {
		Repeat(s string, count int) ([]string, error)
	}
	
	//go:embed demo.json
	var demoNames []byte
	
	func NewDemoCaller() (*jsoncall.Caller, error) {
		var api *Demo
		return jsoncall.NewCaller(api, demoNames)
	}

## server code
	import (
		"example"
	)

	type DemoImpl struct {}
	func (t *DemoImpl) Repeat(s string, count int) ([]string, error) {
		if count < 0 {
			return nil, fmt.Errorf("negative count: %d", count)
		}
		list := make([]string, count)
		for i := 0; i < count; i++ {
			list[i] = s
		}
		return list, nil
	}

	type Server struct {}
	
	func (t *Server) DemoReceiver(w http.ResponseWriter, r *http.Request) interface{} {
		return &example.DemoImpl{}
	}

	func (t *Server) Run() error {
		demoCaller, err := example.NewDemoCaller()
		if err != nil {
			return err
		}
		mux := http.NewServeMux()
		mux.Handle("/demo/", demoCaller.NewHttpHandler(t.DemoReceiver)
	
		return http.ListenAndServe(":8080", mux)
	}
		
## client code
	func NewDemoClient() (example.Demo, error) {
		caller, err := example.NewDemoCaller()
		if err != nil {
			return nil, err
		}
		c := caller.NewHttpClient("http://localhost:8081/demo/")
		return generated.NewDemoClient(c), nil
	}


# Generating a Go client stub
	package main
	
	import (		
		"example"
		"melato.org/jsoncall/generate"
	)

	func GenerateStub() error {
		var g generate.Generator
		g.Init()
		g.Package = "generated"
		g.Type = "exampleClient"
		g.Func = "NewDemoClient"
		g.OutputFile = "../generated/generated_example.go"
		g.Package = "generated"
		//g.Imports = []string{"example/a", "example/b"}
		caller, err := example.NewExampleCaller()
		if err != nil {
			return err
		}
		return g.Output(g.GenerateClient(caller))
	}
	
# Updating the JSON API descriptor
The API descriptor is an optional JSON file.
When the API changes, by adding or removing methods, you need to update the API descriptor by hand.
To facilitate updating, there is a tool that automatically updates the JSON file.
It removes methods that are no longer in the API interface.
It adds new methods, using default naming conventions for the method names and their parameters.
You can then edit the file by hand to replace the default names with custom names.

	package main
	import (
		"melato.org/jsoncall/generate"		
	)
	func Main() {
		var api *example.Example
		err := generate.UpdateDescriptor(api, "../api.json")
		if err != nil {
			fmt.Println(err)
		}		
	}

To create an initial API descriptor, you can create an empty api.json file,
and then update it with the code above.

# API descriptor example
	[
	 {
	  "method": "Repeat",
	  "path": "repeat",
	  "in": [
	   "s",
	   "count"
	  ],
	  "out": [
	   "result",
	   "error"
	  ]
	 }
    ]


## POST request
A client can call this server by making an HTTP POST request to 
http://localhost:8080/demo/repeat with body:

	{"s":"a","count":2}
	
For example, using curl:

	curl --data-binary '{"s":"a","count":2}' http://localhost:8080/repeat

It should get a response like this:

	{"result":["a","a"]}
	

# JSON protocol
The chosen protocol has some similarities with JSON-RPC 2.0, but does not follow it entirely.

The default output keys "result" and "error" follow the JSON-RPC conventions, when there is a single output.

Putting the method name in the url allows us to know how to unmarshal the body of the request into the arguments.  If we encoded the method in the JSON request, as per JSON-RPC, we couldn't unmarshal the request with a single pass.

Putting the method parameters in a map allows us to unmarshal them into a struct (generated via reflection) that has a field for each parameter.

We use a map for the parameters, so we can unmarshal them into their different types.  If we had put the parameters in an array, we could only unmarshal them into one array type, such as []any.  

We could have placed the parameters into a "params" map, as JSON-RPC specifies.
Perhaps we'll add an option to do this.
