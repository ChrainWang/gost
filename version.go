package gost

const VERSION string = "0.0.1"

/*
******** Ver 0.0.1 *********
1. Gost router with linkedtree data structure
	1) every endpoints in the path are stored as Endpoint type instaces
	2) endpoints on differenct layer are linked up bi-directedly
	3) endpoints on the same layer are linked up in one direction, sorted by endpoint names(prefixed rune ':' would be ignored)
	4) in this version only GET, POST, PUT, DELETE, OPTIONS are supported (OPTIONS is not suggested to be used).
	5) A default handler can be registered for all the methods above. If so, other handlers would be override

2. Api register on path root
	1) path root is "/"
	2) endpoints closer to path root are on higher layer

3. Api register on endpoints
	1) if user registers "/a" on endpont "/b", it equals to registering "/b/a" on path root

4. middleware register
	1) middlewares would be run one by one if any descendants of this endpoint is accessed
	2) example:
		uri path: "/a/b/c/d"
		middlewares registerd on "/a" and "/a/b/c"
		handlers will be run one by one:
		MiddlewareA -> MiddlewareC -> HandlerD
*/
