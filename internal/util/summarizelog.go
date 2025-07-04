package util

// app/catalog_service.go
import (
	"strings"
	"unicode/utf8"
)

// SummarizeLog 返回 chunk 摘要
func SummarizeLog(content string, length int) string {
	if content == "" {
		return "[空]"
	}

	runes := []rune(content)
	if len(runes) <= length {
		return string(runes)
	}

	// 截取前 length 个字符，并确保是完整的 Unicode 字符
	truncated := string(runes[:length])
	if !utf8.ValidString(truncated) {
		truncated = string([]rune(truncated)[:len(truncated)-1]) // 去掉可能不完整的最后一个字符
	}

	// 可选：避免在句号、逗号等符号后截断
	if lastPunct := strings.LastIndexAny(truncated, "。.，,；;"); lastPunct > 0 {
		truncated = truncated[:lastPunct+1]
	}

	return truncated + "..."
}
