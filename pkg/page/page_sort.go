package page

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
)

// 分页 排序 帮助函数

type Page struct {
	TotalCount int         `json:"count"`
	List       interface{} `json:"list"`
}

func Offset(pageNo int, pageSize int) int {
	return (pageNo - 1) * pageSize
}

func NewPage(count int, list interface{}) *Page {
	page := Page{
		TotalCount: count,
		List:       list,
	}
	return &page
}

// PageRequest 分页加筛选
type PageRequest struct {
	PageIndex   int                    `json:"page_index"`
	PageSize    int                    `json:"page_size"`
	SortBy      string                 `json:"sort"`
	Filter      map[string]interface{} `json:"filter"`
	AllowFields []string               `json:"-"`
}

func NewPageRequest(pi, ps int, sort string, filter map[string]interface{}) *PageRequest {
	pr := &PageRequest{
		PageIndex: pi,
		PageSize:  ps,
		SortBy:    sort,
		Filter:    filter,
	}
	return pr
}

func NewPageFromRequest(rForm url.Values) *PageRequest {
	pi := cast.ToInt(rForm.Get("pageIndex"))
	ps := cast.ToInt(rForm.Get("pageSize"))
	sort := cast.ToString(rForm.Get("sort"))
	req := NewPageRequest(pi, ps, sort, nil)
	req.checkField()

	// 收集Filter
	for k, v := range rForm {
		if len(v) > 0 {
			vv := v[0]
			switch k {
			case "pageIndex", "pageSize", "sort":
			default:
				req.Filter[k] = vv
			}
		}
	}

	return req
}

func NewPageFromRequestJSON(requestBody []byte) (*PageRequest, error) {
	var req PageRequest

	if !json.Valid(requestBody) {
		return nil, errors.New("json format error")
	}

	gj := gjson.ParseBytes(requestBody)
	req.PageIndex = int(gj.Get("pageIndex").Int())
	req.PageSize = int(gj.Get("pageSize").Int())
	req.SortBy = gj.Get("sort").String()
	req.checkField()

	for k, v := range gj.Map() {
		switch k {
		case "pageIndex", "pageSize", "sort":
		default:
			switch v.Type {
			case gjson.JSON, gjson.Null:
				// ignore complex json params.
				continue
			}

			req.Filter[k] = v.Value()
		}
	}
	return NewPageRequest(req.PageIndex, req.PageSize, req.SortBy, req.Filter), nil
}

// checkField 检查参数是否正确
func (self *PageRequest) checkField() {
	if self.PageIndex == 0 {
		self.PageIndex = 1
	}

	if self.PageSize == 0 {
		self.PageSize = 10
	}

	if self.Filter == nil {
		self.Filter = make(map[string]interface{})
	}
}

// AddAllowSortField 排序字段白名单
func (self *PageRequest) AddAllowSortField(fieldName ...string) {
	if self.AllowFields == nil {
		self.AllowFields = make([]string, 0)
	}
	self.AllowFields = append(self.AllowFields, fieldName...)
}

func (self *PageRequest) GetLimitAndOffset() (limit, offset int) {
	return self.PageSize, Offset(self.PageIndex, self.PageSize)
}

func (self *PageRequest) Sort() (sort string, ok bool) {
	if len(self.SortBy) > 0 {
		// 末尾带+表示升序（默认）
		// 末尾带-表示降序
		sorts := strings.Split(self.SortBy, ",")

		var dbSorts []string
		for _, v := range sorts {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			var desc bool

			if strings.HasSuffix(v, "-") {
				v = strings.ReplaceAll(v, "-", "")
				desc = true
			}
			if strings.HasSuffix(v, "+") {
				v = strings.ReplaceAll(v, "+", "")
			}

			// 是否在 AllowFields 中
			if containString(self.AllowFields, v) {
				if desc {
					dbSorts = append(dbSorts, v+" DESC")
				} else {
					dbSorts = append(dbSorts, v)
				}
			}
		}
		if len(dbSorts) != 0 {
			return strings.Join(dbSorts, ","), true
		}
	}
	return "", false
}

func containString(ss []string, s string) bool {
	for _, e := range ss {
		if strings.TrimSpace(e) == strings.TrimSpace(s) {
			return true
		}
	}
	return false
}
