/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package openvpn

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestParam_Factory(t *testing.T) {
	option := OptionParam("very-value", "1234")
	assert.NotNil(t, option)
}

func TestParam_GetName(t *testing.T) {
	option := OptionParam("very-value", "1234")
	assert.Equal(t, "very-value", option.getName())
}

func TestParam_ToCli(t *testing.T) {
	option := OptionParam("very-value", "1234")

	optionValue, err := option.toCli()
	assert.NoError(t, err)
	assert.Equal(t, "--very-value 1234", optionValue)
}

func TestParam_ToFile(t *testing.T) {
	option := OptionParam("very-value", "1234")

	optionValue, err := option.toFile()
	assert.NoError(t, err)
	assert.Equal(t, "very-value 1234", optionValue)
}

func TestSpacesDontSplitIntoTwoParams(t *testing.T) {
	optionWithSpaces := OptionParam("name", "value with spaces")
	args, err := optionWithSpaces.toCli()
	assert.NoError(t, err)

	cmd := exec.Command("go", "run", filepath.Join("testdata", "args.go"), args)
	buff := bytes.Buffer{}
	cmd.Stdout = &buff
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	assert.NoError(t, err)

	lineReader := bufio.NewReader(&buff)
	argValue, err := lineReader.ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, "--name value with spaces\n", argValue)

	outputEnd, err := lineReader.ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, "END\n", outputEnd)

}
