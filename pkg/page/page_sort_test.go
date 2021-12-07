package page

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPageFromRequest(t *testing.T) {
	t.Run("TestSimpleForm", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/xx", nil)
		r.Form = make(url.Values)
		r.Form.Set("pageIndex", strconv.Itoa(1))
		r.Form.Set("pageSize", strconv.Itoa(10))
		r.Form.Set("sort", "create_time+")
		r.Form.Set("user_id", strconv.Itoa(1))
		pr := NewPageFromRequest(r.Form)

		assert.Equal(t, 1, pr.PageIndex)
		assert.Equal(t, 10, pr.PageSize)
		assert.Equal(t, "create_time+", pr.SortBy)
		assert.Equal(t, "1", pr.Filter["user_id"])
	})

	t.Run("TestEmptyPR", func(t *testing.T) {
		pr := PageRequest{}
		pr.checkField()

		assert.Equal(t, 1, pr.PageIndex)
		assert.Equal(t, 10, pr.PageSize)
		assert.NotNil(t, pr.Filter)
	})

}

func TestNewPageFromRequestJson(t *testing.T) {
	data := make(map[string]interface{})
	data["pageIndex"] = 1
	data["pageSize"] = 10
	data["sort"] = "create_time"
	data["user_id"] = 1
	data["create_time"] = "2022-02-12"
	data["this_will_ignore"] = map[string]interface{}{
		"ca": 1,
	}
	dataJson, _ := json.Marshal(data)

	pr, err := NewPageFromRequestJSON(dataJson)
	if err != nil {
		fmt.Println(err)
		return
	}

	assert.Equal(t, 1, pr.PageIndex)
	assert.Equal(t, 10, pr.PageSize)
	assert.Equal(t, "create_time", pr.SortBy)
	assert.Equal(t, 1.0, pr.Filter["user_id"])
	assert.Equal(t, "2022-02-12", pr.Filter["create_time"])

	_, ok := pr.Filter["this_will_ignore"]
	assert.Equal(t, false, ok)
}

func TestPageRequest_Sort(t *testing.T) {
	t.Run("TestAllowFields", func(t *testing.T) {
		pr := PageRequest{}
		pr.SortBy = "create_time,user_id,id"
		pr.AddAllowSortField("user_id", "id")

		sort, ok := pr.Sort()
		assert.Equal(t, true, ok)

		assert.NotContains(t, "create_time", sort)
	})

	t.Run("TestSort", func(t *testing.T) {
		pr := PageRequest{}
		pr.SortBy = "create_time-,user_id+,id"
		pr.AddAllowSortField("create_time", "user_id", "id")

		sort, ok := pr.Sort()
		assert.Equal(t, true, ok)

		assert.Equal(t, "create_time DESC,user_id,id", sort)
	})
}
