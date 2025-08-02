package utils

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	IRCBold      = "\x02"
	IRCItalic    = "\x1D"
	IRCUnderline = "\x1F"
	IRCColor     = "\x03"
	IRCReset     = "\x0F"
)

var colorMap = map[string]string{
	"white":       "00",
	"black":       "01",
	"blue":        "02",
	"green":       "03",
	"red":         "04",
	"brown":       "05",
	"purple":      "06",
	"orange":      "07",
	"yellow":      "08",
	"light green": "09",
	"cyan":        "10",
	"light cyan":  "11",
	"light blue":  "12",
	"pink":        "13",
	"grey":        "14",
	"light grey":  "15",
}

func FormatMOTD(text string) string {
	text = strings.ReplaceAll(text, "$$", "\x00ESCAPED_DOLLAR\x00")

	boldRegex := regexp.MustCompile(`\$b([^$]*)\$r`)
	text = boldRegex.ReplaceAllString(text, IRCBold+"$1"+IRCReset)

	italicRegex := regexp.MustCompile(`\$i([^$]*)\$r`)
	text = italicRegex.ReplaceAllString(text, IRCItalic+"$1"+IRCReset)

	underlineRegex := regexp.MustCompile(`\$u([^$]*)\$r`)
	text = underlineRegex.ReplaceAllString(text, IRCUnderline+"$1"+IRCReset)

	colorNameRegex := regexp.MustCompile(`\$c\[([^]]+)]([^$]*)\$c`)
	text = colorNameRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := colorNameRegex.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		colorSpec := parts[1]
		content := parts[2]

		colors := strings.Split(colorSpec, ",")
		var ircColorCode string

		if len(colors) >= 1 {
			if fg, exists := colorMap[strings.TrimSpace(colors[0])]; exists {
				ircColorCode = IRCColor + fg
			}
		}

		if len(colors) >= 2 {
			if bg, exists := colorMap[strings.TrimSpace(colors[1])]; exists {
				ircColorCode += "," + bg
			}
		}

		return ircColorCode + content + IRCReset
	})

	colorNumRegex := regexp.MustCompile(`\$c(\d+)(?:,(\d+))?`)
	text = colorNumRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := colorNumRegex.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		fg := parts[1]
		ircColorCode := IRCColor + formatColorCode(fg)

		if len(parts) > 2 && parts[2] != "" {
			bg := parts[2]
			ircColorCode += "," + formatColorCode(bg)
		}

		return ircColorCode
	})

	text = strings.ReplaceAll(text, "$r", IRCReset)

	text = strings.ReplaceAll(text, "\x00ESCAPED_DOLLAR\x00", "$")

	return text
}

func formatColorCode(code string) string {
	if num, err := strconv.Atoi(code); err == nil && num < 10 {
		return "0" + code
	}
	return code
}
