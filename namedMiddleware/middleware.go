package namedMiddleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
)

// WrapKratosMiddleware 封装kratos的中间件为命名中间件
// wrap kratos middleware to named middleware
func WrapKratosMiddleware(name string, mw middleware.Middleware) middleware.Middleware {
	return func(nextHandler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if caller := match(ctx, name); caller != nil {
				return mw(nextHandler)(ctx, req)
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

// Handler 命名kratos中间件的处理函数(带参数)
// the handler function of named kratos middleware with arguments
type Handler func(ctx context.Context, req interface{}, arguments ...string) (interface{}, error)

// MiddlewareWithArguments 注册一个带参数的命名kratos中间件
// register a named kratos middleware with arguments
func MiddlewareWithArguments(name string, handler Handler) middleware.Middleware {
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
