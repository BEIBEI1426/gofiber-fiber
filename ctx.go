package fiber

import (
	"strings"
	"github.com/gofiber/utils/v2"
)

type Config struct {
	ProxyHeader        string
	EnableIPValidation bool
}

type Ctx struct {
	config  *Config
	handler RequestHandler
}

type RequestHandler interface {
	Get(key string) string
	RemoteIP() IPAddr
}

type IPAddr interface {
	String() string
}

type MockRequestHandler struct {
	headers  map[string]string
	remoteIP string
}

func (m *MockRequestHandler) Get(key string) string {
	if m.headers != nil {
		return m.headers[key]
	}
	return ""
}

func (m *MockRequestHandler) RemoteIP() IPAddr {
	return &mockIP{ip: m.remoteIP}
}

type mockIP struct {
	ip string
}

func (m *mockIP) String() string {
	if m.ip == "" {
		return "127.0.0.1"
	}
	return m.ip
}

func (c *Ctx) IP() string {
	if c.config.EnableIPValidation && len(c.config.ProxyHeader) > 0 {
		return c.extractIPFromHeader(c.config.ProxyHeader)
	}
	return c.handler.RemoteIP().String()
}

func (c *Ctx) Get(key string) string {
	return c.handler.Get(key)
}

func (c *Ctx) extractIPFromHeader(header string) string {
	if c.config.EnableIPValidation {
		headerValue := c.Get(header)

		i := 0
		j := -1

	iploop:
		for {
			var v4, v6 bool

			i, j = j+1, j+2

			if j > len(headerValue) {
				break
			}

			for j < len(headerValue) && headerValue[j] != ',' {
				if headerValue[j] == ':' {
					v6 = true
				} else if headerValue[j] == '.' {
					v4 = true
				}
				j++
			}

			for i < j && headerValue[i] == ' ' {
				i++
			}

			s := strings.TrimRight(headerValue[i:j], " ")

			if s == "" {
				continue iploop
			}

			if (!v6 && !v4) || (v6 && !utils.IsIPv6(s)) || (v4 && !utils.IsIPv4(s)) {
				continue iploop
			}

			return s
		}

		return c.handler.RemoteIP().String()
	}

	return c.Get(header)
}