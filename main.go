package main

func main2() {
	storage := NewPostgresStore()
	api := NewApi(storage, ":3000")
	api.Run()
}