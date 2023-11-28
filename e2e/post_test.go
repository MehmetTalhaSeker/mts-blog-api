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

var _ = Describe("post", Ordered, func() {
	ctx := context.Background()

	users := make([]*model.User, 0, 3)
	user := e2e.CreateUserModel(45, types.Registered)
	modUser := e2e.CreateUserModel(57, types.Mod)
	adminUser := e2e.CreateUserModel(58, types.Admin)

	users = append(users, user, modUser, adminUser)

	posts := e2e.CreatePostModels(30)

	var postDTOs []dto.PostResponse
	for _, p := range posts {
		postDTOs = append(postDTOs, *p.ToDTO())
	}

	BeforeAll(func() {
		testutils.InsertUsers(apputils.ToSliceOfAny(users), store.GetInstance())
	})

	AfterAll(func() {
		testutils.DeleteUsers(store.GetInstance())
	})

	BeforeEach(func() {
		testutils.InsertPosts(apputils.ToSliceOfAny(posts), store.GetInstance())
	})

	AfterEach(func() {
		testutils.DeletePosts(store.GetInstance())
	})

	Context("create", func() {
		testCases := []struct {
			when     string
			it       string
			json     string
			authUser *model.User
			want     string
			wantCode int
			wantErrs *errorutils.APIErrors
			wantErr  *errorutils.APIError
		}{
			{
				when:     "valid data (admin)",
				it:       "should succeed",
				json:     `{ "title": "SHOW", "body":"BODY-BODY" }`,
				authUser: adminUser,
				want:     "OK",
				wantCode: http.StatusCreated,
			},
			{
				when:     "valid data (mod)",
				it:       "should succeed",
				json:     `{ "title": "SHOW", "body":"BODY-BODY" }`,
				authUser: modUser,
				want:     "OK",
				wantCode: http.StatusCreated,
			},
			{
				when:     "invalid json",
				it:       "should fail",
				json:     `{ "title: "samil-show" }`,
				authUser: adminUser,
				wantCode: http.StatusBadRequest,
				wantErr:  errorutils.New(errorutils.ErrBinding, nil),
			},
			{
				when:     "title short",
				it:       "should fail",
				json:     `{ "title": "k", "body":"bodyBody" }`,
				authUser: adminUser,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Min("Title"), nil),
				}},
			},
			{
				when:     "title long",
				it:       "should fail",
				json:     `{ "title": "samil-samil-samil-samil-samil-samil-samil-samil-samil-samil", "body":"bodyBody" }`,
				authUser: adminUser,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Max("Title"), nil),
				}},
			},
			{
				when:     "title short and body missing",
				it:       "should fail",
				json:     `{ "title": "k" }`,
				authUser: adminUser,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Min("Title"), nil),
					errorutils.New(errorutils.Required("Body"), nil),
				}},
			},
			{
				when:     "all fields empty",
				it:       "should fail",
				json:     `{}`,
				authUser: adminUser,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Title"), nil),
					errorutils.New(errorutils.Required("Body"), nil),
				}},
			},
			{
				when:     "title long, valid body",
				it:       "should fail",
				json:     `{ "title":"samil-samil-samil-samil-samil-samil-samil-samil-samil-samil-samil-samil-samil-samil-samil-samil-samil" ,"body": "samil" }`,
				authUser: adminUser,
				wantCode: http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Max("Title"), nil),
				}},
			},
			{
				when:     "registered user",
				it:       "should fail",
				authUser: user,
				json:     `{ "title": "samilShow", "body":"BODY_body" }`,
				wantCode: http.StatusUnauthorized,
				wantErr:  errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:     "unauthenticated user",
				it:       "should fail",
				json:     `{ "title": "samil-show", "body":"BODY_body" }`,
				wantCode: http.StatusUnauthorized,
				wantErr:  errorutils.New(errorutils.ErrUnauthorized, nil),
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

					code, body, _, err := e2e.Post(ctx, "/posts", []byte(tc.json))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != "" {
						got := new(string)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, *got); diff != "" {
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

	Context("read", func() {
		testCases := []struct {
			when     string
			it       string
			id       uint64
			authUser *model.User
			want     *dto.PostResponse
			wantCode int
			wantErr  *errorutils.APIError
			wantErrs *errorutils.APIErrors
		}{
			{
				when:     "admin read post",
				it:       "should success",
				authUser: adminUser,
				id:       posts[0].ID,
				want:     posts[0].ToDTO(),
				wantCode: http.StatusOK,
			},
			{
				when:     "mod read post",
				it:       "should success",
				authUser: modUser,
				id:       posts[0].ID,
				want:     posts[0].ToDTO(),
				wantCode: http.StatusOK,
			},
			{
				when:     "read post",
				it:       "should fail",
				authUser: user,
				id:       posts[0].ID,
				want:     posts[0].ToDTO(),
				wantCode: http.StatusOK,
			},
			{
				when:     "with non existing id",
				it:       "should fail",
				authUser: adminUser,
				id:       20000,
				wantCode: http.StatusNotFound,
				wantErr:  errorutils.New(errorutils.ErrPostNotFound, nil),
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

					code, body, _, err := e2e.Get(ctx, "/posts/"+strconv.FormatUint(tc.id, 10))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						got := new(dto.PostResponse)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreTypes(time.Time{}), cmpopts.IgnoreFields(dto.PostResponse{}, "CreatedAt", "UpdatedAt", "DeletedAt", "UpdatedBy", "CreatedBy")); diff != "" {
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
			want        []dto.PostResponse
			wantCode    int
			wantErr     *errorutils.APIError
			wantHeaders map[string]string
		}{
			{
				when:     "sort ascendant",
				it:       "return first 20",
				authUser: adminUser,
				path:     "?sort=createdAt,asc",
				want:     postDTOs[0:20],
				wantCode: http.StatusOK,
				wantHeaders: map[string]string{
					pagination.HeaderLink:        `</v1/posts?page=1&size=20&sort=createdAt%2Casc>; rel="first" </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="next", </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="last"`,
					pagination.HeaderXHasNext:    "true",
					pagination.HeaderXTotalCount: "30",
					pagination.HeaderXTotalPage:  "2",
				},
			},
			{
				when:     "page 4 size 5 sort asc",
				it:       "return third 5",
				authUser: adminUser,
				path:     "?sort=createdAt,asc&page=4&size=5",
				want:     postDTOs[15:20],
				wantCode: http.StatusOK,
				wantHeaders: map[string]string{
					pagination.HeaderLink:        `</v1/posts?page=1&size=5&sort=createdAt%2Casc>; rel="first" </v1/posts?page=5&size=5&sort=createdAt%2Casc>; rel="next", </v1/posts?page=3&size=5&sort=createdAt%2Casc>; rel="prev", </v1/posts?page=6&size=5&sort=createdAt%2Casc>; rel="last"`,
					pagination.HeaderXHasNext:    "true",
					pagination.HeaderXTotalCount: "30",
					pagination.HeaderXTotalPage:  "6",
				},
			},
			{
				when:     "by mod",
				it:       "return first 20",
				authUser: modUser,
				path:     "?sort=createdAt,asc",
				want:     postDTOs[0:20],
				wantCode: http.StatusOK,
				wantHeaders: map[string]string{
					pagination.HeaderLink:        `</v1/posts?page=1&size=20&sort=createdAt%2Casc>; rel="first" </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="next", </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="last"`,
					pagination.HeaderXHasNext:    "true",
					pagination.HeaderXTotalCount: "30",
					pagination.HeaderXTotalPage:  "2",
				},
			},
			{
				when:     "by registered",
				it:       "return first 20",
				authUser: modUser,
				path:     "?sort=createdAt,asc",
				want:     postDTOs[0:20],
				wantCode: http.StatusOK,
				wantHeaders: map[string]string{
					pagination.HeaderLink:        `</v1/posts?page=1&size=20&sort=createdAt%2Casc>; rel="first" </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="next", </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="last"`,
					pagination.HeaderXHasNext:    "true",
					pagination.HeaderXTotalCount: "30",
					pagination.HeaderXTotalPage:  "2",
				},
			},
			{
				when:     "by anonymous user",
				it:       "return first 20",
				path:     "?sort=createdAt,asc",
				want:     postDTOs[0:20],
				wantCode: http.StatusOK,
				wantHeaders: map[string]string{
					pagination.HeaderLink:        `</v1/posts?page=1&size=20&sort=createdAt%2Casc>; rel="first" </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="next", </v1/posts?page=2&size=20&sort=createdAt%2Casc>; rel="last"`,
					pagination.HeaderXHasNext:    "true",
					pagination.HeaderXTotalCount: "30",
					pagination.HeaderXTotalPage:  "2",
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

					code, body, header, err := e2e.Get(ctx, "/posts"+tc.path)
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						var got []dto.PostResponse
						err = json.Unmarshal(body, &got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreTypes(time.Time{}), cmpopts.IgnoreFields(dto.PostResponse{}, "CreatedAt", "UpdatedAt", "DeletedAt")); diff != "" {
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
			want       *dto.PostResponse
			wantCode   int
			wantErr    *errorutils.APIError
			wantErrs   *errorutils.APIErrors
		}{
			{
				when:       "by admin",
				it:         "should success",
				authUser:   adminUser,
				id:         posts[0].ID,
				updateJSON: `{ "title": "TAIL", "body":"12312312^123123123" }`,
				want: &dto.PostResponse{
					ID:    posts[0].ID,
					Title: "TAIL",
					Body:  "12312312^123123123",
				},
				wantCode: http.StatusOK,
			},
			{
				when:       "by anonymous user",
				it:         "should fail",
				id:         posts[0].ID,
				updateJSON: `{ "title": "TAIL" }`,
				wantCode:   http.StatusUnauthorized,
				wantErr:    errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:       "by registered user",
				it:         "should fail",
				authUser:   user,
				id:         posts[0].ID,
				updateJSON: `{ "title": "TAIL" }`,
				wantCode:   http.StatusUnauthorized,
				wantErr:    errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:       "by mod",
				it:         "should success",
				authUser:   modUser,
				id:         posts[0].ID,
				updateJSON: `{ "title": "TAIL", "body":"12312312^123123123" }`,
				want: &dto.PostResponse{
					ID:    posts[0].ID,
					Title: "TAIL",
					Body:  "12312312^123123123",
				},
				wantCode: http.StatusOK,
			},
			{
				when:       "empty fields",
				it:         "should fail",
				authUser:   modUser,
				id:         posts[0].ID,
				updateJSON: `{}`,
				wantCode:   http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Required("Title"), nil),
				}},
			},
			{
				when:       "title too short",
				it:         "should fail",
				authUser:   modUser,
				id:         posts[0].ID,
				updateJSON: `{ "title": "k" }`,
				wantCode:   http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Min("Title"), nil),
				}},
			},
			{
				when:       "title too long",
				it:         "should fail",
				authUser:   modUser,
				id:         posts[0].ID,
				updateJSON: `{ "title": "titleTitleTitleTitleTitleTitle TITLE title" }`,
				wantCode:   http.StatusBadRequest,
				wantErrs: &errorutils.APIErrors{Errors: []*errorutils.APIError{
					errorutils.New(errorutils.Max("Title"), nil),
				}},
			},
			{
				when:       "invalid json",
				it:         "should fail",
				authUser:   modUser,
				id:         user.ID,
				updateJSON: `{ "title": "as das }`,
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

					code, body, _, err := e2e.Put(ctx, "/posts/"+strconv.FormatUint(tc.id, 10), []byte(tc.updateJSON))
					Expect(err).ToNot(HaveOccurred())
					Expect(code).To(Equal(tc.wantCode))

					if tc.want != nil {
						got := new(dto.PostResponse)
						err = json.Unmarshal(body, got)
						Expect(err).ToNot(HaveOccurred())

						if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(dto.PostResponse{}, "CreatedAt", "UpdatedAt", "DeletedAt")); diff != "" {
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
				when:     "by registered",
				it:       "should fail",
				id:       posts[0].ID,
				authUser: user,
				wantCode: http.StatusUnauthorized,
				wantErr:  errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:     "by mod",
				it:       "should fail",
				id:       posts[0].ID,
				authUser: modUser,
				wantCode: http.StatusUnauthorized,
				wantErr:  errorutils.New(errorutils.ErrUnauthorized, nil),
			},
			{
				when:     "by admin",
				it:       "should success",
				id:       posts[0].ID,
				authUser: adminUser,
				want:     &dto.ResponseWithID{ID: strconv.FormatUint(posts[0].ID, 10)},
				wantCode: http.StatusOK,
			},
			{
				when:     "valid id but no data",
				it:       "should fail",
				id:       20000,
				authUser: adminUser,
				wantCode: http.StatusNotFound,
				wantErr:  errorutils.New(errorutils.ErrPostNotFound, nil),
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

					code, body, _, err := e2e.Delete(ctx, "/posts/"+strconv.FormatUint(tc.id, 10))
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
