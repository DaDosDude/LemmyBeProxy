package http

import "testing"

func TestClientIP(t *testing.T) {
	cases := []struct {
		name, realIP, remoteAddr, want string
	}{
		{"x-real-ip wins", "1.2.3.4", "172.17.0.2:41234", "1.2.3.4"},
		{"x-real-ip trimmed", " 1.2.3.4 ", "172.17.0.2:41234", "1.2.3.4"},
		{"ipv6 x-real-ip", "2001:db8::1", "172.17.0.2:41234", "2001:db8::1"},
		{"absent falls back to caller", "", "172.17.0.2:41234", "172.17.0.2"},
		{"garbage falls back to caller", "not-an-ip", "172.17.0.2:41234", "172.17.0.2"},
		{"list is not an ip", "1.2.3.4, 5.6.7.8", "172.17.0.2:41234", "172.17.0.2"},
		{"ipv6 caller", "", "[2001:db8::2]:443", "2001:db8::2"},
	}
	for _, c := range cases {
		if got := ClientIP(c.realIP, c.remoteAddr); got != c.want {
			t.Errorf("%s: ClientIP(%q, %q) = %q, want %q", c.name, c.realIP, c.remoteAddr, got, c.want)
		}
	}
}

func TestSetClientIPReplacesSpoofableHeaders(t *testing.T) {
	h := Headers{
		HeaderRealIP:       "1.2.3.4",
		HeaderForwardedFor: "6.6.6.6, 7.7.7.7",
		HeaderForwarded:    "for=6.6.6.6",
		"Authorization":    "Bearer x",
	}
	h.SetClientIP("1.2.3.4")

	if h[HeaderRealIP] != "1.2.3.4" || h[HeaderForwardedFor] != "1.2.3.4" {
		t.Errorf("got X-Real-Ip=%q X-Forwarded-For=%q, want both 1.2.3.4", h[HeaderRealIP], h[HeaderForwardedFor])
	}
	if _, ok := h[HeaderForwarded]; ok {
		t.Errorf("Forwarded header should have been removed")
	}
	if h["Authorization"] != "Bearer x" {
		t.Errorf("unrelated headers must be left alone")
	}
}
