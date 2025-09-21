package tdms

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type InvalidPathError struct {
	Message string
	Offset  int
}

func NewInvalidPathError(offset int, msg string, args ...any) *InvalidPathError {
	return &InvalidPathError{
		Message: fmt.Sprintf(msg, args...),
		Offset:  offset,
	}
}

func NewInvalidPathErrorExpectingRune(offset int, c rune) *InvalidPathError {
	return &InvalidPathError{
		Message: fmt.Sprintf("invalid path, expecting \"%c\" at offset %d", c, offset),
		Offset:  offset,
	}
}

func (err InvalidPathError) Error() string {
	return err.Message
}

type ObjectPath struct {
	Group   string
	Channel string
}

func ObjectPathFromString(s string) (ObjectPath, error) {
	r := bufio.NewReader(strings.NewReader(s))
	var parts []string
	var part string
	var offset = 0
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return ObjectPath{}, err
			}
		}
		if c != '/' {
			return ObjectPath{}, NewInvalidPathErrorExpectingRune(offset, '/')
		}
		offset++
		part = ""
		c, _, err = r.ReadRune()
		if err != nil {
			if err == io.EOF {
				parts = append(parts, part)
				break
			} else {
				return ObjectPath{}, err
			}
		}
		if c != '\'' {
			return ObjectPath{}, NewInvalidPathErrorExpectingRune(offset, '\'')
		}
		offset++
		for {
			c, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return ObjectPath{}, NewInvalidPathErrorExpectingRune(offset, '\'')
				} else {
					return ObjectPath{}, err
				}
			}
			offset++
			if c == '\'' {
				parts = append(parts, part)
				break
			}
			part += string(c)
		}
	}
	switch len(parts) {
	case 0:
		return ObjectPath{}, NewInvalidPathError(offset, "invalid path, empty")
	case 1:
		return ObjectPath{
			Group: parts[0],
		}, nil
	case 2:
		return ObjectPath{
			Group:   parts[0],
			Channel: parts[1],
		}, nil
	default:
		return ObjectPath{}, NewInvalidPathError(offset, "invalid path, too many parts")
	}
}

func (path ObjectPath) IsRoot() bool {
	return (path.Group == "") && (path.Channel == "")
}

func (path ObjectPath) IsGroup() bool {
	return (path.Group != "") && (path.Channel == "")
}

func (path ObjectPath) IsChannel() bool {
	return (path.Group != "") && (path.Channel != "")
}

func (path ObjectPath) String() string {
	if path.IsRoot() {
		return "/"
	} else if path.IsGroup() {
		return "/'" + path.Group + "'"
	} else if path.IsChannel() {
		return "/'" + path.Group + "'/'" + path.Channel + "'"
	} else {
		return ""
	}
}
