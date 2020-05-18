// +build apitests

/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package tests

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/netfoundry/ziti-edge/controller/model"
	"github.com/netfoundry/ziti-foundation/common/constants"
	nfpem "github.com/netfoundry/ziti-foundation/util/pem"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"testing"
)

func Test_Authenticate_Cert(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.teardown()
	ctx.startServer()

	ctx.requireAdminLogin()

	_, certAuthenticator := ctx.AdminSession.requireCreateIdentityOttEnrollment("test", false)

	var tests = &authCertTests{
		ctx:               ctx,
		certAuthenticator: certAuthenticator,
	}

	t.Run("cert authenticator has full pem for newly created identities", tests.testAuthenticateCertStoresAndFillsFullCert)
	t.Run("login with valid certificate and no client info", tests.testAuthenticateValidCertEmptyBody)
	t.Run("login with valid certificate and client info", tests.testAuthenticateValidCertValidClientInfoBody)
	t.Run("login with valid certificate and invalid JSON client info", tests.testAuthenticateValidCertInvalidJson)
	t.Run("login with valid certificate and client info with extra properties", tests.testAuthenticateValidCertValidClientInfoWithExtraProperties)
	t.Run("login with invalid certificate and no client info", tests.testAuthenticateInvalidCert)
}

type authCertTests struct {
	ctx               *TestContext
	certAuthenticator *certAuthenticator
}

func (test *authCertTests) testAuthenticateCertStoresAndFillsFullCert(t *testing.T) {

	t.Run("newly created cert authenticators have full cert stored as PEM", func(t *testing.T) {
		r := require.New(t)
		authenticator, err := test.ctx.EdgeController.AppEnv.Handlers.Authenticator.ReadByFingerprint(test.certAuthenticator.Fingerprint())

		r.NoError(err)

		certAuth, ok := authenticator.SubType.(*model.AuthenticatorCert)

		r.True(ok, "authenticator was not a certificate type, got: %s", reflect.TypeOf(authenticator.SubType))

		r.NotEmpty(certAuth.Pem, "cert authenticator pem was empty/blank")
	})

	t.Run("cert authenticators with blank pem is stored on authenticate", func(t *testing.T) {
		r := require.New(t)
		authenticator, err := test.ctx.EdgeController.AppEnv.Handlers.Authenticator.ReadByFingerprint(test.certAuthenticator.Fingerprint())

		r.NoError(err)

		certAuth, ok := authenticator.SubType.(*model.AuthenticatorCert)

		r.True(ok, "authenticator was not a certificate type, got: %s", reflect.TypeOf(authenticator.SubType))

		certAuth.Pem = ""

		err = test.ctx.EdgeController.AppEnv.Handlers.Authenticator.Update(authenticator)
		r.NoError(err)

		authenticator, err = test.ctx.EdgeController.AppEnv.Handlers.Authenticator.ReadByFingerprint(test.certAuthenticator.Fingerprint())

		r.NoError(err)

		certAuth, ok = authenticator.SubType.(*model.AuthenticatorCert)

		r.True(ok, "authenticator was not a certificate type, got: %s", reflect.TypeOf(authenticator.SubType))

		r.Empty(certAuth.Pem, "cert authenticator pem was not set to empty/blank")

		testClient, _, transport := test.ctx.NewClientComponents()

		transport.TLSClientConfig.Certificates = test.certAuthenticator.TLSCertificates()
		resp, err := testClient.NewRequest().
			SetHeader("Content-Type", "application/json").
			Post("/authenticate?method=cert")

		standardJsonResponseTests(resp, http.StatusOK, t)

		authenticator, err = test.ctx.EdgeController.AppEnv.Handlers.Authenticator.ReadByFingerprint(test.certAuthenticator.Fingerprint())

		r.NoError(err)

		certAuth, ok = authenticator.SubType.(*model.AuthenticatorCert)

		r.True(ok, "authenticator was not a certificate type, got: %s", reflect.TypeOf(authenticator.SubType))

		r.NotEmpty(certAuth.Pem, "cert authenticator pem was empty/blank after authenticating")
	})
}

func (test *authCertTests) testAuthenticateValidCertEmptyBody(t *testing.T) {
	testClient, _, transport := test.ctx.NewClientComponents()

	transport.TLSClientConfig.Certificates = test.certAuthenticator.TLSCertificates()
	resp, err := testClient.NewRequest().
		SetHeader("Content-Type", "application/json").
		Post("/authenticate?method=cert")

	t.Run("returns without error", func(t *testing.T) {
		require.New(t).NoError(err)
	})

	standardJsonResponseTests(resp, http.StatusOK, t)

	t.Run("returns a session token HTTP headers", func(t *testing.T) {
		require.New(t).NotEmpty(resp.Header().Get(constants.ZitiSession), fmt.Sprintf("HTTP header %s is empty", constants.ZitiSession))
	})

	t.Run("returns a session token in body", func(t *testing.T) {
		r := require.New(t)
		data, err := gabs.ParseJSON(resp.Body())

		r.NoError(err)

		r.True(data.ExistsP("data.token"), "session token property in 'data.token' as not found")
		r.NotEmpty(data.Path("data.token").String(), "session token property in 'data.token' is empty")
	})

	t.Run("body session token matches HTTP header token", func(t *testing.T) {
		r := require.New(t)
		data, err := gabs.ParseJSON(resp.Body())

		r.NoError(err)

		bodyToken := data.Path("data.token").Data().(string)
		headerToken := resp.Header().Get(constants.ZitiSession)
		r.Equal(bodyToken, headerToken)
	})

	t.Run("returns an identity", func(t *testing.T) {
		r := require.New(t)
		data, err := gabs.ParseJSON(resp.Body())

		r.NoError(err)

		r.True(data.ExistsP("data.identity"), "session token property in 'data.token' as not found")

		_, err = data.ObjectP("data.identity")
		r.NoError(err, "session token property in 'data.token' is empty")
	})
}

func (test *authCertTests) testAuthenticateValidCertValidClientInfoBody(t *testing.T) {
	testClient, _, transport := test.ctx.NewClientComponents()

	transport.TLSClientConfig.Certificates = test.certAuthenticator.TLSCertificates()

	bodyJson := `{
  "envInfo": {"os": "windows", "arch": "amd64", "osRelease": "6.2.9200", "osVersion": "6.2.9200"},
  "sdkInfo": {"type": "ziti-sdk-golang", "branch": "unknown", "version": "0.0.0", "revision": "unknown"}
}`
	resp, err := testClient.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyJson).
		Post("/authenticate?method=cert")

	t.Run("returns without error", func(t *testing.T) {
		require.New(t).NoError(err)
	})

	t.Run("returns 200", func(t *testing.T) {
		require.New(t).Equal(http.StatusOK, resp.StatusCode())
	})

	standardJsonResponseTests(resp, http.StatusOK, t)

	t.Run("returns a session token HTTP headers", func(t *testing.T) {
		require.New(t).NotEmpty(resp.Header().Get(constants.ZitiSession), fmt.Sprintf("HTTP header %s is empty", constants.ZitiSession))
	})

	t.Run("returns a session token in body", func(t *testing.T) {
		r := require.New(t)
		data, err := gabs.ParseJSON(resp.Body())

		r.NoError(err)

		r.True(data.ExistsP("data.token"), "session token property in 'data.token' as not found")
		r.NotEmpty(data.Path("data.token").String(), "session token property in 'data.token' is empty")
	})

	t.Run("body session token matches HTTP header token", func(t *testing.T) {
		r := require.New(t)
		data, err := gabs.ParseJSON(resp.Body())

		r.NoError(err)

		bodyToken := data.Path("data.token").Data().(string)
		headerToken := resp.Header().Get(constants.ZitiSession)
		r.Equal(bodyToken, headerToken)
	})

	t.Run("returns an identity", func(t *testing.T) {
		r := require.New(t)
		data, err := gabs.ParseJSON(resp.Body())

		r.NoError(err)

		r.True(data.ExistsP("data.identity"), "session token property in 'data.token' as not found")

		_, err = data.ObjectP("data.identity")
		r.NoError(err, "session token property in 'data.token' is empty")
	})

	t.Run("client info is set on the identity", func(t *testing.T) {
		test.ctx.testContextChanged(t)
		r := test.ctx.req

		data, err := gabs.ParseJSON(resp.Body())

		r.NoError(err)

		r.True(data.ExistsP("data.identity.id"), "identity id not found")
		identityId := data.Path("data.identity.id").Data().(string)
		r.NotEmpty(identityId)

		r.True(data.ExistsP("data.token"), "token not found")
		token := data.Path("data.token").Data().(string)
		r.NotEmpty(token)

		resp, err := test.ctx.AdminSession.newRequest(test.ctx).Get("identities/" + identityId)
		r.NoError(err)

		r.Equal(http.StatusOK, resp.StatusCode())

		identity, err := gabs.ParseJSON(resp.Body())
		r.NoError(err)

		sentInfo, err := gabs.ParseJSON([]byte(bodyJson))
		r.NoError(err)

		sentEnvInfo := sentInfo.Path("envInfo").Data().(map[string]interface{})
		sentSdkInfo := sentInfo.Path("sdkInfo").Data().(map[string]interface{})

		envInfo := identity.Path("data.envInfo").Data().(map[string]interface{})
		r.Equal(sentEnvInfo, envInfo)

		sdkInfo := identity.Path("data.sdkInfo").Data().(map[string]interface{})
		r.Equal(sentSdkInfo, sdkInfo)
	})

	t.Run("client info is updated on the identity", func(t *testing.T) {
		test.ctx.testContextChanged(t)
		r := test.ctx.req

		secondInfo := `{
  "envInfo": {"os": "updatedValueOs", "arch": "updatedValueArch", "osRelease": "updatedValueRelease", "osVersion": "updatedValueOsRelease"},
  "sdkInfo": {"type": "updatedValueType", "branch": "updatedValueBranch", "version": "updatedValueVersion", "revision": "updatedValueRevision"}
}`
		authResp, err := testClient.NewRequest().
			SetHeader("Content-Type", "application/json").
			SetBody(secondInfo).
			Post("/authenticate?method=cert")
		r.NoError(err)
		r.Equal(http.StatusOK, authResp.StatusCode())

		authData, err := gabs.ParseJSON(authResp.Body())
		r.NoError(err)

		identityId := authData.Path("data.identity.id").Data().(string)

		resp, err := test.ctx.AdminSession.newRequest(test.ctx).Get("identities/" + identityId)
		r.NoError(err)

		r.Equal(http.StatusOK, resp.StatusCode())

		identity, err := gabs.ParseJSON(resp.Body())
		r.NoError(err)

		sentInfo, err := gabs.ParseJSON([]byte(secondInfo))
		r.NoError(err)

		sentEnvInfo := sentInfo.Path("envInfo").Data().(map[string]interface{})
		sentSdkInfo := sentInfo.Path("sdkInfo").Data().(map[string]interface{})

		envInfo := identity.Path("data.envInfo").Data().(map[string]interface{})
		r.Equal(sentEnvInfo, envInfo)

		sdkInfo := identity.Path("data.sdkInfo").Data().(map[string]interface{})
		r.Equal(sentSdkInfo, sdkInfo)
	})
}

func (test *authCertTests) testAuthenticateValidCertInvalidJson(t *testing.T) {
	testClient, _, transport := test.ctx.NewClientComponents()

	transport.TLSClientConfig.Certificates = test.certAuthenticator.TLSCertificates()

	bodyJson := "i will not parse"
	resp, err := testClient.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyJson).
		Post("/authenticate?method=cert")

	t.Run("returns without error", func(t *testing.T) {
		require.New(t).NoError(err)
	})

	standardErrorJsonResponseTests(resp, "COULD_NOT_PARSE_BODY", http.StatusBadRequest, t)

	t.Run("returns without a ziti session header", func(t *testing.T) {
		require.New(t).Equal("", resp.Header().Get(constants.ZitiSession))
	})
}

func (test *authCertTests) testAuthenticateValidCertValidClientInfoWithExtraProperties(t *testing.T) {
	testClient, _, transport := test.ctx.NewClientComponents()

	transport.TLSClientConfig.Certificates = test.certAuthenticator.TLSCertificates()

	bodyJson := `{"envInfo": {"os": "windows", "arch": "amd64", "osRelease": "6.2.9200", "osVersion": "6.2.9200", "extraProp1":"extraVal1"},
  "sdkInfo": {"type": "ziti-sdk-golang", "branch": "unknown", "version": "0.0.0", "revision": "unknown", "extraProp2":"extraVal2"},
  "extraProp3": "extraVal3"}`
	resp, err := testClient.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyJson).
		Post("/authenticate?method=cert")

	t.Run("returns without error", func(t *testing.T) {
		require.New(t).NoError(err)
	})

	standardJsonResponseTests(resp, http.StatusOK, t)
}

func (test *authCertTests) testAuthenticateInvalidCert(t *testing.T) {
	r := require.New(t)

	testClient, _, transport := test.ctx.NewClientComponents()

	certAndKeyPem := `-----BEGIN CERTIFICATE-----
MIICyjCCAlCgAwIBAgIRAMbo6szcFH+1lrByi/UvSiMwCgYIKoZIzj0EAwIwZDEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAk5ZMR8wHQYDVQQKDBZTb21lRmFrZUNvcnAg
M3JkIFBhcnR5MScwJQYDVQQDDB5Tb21lRmFrZUNvcnAgM3JkIFBhcnR5IFJvb3Qg
Q0EwHhcNMTkwNTA3MTIzMTU4WhcNMjAwNTE2MTIzMTU4WjBkMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCTlkxHzAdBgNVBAoMFlNvbWVGYWtlQ29ycCAzcmQgUGFydHkx
JzAlBgNVBAMMHlNvbWVGYWtlQ29ycCAzcmQgUGFydHkgUm9vdCBDQTB2MBAGByqG
SM49AgEGBSuBBAAiA2IABN79IXogRQMB3Q3w6JRXXjNcr75UnSnDKY1xfhnyB4zT
WRtcgMeI+FvqFBekFI9iihgNQVMQ1EdIWz9ThzfELKW2tnDhs+fYY7Gu/UJdEQJ8
eCVj687QIW6ZA5Xlt5zZH6OBxTCBwjAJBgNVHRMEAjAAMBEGCWCGSAGG+EIBAQQE
AwIFoDAzBglghkgBhvhCAQ0EJhYkT3BlblNTTCBHZW5lcmF0ZWQgQ2xpZW50IENl
cnRpZmljYXRlMB0GA1UdDgQWBBRYn3XCinkndBO/YsvYiAR+DPa7UzAfBgNVHSME
GDAWgBQf+8pahQr268huQRqUsxmS4LbIVDAOBgNVHQ8BAf8EBAMCBeAwHQYDVR0l
BBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMEMAoGCCqGSM49BAMCA2gAMGUCMQDHyoN8
Y7ZF604e2V/c+S9OZb1JG6x3ZoPsNoHlR/sRr6JNNvqOk89U1uZ8huXJ5eUCMDTh
97wUwCPC3Se2xMm6eHcc+q/EqFFadDQSGIsUm7Pt1Af6S7c9LCVD9keTM5DGcg==
-----END CERTIFICATE-----
-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDAPgK7rxXOfOIqTAfSfeJDYKeIsa5keKS7XFhy/OnsEARUNQrALCniy
ccbzsr2ti0KgBwYFK4EEACKhZANiAATe/SF6IEUDAd0N8OiUV14zXK++VJ0pwymN
cX4Z8geM01kbXIDHiPhb6hQXpBSPYooYDUFTENRHSFs/U4c3xCyltrZw4bPn2GOx
rv1CXRECfHglY+vO0CFumQOV5bec2R8=
-----END EC PRIVATE KEY-----`

	blocks := nfpem.DecodeAll([]byte(certAndKeyPem))
	r.Len(blocks, 3, "cert & key pair pem blocks did not parse, expected 2 blocks, got: %d", len(blocks))

	cert, err := x509.ParseCertificate(blocks[0].Bytes)
	r.NoError(err)

	key, err := x509.ParseECPrivateKey(blocks[2].Bytes)
	r.NoError(err)

	transport.TLSClientConfig.Certificates = []tls.Certificate{
		{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key,
		},
	}

	resp, err := testClient.NewRequest().
		SetHeader("Content-Type", "application/json").
		Post("/authenticate?method=cert")

	t.Run("returns without error", func(t *testing.T) {
		require.New(t).NoError(err)
	})

	standardErrorJsonResponseTests(resp, "INVALID_AUTH", http.StatusUnauthorized, t)

	t.Run("returns without a ziti session header", func(t *testing.T) {
		require.New(t).Equal("", resp.Header().Get(constants.ZitiSession))
	})
}
