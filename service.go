package gost

type HandlerFunc func(*Context)

type _Handler struct {
	handler HandlerFunc
	next    *_Handler
}

func (self *_Handler) findLast() *_Handler {
	if self.next == nil {
		return self
	}
	return self.next.findLast()
}

func (self *_Handler) append(handler *_Handler) {
	self.findLast().next = handler
}

type Service struct {
	method   string
	handlers *_Handler
}

func NewService(method string, handlerFuncs ...HandlerFunc) *Service {
	service := new(Service)
	service.method = method
	current := service.handlers
	for _, handlerFunc := range handlerFuncs {
		if handlerFunc == nil {
			break
		}
		handler := new(_Handler)
		handler.handler = handlerFunc
		if current == nil {
			service.handlers = handler // first one
			current = handler
			continue
		}
		current.next = handler
		current = handler
	}
	return service
}
