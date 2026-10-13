package lemmy

import (
	lemmyRequest "LemmyBeProxy/dto/request/lemmy"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	appHttp "LemmyBeProxy/http"
	"net/url"
	"regexp"
	"strings"
)

// communityPath matches a Lemmy-style community actor path: /c/<name>.
var communityPath = regexp.MustCompile(`^/c/([^/]+)/?$`)

// parseCommunityURL returns the host and community name for a
// q like "https://lemmy.world/c/technology". ok is false for anything
// else (posts, comments, users, @-style identifiers, bare text).
func parseCommunityURL(q string) (host string, name string, ok bool) {
	u, err := url.Parse(strings.TrimSpace(q))
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" ||
		u.RawQuery != "" || u.Fragment != "" {
		return "", "", false
	}
	m := communityPath.FindStringSubmatch(u.EscapedPath())
	if m == nil {
		return "", "", false
	}
	name, err = url.PathUnescape(m[1])
	if err != nil || name == "" || strings.ContainsAny(name, "@/") {
		return "", "", false
	}
	return strings.ToLower(u.Host), name, true
}

// sameActor reports whether a community's actor_id is the URL that was
// asked for (ignoring scheme-irrelevant differences like case in the
// host and a trailing slash). Guards against name@host matching a
// different actor than the URL refers to.
func sameActor(actorId string, q string) bool {
	norm := func(s string) string {
		return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(s)), "/")
	}
	return norm(actorId) == norm(q)
}

// anonymousHeaders returns a copy of headers with every credential
// removed, keeping the rest (notably the client-IP headers).
func anonymousHeaders(headers appHttp.Headers) appHttp.Headers {
	anon := make(appHttp.Headers, len(headers))
	for key, value := range headers {
		if key == "Authorization" || key == "Cookie" {
			continue
		}
		anon[key] = value
	}
	return anon
}

// resolveCommunityLocally answers a community-URL resolve_object from
// Lemmy's own database, never from the community's home instance.
//
// Real resolve_object — and GET /community?name=... too — re-fetches a
// remote community via webfinger whenever it's more than a day stale and
// the request is logged in (Lemmy 0.19 resolve_actor_identifier). For a
// slow or dead home instance that's a 10+ second stall per community,
// repeated on every logged-in lemmyBB index view (~79 resolves).
//
// GET /community?id=... reads straight from the database, with per-user
// state when auth is passed, and never fetches remotely. So:
//  1. name@host -> id via an ANONYMOUS name lookup (anonymous lookups
//     never refetch; they return what's in the database). The id is
//     cached — community ids never change.
//  2. GET /community?id=<id> with the caller's own headers, so
//     subscribed/blocked reflect the logged-in user.
//
// ok is false when the community isn't in the local database (or the
// match is ambiguous); the caller then falls back to real resolve_object,
// which is still how a not-yet-known community gets federated in.
func (receiver *LemmyBackend) resolveCommunityLocally(q string, headers appHttp.Headers) (*lemmyResponse.ResolveObjectResponse, bool) {
	host, name, ok := parseCommunityURL(q)
	if !ok {
		return nil, false
	}
	identifier := name + "@" + host
	key := strings.ToLower(identifier)

	id, known := receiver.communityIDs.Get(key)
	if !known {
		byName, err := receiver.GetCommunity(&lemmyRequest.GetCommunityRequest{Name: &identifier}, anonymousHeaders(headers))
		if err != nil || !sameActor(byName.CommunityView.Community.ActorId, q) {
			return nil, false
		}
		id = byName.CommunityView.Community.Id
		receiver.communityIDs.Set(key, id)

		// An anonymous caller would get exactly this response from the
		// by-id lookup below — skip the second call.
		if headers["Authorization"] == "" && headers["Cookie"] == "" {
			return &lemmyResponse.ResolveObjectResponse{Community: &byName.CommunityView}, true
		}
	}

	byId, err := receiver.GetCommunity(&lemmyRequest.GetCommunityRequest{Id: &id}, headers)
	if err != nil {
		// E.g. the community was purged after its id was cached.
		receiver.communityIDs.Delete(key)
		return nil, false
	}
	return &lemmyResponse.ResolveObjectResponse{Community: &byId.CommunityView}, true
}
