package units

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// FileSize example: B/KB/MB/GB/TB/PB, case insensitive.
type FileSize int64

// Common FileSizes.
const (
	B  FileSize = 1
	KB          = 1024 * B
	MB          = 1024 * KB
	GB          = 1024 * MB
	TB          = 1024 * GB
	PB          = 1024 * TB
	//EB          = 1024 * PB
	//ZB          = 1024 * EB
	//YB          = 1024 * ZB
)

// mapping must be sorted from big to small.
var mapping = []struct {
	baseSize FileSize
	unit     string
}{
	//{YB, "YB"},
	//{ZB, "ZB"},
	//{EB, "EB"},
	{PB, "PB"},
	{TB, "TB"},
	{GB, "GB"},
	{MB, "MB"},
	{KB, "KB"},
	{B, "B"},
}

// NewFileSize return a FileSize from string.
func NewFileSize(s string) (FileSize, error) {
	var fs FileSize
	err := fs.unmarshal([]byte(s))
	return fs, err
}

// String return the string format of FileSize.
func (u FileSize) String() string {
	return u.marshal()
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (u *FileSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return u.unmarshal([]byte(s))
}

// MarshalYAML implements the yaml.Marshaler interface.
func (u FileSize) MarshalYAML() (interface{}, error) {
	return u.marshal(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (u *FileSize) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("bad quoted string")
	}

	var in []byte
	if len(b) > 2 {
		in = b[1 : len(b)-1]
	}
	return u.unmarshal(in)
}

// MarshalJSON implements the json.Marshaler interface.
func (u FileSize) MarshalJSON() ([]byte, error) {
	s := u.marshal()
	return []byte(`"` + s + `"`), nil
}

func (u *FileSize) unmarshal(in []byte) error {
	i := len(in)
	if i == 0 {
		*u = 0
		return nil
	}

	for i--; i >= 0; i-- {
		if isDigital(in[i]) {
			break
		}
	}

	s := string(in[:i+1])
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid number in FileSize:`%s`", in)
	}

	unit := strings.ToUpper(string(in[i+1:]))
	if unit == "" {
		*u = FileSize(n)
		return nil
	}
	if len(unit) == 1 && unit != "B" {
		unit += "B"
	}

	for _, v := range mapping {
		if unit == v.unit {
			*u = FileSize(n) * v.baseSize
			return nil
		}
	}
	return fmt.Errorf("invalid unit in FileSize:`%s`", in)
}

func (u FileSize) marshal() string {
	if u == 0 {
		return "0B"
	}

	positive := FileSize(1)
	if u < 0 {
		positive = -1
		u *= -1
	}

	for _, v := range mapping {
		if u%v.baseSize == 0 {
			return fmt.Sprintf("%d%s", positive*u/v.baseSize, v.unit)
		}
	}
	return fmt.Sprintf("%dB", positive*u)
}

func isDigital(b byte) bool {
	return b >= '0' && b <= '9'
}
