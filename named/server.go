package named

import (
	"context"
	"net/http"
)

// EnableNamedMiddleware 在kratos的http.Server中启用命名中间件
//
//	http.Server(
//	 http.Filter(named.EnableNamedMiddleware),
//	 http.Handler(named.KratosMiddleware("auth", yourMiddleware)),
//	)
func EnableNamedMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		request = request.WithContext(newContext(request.Context(), &namedMiddleware{}))
		handler.ServeHTTP(writer, request)
	})
}

// DispatchMiddleware 给当前路由的添加待执行的命名中间件
func DispatchMiddleware(ctx context.Context, name string, arguments ...string) {
	nm := fromContext(ctx)
	if nm == nil {
		return
	}

	nm.dispatch(name, arguments...)
}
