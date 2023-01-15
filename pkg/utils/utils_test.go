package utils

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestWriteFile(t *testing.T) {
	testContent := func() string {
		s := ""
		for i := 1; i < 10000; i++ {
			s += strconv.Itoa(i) + "\n"
		}
		return s
	}()
	testFilename := "test.m3u8"
	err := WriteFile(testContent, testFilename)
	if err != nil {
		require.Equalf(t, testFilename, err, "写入失败 %v", err)
		return
	}
	require.Equalf(t, testFilename, err, "写入成功 %v", err)

}

func TestFirstLower(t *testing.T) {
	TestCoent := []string{"GENDER", "COUNTRY"}
	for _, tc := range TestCoent {
		Tc := FirstLower(tc)
		require.Equalf(t, Tc, nil, "写入成功 ")
	}

}
func TestFirstUpper(t *testing.T) {
	TestCoent := []string{"gender", "country"}
	for _, tc := range TestCoent {
		Tc := FirstUpper(tc)
		require.Equalf(t, Tc, nil, "写入成功 ")
	}

}
