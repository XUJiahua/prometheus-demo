package cybersource

type Service interface {
	Call(request *Request) error
}
