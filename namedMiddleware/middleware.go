package namedMiddleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
)

var MiddlewareWithArguments = struct{}{}

// WrapKratosMiddleware 封装kratos的中间件为命名中间件
// wrap kratos middleware to named middleware
func WrapKratosMiddleware(name string, mw middleware.Middleware) middleware.Middleware {
	return func(nextHandler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if caller := match(ctx, name); caller != nil {
				// add caller to context, and you can only get the caller in the current middleware frame
				ctx1 := newCallerContext(ctx, caller)
				return mw(nextHandler)(ctx1, req)
			}
			return nextHandler(ctx, req)
		}
	}
}

// Middleware 注册一个命名kratos中间件
// register a named kratos middleware
func Middleware(name string, handler middleware.Handler) middleware.Middleware {
	return func(nextHandler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if caller := match(ctx, name); caller != nil {
				// add caller to context, and you can only get the caller in the current middleware frame
				ctx1 := newCallerContext(ctx, caller)
				newReq, err := handler(ctx1, req)
				if err != nil {
					return nil, err
				}
				return nextHandler(ctx, newReq)
			}
			return nextHandler(ctx, req)
		}
	}
}
