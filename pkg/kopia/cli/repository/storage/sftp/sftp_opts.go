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
	"github.com/kanisterio/safecli/command"

	"github.com/kanisterio/kanister/pkg/kopia/cli"
)

const (
	defaultSFTPPort = "22"
)

var (
	subcmdSFTP = command.NewArgument("sftp")
)

// optHost creates a new host option with a given hostname.
// If the hostname is empty, it returns ErrInvalidHostname.
func optHost(host string) command.Applier {
	if host == "" {
		return command.NewErrorArgument(cli.ErrInvalidHostname)
	}
	return command.NewOptionWithArgument("--host", host)
}

// optPort creates a new port option with a given port.
// If the port is empty, it defaults to 22.
func optPort(port string) command.Applier {
	if port == "" {
		port = defaultSFTPPort
	}
	return command.NewOptionWithArgument("--port", port)
}

// optPath creates a new path option with a given path.
// If the path is empty, it returns ErrInvalidRepoPath.
func optPath(path string) command.Applier {
	if path == "" {
		return command.NewErrorArgument(cli.ErrInvalidRepoPath)
	}
	return command.NewOptionWithArgument("--path", path)
}

// optUsername creates a new username option with a given username.
// If the username is empty, it returns an error.
func optUsername(username string) command.Applier {
	if username == "" {
		return command.NewErrorArgument(cli.ErrInvalidUsername)
	}
	return command.NewOptionWithArgument("--username", username)
}

// optPassword creates a new password option with a given password.
// The password is redacted in logs.
// If the password is empty, the password option is not set.
func optPassword(password string) command.Applier {
	if password == "" {
		return command.NewNoopArgument()
	}
	return command.NewOptionWithRedactedArgument("--sftp-password", password)
}

// optKeyFile creates a new keyfile option with a given path.
// If the path is empty, the keyfile option is not set.
func optKeyFile(keyfile string) command.Applier {
	if keyfile == "" {
		return command.NewNoopArgument()
	}
	return command.NewOptionWithArgument("--keyfile", keyfile)
}

// optKnownHosts creates a new known-hosts option with a given path.
// If the path is empty, the known-hosts option is not set.
func optKnownHosts(knownHosts string) command.Applier {
	if knownHosts == "" {
		return command.NewNoopArgument()
	}
	return command.NewOptionWithArgument("--known-hosts", knownHosts)
}
