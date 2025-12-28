// Copyright 2025 The Kanister Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secrets

import (
	"github.com/kanisterio/errkit"
	corev1 "k8s.io/api/core/v1"

	secerrors "github.com/kanisterio/kanister/pkg/secrets/errors"
)

const (
	// SFTPSecretType represents the secret type for SFTP credentials.
	SFTPSecretType string = "secrets.kanister.io/sftp"

	// SFTPUsername is the key for SFTP username.
	SFTPUsername string = "username"
	// SFTPPassword is the key for SFTP password (optional if using private key).
	SFTPPassword string = "password"
	// SFTPPrivateKey is the key for SFTP private key (optional if using password).
	SFTPPrivateKey string = "privateKey"
)

// ValidateSFTPCredentials validates that the secret has the necessary information
// for SFTP credentials. Either password or privateKey must be provided.
//
// Required field:
// - username
//
// At least one of these must be provided:
// - password
// - privateKey
func ValidateSFTPCredentials(secret *corev1.Secret) error {
	if secret == nil {
		return errkit.New("Nil secret")
	}

	if string(secret.Type) != SFTPSecretType {
		return errkit.Wrap(secerrors.ErrValidate, secerrors.IncompatibleSecretTypeErrorMsg, SFTPSecretType, secret.Namespace, secret.Name)
	}

	// Count known fields
	count := 0
	hasUsername := false
	if _, ok := secret.Data[SFTPUsername]; ok {
		hasUsername = true
		count++
	}
	if _, ok := secret.Data[SFTPPassword]; ok {
		count++
	}
	if _, ok := secret.Data[SFTPPrivateKey]; ok {
		count++
	}

	// Check for unknown fields
	if len(secret.Data) > count {
		return errkit.New("Secret has an unknown field")
	}

	// Username is required
	if !hasUsername {
		return errkit.New("SFTP secret must contain 'username' field")
	}

	// Either password or privateKey must be present
	hasPassword := len(secret.Data[SFTPPassword]) > 0
	hasPrivateKey := len(secret.Data[SFTPPrivateKey]) > 0

	if !hasPassword && !hasPrivateKey {
		return errkit.New("SFTP secret must contain either 'password' or 'privateKey'")
	}

	return nil
}

// SFTPCredentials represents SFTP authentication credentials.
type SFTPCredentials struct {
	Username   string
	Password   string
	PrivateKey string
}

// ExtractSFTPCredentials extracts SFTP credentials from a secret.
func ExtractSFTPCredentials(secret *corev1.Secret) (*SFTPCredentials, error) {
	if err := ValidateSFTPCredentials(secret); err != nil {
		return nil, err
	}

	creds := &SFTPCredentials{
		Username:   string(secret.Data[SFTPUsername]),
		Password:   string(secret.Data[SFTPPassword]),
		PrivateKey: string(secret.Data[SFTPPrivateKey]),
	}

	return creds, nil
}
