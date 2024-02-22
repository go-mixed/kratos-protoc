package namedMiddleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Handler 命名kratos中间件的处理函数(带参数)
// the handler function of named kratos middleware with arguments
type Handler func(ctx context.Context, req interface{}, arguments ...string) (interface{}, error)

// HandlerWithArguments 将一个named.Handler转换带参数的命名kratos中间件
// turn a named.Handler to a named kratos middleware with arguments
func HandlerWithArguments(name string, handler Handler) middleware.Middleware {
	return func(nextHandler middleware.Handler) middleware.Handler {
		// ctx 为当前请求的上下文
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 如果当前请求的上下文的中间件包含当前中间件，取出参数执行
			if caller := match(ctx, name); caller != nil {
				res, err := handler(ctx, req, caller.arguments...)
				if err != nil {
					return nil, err
				}
				return nextHandler(ctx, res)
			}
			// 否则执行下一个中间件
			return nextHandler(ctx, req)
		}
	}
}

// KratosMiddleware 将一个kratos的中间件转换为命名kratos的中间件
// turn a kratos middleware to a named kratos middleware
func KratosMiddleware(name string, mw middleware.Middleware) middleware.Middleware {
	return func(nextHandler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if caller := match(ctx, name); caller != nil {
				return mw(nextHandler)(ctx, req)
			}
			return nextHandler(ctx, req)
		}
	}
}

// KratosHandler 将一个kratos的Handler转换为命名kratos的中间件
// turn a kratos.Handler into a named kratos middleware
func KratosHandler(name string, handler middleware.Handler) middleware.Middleware {
	return func(nextHandler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if caller := match(ctx, name); caller != nil {
				res, err := handler(ctx, req)
				if err != nil {
					return nil, err
				}
				return nextHandler(ctx, res)
			}
			return nextHandler(ctx, req)
		}
	}
}
