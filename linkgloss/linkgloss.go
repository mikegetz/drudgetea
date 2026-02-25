package linkgloss

import "github.com/charmbracelet/lipgloss"

const (
	hrefStart   = "\033]8;;"
	hrefEnd     = "\033\\"
	hrefTextEnd = "\033]8;;\033\\"
)

func hyperlink(url, text string) string {
	return hrefStart + url + hrefEnd + text + hrefTextEnd
}

type Style struct {
	lipgloss.Style
	url      string
	disabled bool
}

func New(s ...lipgloss.Style) Style {
	if len(s) > 0 {
		return Style{Style: s[0]}
	}
	return Style{Style: lipgloss.NewStyle()}
}

// mimic lipgloss.Style methods
func (s Style) Foreground(c lipgloss.TerminalColor) Style {
	s.Style = s.Style.Foreground(c)
	return s
}

// mimic lipgloss.Style methods
func (s Style) Bold(v bool) Style {
	s.Style = s.Style.Bold(v)
	return s
}

// mimic lipgloss.Style methods
func (s Style) Width(v int) Style {
	s.Style = s.Style.Width(v)
	return s
}

// mimic lipgloss.Style methods
func (s Style) Align(a lipgloss.Position) Style {
	s.Style = s.Style.Align(a)
	return s
}

// mimic lipgloss.Style methods
func (s Style) Underline(v bool) Style {
	s.Style = s.Style.Underline(v)
	return s
}

// the linkgloss
func (s Style) URL(url string) Style {
	s.url = url
	return s
}

func (s Style) UnsetURL() Style {
	s.url = ""
	return s
}

// mimic lipgloss.Style methods with hyperlink wrapper
func (s Style) Render(str string) string {
	out := s.Style.Render(str)
	if s.url == "" {
		return out
	}
	return hyperlink(s.url, out)
}
