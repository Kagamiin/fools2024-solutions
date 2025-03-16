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

package maps

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
)

type TilesetAttributes struct {
	Id uint8
	Name string
	ConstName string
	MetatileFileName string
}

func ReadTilesets(tilesetConstantsInfile io.ReadSeekCloser, tilesetsInfile io.ReadSeekCloser, tilesetAttrChan chan<- TilesetAttributes) {
	defer func() {
		close(tilesetAttrChan)
		_ = tilesetsInfile.Close()
		_ = tilesetConstantsInfile.Close()
	}()
	
	tilesetConstantsInfile.Seek(0, io.SeekStart)
	tilesetsInfile.Seek(0, io.SeekStart)
	tilesetConstantsLineReader := bufio.NewReader(tilesetConstantsInfile)
	tilesetsLineReader := bufio.NewReader(tilesetsInfile)
	
	var err error
	
	constRe := regexp.MustCompile(`^\s*const (TILESET_\S+)\s*; ([0-9a-f]+)$`)
	tilesetRe := regexp.MustCompile(`^\s*tileset (Tileset\S+)$`)
	
	tilesetAttrMap := make(map[uint8]TilesetAttributes)
	tilesetAttrMap[0] = TilesetAttributes{
		Id: 0,
	}
	
	for ; err == nil; {
		var line []byte
		var isPrefix bool
		line, isPrefix, err = tilesetConstantsLineReader.ReadLine()
		for ; isPrefix && err == nil; {
			var frag []byte
			frag, isPrefix, err = tilesetConstantsLineReader.ReadLine()
			line = append(line, frag...)
		}
		if err != nil {
			break
		}
		lineS := string(line)
		res := constRe.FindStringSubmatch(lineS)
		if len(res) >= 3 {
			tilesetIdx, err := strconv.ParseUint(res[2], 16, 8)
			handle(err)
			tilesetAttrMap[uint8(tilesetIdx)] = TilesetAttributes{
				Id: uint8(tilesetIdx),
				ConstName: res[1],
			}
		}
	}
	
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
	
	err = nil
	tilesetIdx := uint8(0)
	
	for ; err == nil; {
		var line []byte
		var isPrefix bool
		line, isPrefix, err = tilesetsLineReader.ReadLine()
		for ; isPrefix && err == nil; {
			var frag []byte
			frag, isPrefix, err = tilesetsLineReader.ReadLine()
			line = append(line, frag...)
		}
		if err != nil {
			break
		}
		lineS := string(line)
		res := tilesetRe.FindStringSubmatch(lineS)
		if len(res) >= 2 {
			if _, ok := tilesetAttrMap[tilesetIdx]; !ok {
				tilesetIdx++
				continue
			}
			tAttr := tilesetAttrMap[tilesetIdx]
			tAttr.Name = res[1]
			tilesetAttrMap[tilesetIdx] = tAttr
			tilesetAttrChan <- tAttr
			tilesetIdx++
		}
	}
	
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
}

func ReadMetatileFileNames(infile io.ReadSeekCloser, tilesetAttrChan <-chan TilesetAttributes, tilesetMetatileChan chan<- TilesetAttributes) {
	defer func() {
		close(tilesetMetatileChan)
		_ = infile.Close()
	}()
	
	infile.Seek(0, io.SeekStart)
	lineReader := bufio.NewReader(infile)
	
	var err error
	
	metaLabelRe := regexp.MustCompile(`^(Tileset\S+)Meta::$`)
	metaIncbinRe := regexp.MustCompile(`^INCBIN "(data\/tilesets\/\S+_metatiles\.bin)"$`)
	
	tilesetAttrMap := make(map[string]TilesetAttributes)
	
	for ;; {
		tAttr, ok := <-tilesetAttrChan
		if ok {
			tilesetAttrMap[tAttr.Name] = tAttr
		} else {
			break
		}
	}
	
	
	
	currentMetaLabel := "invalid"
	
	for ; err == nil; {
		var line []byte
		var isPrefix bool
		line, isPrefix, err = lineReader.ReadLine()
		for ; isPrefix && err == nil; {
			var frag []byte
			frag, isPrefix, err = lineReader.ReadLine()
			line = append(line, frag...)
		}
		if err != nil {
			break
		}
		lineS := string(line)
		res := metaLabelRe.FindStringSubmatch(lineS)
		if len(res) >= 2 {
			currentMetaLabel = res[1]
			continue
		}
		res = metaIncbinRe.FindStringSubmatch(lineS)
		if len(res) >= 2 {
			if _, ok := tilesetAttrMap[currentMetaLabel]; !ok {
				continue
			}
			tAttr := tilesetAttrMap[currentMetaLabel]
			tAttr.MetatileFileName = "extern/pokecrystal/" + res[1]
			tilesetAttrMap[currentMetaLabel] = tAttr
			tilesetMetatileChan <- tAttr
		}
	}
	
	tilesetAttrMap["TilesetCave"] = TilesetAttributes{
		Id: 0x18,
		Name: "TilesetCave",
		ConstName: "TILESET_CAVE",
		MetatileFileName: "extern/pokecrystal/data/tilesets/cave_collision.asm",
	}
	tilesetMetatileChan <- tilesetAttrMap["TilesetCave"]
	
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
}

func ServeTilesetAttributes(
	tilesetMetatileChan <-chan TilesetAttributes,
	requestChan <-chan MapAttributes,
	responseChan chan<- MapAttributes,
) {
	defer func() {
		close(responseChan)
	}()
	
	tilesetAttrMap := make(map[string]TilesetAttributes)
	
	for ;; {
		tAttr, ok := <-tilesetMetatileChan
		if ok {
			tilesetAttrMap[tAttr.ConstName] = tAttr
		} else {
			break
		}
		//log.Println(tAttr)
	}
	
	for ;; {
		mapAttr, ok := <-requestChan
		if ok {
			mapAttr.Tileset = tilesetAttrMap[mapAttr.Tileset.ConstName]
			responseChan <- mapAttr
		} else {
			break
		}
	}
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
