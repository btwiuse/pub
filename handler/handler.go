package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/rclone/rclone/backend/local"
	"github.com/rclone/rclone/cmd"
	"github.com/rclone/rclone/fs"
	libhttp "github.com/rclone/rclone/lib/http"
	"github.com/rclone/rclone/lib/http/serve"
	"github.com/rclone/rclone/vfs"
	"github.com/rclone/rclone/vfs/vfsflags"
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
	_, err := url.ParseRequestURI(toTest)
	return err == nil
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

func serveDir(s string) http.Handler {
	VFS := vfs.New(cmd.NewFsSrc([]string{s}), &vfsflags.Opt)
	t, _ := libhttp.GetTemplate("")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remote := strings.TrimPrefix(r.URL.Path, "/")
		// List the directory
		node, err := VFS.Stat(r.URL.Path)
		if err == vfs.ENOENT {
			http.Error(w, "Directory not found", http.StatusNotFound)
			return
		} else if err != nil {
			serve.Error(s, w, "Failed to list directory", err)
			return
		}
		if node.IsDir() && !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", 302)
			return
		}
		if node.IsFile() {
			entry := node.DirEntry()
			if entry == nil {
				http.Error(w, "Can't open file being written", http.StatusNotFound)
				return
			}
			file := node.(*vfs.File)

			// Set content length if we know how long the object is
			knownSize := file.Size() >= 0
			if knownSize {
				w.Header().Set("Content-Length", strconv.FormatInt(node.Size(), 10))
			}

			// Set the Last-Modified header to the timestamp
			w.Header().Set("Last-Modified", file.ModTime().UTC().Format(http.TimeFormat))

			// If HEAD no need to read the object since we have set the headers
			if r.Method == "HEAD" {
				return
			}

			// open the object
			in, err := file.Open(os.O_RDONLY)
			if err != nil {
				serve.Error(remote, w, "Failed to open file", err)
				return
			}
			defer func() {
				err := in.Close()
				if err != nil {
					fs.Errorf(remote, "Failed to close file: %v", err)
				}
			}()

			// http.ServeContent can't serve unknown length files
			if !knownSize {
				if rangeRequest := r.Header.Get("Range"); rangeRequest != "" {
					http.Error(w, "Can't use Range: on files of unknown length", http.StatusRequestedRangeNotSatisfiable)
					return
				}
				n, err := io.Copy(w, in)
				if err != nil {
					fmt.Errorf("Didn't finish writing GET request (wrote %d/unknown bytes): %v", n, err)
					return
				}
			}

			// Serve the file
			http.ServeContent(w, r, remote, node.ModTime(), in)

			return
		}

		if !node.IsDir() {
			http.Error(w, "Not a directory", http.StatusNotFound)
			return
		}

		dir := node.(*vfs.Dir)
		dirEntries, err := dir.ReadDirAll()
		if err != nil {
			serve.Error(s, w, "Failed to list directory", err)
			return
		}

		// Make the entries for display
		directory := serve.NewDirectory(remote, t)
		for _, node := range dirEntries {
			if vfsflags.Opt.NoModTime {
				directory.AddHTMLEntry(node.Path(), node.IsDir(), node.Size(), time.Time{})
			} else {
				directory.AddHTMLEntry(node.Path(), node.IsDir(), node.Size(), node.ModTime().UTC())
			}
		}

		sortParm := r.URL.Query().Get("sort")
		orderParm := r.URL.Query().Get("order")
		directory.ProcessQueryParams(sortParm, orderParm)

		// Set the Last-Modified header to the timestamp
		w.Header().Set("Last-Modified", dir.ModTime().UTC().Format(http.TimeFormat))

		directory.Serve(w, r)
	})
}

func ResourceHandler(s string) http.Handler {
	switch {
	case isFile(s):
		return serveFile(s)
	case pathExists(s):
		return serveDir(s)
	case isPort(s):
		return utils.ReverseProxy(s)
	case isHostPort(s):
		return utils.ReverseProxy(s)
	case isValidURL(s):
		return utils.ReverseProxy(s)
	default:
		return utils.ReverseProxy(s)
	}
}
