package vultr

import (
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

func structToMap(s interface{}) map[string]string {
	st := structs.New(s)
	st.TagName = "json"
	m := st.Map()
	m2 := make(map[string]string)
	for k, v := range m {
		if val, ok := v.(string); ok {
			m2[k] = val
		}
	}
	return m2
}

// Filter is a simple interface to describe type that accepts a
// map[string]string and returns a boolean saying whether or not it
// matches the filter.
type Filter interface {
	F(map[string]string) bool
}

type filter struct {
	name   string
	values []string
}

func (f *filter) F(m map[string]string) bool {
	v, ok := m[f.name]
	if !ok {
		return false
	}
	for _, value := range f.values {
		if v == value {
			return true
		}
	}
	return false
}

type multiFilter []filter

func (f multiFilter) F(m map[string]string) bool {
	for _, filter := range f {
		if !filter.F(m) {
			return false
		}
	}
	return true
}

func filterFromSet(set *schema.Set) Filter {
	var filters []filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string
		for _, value := range m["values"].([]interface{}) {
			values = append(values, value.(string))
		}
		filters = append(filters, filter{
			name:   m["name"].(string),
			values: values,
		})
	}
	return multiFilter(filters)
}
