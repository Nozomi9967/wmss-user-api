package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RpcMetaInterceptor 自动将 context 中的 userId 和 username 传递到 RPC 服务的 metadata
func RpcMetaInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 从 context 中提取 API 中间件存入的 userId 和 username
		userId, _ := ctx.Value("user_id").(string)
		username, _ := ctx.Value("username").(string)

		// 构建 metadata（key 用小写，gRPC 会自动转为小写）
		md := metadata.New(map[string]string{})
		if userId != "" {
			md.Set("user-id", userId)
		}
		if username != "" {
			md.Set("username", username)
		}

		// 将 metadata 附加到 outgoing context
		ctxWithMD := metadata.NewOutgoingContext(ctx, md)

		// 继续执行 RPC 调用
		return invoker(ctxWithMD, method, req, reply, cc, opts...)
	}
}
