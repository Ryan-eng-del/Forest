package inter


const HelloServiceName = "Hello"
const HelloServiceMethod = "Hello.HelloWorld"

type HelloService interface {
	HelloWorld(arg string, reply *string) error
}