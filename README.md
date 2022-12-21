jsoncall is a Go module that facilitates creating HTTP web services from Go interfaces,
using reflection to marshal/unmarshal Go method inputs/outputs to/from a single JSON object.

# Example server
	import "net/http"
	
	type Example struct {}
	func (t *Example) A(s string, d int) (string, error) {...}
	func (t *Example) B() string {...}

	func ExampleServer() error {
		var handler http.Handler
		handler, err := jsoncall.NewHttpHandler(&Example{})
		if err != nil {
			return err
		}
		return http.ListenAndServe(":8080", handler)
	}

## POST request
A client can call this server by making an HTTP POST request to 
http://localhost:8080/A with body:

	{"p1":"a","p2":2}
	
For example, using curl:

	curl --data-binary '{"p1":"a","p2":2}' http://localhost:8080/A

It should get a response like this:

	{"result":"a:2"}
	

## The server can specify the served methods with an interface

	type ExampleInterface interface {
		A(s string, d int) (string, error)
	}

	var api *ExampleInterface
	handler, err := jsoncall.NewHttpHandler(api)
	handler.SetReceiver(&Example{})

## The server can use a function that returns a receiver for each request
	func ExampleProvider(w http.ResponseWriter, r *http.Request) {
		...
		return &Example{}
	}

	handler.SetReceiverFunc(ExampleProvider)

This function can use the request headers to authenticate the request,
and either incorporate authentication info in the returned receiver,
or write an authentication error to the response and return nil.

# Go client, without using generated code:
	var example *ExampleInterface
	caller, err := jsoncall.NewCaller(example, nil)
	client := caller.NewHttpClient("http://localhost:8080/")

	var response map[string]any
	err = client.Call(&response, "A", "hello", 2)

# Go client, with generated code:
	import "{your path}/generated"
	client, err := generated.NewExampleClient()
	if err != nil {
		return err
	}
	s, err := client.A("hello", 7)

# Generating a Go client stub for an interface
	import "melato.org/jsoncall/generate"

	g := generate.NewGenerator()
	g.Func = "NewExampleClient"
	g.OutputFile = "../generated/example.go"
	g.Package = "generated"

	var example *ExampleInterface
	caller, err := jsoncall.NewCaller(example, nil)
	if err != nil {
		return err
	}
	return g.Output(g.GenerateClient(caller))

# melato.org/jsoncall/generate
The package melato.org/jsoncall/generate
- Generates or updates .json API descriptor files with new methods.
- Generates Go client code that implements the Go interface used by the server
and makes the corresponding requests.


# method names and input/output keys
The server uses the last component of the URL path to select the Method to use.
The name of the method in the URL, and the JSON keys for the method inputs and outputs can be specified with an ApiDescriptor,
which can be provided as a JSON file, which can be embedded in the server.
If there is no explicit API descriptor provided, a default descriptor is used:
- The last url path component is the name of the Go method.
- The input parameters of each method are named "p1", "p2", ...
- If the method has exactly one non-error output (any type other than "error"), it is named "result".
- Otherwise, non-error outputs are named "r1", "r2", ....
- If the method has exactly one error output, it is named "error".
- Otherwise, error outputs are named "e1", "e2", ....

Providing an API descriptor ensures that the Go method names and input/output order can change without affecting any clients.

# Code generator
The package melato.org/jsoncall/generate
- Generates or updates .json API descriptor files with new methods.
- Generates Go client code that implements the Go interface used by the server,
and makes the corresponding requests.

# Goals/Features
- Use reflection to automatically marshal/unmarshal Go method inputs/outputs.
- Usable from Go or from other languages.
- Encapsulate all method parameters into a single JSON object
- Encapsulate all method outputs into a single JSON object
- Support any input parameter or output value that can be marshalled/unmarshalled to/from JSON.
- Marshal/unmarshal inputs/outputs, using encoding/json with a single call to json.Marshal, json.Unmarshal.

# JSON protocol
The chosen protocol has some similarities with JSON-RPC 2.0, but does not follow it entirely.

The default output keys "result" and "error" follow the JSON-RPC conventions, when there is a single output.

Putting the method name in the url allows us to know how to unmarshal the body of the request into the arguments.  If we encoded the method in the JSON request, as per JSON-RPC, we couldn't unmarshal the request with a single pass.

Putting the method parameters in a map allows us to unmarshal them into a struct (generated via reflection) that has a field for each 

We use a map for the parameters, so we can unmarshal them into their different types.  If we had put the parameters in an array, we could only unmarshal them into one array type, such as []any.  

We could have placed the parameters into a "params" map, as JSON-RPC specifies, add a level of indirection.  This would allow data other than method parameters to be included in the request, but since we use HTTP, we can already add such data in the HTTP headers.