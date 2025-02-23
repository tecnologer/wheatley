package commands

import (
	"fmt"
	"strings"
)

const (
	ErrEmoji         = "❌"
	MissingArgsEmoji = "❓"
)

type ResponseOption func(*Response)

type Response struct {
	msg           string
	isMissingArgs bool
	isError       bool
	cmd           *Command
}

func (r *Response) Message() string {
	var msg strings.Builder

	if r.isError {
		msg.WriteString(ErrEmoji)
	}

	if r.isMissingArgs {
		msg.WriteString(MissingArgsEmoji)
	}

	if msg.Len() > 0 {
		msg.WriteString(" ")
	}

	if r.msg != "" {
		msg.WriteString(r.msg)
	}

	if r.HasError() && r.cmd.HasHelp() {
		msg.WriteString("\n\n")
		msg.WriteString(r.cmd.Help())
	}

	return msg.String()
}

func (r *Response) HasError() bool {
	return r.isError || r.isMissingArgs
}

func NewResponse(opts ...ResponseOption) *Response {
	response := &Response{}

	for _, opt := range opts {
		opt(response)
	}

	return response
}

func (r *Response) SetError(format string, args ...any) *Response {
	r.isMissingArgs = false
	r.isError = true
	r.msg = fmt.Sprintf(format, args...)

	return r
}

func (r *Response) SetMissingArgs(format string, args ...any) *Response {
	r.isMissingArgs = true
	r.isError = true
	r.msg = fmt.Sprintf(format, args...)

	return r
}

func (r *Response) SetMessage(format string, args ...any) *Response {
	r.isMissingArgs = false
	r.isError = false
	r.msg = fmt.Sprintf(format, args...)

	return r
}

func WithMessage(msg string) ResponseOption {
	return func(r *Response) {
		r.msg = msg
	}
}

func WithMessagef(format string, args ...any) ResponseOption {
	return func(r *Response) {
		r.msg = fmt.Sprintf(format, args...)
	}
}

func WithMissingArgs() ResponseOption {
	return func(r *Response) {
		r.isMissingArgs = true
	}
}

func WithIsError() ResponseOption {
	return func(r *Response) {
		r.isError = true
	}
}

func WithCommand(cmd *Command) ResponseOption {
	return func(r *Response) {
		r.cmd = cmd
	}
}
