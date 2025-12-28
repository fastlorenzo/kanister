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
	"os"
	"testing"

	"github.com/kanisterio/safecli/test"
	"gopkg.in/check.v1"

	"github.com/kanisterio/kanister/pkg/kopia/cli"
	"github.com/kanisterio/kanister/pkg/kopia/cli/internal"
	intlog "github.com/kanisterio/kanister/pkg/kopia/cli/internal/log"
	inttest "github.com/kanisterio/kanister/pkg/kopia/cli/internal/test"
	"github.com/kanisterio/kanister/pkg/log"
)

func TestNewSFTP(t *testing.T) { check.TestingT(t) }

//nolint:unparam
func newLocation(host, port, path, knownHosts string) internal.Location {
	// Set up test environment variables for SFTP credentials
	// These would normally be injected by the repository server controller
	_ = os.Setenv("SFTP_USERNAME", "sftpuser")
	_ = os.Setenv("SFTP_PASSWORD", "sftppass123")

	loc := internal.Location{
		"host": []byte(host),
		"port": []byte(port),
		"path": []byte(path),
	}
	if knownHosts != "" {
		loc["knownHosts"] = []byte(knownHosts)
	}
	return loc
}

// sftptest is a test case for NewSFTP.
type sftptest struct {
	Name        string
	Location    internal.Location
	RepoPath    string
	ExpectedCLI []string
	ExpectedLog string
	ExpectedErr error
	Logger      log.Logger
	LoggerRegex []string
}

// newSFTPTest creates a new test case for NewSFTP.
func newSFTPTest(st sftptest) inttest.ArgumentTest {
	return inttest.ArgumentTest{
		ArgumentTest: test.ArgumentTest{
			Name:        st.Name,
			Argument:    New(st.Location, st.RepoPath, st.Logger),
			ExpectedCLI: st.ExpectedCLI,
			ExpectedLog: st.ExpectedLog,
			ExpectedErr: st.ExpectedErr,
		},
		Logger:      st.Logger,
		LoggerRegex: st.LoggerRegex,
	}
}

// toArgTests converts a list of sftptests to a list of ArgumentTests.
func toArgTests(sftptests []sftptest) []inttest.ArgumentTest {
	argTests := make([]inttest.ArgumentTest, len(sftptests))
	for i, st := range sftptests {
		argTests[i] = newSFTPTest(st)
	}
	return argTests
}

var _ = check.Suite(&inttest.ArgumentSuite{Cmd: "cmd", Arguments: toArgTests([]sftptest{
	{
		Name:     "NewSFTP with password from env",
		Location: newLocation("sftp.example.com", "22", "/backup", ""),
		RepoPath: "repoPath",
		ExpectedCLI: []string{"cmd", "sftp",
			"--host=sftp.example.com",
			"--port=22",
			"--path=/backup/repoPath/",
			"--username=sftpuser",
			"--sftp-password=sftppass123",
		},
		ExpectedLog: "cmd sftp --host=sftp.example.com --port=22 --path=/backup/repoPath/ --username=sftpuser --sftp-password=<****>",
		Logger:      &intlog.StringLogger{},
		LoggerRegex: []string{""},
	},
	{
		Name:     "NewSFTP with custom port",
		Location: newLocation("sftp.example.com", "2222", "/backup", ""),
		RepoPath: "repoPath",
		ExpectedCLI: []string{"cmd", "sftp",
			"--host=sftp.example.com",
			"--port=2222",
			"--path=/backup/repoPath/",
			"--username=sftpuser",
			"--sftp-password=sftppass123",
		},
		ExpectedLog: "cmd sftp --host=sftp.example.com --port=2222 --path=/backup/repoPath/ --username=sftpuser --sftp-password=<****>",
		Logger:      &intlog.StringLogger{},
		LoggerRegex: []string{""},
	},
	{
		Name:     "NewSFTP with minimal config",
		Location: newLocation("sftp.example.com", "", "/backup", ""),
		RepoPath: "",
		ExpectedCLI: []string{"cmd", "sftp",
			"--host=sftp.example.com",
			"--port=22", // defaults to 22 when empty
			"--path=/backup/",
			"--username=sftpuser",
			"--sftp-password=sftppass123",
		},
		ExpectedLog: "cmd sftp --host=sftp.example.com --port=22 --path=/backup/ --username=sftpuser --sftp-password=<****>",
		Logger:      &intlog.StringLogger{},
		LoggerRegex: []string{""},
	},
	{
		Name:        "NewSFTP with empty hostname",
		Location:    newLocation("", "22", "/backup", ""),
		RepoPath:    "repoPath",
		ExpectedCLI: nil,
		ExpectedErr: cli.ErrInvalidHostname,
		Logger:      &intlog.StringLogger{},
		LoggerRegex: []string{""},
	},
	{
		Name:        "NewSFTP with empty path",
		Location:    newLocation("sftp.example.com", "22", "", ""),
		RepoPath:    "", // both path and repoPath empty
		ExpectedCLI: nil,
		ExpectedErr: cli.ErrInvalidRepoPath,
		Logger:      &intlog.StringLogger{},
		LoggerRegex: []string{""},
	},
})})
