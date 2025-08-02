package utils

import (
	"net"
	"regexp"
	"strings"
)

type HostMatcher struct {
	pattern string
	regex   *regexp.Regexp
}

func NewHostMatcher(pattern string) *HostMatcher {
	return &HostMatcher{
		pattern: pattern,
		regex:   compileHostPattern(pattern),
	}
}

func (hm *HostMatcher) Matches(host string) bool {
	if hm.regex != nil {
		return hm.regex.MatchString(host)
	}
	return false
}

func MatchesHost(host, pattern string) bool {
	matcher := NewHostMatcher(pattern)
	return matcher.Matches(host)
}

func MatchesMask(mask, pattern string) bool {
	if pattern == "*" || pattern == "*!*@*" {
		return true
	}

	if !strings.Contains(pattern, "!") && !strings.Contains(pattern, "@") {
		if idx := strings.LastIndex(mask, "@"); idx != -1 {
			host := mask[idx+1:]
			return MatchesHost(host, pattern)
		}
		return MatchesHost(mask, pattern)
	}

	matcher := NewHostMatcher(pattern)
	return matcher.Matches(mask)
}

func compileHostPattern(pattern string) *regexp.Regexp {
	if pattern == "" {
		return nil
	}

	if pattern == "*" {
		return regexp.MustCompile(".*")
	}

	escaped := regexp.QuoteMeta(pattern)

	escaped = strings.ReplaceAll(escaped, "\\*", ".*")
	escaped = strings.ReplaceAll(escaped, "\\?", ".")

	regexPattern := "^" + escaped + "$"

	regex, err := regexp.Compile("(?i)" + regexPattern)
	if err != nil {
		exactPattern := "^" + regexp.QuoteMeta(pattern) + "$"
		regex, _ = regexp.Compile("(?i)" + exactPattern)
	}

	return regex
}

func ValidateHostPattern(pattern string) bool {
	if pattern == "" {
		return false
	}

	if strings.Contains(pattern, "**") {
		return false
	}

	return compileHostPattern(pattern) != nil
}

func NormalizeHostPattern(pattern string) string {
	if pattern == "" {
		return "*"
	}

	pattern = strings.ToLower(pattern)

	if pattern == "localhost" {
		return "127.0.0.1"
	}

	return pattern
}

func IsIPAddress(host string) bool {
	return net.ParseIP(host) != nil
}

func IsValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}

	validHostname := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	return validHostname.MatchString(hostname)
}

func ExpandHostPattern(pattern string) []string {
	switch pattern {
	case "local", "localhost":
		return []string{"127.0.0.1", "::1", "localhost"}
	case "private":
		return []string{"192.168.*", "10.*", "172.16.*", "172.17.*", "172.18.*", "172.19.*", "172.20.*", "172.21.*", "172.22.*", "172.23.*", "172.24.*", "172.25.*", "172.26.*", "172.27.*", "172.28.*", "172.29.*", "172.30.*", "172.31.*"}
	default:
		return []string{pattern}
	}
}

func MatchesAnyHost(host string, patterns []string) bool {
	for _, pattern := range patterns {
		if MatchesHost(host, pattern) {
			return true
		}
	}
	return false
}

func ParseMask(mask string) (nick, user, host string) {
	nick = "*"
	user = "*"
	host = "*"

	exclamPos := strings.Index(mask, "!")
	atPos := strings.LastIndex(mask, "@")

	if exclamPos == -1 && atPos == -1 {
		host = mask
	} else if exclamPos == -1 && atPos != -1 {
		user = mask[:atPos]
		host = mask[atPos+1:]
	} else if exclamPos != -1 && atPos == -1 {
		nick = mask[:exclamPos]
		user = mask[exclamPos+1:]
	} else if exclamPos < atPos {
		nick = mask[:exclamPos]
		user = mask[exclamPos+1 : atPos]
		host = mask[atPos+1:]
	} else {
		host = mask
	}

	return nick, user, host
}

func FormatMask(nick, user, host string) string {
	if nick == "" {
		nick = "*"
	}
	if user == "" {
		user = "*"
	}
	if host == "" {
		host = "*"
	}

	return nick + "!" + user + "@" + host
}

func MatchesBanMask(clientMask, banMask string) bool {
	clientNick, clientUser, clientHost := ParseMask(clientMask)
	banNick, banUser, banHost := ParseMask(banMask)

	nickMatch := MatchesHost(clientNick, banNick)
	userMatch := MatchesHost(clientUser, banUser)
	hostMatch := MatchesHost(clientHost, banHost)

	return nickMatch && userMatch && hostMatch
}
