// Copyright (c) 2021 Asppj  <asppj@foxmail.com>.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package goconfig

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFile_FileNotExist_Default(t *testing.T) {
	require.NoError(t, parseFile(&setup{
		configFilePath: "/doesntexist.conf",
	}, *NewSampleTagOption()))
}

func TestParseFile_FileNotExist_Custom(t *testing.T) {
	require.Error(t, parseFile(&setup{
		configFilePath:   "/doesntexist.conf",
		customConfigFile: true,
	}, *NewSampleTagOption()))
}

func TestParseFile_InvalidJSON(t *testing.T) {
	file, err := ioutil.TempFile("", "goconfig")
	require.NoError(t, err)

	_, err = file.WriteString(`{
		"key": "value",
	}`)
	require.NoError(t, err)

	require.Error(t, parseFile(&setup{
		configFilePath: file.Name(),
		conf: &Conf{
			FileDecoder: DecoderJSON,
		},
	}, *NewSampleTagOption()))
}

func TestParseFile_InvalidYAML(t *testing.T) {
	file, err := ioutil.TempFile("", "goconfig")
	require.NoError(t, err)

	_, err = file.WriteString("test: \"value\n")
	require.NoError(t, err)

	require.Error(t, parseFile(&setup{
		configFilePath: file.Name(),
		conf: &Conf{
			FileDecoder: DecoderYAML,
		},
	}, *NewSampleTagOption()))
}

func TestParseFile_InvalidTOML(t *testing.T) {
	file, err := ioutil.TempFile("", "goconfig")
	require.NoError(t, err)

	_, err = file.WriteString("test = value\n")
	require.NoError(t, err)

	require.Error(t, parseFile(&setup{
		configFilePath: file.Name(),
		conf: &Conf{
			FileDecoder: DecoderTOML,
		},
	}, *NewSampleTagOption()))
}

func TestParseFile_InvalidAny(t *testing.T) {
	file, err := ioutil.TempFile("", "goconfig")
	require.NoError(t, err)

	_, err = file.WriteString("&$_@")
	require.NoError(t, err)

	require.Error(t, parseFile(&setup{
		configFilePath: file.Name(),
		conf: &Conf{
			FileDecoder: DecoderTryAll,
		},
	}, *NewSampleTagOption()))
}

func TestParseFile_MultiDecoder(t *testing.T) {
	file, err := ioutil.TempFile("", "goconfig")
	require.NoError(t, err)

	_, err = file.WriteString("test = \"value\"\n")
	require.NoError(t, err)

	require.NoError(t, parseFile(&setup{
		configFilePath: file.Name(),
		conf: &Conf{
			FileDecoder: NewMultiFileDecoder([]FileDecoderFn{
				DecoderJSON,
				DecoderYAML,
				DecoderTOML,
			}),
		},
	}, *NewSampleTagOption()))
}
