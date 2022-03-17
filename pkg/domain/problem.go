package domain

// Implements rfc7807 (https://datatracker.ietf.org/doc/html/rfc7807)

type Problem struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func NewProblem(title string) *Problem {
	p := &Problem{
		Type:  "about:blank",
		Title: title,
	}

	return p
}
