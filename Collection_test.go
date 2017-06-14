package api

import (
	"fmt"
	"reflect"
	"testing"
)

type SomeDatabase struct{}

func (db *SomeDatabase) Get(table string, id string) (interface{}, error) {
	return nil, nil
}

func (db *SomeDatabase) Set(table string, id string, obj interface{}) error {
	return nil
}

func (db *SomeDatabase) Delete(table string, id string) (bool, error) {
	return false, nil
}

func (db *SomeDatabase) Type(table string) reflect.Type {
	return reflect.TypeOf((*SomeList)(nil)).Elem()
}

func (db *SomeDatabase) Types() map[string]reflect.Type {
	return map[string]reflect.Type{
		"SomeList": reflect.TypeOf((*SomeList)(nil)).Elem(),
	}
}

type SomeList struct {
	Items []SomeListItem `json:"items"`
}

type SomeListItem struct {
	Name string `json:"name"`
}

func (list *SomeList) Add(element interface{}) error {
	return nil
}

func (list *SomeList) Remove(element interface{}) error {
	return nil
}

func (list *SomeList) Contains(element interface{}) bool {
	return false
}

func (list *SomeList) Save() error {
	return nil
}

func TestGet(t *testing.T) {
	db := &SomeDatabase{}
	api := New("/api/", db)
	route, handler := api.Get("SomeList")
	fmt.Println(route, handler)
}
