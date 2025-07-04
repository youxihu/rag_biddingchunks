package util

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const ClientIPKey contextKey = "client_ip"

// WrapMCPHandler 是一个封装函数，接收原始的 MCP HTTP Handler，并返回一个带 IP 记录的新 Handler
func WrapMCPHandler(original http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 提取客户端 IP
		ip := getClientIP(r)

		// 把 IP 存入 Context
		ctx := context.WithValue(r.Context(), ClientIPKey, ip)
		r = r.WithContext(ctx)

		// 调用原始 MCP handler
		original.ServeHTTP(w, r)
	})
}

// getClientIP 获取客户端真实 IP（支持代理）
func getClientIP(r *http.Request) string {
	xfw := r.Header.Get("X-Forwarded-For")
	if xfw != "" {
		ips := strings.Split(xfw, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	ipPort := r.RemoteAddr
	ip := strings.Split(ipPort, ":")[0]
	return ip
}

// LogWithIP 是一个封装后的日志函数，自动带上 IP
func LogWithIP(ctx context.Context, format string, v ...interface{}) {
	ip := "unknown"
	if ctx != nil {
		if val := ctx.Value(ClientIPKey); val != nil {
			ip = val.(string)
		}
	}

	// 自动在日志开头加上 IP
	log.Printf("[来自 IP: %s] "+format, append([]interface{}{ip}, v...)...)
}
