package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/MehmetTalhaSeker/mts-blog-api/e2e"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/apputils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/testutils"
)

var _ = Describe("user", Ordered, func() {
	ctx := context.Background()

	users := e2e.CreateUserModels(27)

	user := e2e.CreateUserModel(45, types.Registered)
	anotherUser := e2e.CreateUserModel(60, types.Registered)
	modUser := e2e.CreateUserModel(57, types.Mod)
	adminUser := e2e.CreateUserModel(58, types.Admin)

	users = append(users, user, modUser, adminUser)

	var userDTOs []dto.UserResponse
	for _, u := range users {
		userDTOs = append(userDTOs, *u.ToDTO())
	}

	BeforeEach(func() {
		testutils.InsertUsers(apputils.ToSliceOfAny(users), store.GetInstance())
	})

	AfterEach(func() {
		testutils.DeleteUsers(store.GetInstance())
	})

	Context("read", func() {
		testCases := []struct {
			when     string
			it       string
			id       uint64
			authUser *model.User
			want     *dto.UserResponse
			wantCode int
			wantErr  *errorutils.APIError
			wantErrs *errorutils.APIErrors
		}{
			{
				when:     "admin read user data",
				it:       "should success",
				authUser: adminUser,
				id:       user.ID,
				want:     user.ToDTO(),
				wantCode: http.StatusOK,
			},
			{
				when:     "mod read user data",
				it:       "should success",
				authUser: modUser,
				id:       user.ID,
				want:     user.ToDTO(),
				wantCode: http.StatusOK,
			},
			{
				when:     "read other user's data",
				it:       "should fail",
				authUser: user,
				id:       anotherUser.ID,
				wantCode: http.StatusUnauthorized,
				wantErr:  errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:     "user data with non exist id",
				it:       "should fail",
				authUser: adminUser,
				id:       20000,
				wantCode: http.StatusNotFound,
				wantErr:  errorutils.New(errorutils.ErrUserNotFound, nil),
			},
		}

		for _, tc := range testCases {
			tc := tc
			When(tc.when, func() {
				AfterEach(func() {
					e2e.ClearAuthMidUser(e)
				})
				It(tc.it, func() {
					if tc.authUser != nil {
						e2e.AuthMidUser(e, tc.authUser)
					}

					code, body, _, err := e2e.Get(ctx, "/users/"+strconv.FormatUint(tc.id, 10))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						got := new(dto.UserResponse)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreTypes(time.Time{}), cmpopts.IgnoreFields(dto.UserResponse{}, "CreatedAt", "UpdatedAt", "DeletedAt", "Status", "UpdatedBy", "CreatedBy")); diff != "" {
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

	Context("reads", func() {
		testCases := []struct {
			when        string
			it          string
			path        string
			authUser    *model.User
			want        []dto.UserResponse
			wantCode    int
			wantErr     *errorutils.APIError
			wantHeaders map[string]string
		}{
			{
				when:     "sort ascendant",
				it:       "return first 20 users",
				authUser: adminUser,
				path:     "?sort=createdAt,asc",
				want:     userDTOs[0:20],
				wantCode: http.StatusOK,
				wantHeaders: map[string]string{
					pagination.HeaderLink:        `</v1/users?page=1&size=20&sort=createdAt%2Casc>; rel="first" </v1/users?page=2&size=20&sort=createdAt%2Casc>; rel="next", </v1/users?page=2&size=20&sort=createdAt%2Casc>; rel="last"`,
					pagination.HeaderXHasNext:    "true",
					pagination.HeaderXTotalCount: "30",
					pagination.HeaderXTotalPage:  "2",
				},
			},
			{
				when:     "page 4 size 5 sort asc",
				it:       "return third 5 users",
				authUser: adminUser,
				path:     "?sort=createdAt,asc&page=4&size=5",
				want:     userDTOs[15:20],
				wantCode: http.StatusOK,
				wantHeaders: map[string]string{
					pagination.HeaderLink:        `</v1/users?page=1&size=5&sort=createdAt%2Casc>; rel="first" </v1/users?page=5&size=5&sort=createdAt%2Casc>; rel="next", </v1/users?page=3&size=5&sort=createdAt%2Casc>; rel="prev", </v1/users?page=6&size=5&sort=createdAt%2Casc>; rel="last"`,
					pagination.HeaderXHasNext:    "true",
					pagination.HeaderXTotalCount: "30",
					pagination.HeaderXTotalPage:  "6",
				},
			},
		}

		for _, tc := range testCases {
			tc := tc
			When(tc.when, func() {
				AfterEach(func() {
					e2e.ClearAuthMidUser(e)
				})
				It(tc.it, func() {
					if tc.authUser != nil {
						e2e.AuthMidUser(e, tc.authUser)
					}

					code, body, header, err := e2e.Get(ctx, "/users"+tc.path)
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						var got []dto.UserResponse
						err = json.Unmarshal(body, &got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreTypes(time.Time{}), cmpopts.IgnoreFields(dto.UserResponse{}, "CreatedAt", "UpdatedAt", "DeletedAt", "Status", "UpdatedBy", "CreatedBy")); diff != "" {
							Expect(diff).To(BeEmpty())
						}

						for key, value := range tc.wantHeaders {
							Expect(header.Get(key)).To(Equal(value))
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
				})
			})
		}
	})

	Context("update", func() {
		testCases := []struct {
			when       string
			it         string
			id         uint64
			updateJSON string
			authUser   *model.User
			want       *dto.UserResponse
			wantCode   int
			wantErr    *errorutils.APIError
			wantErrs   *errorutils.APIErrors
		}{
			{
				when:       "user update own data",
				it:         "should success",
				authUser:   user,
				id:         user.ID,
				updateJSON: `{ "username": "samil" }`,
				want: &dto.UserResponse{
					ID:       user.ID,
					Username: "samil",
					Role:     types.Registered,
					Status:   types.Active,
				},
				wantCode: http.StatusOK,
			},
			{
				when:       "user tries to update other user's data",
				it:         "should fail",
				authUser:   user,
				id:         adminUser.ID,
				updateJSON: `{ "username": "samil"}`,
				wantCode:   http.StatusUnauthorized,
				wantErr:    errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:       "mod update other user's data",
				it:         "should fail",
				wantErr:    errorutils.New(errorutils.ErrUnauthorized, nil),
				authUser:   modUser,
				id:         user.ID,
				updateJSON: `{ "username": "samil-mod" }`,
				wantCode:   http.StatusUnauthorized,
			},
			{
				when:       "admin update other user's data",
				it:         "should success",
				authUser:   adminUser,
				id:         user.ID,
				updateJSON: `{ "username": "samil-admin" }`,
				want: &dto.UserResponse{
					ID:       user.ID,
					Username: "samil-admin",
					Role:     types.Registered,
				},
				wantCode: http.StatusOK,
			},
			{
				when:       "empty fields",
				it:         "should fail",
				authUser:   user,
				id:         user.ID,
				updateJSON: `{}`,
				wantCode:   http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.ErrUsernameRequired, nil),
				}},
			},
			{
				when:       "username too short",
				it:         "should fail",
				authUser:   user,
				id:         user.ID,
				updateJSON: `{ "username": "k"}`,
				wantCode:   http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.ErrShortUsername, nil),
				}},
			},
			{
				when:       "username too long",
				it:         "should fail",
				authUser:   user,
				id:         user.ID,
				updateJSON: `{ "username": "samilsamilsamilsamilsamilsamilsamil"}`,
				wantCode:   http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.ErrLongUsername, nil),
				}},
			},
			{
				when:       "existing username",
				it:         "should fail",
				authUser:   user,
				id:         user.ID,
				updateJSON: `{ "username": "username-10" }`,
				wantCode:   http.StatusBadRequest,
				wantErr:    errorutils.New(errorutils.ErrUsernameAlreadyTaken, nil),
			},
			{
				when:       "invalid json",
				it:         "should fail",
				authUser:   user,
				id:         user.ID,
				updateJSON: `{ "username": "samil }`,
				wantCode:   http.StatusBadRequest,
				wantErr:    errorutils.New(errorutils.ErrBinding, nil),
			},
		}

		for _, tc := range testCases {
			tc := tc
			When(tc.when, func() {
				AfterEach(func() {
					e2e.ClearAuthMidUser(e)
				})
				It(tc.it, func() {
					if tc.authUser != nil {
						e2e.AuthMidUser(e, tc.authUser)
					}

					code, body, _, err := e2e.Put(ctx, "/users/"+strconv.FormatUint(tc.id, 10), []byte(tc.updateJSON))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						got := new(dto.UserResponse)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(dto.UserResponse{}, "CreatedAt", "UpdatedAt", "DeletedAt", "Email", "Status")); diff != "" {
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

	Context("delete", func() {
		testCases := []struct {
			when     string
			it       string
			id       uint64
			authUser *model.User
			want     *dto.ResponseWithID
			wantCode int
			wantErrs *errorutils.APIErrors
			wantErr  *errorutils.APIError
		}{
			{
				when:     "user deletes user",
				it:       "should fail",
				id:       userDTOs[0].ID,
				authUser: user,
				wantCode: http.StatusUnauthorized,
				wantErr:  errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:     "mod deletes user",
				it:       "should success",
				id:       userDTOs[0].ID,
				authUser: modUser,
				wantCode: http.StatusUnauthorized,
				wantErr:  errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:     "admin deletes user",
				it:       "should success",
				id:       userDTOs[0].ID,
				authUser: adminUser,
				want:     &dto.ResponseWithID{ID: strconv.FormatUint(userDTOs[0].ID, 10)},
				wantCode: http.StatusOK,
			},
			{
				when:     "valid id but no data",
				it:       "should fail",
				id:       20000,
				authUser: adminUser,
				wantCode: http.StatusNotFound,
				wantErr:  errorutils.New(errorutils.ErrUserNotFound, nil),
			},
		}

		for _, tc := range testCases {
			tc := tc
			When(tc.when, func() {
				AfterEach(func() {
					e2e.ClearAuthMidUser(e)
				})
				It(tc.it, func() {
					if tc.authUser != nil {
						e2e.AuthMidUser(e, tc.authUser)
					}

					code, body, _, err := e2e.Delete(ctx, "/users/"+strconv.FormatUint(tc.id, 10))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						got := new(dto.ResponseWithID)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got); diff != "" {
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
