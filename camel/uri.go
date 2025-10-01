package camel

import (
	"net/url"
	"strconv"
	"strings"
)

type URI struct {
	rawUri    string
	component string
	path      string
	host      string
	port      string
	fragment  string
	username  string
	password  string
	params    map[string]string
}

func (u *URI) Raw() string {
	return u.rawUri
}

func (u *URI) Component() string {
	return u.component
}

func (u *URI) Host() string {
	return u.host
}

func (u *URI) Port() string {
	return u.port
}

func (u *URI) Fragment() string {
	return u.fragment
}

func (u *URI) Path() string {
	return u.path
}

func (u *URI) Username() string {
	return u.username
}

func (u *URI) Password() string {
	return u.password
}

func (u *URI) Params() map[string]string {
	return u.params
}

func (u *URI) HasParam(name string) bool {
	_, exists := u.params[name]
	return exists
}

func (u *URI) HasParams(names ...string) bool {
	for _, name := range names {
		if !u.HasParam(name) {
			return false
		}
	}
	return true
}

func (u *URI) ParamOrDef(name, def string) string {
	if v, exists := u.params[name]; exists {
		return v
	}
	return def
}

func (u *URI) Param(name string) (string, bool) {
	if v, exists := u.params[name]; exists {
		return v, true
	}
	return "", false
}

func (u *URI) MustParam(name string) string {
	if v, exists := u.Param(name); exists {
		return v
	}
	panic("uri param not found: " + name)
}

func (u *URI) ParamInt(name string) (int, error) {
	s := u.ParamOrDef(name, "0")
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (u *URI) MustParamInt(name string) int {
	v, err := u.ParamInt(name)
	if err != nil {
		panic(err)
	}
	return v
}

func (u *URI) ParamBool(name string) (bool, error) {
	s := u.ParamOrDef(name, "false")
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return b, nil
}

func (u *URI) MustParamBool(name string) bool {
	v, err := u.ParamBool(name)
	if err != nil {
		panic(err)
	}
	return v
}

func (u *URI) String() string {
	return u.rawUri
}

type ParseOptions struct {
	// If TRUE - takes the last key.
	// Default: true
	LastWins bool
	// Prefix for system properties (scheme/host/port/path/fragment/username/password).
	// Default: "_"
	MetaPrefix string
}

func Parse(uri string, opts *ParseOptions) (*URI, error) {
	if opts == nil {
		opts = &ParseOptions{LastWins: true, MetaPrefix: "_"}
	}

	parsed, err := parse(uri, opts)
	if err != nil {
		return nil, err
	}

	u := &URI{
		rawUri:    uri,
		component: parsed["component"],
		params:    map[string]string{},
	}

	delete(parsed, "component")

	for k, v := range parsed {
		if strings.HasPrefix(k, opts.MetaPrefix) {
			ck := strings.TrimPrefix(k, opts.MetaPrefix)
			switch ck {
			case "path":
				u.path = v
			case "host":
				u.host = v
			case "port":
				u.port = v
			case "fragment":
				u.fragment = v
			case "username":
				u.username = v
			case "password":
				u.password = v
			}
		} else {
			u.params[k] = v
		}
	}

	return u, nil
}

// parse decodes Camel-like URI and returns map[string]string.
// Input examples:
//   - "timer:foo?period=1000"
//   - "kafka:topic?brokers=localhost:9092&acks=all"
//   - "file:/var/log?recursive=true"
//   - "http://user:pass@host:8080/a/b?x=1#frag"
func parse(uri string, opts *ParseOptions) (map[string]string, error) {
	if opts.MetaPrefix == "" {
		opts.MetaPrefix = "_"
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string)
	meta := func(k, v string) {
		if v != "" {
			out[opts.MetaPrefix+k] = v
		}
	}

	if u.Scheme != "" {
		out["component"] = u.Scheme
	}

	// Camel-style: scheme:opaque
	if u.Opaque != "" && u.Host == "" && u.Path == "" {
		opaque := strings.TrimSpace(u.Opaque)

		ssp, query, _ := strings.Cut(opaque, "?")
		if ssp != "" {
			meta("path", ssp)
		}
		if query == "" {
			parseQueryInto(out, u.RawQuery, opts.LastWins)
		} else {
			parseQueryInto(out, query, opts.LastWins)
		}
		return out, nil
	}

	// Regular URL
	if u.User != nil {
		meta("username", u.User.Username())
		if pw, ok := u.User.Password(); ok {
			meta("password", pw)
		}
	}
	if h := u.Hostname(); h != "" {
		meta("host", h)
	}
	if p := u.Port(); p != "" {
		meta("port", p)
	}
	if p := u.Path; p != "" {
		meta("path", p)
	}
	if frag := u.Fragment; frag != "" {
		meta("fragment", frag)
	}
	if u.RawQuery != "" {
		parseQueryInto(out, u.RawQuery, opts.LastWins)
	}

	return out, nil
}

func parseQueryInto(out map[string]string, rawQuery string, lastWins bool) {
	values, err := url.ParseQuery(rawQuery)

	if err != nil {
		for _, pair := range strings.Split(rawQuery, "&") {
			if pair == "" {
				continue
			}
			k, v, _ := strings.Cut(pair, "=")
			kk, _ := url.QueryUnescape(k)
			vv, _ := url.QueryUnescape(v)
			if lastWins || out[kk] == "" {
				out[kk] = vv
			}
		}
		return
	}

	for k, vals := range values {
		if len(vals) == 0 {
			continue
		}
		if lastWins {
			out[k] = vals[len(vals)-1]
		} else {
			out[k] = vals[0]
		}
	}
}
