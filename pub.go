package pub

import (
	"expvar"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"

	"github.com/btwiuse/pub/handler"
	"github.com/webteleport/utils"
	"github.com/webteleport/wtf"
	"k0s.io/pkg/rng"
)

type Rule struct {
	Resource string
	Path     string
	Prefix   string
}

func InferPrefix(s string) string {
	switch {
	case strings.HasSuffix(s, "/"):
		return strings.TrimSuffix(s, "/")
	default:
		return ""
	}
}

func SplitPathPrefix(pp string) (path, pfx string) {
	path = pp
	if strings.Contains(pp, "#") {
		parts := strings.SplitN(pp, "#", 2)
		path = parts[0]
		pfx = parts[1]
	} else if pfx == "" {
		pfx = InferPrefix(path)
	}
	return
}

func NewRule(res, path_with_prefix string) Rule {
	path, pfx := SplitPathPrefix(path_with_prefix)
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
		res := rule.Resource
		path := rule.Path
		pfx := rule.Prefix
		emoji := handler.ResourceEmoji(res)
		info := fmt.Sprintf("%s %s 🌐 %s", emoji, res, path)
		if pfx != "" {
			info = fmt.Sprintf("%s (stripping prefix: %s)", info, pfx)
		}
		slog.Info(info)
		handlr := http.StripPrefix(pfx, handler.ResourceHandler(res))
		mux.Handle(rule.Path, handlr)
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
	name := strings.ReplaceAll(rng.New(), "_", "-")
	numb := rand.Intn(9000) + 1000
	addr := fmt.Sprintf("https://pub.webtransport.fun/%s-%s", name, numb)
	return wtf.Serve(addr, handler)
}
