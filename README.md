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
	func ExampleReceiver(w http.ResponseWriter, r *http.Request) {
		...
		return &Example{}
	}

	handler.SetReceiverFunc(ExampleReceiver)


# Go client, without using generated code:
	var example *ExampleInterface
	caller, err := jsoncall.NewCaller(example, nil)
	client := caller.NewHttpClient("http://localhost:8080/")

	var response map[string]any
	err = client.Call(&response, "A", "hello", 2)

# Go client, with generated code:
An included code generator, generates client code that implements
the same interface used by the server.

Generated code is necessary, because Go does not have a mechanism
to implement an interface at runtime using reflection.

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

# ApiDescriptor
The name of the method in the URL, and the JSON keys
for the method inputs and outputs can be specified with an ApiDescriptor.
If there is no explicit API descriptor provided, a default descriptor is used.
An api descriptor is typically specified as an embedded JSON file,
which can be generated as follows:

```
var example *ExampleInterface
desc := generate.GenerateDescriptor(example)
json.Marshal(desc)
```
The resulting JSON is:
```
[
 {
  "in": [
   "p1",
   "p2"
  ],
  "method": "A",
  "out": [
   "result",
   "error"
  ],
  "path": "A"
 },
 {
  "method": "B",
  "out": [
   "result"
  ],
  "path": "B"
 }
]
```

# Code generator
The package melato.org/jsoncall/generate
- Generates .json API descriptor files or updates then with new methods.
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

Putting the method parameters in a map allows us to unmarshal them into a struct (generated via reflection) that has a field for each parameter.

We use a map for the parameters, so we can unmarshal them into their different types.  If we had put the parameters in an array, we could only unmarshal them into one array type, such as []any.  

We could have placed the parameters into a "params" map, as JSON-RPC specifies.
Perhaps we'll add an option to do this.
