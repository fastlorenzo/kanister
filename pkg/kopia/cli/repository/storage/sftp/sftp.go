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

// Package sftp provides functionality for interacting with SFTP storage backends
// in the context of the Kopia repository.
package sftp

import (
	"os"

	"github.com/kanisterio/safecli/command"

	"github.com/kanisterio/kanister/pkg/kopia/cli/internal"
	"github.com/kanisterio/kanister/pkg/log"
)

const (
	// Environment variable names for SFTP credentials
	sftpUsernameEnv   = "SFTP_USERNAME"
	sftpPasswordEnv   = "SFTP_PASSWORD"
	sftpPrivateKeyEnv = "SFTP_PRIVATE_KEY"
)

// New creates a new subcommand for the SFTP storage.
// Credentials (username, password, privateKey) are expected to be available
// via environment variables set by GenerateEnvSpecFromCredentialSecret.
func New(location internal.Location, repoPathPrefix string, logger log.Logger) command.Applier {
	_ = logger // logger parameter required by interface but not currently used
	path := internal.GenerateFullRepoPath(location.SFTPPath(), repoPathPrefix)

	// Get credentials from environment variables
	username := os.Getenv(sftpUsernameEnv)
	password := os.Getenv(sftpPasswordEnv)
	// TODO: privateKey := os.Getenv(sftpPrivateKeyEnv)

	// Build option list
	opts := []command.Applier{
		optHost(location.SFTPHost()),
		optPort(location.SFTPPort()),
		optPath(path),
		optUsername(username),
	}

	// Add password if present
	if password != "" {
		opts = append(opts, optPassword(password))
	}

	// TODO: Add private key support - requires temp file creation
	// if privateKey != "" {
	//     tmpFile := createTempFile(privateKey)
	//     opts = append(opts, optKeyFile(tmpFile))
	// }

	// TODO: Add known hosts support - requires temp file creation
	// knownHosts := location.SFTPKnownHosts()
	// if knownHosts != "" {
	//     tmpFile := createTempFile(knownHosts)
	//     opts = append(opts, optKnownHosts(tmpFile))
	// }

	// Combine subcmd and options
	allArgs := append([]command.Applier{subcmdSFTP}, opts...)
	return command.NewArguments(allArgs...)
}
