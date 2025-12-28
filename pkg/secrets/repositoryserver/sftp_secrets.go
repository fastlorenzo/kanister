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

package repositoryserver

import (
	"github.com/kanisterio/errkit"
	corev1 "k8s.io/api/core/v1"

	secerrors "github.com/kanisterio/kanister/pkg/secrets/errors"
)

var _ Secret = &SFTP{}

type SFTP struct {
	storageLocation *corev1.Secret
}

func NewSFTPLocation(secret *corev1.Secret) *SFTP {
	return &SFTP{
		storageLocation: secret,
	}
}

func (s *SFTP) Validate() error {
	if s.storageLocation == nil {
		return errkit.Wrap(secerrors.ErrValidate, secerrors.NilSecretErrorMessage)
	}
	if len(s.storageLocation.Data) == 0 {
		return errkit.Wrap(secerrors.ErrValidate, secerrors.EmptySecretErrorMessage, s.storageLocation.Namespace, s.storageLocation.Name)
	}

	// Required fields: host and path
	if _, ok := s.storageLocation.Data[HostKey]; !ok {
		return errkit.Wrap(secerrors.ErrValidate, secerrors.MissingRequiredFieldErrorMsg, HostKey, s.storageLocation.Namespace, s.storageLocation.Name)
	}
	if _, ok := s.storageLocation.Data[PathKey]; !ok {
		return errkit.Wrap(secerrors.ErrValidate, secerrors.MissingRequiredFieldErrorMsg, PathKey, s.storageLocation.Namespace, s.storageLocation.Name)
	}

	// Optional fields: port, knownHosts
	// Port defaults to 22 if not provided
	// KnownHosts is optional for strict host key checking

	return nil
}
