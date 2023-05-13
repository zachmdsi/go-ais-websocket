package main

func main() {
	srv := NewServer("data/ais-sample-data.csv")
	srv.Start(":8080")
}
