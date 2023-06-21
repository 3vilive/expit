package main

import "log"

func main() {
	closed := make(chan bool)

	select {
	case <-closed:
		log.Printf("closed\n")
	default:
		log.Printf("default\n")
	}

	close(closed)

	select {
	case <-closed:
		log.Printf("closed\n")
	default:
		log.Printf("default\n")
	}

}
