package main

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	server := NewServer()
	if err := server.Start(serverAddress); err != nil {
		// log error
		//exit
	}

}
