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
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/google/uuid"
	"github.com/openziti/edge/controller/apierror"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_Authenticators_AdminUsingAdminEndpoints(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.teardown()
	ctx.startServer()
	ctx.requireAdminLogin()

	_, _ = ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
	_, _ = ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)

	t.Run("can list authenticators", func(t *testing.T) {
		req := require.New(t)
		resp, err := ctx.AdminSession.newAuthenticatedRequest().Get("/authenticators")

		req.NoError(err)

		standardJsonResponseTests(resp, http.StatusOK, t)

		authenticatorsBody, err := gabs.ParseJSON(resp.Body())

		t.Run("can see three authenticators", func(t *testing.T) {
			req := require.New(t)
			count, err := authenticatorsBody.ArrayCount("data")

			req.NoError(err)
			req.Equal(3, count, "expected three authenticators")
		})

	})
	t.Run("can show details of an authenticator", func(t *testing.T) {
		req := require.New(t)
		listResp, err := ctx.AdminSession.newAuthenticatedRequest().Get("/authenticators")

		req.NoError(err)

		standardJsonResponseTests(listResp, http.StatusOK, t)

		authenticatorsBody, err := gabs.ParseJSON(listResp.Body())

		req.NoError(err)

		authenticatorId := authenticatorsBody.Path("data").Index(0).Path("id").Data().(string)
		req.NotEmpty(authenticatorId)

		detailResp, err := ctx.AdminSession.newAuthenticatedRequest().Get("/authenticators/" + authenticatorId)

		standardJsonResponseTests(detailResp, http.StatusOK, t)
	})

	t.Run("can create updb authenticator for a different identity", func(t *testing.T) {
		ctx.testContextChanged(t)

		identityId := ctx.AdminSession.requireCreateIdentity(uuid.New().String(), false)
		username := uuid.New().String()
		password := uuid.New().String()

		body := gabs.New()
		_, _ = body.Set(identityId, "identityId")
		_, _ = body.Set("updb", "method")
		_, _ = body.Set(username, "username")
		_, _ = body.Set(password, "password")

		resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Post("/authenticators")

		ctx.req.NoError(err)
		standardJsonResponseTests(resp, http.StatusCreated, t)

		t.Run("and the new authenticator can be used for authentication", func(t *testing.T) {
			req := require.New(t)
			authenticator := &updbAuthenticator{
				Username:    username,
				Password:    password,
				ConfigTypes: nil,
			}

			session, err := authenticator.Authenticate(ctx)

			req.NoError(err)
			req.NotEmpty(session.id)

		})
	})

	t.Run("cannot create a updb authenticator for an identity with an existing updb authenticator", func(t *testing.T) {
		ctx.testContextChanged(t)

		identityId, _ := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)

		username := uuid.New().String()
		password := uuid.New().String()

		body := gabs.New()
		_, _ = body.Set(identityId, "identityId")
		_, _ = body.Set("updb", "method")
		_, _ = body.Set(username, "username")
		_, _ = body.Set(password, "password")

		resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Post("/authenticators")

		ctx.req.NoError(err)
		standardErrorJsonResponseTests(resp, apierror.AuthenticatorMethodMaxCode, apierror.AuthenticatorMethodMaxStatus, t)
	})

	t.Run("can update updb authenticator for a different identity", func(t *testing.T) {
		ctx.testContextChanged(t)

		identityId, _ := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)

		result, err := ctx.AdminSession.newAuthenticatedRequest().Get(fmt.Sprintf(`/authenticators?filter=identity="%s"`, identityId))
		ctx.req.NoError(err)

		resultBody, err := gabs.ParseJSON(result.Body())
		ctx.req.NoError(err)

		idContainer := resultBody.Path("data").Index(0).Path("id")
		ctx.req.NotEmpty(idContainer)

		authenticatorId := idContainer.Data().(string)
		ctx.req.NotEmpty(authenticatorId)

		newUsername := uuid.New().String()
		newPassword := uuid.New().String()

		body := gabs.New()
		_, _ = body.Set(newUsername, "username")
		_, _ =
			body.Set(newPassword, "password")

		resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Put("/authenticators/" + authenticatorId)

		ctx.req.NoError(err)
		standardJsonResponseTests(resp, http.StatusOK, t)

		t.Run("newly updated newPassword can be used for authentication", func(t *testing.T) {
			ctx.testContextChanged(t)

			auth := updbAuthenticator{
				Username:    newUsername,
				Password:    newPassword,
				ConfigTypes: nil,
			}

			session, err := auth.Authenticate(ctx)
			ctx.req.NoError(err)
			ctx.req.NotEmpty(session)
		})
	})
	t.Run("can patch updb authenticator for a different identity", func(t *testing.T) {
		t.Run("when patching username only", func(t *testing.T) {
			ctx.testContextChanged(t)
			identityId, authenticator := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)

			result, err := ctx.AdminSession.newAuthenticatedRequest().Get(fmt.Sprintf(`/authenticators?filter=identity="%s"`, identityId))
			ctx.req.NoError(err)

			resultBody, err := gabs.ParseJSON(result.Body())
			ctx.req.NoError(err)

			idContainer := resultBody.Path("data").Index(0).Path("id")
			ctx.req.NotEmpty(idContainer)

			authenticatorId := idContainer.Data().(string)
			ctx.req.NotEmpty(authenticatorId)

			newUsername := uuid.New().String()

			body := gabs.New()
			_, _ = body.Set(newUsername, "username")

			resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Patch("/authenticators/" + authenticatorId)

			ctx.req.NoError(err)
			standardJsonResponseTests(resp, http.StatusOK, t)

			t.Run("newly updated authenticator can be used for authentication", func(t *testing.T) {
				ctx.testContextChanged(t)

				authenticator.Username = newUsername

				session, err := authenticator.Authenticate(ctx)
				ctx.req.NoError(err)
				ctx.req.NotEmpty(session)
			})
		})

		t.Run("when patching password only", func(t *testing.T) {
			ctx.testContextChanged(t)
			identityId, authenticator := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)

			result, err := ctx.AdminSession.newAuthenticatedRequest().Get(fmt.Sprintf(`/authenticators?filter=identity="%s"`, identityId))
			ctx.req.NoError(err)

			resultBody, err := gabs.ParseJSON(result.Body())
			ctx.req.NoError(err)

			idContainer := resultBody.Path("data").Index(0).Path("id")
			ctx.req.NotEmpty(idContainer)

			authenticatorId := idContainer.Data().(string)
			ctx.req.NotEmpty(authenticatorId)

			newPassword := uuid.New().String()

			body := gabs.New()
			_, _ = body.Set(newPassword, "password")

			resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Patch("/authenticators/" + authenticatorId)

			ctx.req.NoError(err)
			standardJsonResponseTests(resp, http.StatusOK, t)

			t.Run("newly patched authenticator can be used for authentication", func(t *testing.T) {
				ctx.testContextChanged(t)

				authenticator.Password = newPassword

				session, err := authenticator.Authenticate(ctx)
				ctx.req.NoError(err)
				ctx.req.NotEmpty(session)
			})
		})

	})
	t.Run("can delete updb authenticator for a different identity", func(t *testing.T) {
		ctx.testContextChanged(t)
		identityId, authenticator := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)

		result, err := ctx.AdminSession.newAuthenticatedRequest().Get(fmt.Sprintf(`/authenticators?filter=identity="%s"`, identityId))
		ctx.req.NoError(err)

		resultBody, err := gabs.ParseJSON(result.Body())
		ctx.req.NoError(err)

		idContainer := resultBody.Path("data").Index(0).Path("id")
		ctx.req.NotEmpty(idContainer)

		authenticatorId := idContainer.Data().(string)
		ctx.req.NotEmpty(authenticatorId)

		resp, err := ctx.AdminSession.newAuthenticatedRequest().Delete("/authenticators/" + authenticatorId)

		ctx.req.NoError(err)

		standardJsonResponseTests(resp, http.StatusOK, t)

		t.Run("identity can not longer authenticate", func(t *testing.T) {
			ctx.testContextChanged(t)
			session, err := authenticator.Authenticate(ctx)

			ctx.req.Error(err)
			ctx.req.Empty(session)
		})
	})

	//cert
	t.Run("can create cert authenticator for a different identity", func(t *testing.T) {
		ctx.testContextChanged(t)

		//used to receive a cert that can be tested with
		unusedIdentityId, unusedCertAuth := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)
		ctx.AdminSession.deleteEntityOfType("identity", unusedIdentityId)

		identityId := ctx.AdminSession.requireCreateIdentity(uuid.New().String(), false)

		body := gabs.New()
		_, _ = body.Set(identityId, "identityId")
		_, _ = body.Set("cert", "method")
		_, _ = body.Set(unusedCertAuth.certPem, "certPem")

		resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Post("/authenticators")

		ctx.req.NoError(err)
		standardJsonResponseTests(resp, http.StatusCreated, t)

		t.Run("and the new authenticator can be used for authentication", func(t *testing.T) {
			req := require.New(t)

			session, err := unusedCertAuth.Authenticate(ctx)

			req.NoError(err)
			req.NotEmpty(session.id)

		})
	})

	t.Run("cannot create a cert authenticator for an identity with an existing cert authenticator", func(t *testing.T) {
		ctx.testContextChanged(t)

		identityId, _ := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)

		body := gabs.New()
		_, _ = body.Set(identityId, "identityId")
		_, _ = body.Set("cert", "method")
		_, _ = body.Set("doesnotmatter", "certPem")
		resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Post("/authenticators")

		ctx.req.NoError(err)
		standardErrorJsonResponseTests(resp, apierror.AuthenticatorMethodMaxCode, apierror.AuthenticatorMethodMaxStatus, t)
	})

	t.Run("can update cert authenticator for a different identity", func(t *testing.T) {
		ctx.testContextChanged(t)

		//used to receive a cert that can be tested with
		unusedIdentityId, unusedCertAuth := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)
		ctx.AdminSession.deleteEntityOfType("identity", unusedIdentityId)

		identityId, _ := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)

		result, err := ctx.AdminSession.newAuthenticatedRequest().Get(fmt.Sprintf(`/authenticators?filter=identity="%s"`, identityId))
		ctx.req.NoError(err)

		resultBody, err := gabs.ParseJSON(result.Body())
		ctx.req.NoError(err)

		idContainer := resultBody.Path("data").Index(0).Path("id")
		ctx.req.NotEmpty(idContainer)

		authenticatorId := idContainer.Data().(string)
		ctx.req.NotEmpty(authenticatorId)

		body := gabs.New()
		_, _ = body.Set(unusedCertAuth.certPem, "certPem")

		resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Put("/authenticators/" + authenticatorId)

		ctx.req.NoError(err)
		standardJsonResponseTests(resp, http.StatusOK, t)

		t.Run("newly updated cert can be used for authentication", func(t *testing.T) {
			ctx.testContextChanged(t)

			session, err := unusedCertAuth.Authenticate(ctx)
			ctx.req.NoError(err)
			ctx.req.NotEmpty(session)
		})
	})

	t.Run("can patch cert authenticator for a different identity", func(t *testing.T) {
		ctx.testContextChanged(t)

		//used to receive a cert that can be tested with
		unusedIdentityId, unusedCertAuth := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)
		ctx.AdminSession.deleteEntityOfType("identity", unusedIdentityId)

		identityId, _ := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)

		result, err := ctx.AdminSession.newAuthenticatedRequest().Get(fmt.Sprintf(`/authenticators?filter=identity="%s"`, identityId))
		ctx.req.NoError(err)

		resultBody, err := gabs.ParseJSON(result.Body())
		ctx.req.NoError(err)

		idContainer := resultBody.Path("data").Index(0).Path("id")
		ctx.req.NotEmpty(idContainer)

		authenticatorId := idContainer.Data().(string)
		ctx.req.NotEmpty(authenticatorId)

		body := gabs.New()
		_, _ = body.Set(unusedCertAuth.certPem, "certPem")

		resp, err := ctx.AdminSession.newAuthenticatedJsonRequest(body.String()).Patch("/authenticators/" + authenticatorId)

		ctx.req.NoError(err)
		standardJsonResponseTests(resp, http.StatusOK, t)

		t.Run("newly updated cert can be used for authentication", func(t *testing.T) {
			ctx.testContextChanged(t)

			session, err := unusedCertAuth.Authenticate(ctx)
			ctx.req.NoError(err)
			ctx.req.NotEmpty(session)
		})
	})

	t.Run("can delete cert authenticator for a different identity", func(t *testing.T) {
		ctx.testContextChanged(t)
		identityId, authenticator := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)

		result, err := ctx.AdminSession.newAuthenticatedRequest().Get(fmt.Sprintf(`/authenticators?filter=identity="%s"`, identityId))
		ctx.req.NoError(err)

		resultBody, err := gabs.ParseJSON(result.Body())
		ctx.req.NoError(err)

		idContainer := resultBody.Path("data").Index(0).Path("id")
		ctx.req.NotEmpty(idContainer)

		authenticatorId := idContainer.Data().(string)
		ctx.req.NotEmpty(authenticatorId)

		resp, err := ctx.AdminSession.newAuthenticatedRequest().Delete("/authenticators/" + authenticatorId)

		ctx.req.NoError(err)

		standardJsonResponseTests(resp, http.StatusOK, t)

		t.Run("identity can not longer authenticate", func(t *testing.T) {
			ctx.testContextChanged(t)
			session, err := authenticator.Authenticate(ctx)

			ctx.req.Error(err)
			ctx.req.Empty(session)
		})
	})
}

func Test_Authenticators_NonAdminUsingAdminEndpoints(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.teardown()
	ctx.startServer()
	ctx.requireAdminLogin()

	updbNonAdminUserId, updbNonAdminUserAuthenticator := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
	updbNonAdminSession, err := updbNonAdminUserAuthenticator.Authenticate(ctx)

	if err != nil {
		ctx.req.NoError(err, "expected no error during non-admin updb authentication")
	}

	certNonAdminUserId, certNonAdminUserAuthenticator := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)
	certNonAdminUserSession, err := certNonAdminUserAuthenticator.Authenticate(ctx)

	if err != nil {
		ctx.req.NoError(err, "expected no error during non-admin cert authentication")
	}

	t.Run("cannot list authenticators, receives unauthorized error", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().Get("/authenticators")
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})
	t.Run("cannot show details of an authenticator, receives unauthorized error", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().Get("/authenticators/ba3d3a94-47aa-4be1-bc89-04877d5d5fe4")
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})

	t.Run("cannot create updb authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().
			SetHeader("content-type", "application/json").
			SetBody(map[string]interface{}{
				"identityId": certNonAdminUserId,
				"method":     "updb",
				"password":   uuid.New().String(),
				"username":   uuid.New().String(),
			}).
			Post("/authenticators")
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})

	t.Run("cannot update updb authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().
			SetHeader("content-type", "application/json").
			SetBody(map[string]interface{}{
				"password": uuid.New().String(),
				"username": uuid.New().String(),
			}).
			Put("/authenticators/" + uuid.New().String())
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})

	t.Run("cannot patch updb authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().
			SetHeader("content-type", "application/json").
			SetBody(map[string]interface{}{
				"password": uuid.New().String(),
			}).
			Patch("/authenticators/" + uuid.New().String())
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})

	t.Run("cannot delete updb authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().Delete("/authenticators/" + uuid.New().String())
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})

	t.Run("cannot create cert authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := certNonAdminUserSession.newAuthenticatedRequest().
			SetHeader("content-type", "application/json").
			SetBody(map[string]interface{}{
				"identityId": updbNonAdminUserId,
				"method":     "updb",
				"password":   uuid.New().String(),
				"username":   uuid.New().String(),
			}).
			Post("/authenticators")
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})

	t.Run("cannot update cert authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().
			SetHeader("content-type", "application/json").
			SetBody(map[string]interface{}{
				"certPem": "",
			}).
			Put("/authenticators/" + uuid.New().String())
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})
	t.Run("cannot patch cert authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().
			SetHeader("content-type", "application/json").
			SetBody(map[string]interface{}{
				"certPem": "",
			}).
			Patch("/authenticators/" + uuid.New().String())
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})
	t.Run("cannot delete cert authenticator for a different identity, receives unauthorized", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().Delete("/authenticators/" + uuid.New().String())
		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)
	})
}

func Test_Authenticators_NonAdminUsingSelfServiceEndpoints(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.teardown()
	ctx.startServer()
	ctx.requireAdminLogin()

	_, updbNonAdminUserAuthenticator := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
	updbNonAdminSession, err := updbNonAdminUserAuthenticator.Authenticate(ctx)

	if err != nil {
		ctx.req.NoError(err, "expected no error during non-admin updb authentication")
	}

	_, certNonAdminUserAuthenticator := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)
	certNonAdminUserSession, err := certNonAdminUserAuthenticator.Authenticate(ctx)

	if err != nil {
		ctx.req.NoError(err, "expected no error during non-admin cert authentication")
	}

	t.Run("can access their authenticators", func(t *testing.T) {
		req := require.New(t)
		resp, err := updbNonAdminSession.newAuthenticatedRequest().Get("/current-identity/authenticators")

		req.NoError(err)

		standardJsonResponseTests(resp, http.StatusOK, t)

		updbAuthenticatorListBody, err := gabs.ParseJSON(resp.Body())

		req.NoError(err)

		t.Run("has one authenticator", func(t *testing.T) {
			req := require.New(t)
			array, ok := updbAuthenticatorListBody.Path("data").Data().([]interface{})

			req.True(ok, "could not cast data to array")
			req.Equal(1, len(array), "number of authenticators returned expected to be 1 got %d", len(array))
		})

		t.Run("authenticator returned is updb", func(t *testing.T) {
			req := require.New(t)
			method, ok := updbAuthenticatorListBody.Search("data").Index(0).Path("method").Data().(string)

			req.True(ok)
			req.Equal("updb", method)
		})

		t.Run("authenticator returned is for updb user", func(t *testing.T) {
			req := require.New(t)
			id, ok := updbAuthenticatorListBody.Search("data").Index(0).Path("identityId").Data().(string)

			req.True(ok)
			req.Equal(updbNonAdminSession.identityId, id)
		})

		t.Run("can get the detail of the authenticator", func(t *testing.T) {
			req := require.New(t)

			authenticatorId, ok := updbAuthenticatorListBody.Search("data").Index(0).Path("id").Data().(string)

			req.True(ok)

			resp, err := updbNonAdminSession.newAuthenticatedRequest().Get("/current-identity/authenticators/" + authenticatorId)

			req.NoError(err)

			standardJsonResponseTests(resp, http.StatusOK, t)
		})
	})

	t.Run("can not access authenticators for other identities", func(t *testing.T) {
		req := require.New(t)
		//get updb's authenticator id
		updbAuthenticatorResp, err := updbNonAdminSession.newAuthenticatedRequest().Get("/current-identity/authenticators")
		req.NoError(err)

		updbAuthenticatorListBody, err := gabs.ParseJSON(updbAuthenticatorResp.Body())
		req.NoError(err)

		authenticatorId, ok := updbAuthenticatorListBody.Search("data").Index(0).Path("id").Data().(string)
		req.True(ok)
		req.NotEmpty(authenticatorId)
		_, err = uuid.Parse(authenticatorId)
		req.NoError(err)

		t.Run("for read if the authenticator id is made up", func(t *testing.T) {
			req := require.New(t)

			resp, err := updbNonAdminSession.newAuthenticatedRequest().Get("current-identity/authenticators/" + uuid.New().String())

			req.NoError(err)

			standardErrorJsonResponseTests(resp, apierror.NotFoundCode, http.StatusNotFound, t)
		})

		t.Run("for read if the authenticator id is for another identity", func(t *testing.T) {
			req := require.New(t)
			//access updb's authenticator from cert identity
			resp, err := certNonAdminUserSession.newAuthenticatedRequest().Get("/current-identity/authenticators/" + authenticatorId)

			req.NoError(err)

			standardErrorJsonResponseTests(resp, apierror.NotFoundCode, http.StatusNotFound, t)
		})

		t.Run("for update if the authenticator id is for another identity", func(t *testing.T) {
			//access updb's authenticator from cert identity
			resp, err := certNonAdminUserSession.newAuthenticatedJsonRequest(`{"currentPassword": "123456", "newPassword":"456789", "username":"username123456"}`).
				Put("/current-identity/authenticators/" + authenticatorId)

			req.NoError(err)

			standardErrorJsonResponseTests(resp, apierror.NotFoundCode, http.StatusNotFound, t)
		})

		t.Run("for update if the authenticator id is made up", func(t *testing.T) {
			//access updb's authenticator from cert identity
			resp, err := certNonAdminUserSession.newAuthenticatedJsonRequest(`{"currentPassword": "123456", "newPassword":"456789", "username":"username123456"}`).
				Put("/current-identity/authenticators/" + uuid.New().String())

			req.NoError(err)

			standardErrorJsonResponseTests(resp, apierror.NotFoundCode, http.StatusNotFound, t)
		})

		t.Run("for patch if the authenticator id is for another identity", func(t *testing.T) {
			//access updb's authenticator from cert identity
			resp, err := certNonAdminUserSession.newAuthenticatedJsonRequest(`{"currentPassword": "123456", "newPassword":"456789"}`).
				Patch("/current-identity/authenticators/" + authenticatorId)

			req.NoError(err)

			standardErrorJsonResponseTests(resp, apierror.NotFoundCode, http.StatusNotFound, t)
		})

		t.Run("for patch if the authenticator id is made up", func(t *testing.T) {
			//access updb's authenticator from cert identity
			resp, err := certNonAdminUserSession.newAuthenticatedJsonRequest(`{"currentPassword": "123456", "newPassword":"456789"}`).
				Patch("/current-identity/authenticators/" + uuid.New().String())

			req.NoError(err)

			standardErrorJsonResponseTests(resp, apierror.NotFoundCode, http.StatusNotFound, t)
		})
	})

	t.Run("can not delete as it isn't supported", func(t *testing.T) {
		req := require.New(t)
		resp, err := certNonAdminUserSession.newAuthenticatedRequest().Delete("/current-identity/authenticators/" + uuid.New().String())

		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.MethodNotAllowedCode, http.StatusMethodNotAllowed, t)
	})

	t.Run("can not create as it isn't supported", func(t *testing.T) {
		req := require.New(t)
		resp, err := certNonAdminUserSession.newAuthenticatedJsonRequest("{}").Post("/current-identity/authenticators/")

		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.MethodNotAllowedCode, http.StatusMethodNotAllowed, t)
	})

	t.Run("a non-admin can update their own updb authenticator", func(t *testing.T) {
		req := require.New(t)
		ctx.testContextChanged(t)

		_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
		authSession, err := auth.Authenticate(ctx)
		req.NoError(err)

		authResp, err := authSession.newAuthenticatedRequest().Get("/current-identity/authenticators")
		req.NoError(err)

		authListBody, err := gabs.ParseJSON(authResp.Body())
		req.NoError(err)

		authId, ok := authListBody.Search("data").Index(0).Path("id").Data().(string)
		req.True(ok)
		req.NotEmpty(authId)
		_, err = uuid.Parse(authId)
		req.NoError(err)

		newUsername := uuid.New().String()
		newPassword := uuid.New().String()

		body := fmt.Sprintf(`{"username":"%s", "newPassword":"%s", "currentPassword":"%s"}`, newUsername, newPassword, auth.Password)
		resp, err := authSession.newAuthenticatedJsonRequest(body).
			Put("/current-identity/authenticators/" + authId)

		req.NoError(err)

		standardJsonResponseTests(resp, http.StatusOK, t)

		t.Run("a non-admin can authenticate with updated updb credentials", func(t *testing.T) {
			ctx.testContextChanged(t)

			auth.Username = newUsername
			auth.Password = newPassword

			_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
			_, err := auth.Authenticate(ctx)

			req.NoError(err)
		})
	})

	t.Run("a non-admin can not update their own updb authenticator with an invalid current password", func(t *testing.T) {
		req := require.New(t)
		ctx.testContextChanged(t)

		_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
		authSession, err := auth.Authenticate(ctx)
		req.NoError(err)

		authResp, err := authSession.newAuthenticatedRequest().Get("/current-identity/authenticators")
		req.NoError(err)

		authListBody, err := gabs.ParseJSON(authResp.Body())
		req.NoError(err)

		authId, ok := authListBody.Search("data").Index(0).Path("id").Data().(string)
		req.True(ok)
		req.NotEmpty(authId)
		_, err = uuid.Parse(authId)
		req.NoError(err)

		newUsername := uuid.New().String()
		newPassword := uuid.New().String()

		body := fmt.Sprintf(`{"username":"%s", "newPassword":"%s", "currentPassword":"%s"}`, newUsername, newPassword, uuid.New().String())
		resp, err := authSession.newAuthenticatedJsonRequest(body).
			Put("/current-identity/authenticators/" + authId)

		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)

		t.Run("a non-admin can authenticate with the original updb credentials", func(t *testing.T) {
			ctx.testContextChanged(t)

			_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
			_, err := auth.Authenticate(ctx)

			req.NoError(err)
		})
	})

	t.Run("a non-admin can patch their own updb authenticator", func(t *testing.T) {
		req := require.New(t)
		ctx.testContextChanged(t)

		_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
		authSession, err := auth.Authenticate(ctx)
		req.NoError(err)

		authResp, err := authSession.newAuthenticatedRequest().Get("/current-identity/authenticators")
		req.NoError(err)

		authListBody, err := gabs.ParseJSON(authResp.Body())
		req.NoError(err)

		authId, ok := authListBody.Search("data").Index(0).Path("id").Data().(string)
		req.True(ok)
		req.NotEmpty(authId)
		_, err = uuid.Parse(authId)
		req.NoError(err)

		newPassword := uuid.New().String()

		body := fmt.Sprintf(`{"newPassword":"%s", "currentPassword":"%s"}`, newPassword, auth.Password)
		resp, err := authSession.newAuthenticatedJsonRequest(body).
			Patch("/current-identity/authenticators/" + authId)

		req.NoError(err)

		standardJsonResponseTests(resp, http.StatusOK, t)

		t.Run("a non-admin can authenticate with patched updb credentials", func(t *testing.T) {
			ctx.testContextChanged(t)

			auth.Password = newPassword

			_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
			_, err := auth.Authenticate(ctx)

			req.NoError(err)
		})
	})

	t.Run("a non-admin can not patch their own updb authenticator with an invalid current password", func(t *testing.T) {
		req := require.New(t)
		ctx.testContextChanged(t)

		_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
		authSession, err := auth.Authenticate(ctx)
		req.NoError(err)

		authResp, err := authSession.newAuthenticatedRequest().Get("/current-identity/authenticators")
		req.NoError(err)

		authListBody, err := gabs.ParseJSON(authResp.Body())
		req.NoError(err)

		authId, ok := authListBody.Search("data").Index(0).Path("id").Data().(string)
		req.True(ok)
		req.NotEmpty(authId)
		_, err = uuid.Parse(authId)
		req.NoError(err)

		newPassword := uuid.New().String()

		body := fmt.Sprintf(`{"newPassword":"%s", "currentPassword":"%s"}`, newPassword, uuid.New().String())
		resp, err := authSession.newAuthenticatedJsonRequest(body).
			Patch("/current-identity/authenticators/" + authId)

		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.UnauthorizedCode, http.StatusUnauthorized, t)

		t.Run("a non-admin can authenticate with original updb credentials", func(t *testing.T) {
			ctx.testContextChanged(t)

			_, auth := ctx.AdminSession.requireCreateIdentityWithUpdbEnrollment(uuid.New().String(), uuid.New().String(), false)
			_, err := auth.Authenticate(ctx)

			req.NoError(err)
		})
	})

	t.Run("a non-admin cannot update their own cert authenticator", func(t *testing.T) {
		req := require.New(t)
		ctx.testContextChanged(t)

		_, auth := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)
		authSession, err := auth.Authenticate(ctx)
		req.NoError(err)

		authResp, err := authSession.newAuthenticatedRequest().Get("/current-identity/authenticators")
		req.NoError(err)

		authListBody, err := gabs.ParseJSON(authResp.Body())
		req.NoError(err)

		authId, ok := authListBody.Search("data").Index(0).Path("id").Data().(string)
		req.True(ok)
		req.NotEmpty(authId)
		_, err = uuid.Parse(authId)
		req.NoError(err)

		resp, err := authSession.newAuthenticatedJsonRequest(`{"certPem": "something"}`).
			Put("/current-identity/authenticators/" + authId)

		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.AuthenticatorCanNotBeUpdatedCode, http.StatusConflict, t)
	})

	t.Run("a non-admin cannot patch their own cert authenticator", func(t *testing.T) {
		req := require.New(t)
		ctx.testContextChanged(t)

		_, auth := ctx.AdminSession.requireCreateIdentityOttEnrollment(uuid.New().String(), false)
		authSession, err := auth.Authenticate(ctx)
		req.NoError(err)

		authResp, err := authSession.newAuthenticatedRequest().Get("/current-identity/authenticators")
		req.NoError(err)

		authListBody, err := gabs.ParseJSON(authResp.Body())
		req.NoError(err)

		authId, ok := authListBody.Search("data").Index(0).Path("id").Data().(string)
		req.True(ok)
		req.NotEmpty(authId)
		_, err = uuid.Parse(authId)
		req.NoError(err)

		resp, err := authSession.newAuthenticatedJsonRequest(`{ "certPem": "something"}`).
			Patch("/current-identity/authenticators/" + authId)

		req.NoError(err)

		standardErrorJsonResponseTests(resp, apierror.AuthenticatorCanNotBeUpdatedCode, http.StatusConflict, t)
	})
}
