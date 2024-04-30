package main

func main() {
	storage := MockStorage{}
	api := NewApi(&storage, ":3000")
	api.Run()
}