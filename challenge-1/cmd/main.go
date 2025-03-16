/*
 * fools2024-solutions: source code for Kagamiin's solutions for TheZZAZZGlitch April Fools Event 2024's Security Testing Program.
 * Copyright (C) 2024 Kagamiin~
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Kagamiin/fools2024-solutions/challenge-1/cmd/decomp"
)

// The ROM must be provided separately and is not included with the repository.
//go:embed pokeblue.gb
var pokéRom []byte

//go:embed pokeblue-ram.dmp
var pokéRam []byte

const MissingnoOffset = 0x1900
const MissingnoBaseWidth = 8
const MissingnoBaseHeight = 8

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "help" || os.Args[1] == "--help" {
		fmt.Printf(`usage: 
	%v rest_in_miss_forever_ingno.sav
	%v (--decompress|-d) pokeblue.sav bank addr width height

rest_in_miss_forever_ingno.sav: save file containing the data to be unscrambled.
Generates the following files in the current directory:
- result.bin: contains the best-effort unscrambled data
- unknownbits.bin: contains a bitmap of memory where data was permanently overwritten`, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	
	if os.Args[1] == "--decompress" || os.Args[1] == "-d" {
		decompressSprite()
		return
	}
	
	savData, err := readSavFile(os.Args[1])
	handle(err)
	
	memSpace := make([]byte, 65536)
	
	prepareMemSpace(memSpace, savData)
	
	recording := decomp.RecordDecompressSprite(pokéRom, MissingnoOffset, MissingnoBaseWidth, MissingnoBaseHeight)
	
	unknownBitMap := recording.UndoRecording(&memSpace)

	err = dumpBin("result.bin", &memSpace)
	handle(err)
	
	err = dumpBin("unknownbits.bin", unknownBitMap)
	handle(err)
	
	log.Println("Done.")
}

func decompressSprite() {
	if len(os.Args) < 5 || os.Args[2] == "-h" || os.Args[1] == "help" || os.Args[1] == "--help" {
		fmt.Printf(`usage: 
	%v rest_in_miss_forever_ingno.sav
	%v (--decompress|-d) pokeblue.sav bank addr [width height]

pokeblue.sav: save file containing the source data where the sprite will be decompressed.
bank: ROM bank where the sprite decompression is performed
addr: pointer to the sprite
width: width of the sprite in its base data, in tiles (omit to use the width from the sprite data)
height: height of the sprite in its base data, in tiles (omit to use the height from the sprite data)
Generates the following files in the current directory:
- decompressed.sav: contains the save data after decompression`, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	
	savData, err := readSavFile(os.Args[2])
	handle(err)
	
	bank, err := strconv.ParseUint(os.Args[3], 16, 64)
	handle(err)
	
	addr, err := strconv.ParseUint(os.Args[4], 16, 64)
	handle(err)
	
	var width int64
	var height int64
	
	if len(os.Args) >= 7 {
		width, err = strconv.ParseInt(os.Args[5], 16, 64)
		handle(err)
		
		height, err = strconv.ParseInt(os.Args[6], 16, 64)
		handle(err)
	} else {
		width = -1
		height = -1
	}
	
	memSpace := make([]byte, 65536)
	
	prepareMemSpace(memSpace, savData)
	
	for offs := 0; offs < 0x4000; offs++ {
		if bank == 0 {
			bank = 1
		}
		srcAddr := int(bank) * 0x4000 + offs
		destAddr := 0x4000 + offs
		memSpace[destAddr] = pokéRom[srcAddr]
	}
	
	recording := decomp.RecordDecompressSprite(memSpace, int(addr), int(width), int(height))
	recording.ApplyRecording(&memSpace)
	
	err = dumpBin("decompressed.bin", &memSpace)
	handle(err)
	
	log.Println("Done.")
}

func prepareMemSpace(memSpace []byte, savData []byte) {
	for offs := 0; offs < 0x2000; offs++ {
		srcAddr := offs
		destAddr := 0xa000 + offs
		memSpace[destAddr] = savData[srcAddr]
		
		srcAddr = offs
		destAddr = offs
		memSpace[destAddr] = pokéRom[srcAddr]
		srcAddr = 0x2000 + offs
		destAddr = 0x2000 + offs
		memSpace[destAddr] = pokéRom[srcAddr]
		
		srcAddr = 0xc000 + offs
		destAddr = 0xc000 + offs
		memSpace[destAddr] = pokéRam[srcAddr]
		srcAddr = 0xe000 + offs
		destAddr = 0xe000 + offs
		memSpace[destAddr] = pokéRam[srcAddr]
	}
	for offs := 0; offs < 0x1000; offs++ {
		srcAddr := 0x8000 + offs
		destAddr := 0x8000 + offs
		memSpace[destAddr] = pokéRam[srcAddr]
	}
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
	
	data := make([]byte, 32768)
	numBytes, err := savFile.Read(data)
	if err != nil {
		return nil, err
	}
	if numBytes < 32768 {
		return nil, fmt.Errorf("expected 32768 bytes, got %d", numBytes)
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
