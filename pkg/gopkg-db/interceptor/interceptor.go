package interceptor

import (
	"context"
	db2 "effectivemobile/pkg/gopkg-db"
	"net/http"

	"google.golang.org/grpc"
)

// DBClientProvider is a func for wrapping database
type DBClientProvider func(ctx context.Context) db2.IClient

// NewDBInterceptor wrap endpoint with middleware mixing in db connection
func NewDBInterceptor(dbClientProvider DBClientProvider, option ...db2.Option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(db2.AddToContext(ctx, dbClientProvider(ctx)), req)
	}
}

// NewDBServerMiddleware wrap endpoint with middleware mixing in db connection
func NewDBServerMiddleware(dbClientProvider DBClientProvider, option ...db2.Option) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			r = r.WithContext(db2.AddToContext(ctx, dbClientProvider(ctx)))
			next.ServeHTTP(w, r)
		})
	}
}
