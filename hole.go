package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	Black       = string([]byte{27, 91, 51, 48, 109})
	Red         = string([]byte{27, 91, 51, 49, 109})
	Green       = string([]byte{27, 91, 51, 50, 109})
	Yellow      = string([]byte{27, 91, 51, 51, 109})
	Blue        = string([]byte{27, 91, 51, 52, 109})
	Magenta     = string([]byte{27, 91, 51, 53, 109})
	Cyan        = string([]byte{27, 91, 51, 54, 109})
	White       = string([]byte{27, 91, 51, 55, 109})
	BlackWhite  = string([]byte{27, 91, 51, 48, 59, 52, 55, 109}) // 黑字白底
	YellowBlack = string([]byte{27, 91, 51, 51, 59, 52, 48, 109}) // 黄字黑底

	Reset = string([]byte{27, 91, 48, 109})
)

func main() {

	_, err := os.Stdin.Stat()

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {

		input, err := reader.ReadBytes('\n')
		if err != nil && err == io.EOF {
			break
		}

		fields := make(map[string]interface{})

		decoder := json.NewDecoder(bytes.NewReader(input))

		decoder.UseNumber()

		if err := decoder.Decode(&fields); err != nil {
			fmt.Print(string(input))
		} else {
			jsonOutput(fields)
		}
	}

}

func jsonOutput(fields map[string]interface{}) {
	if fields["method"] != nil && fields["path"] != nil {
		if fields["status"] != nil && fields["latency"] != nil {
			httpRequestEndOutput(fields)

		} else {
			httpRequestOutput(fields)
		}
	} else {
		normalOutput(fields)
	}
}

func newSpecialItem(specialKeys []string, fields map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for _, k := range specialKeys {
		res[k] = "-"
	}

	for _, key := range specialKeys {
		if val, exists := fields[key]; exists {
			res[key] = val

		} else if val, exists := fields["_"+key]; exists {
			res[key] = val
			delete(fields, "_"+key)
		}
	}
	return res
}

func httpRequestOutput(fields map[string]interface{}) {
	specialKeys := []string{"app", "time", "method", "path", "msg", "contentType", "clientIP", "traceId", "spanId", "parentSpanId", "userAgent"}
	excludeKeys := append(specialKeys, "body", "level")

	item := newSpecialItem(specialKeys, fields)

	fmt.Print("[", item["app"], "] ")
	fmt.Printf(
		"[%s][%13s][%s%s%s] [%s] [%s][%s][%s] [%s] [%s]",
		item["time"],
		item["clientIP"],
		colorForMethod(fmt.Sprint(item["method"])), item["method"], Reset, item["path"],
		item["contentType"],
		item["msg"],
		item["traceId"], item["spanId"], item["parentSpanId"],
	)

	fmt.Print("| ", item["userAgent"])
	fmt.Print("\n")

	for k, v := range fields {
		if !inSlice(k, excludeKeys) {
			fmt.Printf("  %-15s:%v\n", k, humanize(v))
		}
	}
	if body, exists := fields["body"]; exists {
		fmt.Printf("  %-15s:\n%v\n\n", "body", body)
	}
}

func httpRequestEndOutput(fields map[string]interface{}) {
	specialKeys := []string{"app", "time", "method", "path", "msg", "status", "start", "end", "latency", "comment"}
	item := newSpecialItem(specialKeys, fields)

	var latency int64

	if v, ok := item["latency"].(json.Number); ok {
		latency, _ = v.Int64()
	} else if v, ok := item["latency"].(float64); ok {
		latency = int64(v)
	}

	msg, _ := item["msg"].(string)
	if msg == "ok" {
		msg = Green + msg + Reset
	} else {
		msg = Red + msg + Reset
	}

	fmt.Print("[", item["app"], "] ")
	fmt.Printf(
		"[%s][%s%v%s][%13v][%s%s%s %s][%s%s%s %s]\n",
		item["time"],
		Green, item["status"], Reset,
		time.Duration(latency),
		colorForMethod(fmt.Sprint(item["method"])), item["method"], Reset, item["path"],
		Yellow, msg, Reset, item["comment"],
	)
}

func normalOutput(fields map[string]interface{}) {

	specialKeys := []string{"app", "time", "level", "line", "msg", "traceId", "spanId", "parentSpanId"}

	excludeKeys := specialKeys

	item := newSpecialItem(specialKeys, fields)

	fmt.Print("[", item["app"], "]")

	fmt.Printf(
		"[%s%s%s][%s][%s][%s%s%s][%s %s %s]\n",
		colorFormLevel(fmt.Sprint(item["level"])), strings.ToUpper(item["level"].(string)), Reset,
		item["time"],
		item["line"],
		colorFormLevel(fmt.Sprint(item["level"])),
		item["msg"], Reset,
		item["traceId"], item["spanId"], item["parentSpanId"],
	)

	for k, v := range fields {
		if !inSlice(k, excludeKeys) {
			fmt.Printf("  %-15s:%v\n", k, humanize(v))
		}
	}
}

func humanize(val interface{}) string {

	switch vv := val.(type) {

	case string:
		return fmt.Sprintf("%s\"%s%s%s\"%s", White, Reset, vv, White, Reset)

	case map[string]interface{}, []map[string]interface{}:
		data, err := json.MarshalIndent(vv, "", " ")
		if err != nil {
			return fmt.Sprint(vv)
		}
		return "\n" + string(data)

	default:
		return fmt.Sprint(vv)
	}
}

func colorFormLevel(level string) string {
	switch level {
	case "debug", "DEBUG":
		return Blue
	case "info", "INFO":
		return Green
	case "warn", "warning", "WARN", "WARNING":
		return Yellow
	case "error", "ERROR":
		return Red
	case "fatal", "FATAL":
		return Magenta
	case "panic", "PANIC":
		return Cyan
	default:
		return Reset
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET", "get":
		return Blue
	case "POST", "post":
		return Cyan
	case "PUT", "put":
		return Yellow
	case "DELETE", "delete":
		return Red
	case "PATCH", "patch":
		return Green
	case "HEAD", "head":
		return Magenta
	case "OPTIONS", "options":
		return BlackWhite
	default:
		return Reset
	}
}

func inSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
