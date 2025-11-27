package security

import (
	"net/url"
	"slices"
	"strings"
)

func ValidateRedirectURL(redirectURL string, allowedHosts []string) (string, bool) {
	if redirectURL == "" {
		return "/", true
	}

	parsed, err := url.Parse(redirectURL)
	if err != nil {
		return "/", false
	}

	if parsed.Scheme != "" || parsed.Host != "" {
		if len(allowedHosts) > 0 && !isAllowedHost(parsed.Host, allowedHosts) {
			return "/", false
		}

		if len(allowedHosts) == 0 {
			return "/", false
		}
	}

	path := parsed.Path
	if !strings.HasPrefix(path, "/") {
		return "/", false
	}

	if strings.HasPrefix(path, "//") {
		return "/", false
	}

	if strings.Contains(redirectURL, "@") {
		return "/", false
	}

	if parsed.RawQuery != "" {
		path = path + "?" + parsed.RawQuery
	}

	if parsed.Fragment != "" {
		path = path + "#" + parsed.Fragment
	}

	return path, true
}

func isAllowedHost(host string, allowedHosts []string) bool {
	return slices.Contains(allowedHosts, host)
}
