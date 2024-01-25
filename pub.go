package pub

import (
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/btwiuse/pub/handler"
	"github.com/webteleport/utils"
	"github.com/webteleport/wtf"
)

type Rule struct {
	Resource string
	Path     string
	Prefix   string
}

func NewRule(res, path_with_prefix string) Rule {
	path := path_with_prefix
	pfx := ""
	if strings.Contains(path_with_prefix, "#") {
		parts := strings.SplitN(path_with_prefix, "#", 2)
		path = parts[0]
		pfx = parts[1]
	}
	/*
		if pfx == "" {
			pfx = handler.InferPrefix(res)
			slog.Info("infer: "+pfx)
		}
	*/
	return Rule{res, path, pfx}
}

type Rules []Rule

func (s *Rules) Push(r Rule) {
	*s = append(*s, r)
}

func Parse(s []string) (rules Rules) {
	for i := 0; i < len(s); i += 2 {
		res := s[i]
		pwf := "/"
		if i+1 < len(s) {
			pwf = s[i+1]
		}
		rules.Push(NewRule(res, pwf))
	}
	return
}

func ApplyRules(mux *http.ServeMux, rules Rules) {
	for _, rule := range rules {
		if rule.Prefix == "" {
			slog.Info(fmt.Sprintf("✅ %s => %s", rule.Path, rule.Resource))
		} else {
			slog.Info(fmt.Sprintf("✅ %s => %s (stripping prefix: %s)", rule.Path, rule.Resource, rule.Prefix))
		}
		mux.Handle(rule.Path, http.StripPrefix(rule.Prefix, handler.ResourceHandler(rule.Resource)))
	}
	mux.HandleFunc("/debug/vars", expvar.Handler().ServeHTTP)
}

func Handler(rules Rules) http.Handler {
	mux := http.NewServeMux()
	ApplyRules(mux, rules)
	return mux
}

func Run(args []string) error {
	rules := Parse(args)
	handler := Handler(rules)
	handler = utils.GzipMiddleware(handler)
	handler = utils.GinLoggerMiddleware(handler)
	return wtf.Serve("https://k0s.io", handler)
}
