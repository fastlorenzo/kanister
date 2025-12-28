// Copyright 2025 The Kanister Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sftp

import (
	"testing"

	"github.com/kanisterio/safecli/command"
	"github.com/kanisterio/safecli/test"
	"gopkg.in/check.v1"

	"github.com/kanisterio/kanister/pkg/kopia/cli"
)

func TestSFTPOptions(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&test.ArgumentSuite{Cmd: "cmd", Arguments: []test.ArgumentTest{
	{
		Name:        "optHost with hostname should return option",
		Argument:    optHost("sftp.example.com"),
		ExpectedCLI: []string{"cmd", "--host=sftp.example.com"},
	},
	{
		Name:        "optHost with empty hostname should return error",
		Argument:    optHost(""),
		ExpectedErr: cli.ErrInvalidHostname,
	},
	{
		Name:        "optPort with port should return option",
		Argument:    optPort("2222"),
		ExpectedCLI: []string{"cmd", "--port=2222"},
	},
	{
		Name:        "optPort with empty port should use default",
		Argument:    optPort(""),
		ExpectedCLI: []string{"cmd", "--port=22"},
	},
	{
		Name:        "optPath with path should return option",
		Argument:    optPath("/backup"),
		ExpectedCLI: []string{"cmd", "--path=/backup"},
	},
	{
		Name:        "optPath with empty path should return error",
		Argument:    optPath(""),
		ExpectedErr: cli.ErrInvalidRepoPath,
	},
	{
		Name:        "optUsername with username should return option",
		Argument:    optUsername("backup-user"),
		ExpectedCLI: []string{"cmd", "--username=backup-user"},
	},
	{
		Name:        "optUsername with empty username should return error",
		Argument:    optUsername(""),
		ExpectedErr: cli.ErrInvalidUsername,
	},
	{
		Name:        "optKeyFile with keyfile should return option",
		Argument:    optKeyFile("/keys/id_rsa"),
		ExpectedCLI: []string{"cmd", "--keyfile=/keys/id_rsa"},
	},
	{
		Name:        "optKeyFile with empty keyfile should return noop",
		Argument:    command.NewArguments(optKeyFile("")),
		ExpectedCLI: []string{"cmd"},
	},
	{
		Name:        "optKnownHosts with known_hosts should return option",
		Argument:    optKnownHosts("/keys/known_hosts"),
		ExpectedCLI: []string{"cmd", "--known-hosts=/keys/known_hosts"},
	},
	{
		Name:        "optKnownHosts with empty known_hosts should return noop",
		Argument:    command.NewArguments(optKnownHosts("")),
		ExpectedCLI: []string{"cmd"},
	},
}})
