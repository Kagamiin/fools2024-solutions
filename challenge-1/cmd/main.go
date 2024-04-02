package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/Kagamiin/fools2024-solutions/challenge-1/cmd/decomp"
)

// The ROM must be provided separately and is not included with the repository.
//go:embed pokeblue.gb
var pokéRom []byte

const MissingnoOffset = 0x1900
const MissingnoBaseWidth = 8
const MissingnoBaseHeight = 8

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "help" || os.Args[1] == "--help" {
		fmt.Printf(`usage: %v rest_in_miss_forever_ingno.sav

rest_in_miss_forever_ingno.sav: save file containing the data to be unscrambled.
Generates the following files in the current directory:
- result.bin: contains the best-effort unscrambled data
- unknownbits.bin: contains a bitmap of memory where data was permanently overwritten`, os.Args[0])
		os.Exit(1)
	}
	
	savData, err := readSavFile(os.Args[1])
	handle(err)
	
	memSpace := make([]byte, 65536)
	
	for offs := 0; offs < 8192; offs++ {
		srcAddr := offs
		destAddr := 0xa000 + offs
		memSpace[destAddr] = savData[srcAddr]
	}
	
	recording := decomp.RecordDecompressSprite(pokéRom, MissingnoOffset, MissingnoBaseWidth, MissingnoBaseHeight)
	
	unknownBitMap := recording.UndoRecording(&memSpace)

	err = dumpBin("result.bin", &memSpace)
	handle(err)
	
	err = dumpBin("unknownbits.bin", unknownBitMap)
	handle(err)
	
	log.Println("Done.")
}

func readSavFile(path string) ([]byte, error) {
	savFile, err := os.Open(path)
	defer func() {
		ec := savFile.Close()
		handle(ec)
	}()
	if err != nil {
		return nil, err
	}
	
	data := make([]byte, 8192)
	numBytes, err := savFile.Read(data)
	if err != nil {
		return nil, err
	}
	if numBytes < 8192 {
		return nil, fmt.Errorf("expected 8192 bytes, got %d", numBytes)
	}
	
	return data, nil
}

func dumpBin(path string, data *[]byte) error {
	f, err := os.Create(path)
	defer func() {
		ec := f.Close()
		handle(ec)
	}()
	if err != nil {
		return err
	}
	fWriter := bufio.NewWriter(f)

	_, err = fWriter.Write(*data)
	if err != nil {
		return err
	}
	err = fWriter.Flush()
	if err != nil {
		return err
	}
	
	return nil
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
