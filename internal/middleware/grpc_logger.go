package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// key untuk context
type ctxKey string

const (
	grpcLogFieldsKey ctxKey = "grpc_log_fields"
	grpcRequestIDKey ctxKey = "grpc_request_id"
)

// AddLogFields memungkinkan middleware lain menambahkan field ke log utama
func AddLogFields(ctx context.Context, fields ...slog.Attr) {
	if v, ok := ctx.Value(grpcLogFieldsKey).(*[]slog.Attr); ok {
		*v = append(*v, fields...)
	}
}

// GetRequestID mengambil request ID dari context (untuk dipakai di handler/service)
func GetGrpcRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(grpcRequestIDKey).(string); ok {
		return v
	}
	return ""
}

type GrpcLoggerInterceptor struct {
	logger *slog.Logger
}

func NewGrpcLoggerInterceptor(logger *slog.Logger) *GrpcLoggerInterceptor {
	return &GrpcLoggerInterceptor{logger: logger}
}

func (i *GrpcLoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. Generate Request ID (UUID sederhana)
		reqID := generateRequestID()
		ctx = context.WithValue(ctx, grpcRequestIDKey, reqID)
		start := time.Now()

		// 2. Siapkan wadah untuk log tambahan (username, warehouse_id, dll)
		var extraFields []slog.Attr
		ctx = context.WithValue(ctx, grpcLogFieldsKey, &extraFields)

		// Panggil handler berikutnya
		resp, err := handler(ctx, req)

		duration := time.Since(start)

		// Siapkan field log dasar
		logFields := []any{
			slog.String("request_id", reqID),
			slog.String("method", info.FullMethod),
			slog.Duration("duration", duration),
		}

		// 3. Masukkan field tambahan dari middleware lain
		for _, attr := range extraFields {
			logFields = append(logFields, attr)
		}

		if err != nil {
			st, _ := status.FromError(err)
			logFields = append(logFields, slog.String("status", st.Code().String()), slog.String("error", err.Error()))
			i.logger.Error("gRPC Request Failed", logFields...)
		} else {
			logFields = append(logFields, slog.String("status", "OK"))
			i.logger.Info("gRPC Request Success", logFields...)
		}

		return resp, err
	}
}

func generateRequestID() string {
	b := make([]byte, 8) // 16 karakter hex
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
