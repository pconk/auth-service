package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// key untuk context
type ctxKey string

const grpcRequestIDKey ctxKey = "grpc_request_id"

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
		// 1. Cek apakah ada Request ID di metadata (dikirim oleh Gateway)
		md, ok := metadata.FromIncomingContext(ctx)
		var reqID string
		if ok {
			if ids := md.Get("x-request-id"); len(ids) > 0 {
				reqID = ids[0]
			}
		}

		// 2. Jika tidak ada di metadata, baru generate sendiri
		if reqID == "" {
			reqID = generateRequestID()
		}

		// Simpan ID ke context gRPC
		ctx = context.WithValue(ctx, grpcRequestIDKey, reqID)
		start := time.Now()

		// 2. Siapkan wadah untuk log tambahan (username, warehouse_id, dll)
		var extraFields []slog.Attr
		ctx = context.WithValue(ctx, LogFieldsKey, &extraFields)

		// 3. Kirim balik Request ID ke Client (Gateway) lewat Header Response
		header := metadata.Pairs("x-request-id", reqID)
		grpc.SendHeader(ctx, header)

		// Panggil handler berikutnya
		resp, err := handler(ctx, req)

		// Tentukan status, level, dan message
		level := slog.LevelInfo
		statusStr := "OK"
		msg := "gRPC Request Success"

		if err != nil {
			st, _ := status.FromError(err)
			statusStr = st.Code().String()
			level = slog.LevelError
			msg = "gRPC Request Failed"
			// Masukkan error ke dalam extra fields agar masuk ke group trace
			extraFields = append(extraFields, slog.String("error", err.Error()))
		}

		// Log dengan struktur Grouping (trace & grpc)
		i.logger.LogAttrs(ctx, level, msg,
			CreateTraceGroup(reqID, extraFields),
			slog.Group("grpc",
				slog.String("method", info.FullMethod),
				slog.String("status", statusStr),
				slog.Float64("duration_ms", DurationToMs(time.Since(start))),
			),
		)

		return resp, err
	}
}

func generateRequestID() string {
	b := make([]byte, 8) // 16 karakter hex
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
