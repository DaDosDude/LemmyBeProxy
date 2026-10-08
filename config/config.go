package config

import (
	"LemmyBeProxy/router"
	"LemmyBeProxy/service"
	"LemmyBeProxy/service/backend"
	"LemmyBeProxy/service/frontend"
	lemmyService "LemmyBeProxy/service/lemmy"
	piefedService "LemmyBeProxy/service/piefed"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var AppHttpPort = 8080
var AppRouter *router.Router
var CertificateFile string
var CertificateKey string
var CorsRegex *regexp.Regexp

var simulateLemmy bool
var activityPub *service.ActivityPub

// activeBackend is what every controller talks to — see
// service/backend/Backend.go. Every controller this proxy implements
// (Post, Comment, Community, User, Search, Site, Upload) is migrated
// onto it; there's no controller left still hardcoded to a Piefed-shaped
// client. Setting BACKEND_TYPE=lemmy genuinely changes every endpoint's
// behavior now, not just some of them.
var activeBackend backend.Backend

// activeFrontend is the wire-format counterpart to activeBackend — see
// service/frontend/Frontend.go. Every endpoint this proxy implements
// works on both wire formats (0.19.x and 0.17.x) except Upload, which
// doesn't need Frontend-axis work at all since pict-rs's own protocol
// is version-agnostic by nature.
var activeFrontend frontend.Frontend

// FrontendVersion is the raw FRONTEND_VERSION value ("0.19" or "0.17"),
// exported so main.go can apply 0.17-specific post-processing (stripping
// timezone suffixes from timestamps — real Lemmy 0.17.x uses
// chrono::NaiveDateTime, which has no timezone concept and rejects one
// if present) without main needing its own copy of this logic or a
// circular import back into this package from a lower-level one.
var FrontendVersion string

// imageHosts is the allowlist UploadController.ServeImage checks decoded
// image URLs against: the backend's own host, plus any EXTRA_IMAGE_HOSTS
// (comma-separated) for a backend that serves media from a separate
// CDN/S3 host.
var imageHosts []string

// backendHost extracts the bare host[:port] from BACKEND_INSTANCE, which
// may be a plain hostname ("retrofed.com") or carry an explicit scheme
// ("http://lemmy-easy-deploy-lemmy-1:8536").
func backendHost(instance string) string {
	if strings.Contains(instance, "://") {
		if parsed, err := url.Parse(instance); err == nil && parsed.Host != "" {
			return parsed.Host
		}
	}
	return strings.TrimSuffix(instance, "/")
}

func init() {
	// Every outbound call (backend API, image upload, ActivityPub actor
	// fetch) goes through http.DefaultClient, which has no timeout at
	// all: one hung backend request held its goroutine and socket open
	// forever. 60s comfortably covers a slow upload or resolve_object.
	http.DefaultClient.Timeout = 60 * time.Second

	if port, exists := os.LookupEnv("PORT"); exists {
		parsed, err := strconv.Atoi(port)
		if err != nil {
			panic(err)
		}

		AppHttpPort = parsed
	}
	if simulate, exists := os.LookupEnv("SIMULATE"); exists {
		var err error
		simulateLemmy, err = strconv.ParseBool(simulate)
		if err != nil {
			panic(err)
		}
	}

	// BACKEND_TYPE and BACKEND_INSTANCE are both required, with no
	// default — this proxy takes no position on which backend you're
	// running, so every deployment must say explicitly. "piefed" or
	// "lemmy" for BACKEND_TYPE; BACKEND_INSTANCE is that backend's
	// hostname (e.g. retrofed.com).
	backendType, hasBackendType := os.LookupEnv("BACKEND_TYPE")
	if !hasBackendType {
		panic("BACKEND_TYPE environment variable not set — must be \"piefed\" or \"lemmy\"")
	}
	backendInstance, hasBackendInstance := os.LookupEnv("BACKEND_INSTANCE")
	if !hasBackendInstance {
		panic("BACKEND_INSTANCE environment variable not set")
	}

	// FRONTEND_VERSION is required too, same "no preference" philosophy.
	// "0.19" for the current Lemmy wire format (what every deployment of
	// this proxy has used until now), "0.17" for the older format
	// lemmyBB and similar older clients actually speak.
	frontendVersion, hasFrontendVersion := os.LookupEnv("FRONTEND_VERSION")
	if !hasFrontendVersion {
		panic("FRONTEND_VERSION environment variable not set — must be \"0.19\" or \"0.17\"")
	}

	activityPub = service.NewActivityPub()

	switch backendType {
	case "piefed":
		activeBackend = piefedService.NewPiefedBackend(piefedService.NewPiefed(backendInstance), activityPub, simulateLemmy)
	case "lemmy":
		activeBackend = lemmyService.NewLemmyBackend(lemmyService.NewLemmy(backendInstance))
	default:
		panic(fmt.Sprintf("unknown BACKEND_TYPE %q — expected \"piefed\" or \"lemmy\"", backendType))
	}

	switch frontendVersion {
	case "0.19":
		activeFrontend = frontend.NewFrontend019()
	case "0.17":
		activeFrontend = frontend.NewFrontend017()
	default:
		panic(fmt.Sprintf("unknown FRONTEND_VERSION %q — expected \"0.19\" or \"0.17\"", frontendVersion))
	}
	FrontendVersion = frontendVersion

	imageHosts = []string{backendHost(backendInstance)}
	if extra, exists := os.LookupEnv("EXTRA_IMAGE_HOSTS"); exists {
		imageHosts = append(imageHosts, strings.Split(extra, ",")...)
	}

	AppRouter = router.NewRouter()

	if certFile, exists := os.LookupEnv("CERTIFICATE_FILE"); exists {
		CertificateFile = certFile
	}
	if certKey, exists := os.LookupEnv("CERTIFICATE_KEY"); exists {
		CertificateKey = certKey
	}
	if regexStr, exists := os.LookupEnv("CORS_REGEX"); exists {
		CorsRegex = regexp.MustCompile(regexStr)
	}
}
