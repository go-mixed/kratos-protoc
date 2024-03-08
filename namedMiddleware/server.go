package namedMiddleware

import (
	"context"
	"net/http"
)

// enableNamedMiddleware handler of http.Server
func enableNamedMiddlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(newContext(request.Context(), &namedMiddleware{}))
		handler.ServeHTTP(writer, request)
	})
}

// EnableNamedMiddleware 在kratos的http.Server中启用命名中间件
// enable named middleware in kratos http.Server
//
//	http.Server(
//	 http.Filter(namedMiddleware.EnableNamedMiddleware()),
//	 http.Handler(namedMiddleware.WrapKratosMiddleware("auth", yourMiddleware)),
//	)
func EnableNamedMiddleware() func(http.Handler) http.Handler {
	return enableNamedMiddlewareHandler
}

// DispatchMiddleware 给当前路由的添加待执行的命名中间件（无需主动调用），将在xxx_http.pb.go的路由中调用
// add the named middleware to the current route (no need to call it actively), it will be called in the route of xxx_http.pb.go
func DispatchMiddleware(ctx context.Context, name string, arguments ...string) {
	nm := fromContext(ctx)
	if nm == nil {
		return
	}

	nm.dispatch(name, arguments...)
}
