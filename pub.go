package pub

import (
	_ "embed"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/btwiuse/pub/handler"
	"github.com/btwiuse/rng"
	"github.com/webteleport/utils"
	"github.com/webteleport/wtf"
)

type Rule struct {
	Resource string
	Pattern  string
	Prefix   string
}

func InferPrefix(pattern string) string {
	parts := strings.Split(pattern, " ")
	pattern = parts[len(parts)-1]
	switch {
	case strings.HasSuffix(pattern, "/"):
		return strings.TrimSuffix(pattern, "/")
	default:
		return ""
	}
}

func SplitPathPrefix(pp string) (pattern, pfx string) {
	pattern = pp
	if strings.Contains(pp, "#") {
		parts := strings.SplitN(pp, "#", 2)
		pattern = parts[0]
		pfx = parts[1]
	} else if pfx == "" {
		pfx = InferPrefix(pattern)
	}
	return
}

func NewRule(res, pattern_with_prefix string) Rule {
	pattern, pfx := SplitPathPrefix(pattern_with_prefix)
	return Rule{res, pattern, pfx}
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
		pattern := rule.Pattern
		pfx := rule.Prefix
		emoji := handler.ResourceEmoji(res)
		info := fmt.Sprintf("%s %s 🌐 %s", emoji, res, pattern)
		if pfx != "" {
			info = fmt.Sprintf("%s (stripping prefix: %s)", info, pfx)
		}
		slog.Info(info)
		handlr := http.StripPrefix(pfx, handler.ResourceHandler(res))
		mux.Handle(rule.Pattern, handlr)
	}
}

func Handler(rules Rules) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /debug/vars", expvar.Handler().ServeHTTP)
	ApplyRules(mux, rules)
	return mux
}

//go:embed README.md
var Usage string

func RelayAddr() string {
	if relay := os.Getenv("RELAY"); relay != "" {
		return relay
	}

	name := rng.NewDockerSepDigits("-", 4)
	return fmt.Sprintf("https://pub.webtransport.fun/%s", name)
}

func Run(args []string) error {
	rules := Parse(args)

	handler := Handler(rules)
	handler = utils.GzipMiddleware(handler)
	handler = utils.GinLoggerMiddleware(handler)
	handler = utils.AllowAllCorsMiddleware(handler)

	return wtf.Serve(RelayAddr(), handler)
}
