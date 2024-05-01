package main

func main() {
	storage := NewPostgresStore()
	api := NewApi(storage, ":3000")
	api.Run()
}