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
	"gopkg.in/check.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	secerrors "github.com/kanisterio/kanister/pkg/secrets/errors"
)

type SFTPSecretSuite struct{}

var _ = check.Suite(&SFTPSecretSuite{})

func (s *SFTPSecretSuite) TestValidateSFTPLocation(c *check.C) {
	for i, tc := range []struct {
		secret        Secret
		errChecker    check.Checker
		expectedError error
	}{
		{ // Valid SFTP Secret with all optional fields
			secret: NewSFTPLocation(&corev1.Secret{
				Type: corev1.SecretType(LocTypeSFTP),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sec",
					Namespace: "ns",
				},
				Data: map[string][]byte{
					HostKey:       []byte("sftp.example.com"),
					PortKey:       []byte("2222"),
					PathKey:       []byte("/backup"),
					KnownHostsKey: []byte("sftp.example.com ssh-rsa AAAA..."),
				},
			}),
			errChecker: check.IsNil,
		},
		{ // Valid SFTP Secret with minimal required fields
			secret: NewSFTPLocation(&corev1.Secret{
				Type: corev1.SecretType(LocTypeSFTP),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sec",
					Namespace: "ns",
				},
				Data: map[string][]byte{
					HostKey: []byte("sftp.example.com"),
					PathKey: []byte("/backup"),
				},
			}),
			errChecker: check.IsNil,
		},
		{ // Valid SFTP Secret without port (defaults to 22)
			secret: NewSFTPLocation(&corev1.Secret{
				Type: corev1.SecretType(LocTypeSFTP),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sec",
					Namespace: "ns",
				},
				Data: map[string][]byte{
					HostKey:       []byte("sftp.example.com"),
					PathKey:       []byte("/backup"),
					KnownHostsKey: []byte("sftp.example.com ssh-rsa AAAA..."),
				},
			}),
			errChecker: check.IsNil,
		},
		{ // Missing required field - Host
			secret: NewSFTPLocation(&corev1.Secret{
				Type: corev1.SecretType(LocTypeSFTP),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sec",
					Namespace: "ns",
				},
				Data: map[string][]byte{
					PathKey: []byte("/backup"),
				},
			}),
			errChecker:    check.NotNil,
			expectedError: errkit.Wrap(secerrors.ErrValidate, secerrors.MissingRequiredFieldErrorMsg, HostKey, "ns", "sec"),
		},
		{ // Missing required field - Path
			secret: NewSFTPLocation(&corev1.Secret{
				Type: corev1.SecretType(LocTypeSFTP),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sec",
					Namespace: "ns",
				},
				Data: map[string][]byte{
					HostKey: []byte("sftp.example.com"),
				},
			}),
			errChecker:    check.NotNil,
			expectedError: errkit.Wrap(secerrors.ErrValidate, secerrors.MissingRequiredFieldErrorMsg, PathKey, "ns", "sec"),
		},
		{ // Empty Secret
			secret: NewSFTPLocation(&corev1.Secret{
				Type: corev1.SecretType(LocTypeSFTP),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sec",
					Namespace: "ns",
				},
			}),
			errChecker:    check.NotNil,
			expectedError: errkit.Wrap(secerrors.ErrValidate, secerrors.EmptySecretErrorMessage, "ns", "sec"),
		},
		{ // Nil Secret
			secret:        NewSFTPLocation(nil),
			errChecker:    check.NotNil,
			expectedError: errkit.Wrap(secerrors.ErrValidate, secerrors.NilSecretErrorMessage),
		},
	} {
		err := tc.secret.Validate()
		c.Check(err, tc.errChecker)
		if err != nil {
			c.Check(err.Error(), check.Equals, tc.expectedError.Error(), check.Commentf("test number: %d", i))
		}
	}
}
