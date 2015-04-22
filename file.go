package gofigure

import (
	"log"
	"strings"
)

type CategoryMap map[string]*Category
type ValueMap map[string]string

type Category struct {
	Name		string
	Parent		*Category
	Categories	CategoryMap
	Values		ValueMap
}

func NewCategory(name string, parent *Category) *Category {
	return &Category {
		Name:       name,
		Parent:     parent,
		Categories: CategoryMap{},
		Values:     ValueMap{},
	}
}

func (c Category) FindOption(o *Option) (string, bool) {
	if o.fileSpec != "" {
		val, ok := c.Find(o.fileSpec)
		return val, ok
	}
	return "", false
}

func (c Category) Find(spec string) (string, bool) {
	fs := strings.Split(spec, ".")
	if len(fs) > 1 {
		next := c.Categories.Get(fs[0])
		if next != nil {
			n, e := next.Find(strings.Join(fs[1:], "."))
			return n, e
		} else {
			return "", false
		}
	} else {
		return c.Values.Get(fs[0]), c.Values.Exists(fs[0])
	}
}

func (c CategoryMap) Exists(key string) bool {
	if _, ok := c[key]; ok {
		return true
	}
	return false
}

func (c CategoryMap) MustGet(key string) *Category {
	if v, ok := c[key]; ok {
		return v
	}
	log.Panicf("Undefined required category %s.", key)
	return nil
}

func (c CategoryMap) Get(key string) *Category {
	if v, ok := c[key]; ok {
		return v
	}
	return nil
}

func (c CategoryMap) Set(key string, value *Category) {
	c[key] = value
}

func (c CategoryMap) Delete(key string) {
	delete(c, key)
}

func (c ValueMap) Exists(key string) bool {
	if _, ok := c[key]; ok {
		return true
	}
	return false
}

func (c ValueMap) MustGet(key string) string {
	if v, ok := c[key]; ok {
		return v
	}
	log.Panicf("Undefined required value %s.", key)
	return ""
}

func (c ValueMap) Get(key string) string {
	if v, ok := c[key]; ok {
		return v
	}
	return ""
}

func (c ValueMap) Set(key, value string) {
	c[key] = value
}

func (c ValueMap) Delete(key string) {
	delete(c, key)
}


type File interface {
	Parse(name string) (*Category, error)
	ParseConfig(name string, specs map[string]*Option) (ValueMap, error)
}