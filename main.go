package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
)

var pngSignature = []byte{
	137, 80, 78, 71, 13, 10, 26, 10,
}

type Chunk struct {
	Length uint32
	Type   string
	Data   []byte
	CRC    uint32
}

func main() {
	data, err := os.ReadFile("frame_5700.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := verifySignature(data); err != nil {
		log.Fatal(err)
	}

	offset := len(pngSignature)

	for {
		chunk, nextOffset, err := readChunk(data, offset)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"Chunk: %-4s Length: %d\n",
			chunk.Type,
			chunk.Length,
		)

		offset = nextOffset

		if chunk.Type == "IEND" {
			break
		}
	}
}

func verifySignature(data []byte) error {
	if len(data) < len(pngSignature) {
		return errors.New("file too small")
	}

	if !slices.Equal(data[:8], pngSignature) {
		return errors.New("invalid PNG signature")
	}

	return nil
}

func readChunk(data []byte, offset int) (Chunk, int, error) {
	if offset+8 > len(data) {
		return Chunk{}, 0, errors.New("unexpected end of file")
	}

	length := binary.BigEndian.Uint32(
		data[offset : offset+4],
	)
	offset += 4

	chunkType := string(
		data[offset : offset+4],
	)
	offset += 4

	endOfData := offset + int(length)

	if endOfData+4 > len(data) {
		return Chunk{}, 0, errors.New("truncated chunk")
	}

	chunkData := data[offset:endOfData]
	offset = endOfData

	crc := binary.BigEndian.Uint32(
		data[offset : offset+4],
	)
	offset += 4

	return Chunk{
		Length: length,
		Type:   chunkType,
		Data:   chunkData,
		CRC:    crc,
	}, offset, nil
}