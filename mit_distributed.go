package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	store := Store{dict: make(map[string][]string)}
	store.SetInitDict()
	fmt.Println("-------sync-go----------")
	syncFetch("first", store, make(map[string]bool))
	fmt.Println("-------async-go----------")
	asyncFetch("first", store, &asyncGo{dict: make(map[string]bool)})
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
		syncFetch(url, store, dict)
	}
}

func asyncFetch(url string, store Store, ass *asyncGo) {
	ass.mu.Lock()
	if _, ok := ass.dict[url]; ok {
		ass.mu.Unlock()
		return
	}
	ass.dict[url] = true
	ass.mu.Unlock()

	res, err := store.Fetch(url)
	if err != nil {
		fmt.Println("error fetching")
	}

	var wg sync.WaitGroup
	wg.Add(len(res))

	for _, url := range res {
		go func(url string) {
			defer wg.Done()
			asyncFetch(url, store, ass)
		}(url)
	}
	wg.Wait()
}

type asyncGo struct {
	mu   sync.Mutex
	dict map[string]bool
}
type Store struct {
	dict map[string][]string
}

func (d Store) Fetch(url string) ([]string, error) {
	if res, ok := d.dict[url]; ok {
		fmt.Printf("found:   %s\n", url)
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
