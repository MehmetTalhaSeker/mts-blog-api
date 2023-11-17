package pagination

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
)

// Pageable defined for pagination.
type Pageable struct {
	Page       *int64 `query:"page" validate:"omitempty,min=1"`
	Size       *int64 `query:"size" validate:"omitempty,min=1,max=100"`
	Sort       string `query:"sort"`
	TotalCount int64
}

func NewPagination() *Pageable {
	defaultPage := int64(1)
	defaultSize := int64(20)

	return &Pageable{Page: &defaultPage, Size: &defaultSize, Sort: "createdAt,desc"}
}

func (p *Pageable) HasNext() bool {
	return *p.Page < p.GetTotalPage()-1
}

func (p *Pageable) GetTotalPage() int64 {
	return int64(math.Ceil(float64(p.TotalCount) / float64(*p.Size)))
}

func (p *Pageable) IsLast() bool {
	return *p.Page == p.GetTotalPage()
}

func (p *Pageable) Limit() *int64 {
	return p.Size
}

func (p *Pageable) Offset() int {
	return int((*p.Page - 1) * *p.Size)
}

func (p *Pageable) Order() string {
	sk, sv := p.GetSortKeyAndValue()

	return fmt.Sprintf("%s %s", strcase.ToSnake(sk), sv)
}

func (p *Pageable) GetSortKeyAndValue() (string, string) {
	split := strings.Split(p.Sort, ",")
	if len(split) == 2 && (split[1] == "desc" || split[1] == "asc") {
		return split[0], split[1]
	}

	return "createdAt", "DESC"
}

func (p *Pageable) GetSortKey() string {
	split := strings.Split(p.Sort, ",")

	return strcase.ToSnake(split[0])
}

func (p *Pageable) GetSortValue() int {
	split := strings.Split(p.Sort, ",")
	if len(split) == 2 {
		switch split[1] {
		case "asc":
			return 1
		case "desc":
			return -1
		}
	}

	return -1
}

func (p *Pageable) PaginationHeader(c echo.Context) (int64, string) {
	u := c.Request().URL
	link := ""

	link += fmt.Sprintf("%s\n", p.prepareLink(u, 1, "first"))

	if *p.Page < p.GetTotalPage() {
		link += fmt.Sprintf("%s,\n", p.prepareLink(u, *p.Page+1, "next"))
	}

	if *p.Page > 1 {
		link += fmt.Sprintf("%s,\n", p.prepareLink(u, *p.Page-1, "prev"))
	}

	link += p.prepareLink(u, p.GetTotalPage(), "last")

	c.Response().Header().Set(HeaderLink, link)
	c.Response().Header().Set(HeaderXTotalCount, fmt.Sprintf("%d", p.TotalCount))
	c.Response().Header().Set(HeaderXTotalPage, fmt.Sprintf("%d", p.GetTotalPage()))
	c.Response().Header().Set(HeaderXHasNext, fmt.Sprintf("%t", p.HasNext()))

	return p.TotalCount, link
}

func (p *Pageable) prepareLink(requestURL *url.URL, page int64, relType string) string {
	u := *requestURL
	q := u.Query()
	q.Set("page", strconv.FormatInt(page, 10))
	q.Set("size", strconv.FormatInt(*p.Size, 10))
	key, value := p.GetSortKeyAndValue()
	q.Set("sort", fmt.Sprintf("%s,%s", key, value))

	u.RawQuery = q.Encode()

	return fmt.Sprintf("<%s>; rel=\"%s\"", u.String(), relType)
}

const (
	HeaderLink        = "Link"
	HeaderXTotalCount = "X-Total-Count"
	HeaderXTotalPage  = "X-Total-Page"
	HeaderXHasNext    = "X-Has-Next"
)
