package cabinet

import (
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// r is a variable of type `*rand.Rand` that is used for
// generating random integers.
var r = rand.New(rand.NewSource(time.Now().UnixMicro()))

// GenerateCode generates and returns a random 4-digit code.
// The code is generated using a random number generator initialized with
// the current Unix time in microseconds.
func GenerateCode() string {
	return strconv.Itoa(r.Intn(9000) + 1000)
}

type SaveHandler struct {
	Ctx echo.Context
	fh  *multipart.FileHeader
	src multipart.File
	dst *os.File
	err error
}

func (s *SaveHandler) FormFile(name string) {
	s.fh, s.err = s.Ctx.FormFile(name)
}

func (s *SaveHandler) OpenFileHeader() {
	if s.err != nil {
		return
	}
	s.src, s.err = s.fh.Open()
}

func (s *SaveHandler) OpenOSFile(name string, flag int, perm os.FileMode) {
	if s.err != nil {
		return
	}
	s.dst, s.err = os.OpenFile(name, flag, perm)
}

func (s *SaveHandler) Save() {
	if s.err != nil {
		return
	}
	_, s.err = io.Copy(s.dst, s.src)
}
func (s *SaveHandler) Err() error {
	return s.err
}
