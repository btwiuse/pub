package handler

import (
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/btwiuse/better"
	"github.com/webteleport/utils"
)

// passing symlink also returns true
func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// File exists but other error occurred
	return false
}

func isPort(s string) bool {
	match, _ := regexp.MatchString(`^:\d{1,5}$`, s)
	return match
}

func isHostPort(s string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9.-]+:\d{1,5}$`, s)
	return match
}

func isValidURL(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	return err == nil && u.Host != ""
}

func serveFile(s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s)
	})
}

func ResourceEmoji(s string) string {
	switch {
	case isFile(s):
		return "📄"
	case pathExists(s):
		return "📁"
	case isPort(s):
		return "🔌"
	case isHostPort(s):
		return "💻"
	case isValidURL(s):
		return "🌐"
	default:
		return "🔀"
	}
}

func ResourceHandler(s string) http.Handler {
	switch {
	case isPort(s), isHostPort(s), isValidURL(s):
		return utils.ReverseProxy(s)
	case isFile(s):
		return serveFile(s)
	case pathExists(s):
		return better.FileServer(http.Dir(s))
	default:
		return serveLazyFS(s)
	}
}

func serveLazyFS(s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case isFile(s):
			http.ServeFile(w, r, s)
		case pathExists(s):
			better.FileServer(http.Dir(s)).ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}
