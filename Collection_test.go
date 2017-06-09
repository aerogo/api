package api

import (
	"fmt"
	"testing"
)

type SomeDatabase struct{}

func (db *SomeDatabase) Get(table string, id string) (interface{}, error) {
	return nil, nil
}

func (db *SomeDatabase) Set(table string, id string, obj interface{}) error {
	return nil
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
	list := &SomeList{}
	db := &SomeDatabase{}
	api := New("/api/", db)
	route, handler := api.Get(list)
	fmt.Println(route, handler)
}
