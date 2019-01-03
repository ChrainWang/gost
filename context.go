package gost

import (
	"bytes"
	"net/http"
)

type _Param struct {
	key   []byte
	value []byte
	next  *_Param
}

func (self *_Param) find(key []byte) (*_Param, bool) {
	switch bytes.Compare(self.key, key) {
	case 0:
		return self, true
	case 1:
		return nil, false
	case -1:
		if self.next == nil {
			return self, false
		}
		target, found := self.next.find(key)
		if target == nil {
			return self, false
		}
		return target, found
	}

	return nil, false
}

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	urlParam *_Param

	handlers    *_Handler
	middlewares *_Handler
}

func (self *Context) Get(key string) (string, bool) {
	for current := self.urlParam; current != nil; current = current.next {
		if bytes.Compare(current.key, []byte(key)) == 0 {
			return string(current.value), true
		}
	}
	return "", false
}

func (self *Context) setUrlParam(key, value []byte) {
	keyBuf := []byte(key)

	if self.urlParam == nil {
		self.urlParam = new(_Param)
		self.urlParam.key = keyBuf
		self.urlParam.value = value
		return
	}

	position, found := self.urlParam.find([]byte(key))
	if found {
		position.value = value
		return
	}

	param := new(_Param)
	param.key = keyBuf
	param.value = value

	if position == nil {
		param.next = self.urlParam
		self.urlParam = param
		return
	}

	param.next = position.next
	position.next = param
}
