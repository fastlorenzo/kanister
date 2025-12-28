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
	"gopkg.in/check.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SFTPCredentialSecretSuite struct{}

var _ = check.Suite(&SFTPCredentialSecretSuite{})

func (s *SFTPCredentialSecretSuite) TestValidateSFTPCredentials(c *check.C) {
	for _, tc := range []struct {
		secret     *corev1.Secret
		errChecker check.Checker
	}{
		{ // Valid SFTP credential secret with password
			secret: &corev1.Secret{
				Type: corev1.SecretType(SFTPSecretType),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sftp-creds",
					Namespace: "default",
				},
				Data: map[string][]byte{
					SFTPUsername: []byte("sftpuser"),
					SFTPPassword: []byte("secret-password"),
				},
			},
			errChecker: check.IsNil,
		},
		{ // Valid SFTP credential secret with private key
			secret: &corev1.Secret{
				Type: corev1.SecretType(SFTPSecretType),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sftp-creds",
					Namespace: "default",
				},
				Data: map[string][]byte{
					SFTPUsername:   []byte("sftpuser"),
					SFTPPrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\n..."),
				},
			},
			errChecker: check.IsNil,
		},
		{ // Valid SFTP credential secret with both password and private key
			secret: &corev1.Secret{
				Type: corev1.SecretType(SFTPSecretType),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sftp-creds",
					Namespace: "default",
				},
				Data: map[string][]byte{
					SFTPUsername:   []byte("sftpuser"),
					SFTPPassword:   []byte("secret-password"),
					SFTPPrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\n..."),
				},
			},
			errChecker: check.IsNil,
		},
		{ // Invalid - missing username
			secret: &corev1.Secret{
				Type: corev1.SecretType(SFTPSecretType),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sftp-creds",
					Namespace: "default",
				},
				Data: map[string][]byte{
					SFTPPassword: []byte("secret-password"),
				},
			},
			errChecker: check.NotNil,
		},
		{ // Invalid - missing both password and private key
			secret: &corev1.Secret{
				Type: corev1.SecretType(SFTPSecretType),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sftp-creds",
					Namespace: "default",
				},
				Data: map[string][]byte{
					SFTPUsername: []byte("sftpuser"),
				},
			},
			errChecker: check.NotNil,
		},
		{ // Invalid - wrong secret type
			secret: &corev1.Secret{
				Type: corev1.SecretType("Opaque"),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sftp-creds",
					Namespace: "default",
				},
				Data: map[string][]byte{
					SFTPUsername: []byte("sftpuser"),
					SFTPPassword: []byte("secret-password"),
				},
			},
			errChecker: check.NotNil,
		},
		{ // Invalid - unknown field
			secret: &corev1.Secret{
				Type: corev1.SecretType(SFTPSecretType),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sftp-creds",
					Namespace: "default",
				},
				Data: map[string][]byte{
					SFTPUsername: []byte("sftpuser"),
					SFTPPassword: []byte("secret-password"),
					"unknown":    []byte("value"),
				},
			},
			errChecker: check.NotNil,
		},
		{ // Invalid - nil secret
			secret:     nil,
			errChecker: check.NotNil,
		},
	} {
		err := ValidateSFTPCredentials(tc.secret)
		c.Check(err, tc.errChecker)
	}
}

func (s *SFTPCredentialSecretSuite) TestExtractSFTPCredentials(c *check.C) {
	// Valid credentials with password
	secret := &corev1.Secret{
		Type: corev1.SecretType(SFTPSecretType),
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sftp-creds",
			Namespace: "default",
		},
		Data: map[string][]byte{
			SFTPUsername: []byte("sftpuser"),
			SFTPPassword: []byte("secret-password"),
		},
	}

	creds, err := ExtractSFTPCredentials(secret)
	c.Assert(err, check.IsNil)
	c.Assert(creds, check.NotNil)
	c.Check(creds.Username, check.Equals, "sftpuser")
	c.Check(creds.Password, check.Equals, "secret-password")
	c.Check(creds.PrivateKey, check.Equals, "")

	// Valid credentials with private key
	secret = &corev1.Secret{
		Type: corev1.SecretType(SFTPSecretType),
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sftp-creds",
			Namespace: "default",
		},
		Data: map[string][]byte{
			SFTPUsername:   []byte("sftpuser"),
			SFTPPrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----\n..."),
		},
	}

	creds, err = ExtractSFTPCredentials(secret)
	c.Assert(err, check.IsNil)
	c.Assert(creds, check.NotNil)
	c.Check(creds.Username, check.Equals, "sftpuser")
	c.Check(creds.Password, check.Equals, "")
	c.Check(creds.PrivateKey, check.Equals, "-----BEGIN RSA PRIVATE KEY-----\n...")
}
