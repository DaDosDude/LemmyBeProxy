package http

import (
	"net"
	goHttp "net/http"
	"strings"
)

// Header names in Go's canonical form — the form incoming header keys
// already have by the time they reach Headers (net/http canonicalizes
// them), so lookups and overwrites hit the same key.
const (
	HeaderRealIP       = "X-Real-Ip"
	HeaderForwardedFor = "X-Forwarded-For"
	HeaderForwarded    = "Forwarded"
)

// ClientIP determines the real visitor's address for an incoming request.
//
// X-Real-IP is the source of truth: the front proxy (Caddy, or mlmym,
// which passes its own caller's X-Real-IP along) sets it, and Caddy
// overwrites any client-supplied value with the actual remote address,
// so it can't be forged from outside. A client-supplied X-Forwarded-For
// is deliberately NOT consulted — anyone calling /api/v3 can put any
// address in it. If X-Real-IP is absent or isn't a valid IP, the direct
// caller's own address (remoteAddr) is used, which is this proxy's
// behaviour before client-IP forwarding existed.
func ClientIP(realIP string, remoteAddr string) string {
	if ip := net.ParseIP(strings.TrimSpace(realIP)); ip != nil {
		return ip.String()
	}

	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

// SetClientIP rewrites the client-IP headers that get forwarded to the
// backend so they carry exactly one address: the real visitor's.
//
// X-Forwarded-For is replaced rather than appended to. Real Lemmy
// (actix-web) takes the client address from the Forwarded header first,
// then the *leftmost* X-Forwarded-For entry — so appending to a
// client-supplied chain would let any caller choose the IP the backend
// rate-limits on. Forwarded is removed for the same reason. Backends that
// read X-Real-IP (or the rightmost X-Forwarded-For entry) get the same
// address either way.
func (h Headers) SetClientIP(ip string) {
	delete(h, HeaderForwarded)
	if ip == "" {
		delete(h, HeaderRealIP)
		delete(h, HeaderForwardedFor)
		return
	}
	h[HeaderRealIP] = ip
	h[HeaderForwardedFor] = ip
}

// ApplyClientIP copies the client-IP headers onto an outgoing request
// built outside the normal sendRequest path (image upload and fetch),
// which don't forward the full incoming header set.
func (h Headers) ApplyClientIP(req *goHttp.Request) {
	for _, key := range []string{HeaderRealIP, HeaderForwardedFor} {
		if value, ok := h[key]; ok && value != "" {
			req.Header.Set(key, value)
		}
	}
}
