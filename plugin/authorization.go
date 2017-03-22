package plugin

import (
	"errors"
	"net/rpc"
	"net/url"
	"strings"

	"github.com/henrylee2cn/rpc2"
)

type (
	// AuthorizationPlugin authorization plugin
	AuthorizationPlugin struct {
		token             string // Authorization token
		tag               string // extra tag for Authorization
		authorizationFunc AuthorizationFunc
	}

	// AuthorizationFunc defines a method type which handles Authorization info
	AuthorizationFunc func(token string, tag string, serviceMethod string) error
)

// NewServerAuthorizationPlugin is by name
func NewServerAuthorizationPlugin(fn AuthorizationFunc) *AuthorizationPlugin {
	return &AuthorizationPlugin{
		authorizationFunc: fn,
	}
}

// NewClientAuthorizationPlugin is by name
func NewClientAuthorizationPlugin(token string, tag string) *AuthorizationPlugin {
	return &AuthorizationPlugin{
		token: token,
		tag:   tag,
	}
}

// Name returns plugin name.
func (auth *AuthorizationPlugin) Name() string {
	return "AuthorizationPlugin"
}

var _ rpc2.IPreWriteRequestPlugin = new(AuthorizationPlugin)

func (auth *AuthorizationPlugin) PreWriteRequest(r *rpc.Request, body interface{}) error {
	s := url.Values{"auth": []string{auth.token + "\x1f" + auth.tag}}.Encode()
	idx := strings.Index(r.ServiceMethod, "?")
	if idx < 0 {
		r.ServiceMethod += "?" + s
	} else {
		r.ServiceMethod += "&" + s
	}
	return nil
}

var _ rpc2.IPostReadRequestHeaderPlugin = new(AuthorizationPlugin)

func (auth *AuthorizationPlugin) PostReadRequestHeader(ctx *rpc2.Context) error {
	if auth.authorizationFunc == nil {
		return nil
	}
	s := ctx.Query().Get("auth")
	a := strings.Split(s, "\x1f")
	if len(a) != 2 {
		return errors.New("The authorization is not formatted correctly: " + s)
	}
	return auth.authorizationFunc(a[0], a[1], ctx.Path())
}
