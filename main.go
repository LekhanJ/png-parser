package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
)

func main() {
	data, err := os.ReadFile("frame_5700.png")
	if err != nil {
		log.Fatal(err)
	}

	_, err = verifySignature(data[:8])
	if err != nil {
		fmt.Println("Error occurred: ", err)
	}

	fmt.Println(data)
}

func readUInt32(bytes []byte, offset int) uint32 {
	return uint32(bytes[offset])<<24 |
		uint32(bytes[offset+1])<<16 |
		uint32(bytes[offset+2])<<8 |
		uint32(bytes[offset+3])
}

func verifySignature(data []byte) (bool, error) {
	signature := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	if !slices.Equal(data, signature) {
		return false, errors.New("Invalid Signature")
	}
	return true, nil
}