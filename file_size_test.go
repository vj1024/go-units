package units

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestNewFileSize(t *testing.T) {
	testCase := []struct {
		expected FileSize
		strFmt   string
	}{
		{0, ""},
		{0, "0B"},
		{1, "1B"},
		{1024, "1KB"},
		{2 * 1024 * 1024, "2MB"},
		{100 * GB, "100GB"},
		{1024 * GB, "1TB"},
		{1026 * GB, "1026GB"},
		{123 * PB, "123PB"},

		{-1 * GB, "-1GB"},
		{-10 * MB, "-10MB"},

		{123 * B, "123B"},
		{-123 * B, "-123B"},
	}

	for _, v := range testCase {
		got, err := NewFileSize(v.strFmt)
		t.Logf("parse from string:%s, got_int:%d, got_str:%s", v.strFmt, int(got), got)
		assert.Nil(t, err)
		assert.Equal(t, v.expected, got)
	}
}
func TestFileSize_MarshalYAML(t *testing.T) {
	testCase := []struct {
		fileSize FileSize
		expected string
	}{
		{0, "0B"},
		{1, "1B"},
		{1024, "1KB"},
		{2 * 1024 * 1024, "2MB"},
		{100 * GB, "100GB"},
		{1024 * GB, "1TB"},
		{1026 * GB, "1026GB"},
		{123 * PB, "123PB"},

		{-1 * GB, "-1GB"},
		{-10 * MB, "-10MB"},

		{123 * B, "123B"},
		{-123 * B, "-123B"},
	}

	for _, v := range testCase {
		got, err := v.fileSize.MarshalYAML()
		//t.Logf("%+v: %s", v.fileSize, got)
		assert.Nil(t, err)
		assert.Equal(t, v.expected, got)

		m := map[string]interface{}{
			"file_size": v.fileSize,
		}
		bs, _ := yaml.Marshal(m)
		t.Logf("yaml: %s", bs)
	}
}

func TestFileSize_UnmarshalYAML(t *testing.T) {
	testCase := []struct {
		expected FileSize
		src      string
	}{
		{0, "0B"},
		{1, "1B"},
		{1024, "1KB"},
		{2 * 1024 * 1024, "2MB"},
		{100 * GB, "100GB"},
		{1024 * GB, "1TB"},
		{1026 * GB, "1026GB"},
		{123 * PB, "123PB"},

		{-1 * GB, "-1GB"},
		{-10 * MB, "-10MB"},

		{123 * B, "123"},
		{-123 * B, "-123"},
	}

	type Tmp struct {
		FileSize FileSize `yaml:"file_size"`
	}

	for _, v := range testCase {
		var tmp Tmp
		err := yaml.Unmarshal([]byte("file_size: "+v.src), &tmp)
		assert.Nil(t, err)
		assert.Equal(t, v.expected, tmp.FileSize)

		t.Logf("unmarshaled: %+v", tmp)
	}
}

func TestFileSize_MarshalJSON(t *testing.T) {
	v := map[string]interface{}{
		"size": 101 * MB,
	}
	bs, err := json.Marshal(&v)
	if err != nil {
		t.Fatal(err)
	}
	expect := `{"size":"101MB"}`
	if string(bs) != expect {
		t.Fatalf("expect:%s, got:%s", expect, string(bs))
	}
	t.Logf("marshal json: %s", bs)
}

func TestFileSize_UnmarshalJSON(t *testing.T) {
	s := `{"size":"101MB"}`
	v := map[string]FileSize{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Fatal(err)
	}

	if v["size"] != 101*MB {
		t.Fatalf("expect:%v, got:%v", 101*MB, v["size"])
	}
	t.Logf("unmarshal json: %+v", v)
}
