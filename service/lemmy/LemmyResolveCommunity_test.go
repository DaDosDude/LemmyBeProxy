package lemmy

import "testing"

func TestParseCommunityURL(t *testing.T) {
	cases := []struct {
		q, host, name string
		ok            bool
	}{
		{"https://lemmy.world/c/technology", "lemmy.world", "technology", true},
		{"https://Lemmy.World/c/Technology/", "lemmy.world", "Technology", true},
		{"https://startrek.website/c/risa", "startrek.website", "risa", true},
		{"http://lemmy:8536/c/general", "lemmy:8536", "general", true},
		{"https://lemmy.world/post/123", "", "", false},
		{"https://lemmy.world/c/technology/p/1", "", "", false},
		{"https://lemmy.world/u/someone", "", "", false},
		{"https://lemmy.world/c/technology?x=1", "", "", false},
		{"https://lemmy.world/c/", "", "", false},
		{"!technology@lemmy.world", "", "", false},
		{"technology", "", "", false},
		{"ftp://lemmy.world/c/technology", "", "", false},
	}
	for _, c := range cases {
		host, name, ok := parseCommunityURL(c.q)
		if ok != c.ok || host != c.host || name != c.name {
			t.Errorf("parseCommunityURL(%q) = (%q, %q, %v), want (%q, %q, %v)", c.q, host, name, ok, c.host, c.name, c.ok)
		}
	}
}

func TestSameActor(t *testing.T) {
	if !sameActor("https://lemmy.world/c/technology", "https://Lemmy.World/c/technology/") {
		t.Error("case and trailing slash must not matter")
	}
	if sameActor("https://lemmy.world/c/technology", "https://lemmy.world/c/tech") {
		t.Error("different actors must not match")
	}
}
