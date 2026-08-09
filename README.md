# LemmyBeProxy

A simple proxy that acts as a compatibility layer for Lemmy apps against Piefed and newer versions of Lemmy.

It sits in front of a configurable backend and speaks Lemmy's API to
whatever's talking to it, regardless of what's actually behind it or how
old the client is. Any Lemmy-compatible app or frontend (mlmym,
Alexandrite, Voyager, official Lemmy apps, lemmyBB, and so on) can talk
to this proxy without knowing or caring what it's really connected to.

The goal is a genuine compatibility layer across Lemmy API generations —
not just Piefed vs. real Lemmy, but old clients against new backends and
vice versa, so nothing breaks outright just because either side moved
forward. That work is ongoing; see Features below for exactly where it
currently stands.

## How it works

```mermaid
sequenceDiagram
    actor User
    participant Client as Lemmy client
    participant Proxy as LemmyBeProxy
    participant Backend as Configured backend

    User->>Client: Visit a page using a Lemmy client
    Client->>Proxy: Send request to proxy
    Proxy->>Backend: Transform to backend's request shape
    Backend->>Proxy: Return backend response
    Proxy->>Client: Transform to the client's expected shape
    Client->>User: User sees content as if on Lemmy
```

The proxy exposes Lemmy's API at `/api/v3/*` and image uploads at
`/pictrs/image`, the same paths a real Lemmy server uses. Internally,
`BACKEND_TYPE` picks what it talks to and `FRONTEND_VERSION` picks the
wire format it accepts from and replies to clients with — the two are
independent, and neither takes a default. Against a real Lemmy backend
on the current wire format, translation is close to a no-op, since this
proxy's internal shapes are already modeled on Lemmy's own API.

## Quick start

```
docker build -t lemmybeproxy .

docker run -d \
  --name lemmybeproxy \
  --restart unless-stopped \
  -p 127.0.0.1:8050:8080 \
  -e BACKEND_TYPE=piefed \
  -e BACKEND_INSTANCE=your-instance.example \
  -e FRONTEND_VERSION=0.19 \
  lemmybeproxy
```

`BACKEND_TYPE` is `piefed` or `lemmy`; `BACKEND_INSTANCE` is that
backend's hostname (or a full internal URL — see Environment variables).
`FRONTEND_VERSION` is `0.19`
(current Lemmy wire format) or `0.17` (older format, e.g. lemmyBB). All
three are required — see Environment variables for the full list.

Point a reverse proxy's `/api/v3/*` and `/pictrs/image*` routes at this
container, and everything else at your actual frontend (e.g. mlmym).
Nginx and Caddy examples, plus the full walkthrough including automatic
updates via Watchtower, are in [DEPLOYMENT.md](DEPLOYMENT.md).

Verify it's working:

```
curl -s "https://your-domain.example/api/v3/site" | head -c 200
```

## Features

**Backends** — what this proxy can translate *to*:
- Piefed, via `/api/alpha/*`
- Real Lemmy, close to a passthrough since the shapes already match

**Frontend / wire compatibility** — what this proxy can accept *from* clients:
- Lemmy 0.19.x (current) — full coverage
- Lemmy 0.18.x authentication convention (`auth` in the request body or
  query string instead of an `Authorization` header) — applied centrally
  to every route, not endpoint-specific
- Lemmy 0.17.x wire format (e.g. lemmyBB) — Post, Comment, Community,
  User, Search, and Site endpoints, all verified against Lemmy's real
  0.17.2 source *and* confirmed end-to-end against the real, official,
  unmodified lemmyBB client — not just source-level checking. `pictrs/image`
  needs no version-specific handling at all and works as-is — confirmed
  directly against lemmyBB's own upload code, which uses the identical
  field name and response shape. That's full 0.17.x coverage of every
  endpoint this proxy implements.

**Endpoints implemented and tested against live Piefed and real Lemmy
instances:** `user/login`, `user/register`, `user/get_captcha`,
`user/unread_count`, `user`, `user/block`,
`user/save_user_settings` (Piefed only supports a subset of fields — see
Features not working yet), `user/mention`, `user/replies`,
`private_message/list` (Piefed has no equivalent for these three —
returns an honest empty list rather than a 404, see Features not
working yet), `site`, `resolve_object`, `post/list`, `post`
(get/create/edit), `post/like`, `post/mark_as_read`, `comment/list`,
`comment` (get/create), `comment/like`, `community`, `community/list`,
`community/follow`, `community/block`, `search`, `pictrs/image` (upload
and fetch, outside `/api/v3` since that's where real clients send them).

## Features not working yet

**Backend pluggability is done — every endpoint this proxy implements
works against either Piefed or real Lemmy.** No controller is left
hardcoded to a Piefed-shaped client.

**Frontend pluggability (0.17.x wire format) is done too — every
endpoint this proxy implements works on both wire formats.** Upload
doesn't need Frontend-axis work at all, since `pictrs/image` is
version-agnostic by nature. Both axes — which backend, and which wire
format — are now genuinely independent and fully covered for the entire
API surface this proxy implements.

**Not implemented in Piefed itself**, confirmed from Piefed's own source
— not a translation gap, there's nothing on Piefed's side to translate
to: report count, password reset/change, TOTP, account
deletion, email verification, admin tools, custom emoji.

**Known limitations:**
- Image thumbnails always load at full resolution — Piefed ignores the
  `?format=jpg&thumbnail=96`-style query params real pict-rs understands.
- `save_user_settings` silently drops fields Piefed has no equivalent
  for (most notably mlmym's "endless scrolling," which can never persist
  server-side against Piefed) rather than failing the whole save.
- `user/mention`, `user/replies`, and `private_message/list` always
  return an empty list against Piefed — Piefed has unread *counts* for
  these but no equivalent endpoint to fetch the actual items in the same
  shape Lemmy has them. Real Lemmy backends return the genuine list.
  Returning empty rather than 404 matters concretely: a 404 here is
  exactly what crashes a real Lemmy client (lemmyBB included), since
  these are fetched on every authenticated page load.
- `resolve_object` (used by lemmyBB's frontpage category configuration,
  among other things) only resolves *remote/federated* objects when the
  request is authenticated — confirmed as genuine Lemmy API behavior,
  not a proxy limitation. A community you're personally subscribed to
  resolves fine for you, but won't resolve for anonymous visitors unless
  it's local to the instance you're querying.
- `user/get_captcha` always returns no captcha required against Piefed —
  Piefed has no image/audio captcha system to forward to, so this is an
  honest "none needed" rather than a faked one. Real Lemmy backends
  return a genuine captcha challenge when the instance has one enabled
  (confirmed against retrolemmy.com, which does).
- 0.17.x `mark_as_read` returns the canonical `{success: bool}` shape
  instead of real 0.17.x's `PostResponse{post_view}` — building the
  latter needs an extra backend round-trip, deferred rather than rushed.
- 0.17.x comment responses don't echo back a client's `form_id` (used
  for optimistic-UI correlation) — accepted on requests, not threaded
  through to responses yet.

**Real interop lessons**, found by testing against actual live
software — a real client (lemmyBB) and a real backend (retrolemmy.com)
— not just checking source-level shapes, which alone wasn't enough to
catch these:
- Real Lemmy's own pict-rs returns `201 Created` for a successful image
  upload, not `200 OK` — `LemmyBackend.UploadImage` only accepted 200
  at first, silently treating every successful upload as a failure
  against a real Lemmy backend specifically (Piefed's own upload
  behavior happened to return 200, which is why this went unnoticed
  until testing against retrolemmy.com directly).
- Timestamps must not carry a timezone suffix for 0.17.x — real Lemmy
  0.17.x uses `chrono::NaiveDateTime` (no timezone concept) throughout,
  and its deserializer errors on one if present. See
  `helper.StripTimezoneSuffixes`.
- `Option<T>` fields in Rust's `#[derive(Deserialize)]` still require
  the JSON *key* to be present (as `null`) unless the struct opts into
  `#[serde(default)]`, which none of Lemmy's real structs do — a missing
  key isn't the same as a `null` value here. Affects every 0.17-shaped
  type with an optional field.
- `RegistrationMode` is the one enum in Lemmy 0.17.2 with
  `#[serde(rename_all = "lowercase")]` — its wire values are lowercase,
  unlike every other enum in the API, which uses standard PascalCase.
- A response's `Content-Length` header needs to be set explicitly —
  without it, Go's server defaults to `Transfer-Encoding: chunked`,
  which Rust's `reqwest` (lemmyBB's HTTP client) read as truncated even
  though `curl` read the identical response correctly. See
  `WriteHttpResponse`.
- A required timestamp field populated with a placeholder empty string
  (rather than real or reasonably-derived data) fails outright on a
  strict Rust deserializer — `Site.LastRefreshedAt` was previously left
  as `""` with a `// todo`; it now derives from the ActivityPub actor's
  own `updated`/`published` timestamp instead. See `ConvertSite`.
- A response field genuinely absent in some real cases needs to be
  modeled as a pointer, not a plain value with an empty-string
  convention — `LoginResponse.Jwt` is `Option<String>` in real Lemmy,
  since registration can succeed without an immediate JWT when an
  application needs admin approval first (retrolemmy.com's own
  `requireapplication` mode is a live example).
- A field's real Rust type isn't always what its name suggests — real
  Lemmy's `PostAggregates.unread_comments` is a signed `i64`, not
  unsigned, and genuinely sends `-1` in a real edge case (confirmed: an
  old post in retrolemmy.com's own `general` community). It's a
  per-user field only calculated when logged in, so an unsigned type
  here crashed with "internal lemmy layer proxy error" specifically for
  authenticated requests — anonymous requests never triggered it,
  making this easy to miss without testing both states live.

## Environment variables

| Variable | Required | Description |
|---|---|---|
| `BACKEND_TYPE` | Yes | `piefed` or `lemmy`. No default. |
| `BACKEND_INSTANCE` | Yes | Hostname of the backend, no trailing slash. Defaults to `https://`. Can instead be given a full `http://` or `https://` URL to reach an internal address directly (e.g. `http://lemmy:8536` on the backend's own Docker network) — bypasses a public domain's reverse proxy and whatever's in front of it (anti-abuse layers, CDNs) entirely. Useful when the proxy runs on the same host/network as the real backend; see Troubleshooting for why this can matter. |
| `FRONTEND_VERSION` | Yes | `0.19` or `0.17`. No default. |
| `PORT` | No | Port to listen on. Defaults to `8080`. |
| `SIMULATE` | No | If set, the proxy presents itself as a Lemmy server when asked for version/software info. |

mlmym (or whatever frontend you run alongside this) has its own env
vars, including `LEMMY_DOMAIN`, `COLLAPSE_MEDIA`, `HIDE_THUMBNAILS` — see
that project's own documentation.

## Troubleshooting

**Requests to the backend fail or get rate-limited for no obvious
reason, even though the exact same request works fine from elsewhere**
— if the proxy runs on the same server as the backend itself, and
`BACKEND_INSTANCE` is set to that server's own public domain, outbound
requests can hairpin through the host's own public IP and get treated
as suspicious by the site's own anti-abuse layer — confirmed
reproducible with plain `curl`, nothing specific to this proxy's
client. Point `BACKEND_INSTANCE` at the backend's internal address
directly instead (e.g. `http://lemmy:8536` on its own Docker network) —
see the `BACKEND_INSTANCE` row above.

**Parse error mentioning an invalid character around byte `0x1f`** — a
gzip-compressed response got handed to the JSON parser undecompressed.
Caused by forwarding a client's own `Accept-Encoding` header onto the
outbound request, which disables Go's automatic gzip handling. This
project already strips that header before forwarding; keep that in mind
if you add new outbound requests elsewhere.

**Frontend is slow on every page load**, despite the API responding
fast when tested directly — the frontend container is likely reaching
its own API calls by round-tripping out to the public internet instead
of staying local. Fixed with `--add-host yourdomain.com:host-gateway`
on the frontend container — see DEPLOYMENT.md.

**A specific sort or listing combination fails or is unusually slow** —
Piefed validates sort values per endpoint, and the valid set differs
between them (community listing accepts fewer values than post
listing, for example). Check Piefed's own error message for what's
actually valid on that specific endpoint.

**Build fails with a `slices` package / GOROOT error** — an old Go
toolchain is active, usually one installed via a distro package manager
instead of the official binary. Needs Go 1.24+; check `go version`.

## Contributing

This project translates between APIs that are similar but not
identical, and the real gaps are usually only discoverable by testing
against a live instance and comparing actual response shapes and
validation rules to what the code assumes — this project's own history
includes cases where a schema's own documentation didn't match what its
API validated at runtime. If you add an endpoint, test it live before
trusting the schema alone.
