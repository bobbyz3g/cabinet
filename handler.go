package cabinet

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var ErrReceiveTimeout = errors.New("wait receive timeout")

// r is a variable of type `*rand.Rand` that is used for
// generating random integers.
var r = rand.New(rand.NewSource(time.Now().UnixMicro()))

// GenerateCode generates and returns a random 4-digit code.
// The code is generated using a random number generator initialized with
// the current Unix time in microseconds.
func GenerateCode() string {
	return strconv.Itoa(r.Intn(9000) + 1000)
}

type PushHandler struct {
	Sessions *Sessions

	err  error
	code Code
	done chan struct{}
}

func (p *PushHandler) Prepare(ctx echo.Context) {
	fh, err := ctx.FormFile("file")
	if err != nil {
		p.err = err
		return
	}

	p.done = make(chan struct{})
	p.code = Code(ctx.FormValue("code"))

	t := &Translator{Name: fh.Filename, Done: p.done}
	t.Reader, p.err = fh.Open()
	p.Sessions.Push(p.code, t)
}

func (p *PushHandler) Flush() error {
	if p.err != nil {
		return p.err
	}

	select {
	case <-p.done:
		return nil
	case <-time.After(time.Minute):
		// remove translator from session, avoid memory leak
		p.Sessions.Pop(p.code)
		return ErrReceiveTimeout
	}
}
