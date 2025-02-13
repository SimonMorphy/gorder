package config

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// CustomFormatter 自定义日志格式
type CustomFormatter struct {
	TimestampFormat string
	ServiceName     string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 时间格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// 获取调用信息
	var caller string
	if entry.HasCaller() {
		caller = fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)
	}

	// 日志级别颜色
	levelColor := getColorByLevel(entry.Level)

	// 服务名
	serviceInfo := ""
	if f.ServiceName != "" {
		serviceInfo = fmt.Sprintf("[%s] ", f.ServiceName)
	}

	// 格式化日志头部
	fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m [\x1b[36m%s\x1b[0m] %s\x1b[33m%-24s\x1b[0m",
		levelColor,
		strings.ToUpper(entry.Level.String()),
		timestamp,
		serviceInfo,
		caller,
	)

	// 添加字段信息
	if len(entry.Data) > 0 {
		fields := make([]string, 0, len(entry.Data))
		for k, v := range entry.Data {
			fields = append(fields, fmt.Sprintf("\x1b[36m%s\x1b[0m=\x1b[35m%v\x1b[0m", k, v))
		}
		fmt.Fprintf(b, " (%s)", strings.Join(fields, ", "))
	}

	// 处理消息内容
	message := entry.Message
	if shouldPrettyPrint(message) {
		// 结构化数据，先输出基本信息行
		fmt.Fprintf(b, " \x1b[90m|\x1b[0m \x1b[32m%s\x1b[0m\n", getMessageSummary(message))
		// 然后输出格式化的详细信息
		fmt.Fprintf(b, "%s\n", prettyPrintMessage(message))
	} else {
		// 普通消息保持在同一行
		fmt.Fprintf(b, " \x1b[90m|\x1b[0m \x1b[32m%s\x1b[0m\n", message)
	}

	return b.Bytes(), nil
}

// shouldPrettyPrint 判断是否需要美化打印
func shouldPrettyPrint(message string) bool {
	return len(message) > 100 && (strings.Contains(message, "{") ||
		strings.Contains(message, "[") ||
		strings.Contains(message, "Items:") ||
		strings.Contains(message, "Response:"))
}

// getMessageSummary 获取消息的简短摘要
func getMessageSummary(message string) string {
	// 提取冒号前的内容作为摘要
	if idx := strings.Index(message, ":"); idx > 0 {
		return message[:idx+1]
	}
	// 如果消息过长，截取一部分
	if len(message) > 50 {
		return message[:50] + "..."
	}
	return message
}

// prettyPrintMessage 美化打印消息
func prettyPrintMessage(message string) string {
	// 处理常见的结构化数据格式
	message = strings.ReplaceAll(message, "Items:{", "\n    Items: {")
	message = strings.ReplaceAll(message, "} Items:{", "}\n    Items: {")
	message = strings.ReplaceAll(message, "ID:", "\n        ID: ")
	message = strings.ReplaceAll(message, "Name:", "\n        Name: ")
	message = strings.ReplaceAll(message, "Quantity:", "\n        Quantity: ")
	message = strings.ReplaceAll(message, "PriceID:", "\n        PriceID: ")

	// 移除多余的换行
	lines := strings.Split(message, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, "    "+strings.TrimSpace(line))
		}
	}

	return strings.Join(result, "\n")
}

// getColorByLevel 根据日志级别返回颜色代码
func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return 36 // cyan
	case logrus.DebugLevel:
		return 34 // blue
	case logrus.InfoLevel:
		return 32 // green
	case logrus.WarnLevel:
		return 33 // yellow
	case logrus.ErrorLevel:
		return 31 // red
	case logrus.FatalLevel, logrus.PanicLevel:
		return 35 // magenta
	default:
		return 37 // white
	}
}

// NewLogrusConfig 初始化日志配置
func NewLogrusConfig(options ...LogOption) {
	// 默认配置
	config := &logConfig{
		level:           logrus.InfoLevel,
		serviceName:     "",
		reportCaller:    true,
		timestampFormat: time.RFC3339,
	}

	// 应用自定义选项
	for _, opt := range options {
		opt(config)
	}

	// 设置日志格式化器
	logrus.SetFormatter(&CustomFormatter{
		TimestampFormat: config.timestampFormat,
		ServiceName:     config.serviceName,
	})

	// 设置基本配置
	logrus.SetLevel(config.level)
	logrus.SetReportCaller(config.reportCaller)
}

// logConfig 日志配置项
type logConfig struct {
	level           logrus.Level
	serviceName     string
	reportCaller    bool
	timestampFormat string
}

// LogOption 配置选项函数类型
type LogOption func(*logConfig)

// WithLevel 设置日志级别
func WithLevel(level logrus.Level) LogOption {
	return func(c *logConfig) {
		c.level = level
	}
}

// WithServiceName 设置服务名
func WithServiceName(name string) LogOption {
	return func(c *logConfig) {
		c.serviceName = name
	}
}

// WithReportCaller 设置是否记录调用信息
func WithReportCaller(report bool) LogOption {
	return func(c *logConfig) {
		c.reportCaller = report
	}
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(format string) LogOption {
	return func(c *logConfig) {
		c.timestampFormat = format
	}
}
