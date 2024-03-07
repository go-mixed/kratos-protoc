package namedMiddleware

import "context"

type namedMiddlewareKey struct{}

type namedMiddleware struct {
	middlewareCallers []*middlewareCaller
}

type middlewareCaller struct {
	name      string
	arguments []string
}

func (nm *namedMiddleware) dispatch(name string, arguments ...string) {
	nm.middlewareCallers = append(nm.middlewareCallers, &middlewareCaller{
		name:      name,
		arguments: arguments,
	})
}

func fromContext(ctx context.Context) *namedMiddleware {
	return ctx.Value(namedMiddlewareKey{}).(*namedMiddleware)
}

func match(ctx context.Context, name string) *middlewareCaller {
	nw := fromContext(ctx)
	if nw == nil {
		return nil
	}
	for _, middleware := range nw.middlewareCallers {
		if middleware.name == name {
			return middleware
		}
	}
	return nil
}

func newContext(ctx context.Context, nm *namedMiddleware) context.Context {
	return context.WithValue(ctx, namedMiddlewareKey{}, nm)
}
