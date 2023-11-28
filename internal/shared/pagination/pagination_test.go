package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/shared/pagination"
)

func TestPagination(t *testing.T) {
	tests := []struct {
		name          string
		pageable      *pagination.Pageable
		wantHasNext   bool
		wantTotalPage int64
	}{
		{
			name:          "first page",
			pageable:      &pagination.Pageable{Page: getInt64Pointer(1), Size: getInt64Pointer(20), TotalCount: 100},
			wantHasNext:   true,
			wantTotalPage: 5,
		},
		{
			name:          "last page",
			pageable:      &pagination.Pageable{Page: getInt64Pointer(5), Size: getInt64Pointer(20), TotalCount: 100},
			wantHasNext:   false,
			wantTotalPage: 5,
		},
		{
			name:          "second page",
			pageable:      &pagination.Pageable{Page: getInt64Pointer(2), Size: getInt64Pointer(10), TotalCount: 50},
			wantHasNext:   true,
			wantTotalPage: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if hasNext := tt.pageable.HasNext(); hasNext != tt.wantHasNext {
				t.Errorf("HasNext() = %v, want %v", hasNext, tt.wantHasNext)
			}

			if totalPage := tt.pageable.GetTotalPage(); totalPage != tt.wantTotalPage {
				t.Errorf("GetTotalPage() = %v, want %v", totalPage, tt.wantTotalPage)
			}
		})
	}
}

func TestPaginationHeaders(t *testing.T) {
	testCases := []struct {
		name        string
		pageable    *pagination.Pageable
		wantCount   int64
		wantPages   int64
		wantHasNext bool
		wantLink    string
	}{
		{
			name: "first page",
			pageable: &pagination.Pageable{
				Size:       getInt64Pointer(20),
				Page:       getInt64Pointer(1),
				TotalCount: 100,
			},
			wantCount:   100,
			wantPages:   5,
			wantHasNext: true,
			wantLink: `</test-pagination?page=1&size=20&sort=createdAt%2CDESC>; rel="first"
</test-pagination?page=2&size=20&sort=createdAt%2CDESC>; rel="next",
</test-pagination?page=5&size=20&sort=createdAt%2CDESC>; rel="last"`,
		},
		{
			name: "last page",
			pageable: &pagination.Pageable{
				Size:       getInt64Pointer(20),
				Page:       getInt64Pointer(5),
				TotalCount: 100,
			},
			wantCount:   100,
			wantPages:   5,
			wantHasNext: false,
			wantLink: `</test-pagination?page=1&size=20&sort=createdAt%2CDESC>; rel="first"
</test-pagination?page=4&size=20&sort=createdAt%2CDESC>; rel="prev",
</test-pagination?page=5&size=20&sort=createdAt%2CDESC>; rel="last"`,
		},
		{
			name: "third page",
			pageable: &pagination.Pageable{
				Size:       getInt64Pointer(20),
				Page:       getInt64Pointer(2),
				TotalCount: 100,
			},
			wantCount:   100,
			wantPages:   5,
			wantHasNext: true,
			wantLink: `</test-pagination?page=1&size=20&sort=createdAt%2CDESC>; rel="first"
</test-pagination?page=3&size=20&sort=createdAt%2CDESC>; rel="next",
</test-pagination?page=1&size=20&sort=createdAt%2CDESC>; rel="prev",
</test-pagination?page=5&size=20&sort=createdAt%2CDESC>; rel="last"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/test-pagination", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			totalCount, link := tc.pageable.PaginationHeader(c)

			if link != tc.wantLink {
				t.Errorf("Expected link to be %s, got %s", tc.wantLink, link)
			}

			if totalCount != tc.wantCount {
				t.Errorf("Expected totalCount to be %d, got %d", tc.wantCount, totalCount)
			}

			if xTotalCount := rec.Header().Get(pagination.HeaderXTotalCount); xTotalCount != strconv.FormatInt(tc.wantCount, 10) {
				t.Errorf("Expected X-Total-Count header to be %d, got %s", tc.wantCount, xTotalCount)
			}

			xTotalPage, _ := strconv.ParseInt(rec.Header().Get(pagination.HeaderXTotalPage), 10, 64)
			if xTotalPage != tc.wantPages {
				t.Errorf("Expected X-Total-Page header to be %d, got %d", tc.wantPages, xTotalPage)
			}

			xHasNext := rec.Header().Get(pagination.HeaderXHasNext) == "true"
			if xHasNext != tc.wantHasNext {
				t.Errorf("Expected X-Has-Next header to be %t, got %t", tc.wantHasNext, xHasNext)
			}
		})
	}
}

func getInt64Pointer(v int64) *int64 {
	return &v
}
