package main

import (
	"errors"
	"fmt"
)

func main() {
	store := Store{dict: make(map[string][]string)}
	store.SetInitDict()

	syncFetch("first", store, make(map[string]bool))
}

func syncFetch(url string, store Store, dict map[string]bool) {
	if _, ok := dict[url]; ok {
		return
	}
	dict[url] = true

	res, err := store.Fetch(url)
	if err != nil {
		fmt.Println("error fetching")
	}

	for _, url := range res {
		fmt.Println(url)
		syncFetch(url, store, dict)
	}

}

type Store struct {
	dict map[string][]string
}

func (d Store) Fetch(url string) ([]string, error) {
	if res, ok := d.dict[url]; ok {
		return res, nil
	}

	return nil, errors.New("incorrect url")
}

func (d Store) SetInitDict() {
	d.dict["first"] = []string{"second"}
	d.dict["second"] = []string{"first", "third"}
	d.dict["third"] = []string{"fourth"}
	d.dict["fourth"] = []string{}
}
