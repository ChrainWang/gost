package gost

import (
	"bytes"
	"errors"
	"strings"
)

var _SPLIT_PATTERNS []rune = []rune{
	'/', '\\',
}

type _RuneRange struct {
	Lo byte
	Hi byte
}

type _LegalRunes struct {
	ranges []_RuneRange
}

func (self _LegalRunes) isLegal(r byte) bool {
	for _, runeRange := range self.ranges {
		if r >= runeRange.Lo && r <= runeRange.Hi {
			return true
		}
	}
	return false
}

var (
	_LEGAL_RUNES_CASESENSITIVE = &_LegalRunes{
		ranges: []_RuneRange{
			_RuneRange{Lo: '-', Hi: '-'},
			_RuneRange{Lo: '_', Hi: '_'},
			_RuneRange{Lo: '0', Hi: '9'},
			_RuneRange{Lo: 'a', Hi: 'z'},
			_RuneRange{Lo: 'A', Hi: 'Z'},
		},
	}
	_LEGAL_RUNES_NON_CASESENSITIVE = &_LegalRunes{
		ranges: []_RuneRange{
			_RuneRange{Lo: '-', Hi: '-'},
			_RuneRange{Lo: '_', Hi: '_'},
			_RuneRange{Lo: '0', Hi: '9'},
			_RuneRange{Lo: 'a', Hi: 'z'},
		},
	}
)

// common errors
var (
	_IlegalRuneErr         = errors.New("Ilegal rune in path")
	_DynamicEpSingletonErr = errors.New("Dynamic endpoint must be the only child of it's parent")
	_MethodNotSupported    = errors.New("Method is not supported")
)

// methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"

	OPTIONS = "OPTIONS"
)

/******* Endpoint *******/
//  Internaly used
type Endpoint struct {
	name []byte

	parent *Endpoint
	siblin *Endpoint
	child  *Endpoint

	casesensitive bool
	dynamic       bool // if the endpoint is prefixed with ':'

	middlewares *_Handler

	any  *_Handler
	get  *_Handler
	post *_Handler
	put  *_Handler
	del  *_Handler

	options *_Handler
}

func (self *Endpoint) NewEndpoint() *Endpoint {
	ep := new(Endpoint)
	ep.casesensitive = self.casesensitive

	return ep
}

func (self *Endpoint) AddService(path string, service *Service) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.addChild(spliter)

	if err != nil {
		return err
	}

	if service.method == "" {
		if endpoint.any == nil {
			endpoint.any = service.handlers
			return err
		}

		endpoint.any.append(service.handlers)
		return err
	}

	switch strings.ToUpper(service.method) {
	case GET:
		if endpoint.get == nil {
			endpoint.get = service.handlers
			break
		}
		endpoint.get.append(service.handlers)
	case POST:
		if endpoint.post == nil {
			endpoint.post = service.handlers
			break
		}
		endpoint.post.append(service.handlers)
	case PUT:
		if endpoint.put == nil {
			endpoint.put = service.handlers
			break
		}
		endpoint.put.append(service.handlers)
	case DELETE:
		if endpoint.del == nil {
			endpoint.del = service.handlers
			break
		}
		endpoint.del.append(service.handlers)
	case OPTIONS:
		if endpoint.options == nil {
			endpoint.options = service.handlers
			break
		}
		endpoint.options.append(service.handlers)
	default:
		return _MethodNotSupported
	}

	return err
}

// path register
func (self *Endpoint) Any(path string, handlers ...HandlerFunc) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.addChild(spliter)
	if err != nil {
		return err
	}

	if endpoint.any == nil {
		endpoint.any = new(_Handler)
		endpoint.any.handler = handlers[0]
		handlers = handlers[1:]
	}

	_handler := endpoint.any.findLast()
	for _, handler := range handlers {
		_handler.next = new(_Handler)
		_handler.next.handler = handler
		_handler = _handler.next
	}

	return nil
}

func (self *Endpoint) Get(path string, handlers ...HandlerFunc) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.addChild(spliter)
	if err != nil {
		return err
	}

	if endpoint.get == nil {
		endpoint.get = new(_Handler)
		endpoint.get.handler = handlers[0]
		handlers = handlers[1:]
	}

	_handler := endpoint.get.findLast()
	for _, handler := range handlers {
		_handler.next = new(_Handler)
		_handler.next.handler = handler
		_handler = _handler.next
	}

	return nil
}

func (self *Endpoint) Post(path string, handlers ...HandlerFunc) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.addChild(spliter)
	if err != nil {
		return err
	}

	if endpoint.post == nil {
		endpoint.post = new(_Handler)
		endpoint.post.handler = handlers[0]
		handlers = handlers[1:]
	}

	_handler := endpoint.post.findLast()
	for _, handler := range handlers {
		_handler.next = new(_Handler)
		_handler.next.handler = handler
		_handler = _handler.next
	}

	return nil
}

func (self *Endpoint) Put(path string, handlers ...HandlerFunc) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.addChild(spliter)
	if err != nil {
		return err
	}

	if endpoint.put == nil {
		endpoint.put = new(_Handler)
		endpoint.put.handler = handlers[0]
		handlers = handlers[1:]
	}

	_handler := endpoint.put.findLast()
	for _, handler := range handlers {
		_handler.next = new(_Handler)
		_handler.next.handler = handler
		_handler = _handler.next
	}

	return nil
}

func (self *Endpoint) Delete(path string, handlers ...HandlerFunc) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.addChild(spliter)
	if err != nil {
		return err
	}

	if endpoint.del == nil {
		endpoint.del = new(_Handler)
		endpoint.del.handler = handlers[0]
		handlers = handlers[1:]
	}

	_handler := endpoint.del.findLast()
	for _, handler := range handlers {
		_handler.next = new(_Handler)
		_handler.next.handler = handler
		_handler = _handler.next
	}

	return nil
}

func (self *Endpoint) Options(path string, handlers ...HandlerFunc) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.addChild(spliter)
	if err != nil {
		return err
	}

	if endpoint.options == nil {
		endpoint.options = new(_Handler)
		endpoint.options.handler = handlers[0]
		handlers = handlers[1:]
	}

	_handler := endpoint.options.findLast()
	for _, handler := range handlers {
		_handler.next = new(_Handler)
		_handler.next.handler = handler
		_handler = _handler.next
	}

	return nil
}

// try to find endpoint with given path
// if endpoint is found, return it
// if not, create one and insert it into the tree, then return it
func (self *Endpoint) addChild(provider *Spliter) (*Endpoint, error) {
	epName, err := provider.Next()
	if err != nil {
		return nil, err
	}

	if epName == nil {
		return self, nil
	}

	var legalRunes *_LegalRunes
	if self.casesensitive {
		legalRunes = _LEGAL_RUNES_CASESENSITIVE
	} else {
		legalRunes = _LEGAL_RUNES_NON_CASESENSITIVE
	}

	/*********************** dynamic endpoint *******************/
	if epName[0] == ':' {
		epName = epName[1:]

		for _, r := range epName {
			if !legalRunes.isLegal(r) {
				return nil, _IlegalRuneErr
			}
		}

		// current endpoint has no child
		if self.child == nil {
			// create a new endpoint
			newEndpoint := self.NewEndpoint()
			newEndpoint.name = epName
			newEndpoint.dynamic = true
			newEndpoint.parent = self
			// the target is the endpoint that will register handlers or middlewares with
			target, err := newEndpoint.addChild(provider)
			if err != nil {
				return nil, err
			}
			// make the new endpoint as the child of current one
			self.child = newEndpoint
			return target, nil
		}

		// current endpoint has children
		if !self.child.dynamic || bytes.Compare(self.child.name, epName) != 0 {
			// current endpoint has a static child
			// or current endpoint has a dynamic child with different name from new endpoint
			// return error
			return nil, _DynamicEpSingletonErr
		}

		return self.child.addChild(provider)

	}

	/********************* static endpoint done *****************/
	for _, r := range epName {
		if !legalRunes.isLegal(r) {
			return nil, _IlegalRuneErr
		}
	}

	// current endpoint has no child
	if self.child == nil {
		// create a new endpoint
		newEndpoint := self.NewEndpoint()
		newEndpoint.name = epName
		newEndpoint.parent = self
		// the target is the endpoint that will register handlers or middlewares with
		target, err := newEndpoint.addChild(provider)
		if err == nil {
			// make the new endpoint as the child of current one
			self.child = newEndpoint
		}
		return target, nil
	}

	// current endpoint has children
	// find the position that the new one shall be inserted at
	predictedPosition, found := self.child.findSiblin(epName)
	if found {
		return predictedPosition.addChild(provider)
	}

	newEndpoint := self.NewEndpoint()
	newEndpoint.name = epName
	newEndpoint.parent = self

	target, err := newEndpoint.addChild(provider)

	if err == nil {
		if predictedPosition == nil {
			newEndpoint.siblin = self.child
			self.child = newEndpoint
		} else {
			newEndpoint.siblin = predictedPosition.siblin
			predictedPosition.siblin = newEndpoint
		}
	}

	return target, err
}

func (self *Endpoint) findSiblin(epName []byte) (*Endpoint, bool) {
	switch bytes.Compare(self.name, epName) {
	case 0: // this is exactly the endpoint that is searching for
		return self, true
	case 1:
		// current endpoint has a bigger name
		// the new endpoint shall be inserted before current endpoint
		return nil, false
	case -1:
		if self.siblin == nil {
			// current endpoint is the end of child endpoint chain
			return self, false
		}
		predictedPosition, found := self.siblin.findSiblin(epName)
		if predictedPosition == nil {
			return self, false
		}
		return predictedPosition, found
	default:
		// unreachable
		return nil, true
	}
}

func (self *Endpoint) route(spliter *Spliter, c *Context) *Endpoint {
	if c.middlewares == nil {
		c.middlewares = self.middlewares
	} else {
		mw := c.middlewares
		for {
			if mw.next == nil {
				break
			}
			mw = mw.next
		}
		mw.next = self.middlewares
	}
	if name, err := spliter.Next(); err == nil {
		if name == nil {
			return self
		}

		if self.child == nil {
			spliter.Close()
			return self
		}

		if self.child.dynamic {
			c.setUrlParam(self.child.name, name)
			return self.child.route(spliter, c)
		}

		if child, found := self.child.findSiblin(name); found {
			return child.route(spliter, c)
		}

		spliter.Close()
	}
	return nil
}

func (self *Endpoint) getChild(spliter *Spliter) *Endpoint {
	if name, err := spliter.Next(); err == nil {
		if name == nil {
			return self
		}
		if self.child == nil {
			spliter.Close()
			return nil
		}

		if self.child.dynamic {
			return self.child.getChild(spliter)
		}

		child, found := self.child.findSiblin(name)
		if found {
			return child.getChild(spliter)
		}
		spliter.Close()
		return nil
	}
	return nil
}

func (self *Endpoint) FullPath() string {
	buffer := bytes.Buffer{}
	if err := self.searchBack(&buffer); err != nil {
		return ""
	}
	return buffer.String()
}

func (self *Endpoint) searchBack(buffer *bytes.Buffer) error {
	if self.parent == nil {
		return nil
	}
	self.parent.searchBack(buffer)
	if err := buffer.WriteByte('/'); err != nil {
		return err
	}
	if self.dynamic {
		if err := buffer.WriteByte(':'); err != nil {
			return err
		}
	}
	if _, err := buffer.Write(self.name); err != nil {
		return err
	}
	return nil
}

/*********** Gost router ********/
type Gost struct {
	root          *Endpoint
	casesensitive bool
}

// Constructor for Gost
func NewGost(casesensitive bool) *Gost {
	g := new(Gost)
	g.casesensitive = casesensitive
	g.root = g.NewEndpoint()
	return g
}

func (self *Gost) NewEndpoint() *Endpoint {
	node := new(Endpoint)
	node.casesensitive = self.casesensitive
	return node
}

// Deprecated
// Only for testing in recent versions
// Will be removed in the future
func (self *Gost) AddService(path string, service *Service) error {
	var handlers []HandlerFunc
	for _handler := service.handlers; _handler != nil; _handler = _handler.next {
		handlers = append(handlers, _handler.handler)
	}

	if len(service.method) == 0 {
		return self.root.Any(path, handlers...)
	}
	switch strings.ToUpper(service.method) {
	case GET:
		return self.root.Get(path, handlers...)
	case POST:
		return self.root.Post(path, handlers...)
	case PUT:
		return self.root.Put(path, handlers...)
	case DELETE:
		return self.root.Delete(path, handlers...)
	case OPTIONS:
		return self.root.Options(path, handlers...)
	default:
		return _MethodNotSupported
	}
}

func (self *Gost) route(c *Context) error {
	spliter := NewSpliter([]byte(c.Request.RequestURI), _SPLIT_PATTERNS)
	spliter.Split()
	if endpoint := self.root.route(spliter, c); endpoint != nil {
		var handler *_Handler
		switch strings.ToUpper(c.Request.Method) {
		case GET:
			handler = endpoint.get
		case POST:
			handler = endpoint.post
		case PUT:
			handler = endpoint.put
		case DELETE:
			handler = endpoint.del
		default:
			c.Writer.WriteHeader(405)
			return _MethodNotSupported
		}
		c.handlers = handler
		return nil
	}

	c.Writer.WriteHeader(405)
	return _MethodNotSupported
}

func (self *Gost) GetEndpoint(path string) *Endpoint {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()
	return self.root.getChild(spliter)
}

func (self *Gost) Any(path string, handlers ...HandlerFunc) error {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()

	endpoint, err := self.root.addChild(spliter)
	if err != nil {
		return err
	}

	if endpoint.any == nil {
		endpoint.any = new(_Handler)
		endpoint.any.handler = handlers[0]
		handlers = handlers[1:]
	}

	_lastHandler := endpoint.any.findLast()
	for _, _handler := range handlers {
		_lastHandler.next = new(_Handler)
		_lastHandler.next.handler = _handler
		_lastHandler = _lastHandler.next
	}

	return nil

}

func (self *Gost) Get(path string, handlers ...HandlerFunc) error {
	return self.root.Get(path, handlers...)
}

func (self *Gost) Post(path string, handlers ...HandlerFunc) error {
	return self.root.Post(path, handlers...)

}

func (self *Gost) Put(path string, handlers ...HandlerFunc) error {
	return self.root.Put(path, handlers...)

}

func (self *Gost) Delete(path string, handlers ...HandlerFunc) error {
	return self.root.Delete(path, handlers...)
}

func (self *Gost) Options(path string, handlers ...HandlerFunc) error {
	return self.root.Options(path, handlers...)
}

func (self *Gost) Group(path string, middlewares ...HandlerFunc) *Endpoint {
	spliter := NewSpliter([]byte(path), _SPLIT_PATTERNS)
	spliter.Split()
	endpoint := self.root.getChild(spliter)

	if endpoint.middlewares == nil {
		endpoint.middlewares = new(_Handler)
		endpoint.middlewares.handler = middlewares[0]
		_middleware := endpoint.middlewares
		for _, middleware := range middlewares[1:] {
			_middleware.next = new(_Handler)
			_middleware.next.handler = middleware
			_middleware = _middleware.next
		}
		return endpoint
	}

	_middleware := endpoint.middlewares
	for _, middleware := range middlewares {
		_middleware.next = new(_Handler)
		_middleware.next.handler = middleware
		_middleware = _middleware.next
	}

	return endpoint
}
