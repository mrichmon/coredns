package health

import (
	"github.com/miekg/coredns/middleware"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("health", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	addr, err := healthParse(c)
	if err != nil {
		return middleware.Error("health", err)
	}

	health := &Health{Addr: addr}
	c.OnStartup(health.Startup)
	c.OnShutdown(health.Shutdown)

	// Don't do AddMiddleware, as health is not *really* a middleware just a separate
	// webserver running.

	return nil
}

func healthParse(c *caddy.Controller) (string, error) {
	addr := ""
	for c.Next() {
		args := c.RemainingArgs()

		switch len(args) {
		case 0:
		case 1:
			addr = args[0]
		default:
			return "", c.ArgErr()
		}
	}
	return addr, nil
}
