package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func printReqHeaders(headers http.Header) string {
	var result strings.Builder
	for h, v := range headers {
		fmt.Fprintf(&result, "%s: %v\n", h, v)
	}
	return result.String()
}

// clientIP returns the originating client address. Requests arrive through
// Cloudflare and then Caddy, so r.RemoteAddr is the proxy's address and the
// first X-Forwarded-For hop is a Cloudflare edge IP, not the attacker's.
// Cloudflare carries the real client address in Cf-Connecting-Ip; fall back to
// X-Forwarded-For, then r.RemoteAddr, when it is absent.
func clientIP(r *http.Request) string {
	if cf := r.Header.Get("Cf-Connecting-Ip"); cf != "" {
		return strings.TrimSpace(cf)
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
	}
	return r.RemoteAddr
}

type LogFormat struct {
	EventName string `json:"eventName"`
	Actor     struct {
		Id string `json:"id"`
	} `json:"actor"`
	Resources struct {
		Headers string `json:"headers"`
	} `json:"resources"`
	ReqProto string `json:"reqProto"`
	ReqIP    string `json:"srcIP"`
	ReqPath  string `json:"reqPath"`
	Request  string `json:"request"`
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := sha256.Sum256([]byte(fmt.Sprintf("%s-%s", clientIP(r), r.UserAgent())))
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("failed to read request body", "error", err)
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		jsonTxt := LogFormat{
			ReqProto:  r.Proto,
			ReqIP:     clientIP(r),
			ReqPath:   r.RequestURI,
			EventName: r.Method,
			Resources: struct {
				Headers string `json:"headers"`
			}{
				Headers: printReqHeaders(r.Header),
			},
			Actor: struct {
				Id string `json:"id"`
			}{
				Id: hex.EncodeToString(userID[:]),
			},
			Request: string(body),
		}
		out, err := json.Marshal(jsonTxt)
		if err != nil {
			slog.Error("failed to marshal request data", "error", err)
		} else if _, err := os.Stdout.Write(append(out, '\n')); err != nil {
			slog.Error("failed to write request log", "error", err)
		}
		next.ServeHTTP(w, r)
	})
}
