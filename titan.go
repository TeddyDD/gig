package gig

import (
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
)

const (
	titanScheme = "titan"
)

type titanParams struct {
	token string
	mime  string
	size  int
}

func newTitanParams(url *url.URL) (p titanParams) {
	fragments := strings.Split(url.Path, ";")
	for i := range fragments {
		kv := strings.SplitN(fragments[i], "=", 2)
		if len(kv) != 2 {
			continue
		}

		switch {
		case kv[0] == "token":
			p.token = kv[1]
		case kv[0] == "mime":
			p.mime = kv[1]
		case kv[0] == "size":
			if v, err := strconv.Atoi(kv[1]); err == nil {
				p.size = v
			}
		}
	}

	return
}

// Titan returns a middleware that implements Titan protocol request parsing
// and validation. To limit size of uploaded files set sizeLimit to value
// greater than 0.
func Titan(sizeLimit int) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			switch c.URL().Scheme {
			case titanScheme:
				c.Set("titan", true)

				// Parameters
				params := newTitanParams(c.URL())

				if params.size <= 0 {
					return c.NoContent(StatusBadRequest, "Size parameter is incorrect or not provided")
				}

				if sizeLimit > 0 && sizeLimit < params.size {
					return c.NoContent(StatusBadRequest, "Request is bigger than allowed %d bytes", sizeLimit)
				}

				c.Set("size", params.size)
				c.Set("token", params.token)
				c.Set("mime", params.mime)
			default:
				c.Set("titan", false)
			}

			return next(c)
		}
	}
}

// TitanReadFull is utility wrapper that allocates new buffer and reads
// Titan request body into it.
//
// To store file on disk directly io.CopyN is preferable.
func TitanReadFull(c Context) ([]byte, error) {
	size := c.Get("size").(int)
	buffer := make([]byte, size)

	var err error
	if r := c.Reader(); r != nil {
		_, err = io.ReadFull(c.Reader(), buffer)
	} else {
		err = errors.New("context reader is nil")
	}

	return buffer, err
}
