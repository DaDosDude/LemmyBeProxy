# LemmyBeProxy

A proxy that lets Lemmy clients (mlmym, lemmyBB, Voyager, Alexandrite,
official apps) talk to a PieFed instance or a different Lemmy version.
It serves Lemmy's API at `/api/v3/*` and `/pictrs/image`, and translates
each request to whatever backend sits behind it.

- **Backends:** PieFed (`/api/alpha`) or real Lemmy.
- **Client wire formats:** Lemmy 0.19 (current) and 0.17 (e.g. lemmyBB),
  plus 0.18-style `auth` in the body or query string.

## Quick start

```
docker build -t lemmybeproxy .

docker run -d --name lemmybeproxy --restart unless-stopped \
  -p 127.0.0.1:8050:8080 \
  -e BACKEND_TYPE=piefed \
  -e BACKEND_INSTANCE=your-instance.example \
  -e FRONTEND_VERSION=0.19 \
  lemmybeproxy
```

Then point your reverse proxy's `/api/v3/*` and `/pictrs/image*` routes
at the container, and everything else at your frontend (e.g. mlmym).
Nginx and Caddy examples are in [DEPLOYMENT.md](DEPLOYMENT.md).

Check it works:

```
curl -s https://your-domain.example/api/v3/site | head -c 200
```

## Configuration

| Variable | Required | Description |
|---|---|---|
| `BACKEND_TYPE` | Yes | `piefed` or `lemmy`. |
| `BACKEND_INSTANCE` | Yes | Backend hostname (`https://` assumed), or a full URL such as `http://lemmy:8536` to reach it directly on its internal Docker network. |
| `FRONTEND_VERSION` | Yes | `0.19` or `0.17` — the wire format your clients speak. |
| `PORT` | No | Listen port. Default `8080`. |
| `SIMULATE` | No | Report itself as a Lemmy server in version info. |
| `RESOLVE_OBJECT_CACHE_TTL` | No | Seconds to cache anonymous `resolve_object` results. Default `600`, `0` disables. |
| `RESOLVE_OBJECT_CACHE_SIZE` | No | Max cached `resolve_object` entries. Default `2000`. |
| `EXTRA_IMAGE_HOSTS` | No | Comma-separated extra hosts allowed for `/pictrs/image` redirects (e.g. a separate media/CDN host). |

## Reverse proxy requirements

- **Set `X-Real-IP` to the visitor's address** (Caddy:
  `header_up X-Real-IP {client_ip}`; nginx:
  `proxy_set_header X-Real-IP $remote_addr;`). The proxy forwards it to
  the backend as `X-Real-IP` and `X-Forwarded-For`, so the backend
  rate-limits per visitor instead of treating everyone as one client.
  Client-supplied `X-Forwarded-For`/`Forwarded` headers are ignored.
- **Don't expose the container port publicly.** Bind it to `127.0.0.1`
  or a private network, so `X-Real-IP` can only come from your proxy.
- **If `BACKEND_INSTANCE` is a public domain**, the reverse proxy in
  front of the backend must trust this proxy's address (Caddy:
  `trusted_proxies`), or it overwrites the visitor IP. Using the
  backend's internal address avoids this.

## Behaviour worth knowing

- **Anonymous `resolve_object` results are cached** (see the TTL
  setting). Requests with an `Authorization` header or a cookie are
  never cached. Cached responses carry `X-Proxy-Cache: HIT`/`MISS`.
- **Lemmy backend:** a `resolve_object` for a community URL that Lemmy
  already knows is answered from its database, without contacting the
  community's home instance. Unknown communities, posts and users still
  use the real `resolve_object`.
- **Images** uploaded through `/pictrs/image` are served by redirecting
  to the backend; with an internal `http://` backend address the proxy
  streams them itself.

## Known limitations

- **PieFed:** `user/mention`, `user/replies` and `private_message/list`
  return empty lists; `save_user_settings` drops fields PieFed doesn't
  have; no captcha, report counts, password reset, TOTP or account
  deletion; image thumbnails load at full size.
- **0.17 clients:** `mark_as_read` returns `{success}` instead of a post;
  `form_id` isn't echoed back.
- **Clients that don't send `X-Real-IP`** (e.g. lemmyBB's server-side
  calls) reach the backend as one client and share its rate limit.
- **Remote objects** resolve only for logged-in users (Lemmy's own rule).

## Troubleshooting

- **Everyone gets HTTP 429 at once:** the backend sees all proxied
  traffic as one client — see Reverse proxy requirements.
- **Backend requests fail or get blocked from the same host:** point
  `BACKEND_INSTANCE` at the backend's internal address instead of its
  public domain.
- **Slow frontend on every page:** the frontend container is reaching
  the API via the public internet; use `--add-host
  yourdomain:host-gateway` (see DEPLOYMENT.md).
- **Build fails with a `slices`/GOROOT error:** needs Go 1.24+.
