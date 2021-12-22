package page

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

// package page
// 分页 排序 帮助函数

// Page 系统分页默认结构
type Page struct {
	TotalCount int64       `json:"count"`
	List       interface{} `json:"list"`
}

// Offset 计算真实偏移量
func Offset(pageNo int, pageSize int) int {
	return (pageNo - 1) * pageSize
}

// NewPage 新建一个 Page
func NewPage(count int64, list interface{}) *Page {
	page := Page{
		TotalCount: count,
		List:       list,
	}
	return &page
}

// PageRequest 封装分页加筛选请求
type PageRequest struct {
	PageIndex   int                    `form:"page_index" json:"page_index"`
	PageSize    int                    `form:"page_size" json:"page_size"`
	SortBy      string                 `form:"sort" json:"sort"`
	Filter      map[string]interface{} `form:"filter" json:"filter"`
	AllowFields []string               `form:"-" json:"-"`
}

// NewPageRequest 新建一个 PageRequest
func NewPageRequest(pi, ps int, sort string, filter map[string]interface{}) *PageRequest {
	pr := &PageRequest{
		PageIndex: pi,
		PageSize:  ps,
		SortBy:    sort,
		Filter:    filter,
	}
	return pr
}

// NewPageFromRequest 从 http.Request 中解析参数
// 使用 request.Form 前，需要 request.ParseForm()
func NewPageFromRequest(rForm url.Values) *PageRequest {
	pi := cast.ToInt(rForm.Get("page_index"))
	ps := cast.ToInt(rForm.Get("page_size"))
	sort := cast.ToString(rForm.Get("sort"))
	req := NewPageRequest(pi, ps, sort, nil)
	req.checkField()

	// 收集Filter
	for k, v := range rForm {
		if len(v) > 0 {
			vv := v[0]
			switch k {
			case "page_index", "page_size", "sort":
			default:
				req.Filter[k] = vv
			}
		}
	}

	return req
}

// NewPageFromRequestJSON 从请求体 JSON 中解析参数
func NewPageFromRequestJSON(requestBody []byte) (*PageRequest, error) {
	var req PageRequest

	if !json.Valid(requestBody) {
		return nil, errors.New("json format error")
	}

	gj := gjson.ParseBytes(requestBody)
	req.PageIndex = int(gj.Get("page_index").Int())
	req.PageSize = int(gj.Get("page_size").Int())
	req.SortBy = gj.Get("sort").String()
	req.checkField()

	for k, v := range gj.Map() {
		switch k {
		case "page_index", "page_size", "sort":
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
func (pr *PageRequest) checkField() {
	if pr.PageIndex == 0 {
		pr.PageIndex = 1
	}

	if pr.PageSize == 0 {
		pr.PageSize = 10
	}

	if pr.Filter == nil {
		pr.Filter = make(map[string]interface{})
	}
}

// AddAllowSortField 排序字段白名单
func (pr *PageRequest) AddAllowSortField(fieldName ...string) {
	if pr.AllowFields == nil {
		pr.AllowFields = make([]string, 0)
	}
	pr.AllowFields = append(pr.AllowFields, fieldName...)
}

// GetLimitAndOffset 获取 limit offset 用于数据库分页
func (pr *PageRequest) GetLimitAndOffset() (limit, offset int) {
	return pr.PageSize, Offset(pr.PageIndex, pr.PageSize)
}

// Sort 获取 SQL 的 Order By 子句
// create_time+ -> create_time
// create_time- -> create_time DESC
// id+,create_time- -> id,create_time DESC
func (pr *PageRequest) Sort() (sort string, ok bool) {
	if len(pr.SortBy) > 0 {
		// 末尾带+表示升序（默认）
		// 末尾带-表示降序
		sorts := strings.Split(pr.SortBy, ",")

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
			if containString(pr.AllowFields, v) {
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
