package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/MehmetTalhaSeker/mts-blog-api/e2e"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/testutils"
)

var _ = Describe("authentication", Ordered, func() {
	ctx := context.Background()

	user := e2e.CreateUserModel(45, types.Registered)

	var users []*model.User
	users = append(users, user)

	BeforeEach(func() {
		testutils.InsertUsers(apputils.ToSliceOfAny(users), store.GetInstance())
	})

	AfterEach(func() {
		testutils.DeleteUsers(store.GetInstance())
	})

	Context("register", func() {
		testCases := []struct {
			when     string
			it       string
			json     string
			want     *dto.WithTokenResponse
			wantCode int
			wantErrs *errorutils.APIErrors
			wantErr  *errorutils.APIError
		}{
			{
				when:     "valid data",
				it:       "should succeed",
				json:     `{ "username": "samil", "email": "samil@samilov.com", "termsOfService": true, "password": "12341234" }`,
				wantCode: http.StatusCreated,
			},
			{
				when:     "all fields empty",
				it:       "should fail",
				json:     `{}`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Email"), nil),
					errorutils.New(errorutils.Required("Password"), nil),
					errorutils.New(errorutils.Required("Username"), nil),
					errorutils.New(errorutils.Required("TermsOfService"), nil),
				}},
			},
			{
				when:     "invalid json",
				it:       "should fail",
				json:     `{ "username": "samil", "email": "samil@samilov.com", "termsOfService": true, "password": "12341234", }`,
				wantCode: http.StatusBadRequest,
				wantErr:  errorutils.New(errorutils.ErrBinding, nil),
			},
			{
				when:     "termsOfService is false",
				it:       "should fail",
				json:     `{ "username": "samil", "email": "samil@samilov.com", "termsOfService": false, "password": "12341234" }`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("TermsOfService"), nil),
				}},
			},
			{
				when:     "empty username",
				it:       "should fail",
				json:     `{ "email": "samil@samilov.com", "termsOfService": true, "password": "12341234" }`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Username"), nil),
				}},
			},
			{
				when:     "username short",
				it:       "should fail",
				json:     `{ "username": "ka", "email": "samil@samilov.com", "termsOfService": true, "password": "12341234" }`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Min("Username"), nil),
				}},
			},
			{
				when:     "username long",
				it:       "should fail",
				json:     `{ "username": "samilsamilsamilsamilsamilsamilsamilsamilsamilsamilsamilsamil", "email": "samil@samilov.com", "termsOfService": true, "password": "12341234" }`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Max("Username"), nil),
				}},
			},
			{
				when:     "empty email",
				it:       "should fail",
				json:     `{ "username": "samil", "termsOfService": true, "password": "12341234" }`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Email"), nil),
				}},
			},
			{
				when:     "same username",
				it:       "should fail",
				json:     fmt.Sprintf(`{ "username": "%s", "email": "samil@samilov.com", "termsOfService": true, "password": "12341234" }`, user.Username),
				wantCode: http.StatusBadRequest,
				wantErr:  errorutils.New(errorutils.ErrUsernameAlreadyTaken, nil),
			},
			{
				when:     "same email",
				it:       "should fail",
				json:     fmt.Sprintf(`{ "username": "samil", "email": "%s", "termsOfService": true, "password": "12341234" }`, user.Email),
				wantCode: http.StatusBadRequest,
				wantErr:  errorutils.New(errorutils.ErrEmailAlreadyTaken, nil),
			},
		}

		for _, tc := range testCases {
			tc := tc
			When(tc.when, func() {
				It(tc.it, func() {
					code, body, _, err := e2e.Post(ctx, "/auth/register", []byte(tc.json))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						got := new(dto.WithTokenResponse)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreTypes(time.Time{})); diff != "" {
							Expect(diff).To(BeEmpty())
						}
					}

					if tc.wantErr != nil {
						got := new(errorutils.APIError)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.wantErr, got); diff != "" {
							Expect(diff).To(BeEmpty())
						}
					}

					if tc.wantErrs != nil {
						got := new(errorutils.APIErrors)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.wantErrs, got); diff != "" {
							Expect(diff).To(BeEmpty())
						}
					}
				})
			})
		}
	})

	Context("login", func() {
		testCases := []struct {
			when     string
			it       string
			json     string
			want     *dto.WithTokenResponse
			wantCode int
			wantErrs *errorutils.APIErrors
			wantErr  *errorutils.APIError
		}{
			{
				when:     "valid data",
				it:       "should succeed",
				json:     fmt.Sprintf(`{ "email": "%s", "password": "12341234" }`, user.Email),
				wantCode: http.StatusOK,
			},
			{
				when:     "all fields empty",
				it:       "should fail",
				json:     `{}`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Email"), nil),
					errorutils.New(errorutils.Required("Password"), nil),
				}},
			},
			{
				when:     "invalid json",
				it:       "should fail",
				json:     `{ "email": "samil@samilov.com", "password": "12341234", }`,
				wantCode: http.StatusBadRequest,
				wantErr:  errorutils.New(errorutils.ErrBinding, nil),
			},
			{
				when:     "empty email",
				it:       "should fail",
				json:     `{ "password": "12341234" }`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Email"), nil),
				}},
			},
			{
				when:     "empty password",
				it:       "should fail",
				json:     `{ "email": "samil@samilov.com" }`,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Password"), nil),
				}},
			},
			{
				when:     "non existing user",
				it:       "should fail",
				json:     `{ "email": "samil@samilov.com", "password": "12341234" }`,
				wantCode: http.StatusBadRequest,
				wantErr:  errorutils.New(errorutils.ErrEmailNotFound, nil),
			},
		}

		for _, tc := range testCases {
			tc := tc
			When(tc.when, func() {
				It(tc.it, func() {
					code, body, _, err := e2e.Post(ctx, "/auth/login", []byte(tc.json))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						got := new(dto.UserResponse)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreTypes(time.Time{})); diff != "" {
							Expect(diff).To(BeEmpty())
						}
					}

					if tc.wantErr != nil {
						got := new(errorutils.APIError)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.wantErr, got); diff != "" {
							Expect(diff).To(BeEmpty())
						}
					}

					if tc.wantErrs != nil {
						got := new(errorutils.APIErrors)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.wantErrs, got); diff != "" {
							Expect(diff).To(BeEmpty())
						}
					}
				})
			})
		}
	})
})
