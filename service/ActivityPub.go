package service

import (
	"LemmyBeProxy/dto/model/ap"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// actorCacheTTL bounds how stale a cached actor may be. The site actor
// (the only one fetched) changes only when an admin edits the site's
// name/description/keys, while /site is requested on nearly every page
// load — previously each of those triggered a fresh outbound HTTPS fetch.
const actorCacheTTL = 10 * time.Minute

type cachedActor struct {
	actor   ap.Actor
	fetched time.Time
}

type ActivityPub struct {
	mu    sync.Mutex
	cache map[string]cachedActor
}

func NewActivityPub() *ActivityPub {
	return &ActivityPub{cache: make(map[string]cachedActor)}
}

func (receiver *ActivityPub) FetchActor(actorId string) (ap.Actor, error) {
	receiver.mu.Lock()
	entry, ok := receiver.cache[actorId]
	receiver.mu.Unlock()
	if ok && time.Since(entry.fetched) < actorCacheTTL {
		return entry.actor, nil
	}

	actor, err := receiver.fetch(actorId)
	if err != nil {
		// Serve a stale copy rather than failing /site outright if a
		// refresh fails transiently.
		if ok {
			return entry.actor, nil
		}
		return actor, err
	}

	receiver.mu.Lock()
	receiver.cache[actorId] = cachedActor{actor: actor, fetched: time.Now()}
	receiver.mu.Unlock()

	return actor, nil
}

func (receiver *ActivityPub) fetch(actorId string) (result ap.Actor, err error) {
	req, err := http.NewRequest(http.MethodGet, actorId, nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/activity+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("actor fetch %s returned HTTP %d", actorId, resp.StatusCode)
		return
	}

	err = json.Unmarshal(body, &result)
	return
}
