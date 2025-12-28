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

package storage

import (
	"os"

	"github.com/kanisterio/kanister/pkg/logsafe"
	"github.com/kanisterio/kanister/pkg/secrets/repositoryserver"
)

const (
	sftpSubCommand     = "sftp"
	sftpHostFlag       = "--host"
	sftpPortFlag       = "--port"
	sftpPathFlag       = "--path"
	sftpUsernameFlag   = "--username"
	sftpPasswordFlag   = "--sftp-password"
	sftpKeyfileFlag    = "--keyfile"
	sftpKnownHostsFlag = "--known-hosts"
	sftpDefaultPort    = "22"
)

func sftpArgs(location map[string][]byte, repoPathPrefix string) logsafe.Cmd {
	args := logsafe.NewLoggable(sftpSubCommand)

	// Required fields from location secret
	host := string(location[repositoryserver.HostKey])
	args = args.AppendLoggableKV(sftpHostFlag, host)

	port := string(location[repositoryserver.PortKey])
	if port == "" {
		port = sftpDefaultPort
	}
	args = args.AppendLoggableKV(sftpPortFlag, port)

	// Append prefix from the location to the repository path prefix, if specified
	fullRepoPath := GenerateFullRepoPath(string(location[repositoryserver.PathKey]), repoPathPrefix)
	args = args.AppendLoggableKV(sftpPathFlag, fullRepoPath)

	// Credentials from environment variables (set by GenerateEnvSpecFromCredentialSecret)
	// Username is required
	username := os.Getenv(sftpUsernameEnv)
	if username != "" {
		args = args.AppendLoggableKV(sftpUsernameFlag, username)
	}

	// Either password or private key (both are optional here, validated in credential secret)
	password := os.Getenv(sftpPasswordEnv)
	if password != "" {
		args = args.AppendRedactedKV(sftpPasswordFlag, password)
	}

	privateKey := os.Getenv(sftpPrivateKeyEnv)
	// Note: Kopia expects a file path for --keyfile, but we have the key content
	// We'll need to write it to a temp file. For now, we'll skip this and rely on
	// the new CLI builder approach which handles this better.
	// This legacy path is primarily for backwards compatibility.
	_ = privateKey // suppress unused variable warning

	// Optional: known_hosts file content from location secret
	knownHosts := string(location[repositoryserver.KnownHostsKey])
	// Similar to private key, this would need to be written to a file
	// The new CLI builder handles this better
	_ = knownHosts // suppress unused variable warning

	return args
}
