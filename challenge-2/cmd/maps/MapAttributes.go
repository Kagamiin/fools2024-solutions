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
	"strings"
)

type MapConnection struct {
	Enabled bool
	DestinationMapName string
	BlockOffset int8
}

type MapAttributes struct {
	Name string
	IdName string
	BorderBlock uint8
	Height uint8
	Width uint8
	BlockFilePath string
	Tileset TilesetAttributes
	ConnectionNorth MapConnection
	ConnectionSouth MapConnection
	ConnectionWest MapConnection
	ConnectionEast MapConnection
}

func ReadBlocksAsm(infile io.ReadSeekCloser, mapAttrChan chan<- MapAttributes) {
	defer func() {
		close(mapAttrChan)
		_ = infile.Close()
	}()
	
	infile.Seek(0, io.SeekStart)
	lineReader := bufio.NewReader(infile)
	_, isPrefix, err := lineReader.ReadLine()
	for ; isPrefix && err == nil; {
		_, isPrefix, err = lineReader.ReadLine()
	}
	if err != nil {
		panic(err)
	}
	
	mapAttrNameList := []MapAttributes{}
	mapAttrRe := regexp.MustCompile(`^(\S*)_Blocks:$`)
	blocksFNameRe := regexp.MustCompile(`^\s*INCBIN "(\S*)"$`)
	
	for ; err == nil; {
		var line []byte
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
		res := mapAttrRe.FindStringSubmatch(lineS)
		if len(res) >= 2 {
			mapAttrNameList = append(mapAttrNameList, MapAttributes{
				Name: res[1],
			})
		} else if len(mapAttrNameList) > 0 {
			res = blocksFNameRe.FindStringSubmatch(lineS)
			if len(res) >= 2 {
				for _, item := range mapAttrNameList {
					item.BlockFilePath = "extern/pokecrystal/" + res[1]
					mapAttrChan <- item
				}
				mapAttrNameList = []MapAttributes{}
			}
		}
	}

	if err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
}

func ReadMapAttributes(infile io.ReadSeekCloser, mapAttrChan <-chan MapAttributes, mapTilesetChan chan<- MapAttributes) {
	defer func() {
		close(mapTilesetChan)
		_ = infile.Close()
	}()
	
	infile.Seek(0, io.SeekStart)
	lineReader := bufio.NewReader(infile)
	
	var err error
	
	mapAttrRe := regexp.MustCompile(`^\s*map_attributes (\S+), (\S+), \$([0-9][0-9]), (.*)$`)
	
	mapAttrMap := make(map[string]MapAttributes)
	
	for ;; {
		mapAttr, ok := <-mapAttrChan
		if ok {
			mapAttrMap[mapAttr.Name] = mapAttr
		} else {
			break
		}
	}
	
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
		res := mapAttrRe.FindStringSubmatch(lineS)
		if len(res) >= 5 {
			if _, ok := mapAttrMap[res[1]]; !ok {
				continue
			}
			mapAttr := mapAttrMap[res[1]]
			borderBlock, err := strconv.ParseUint(res[3], 16, 8)
			if err != nil {
				panic(err)
			}
			
			mapAttr.IdName = res[2]
			mapAttr.BorderBlock = uint8(borderBlock)
			readMapConnections(lineReader, &mapAttr, res[4])
			mapTilesetChan <- mapAttr
		}
	}
	
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
}

func readMapConnections(lineReader *bufio.Reader, mapAttr *MapAttributes, connStr string) {
	mapConnRe := regexp.MustCompile(`^\s*connection (\S+), (\S+), (\S+), (-?[0-9]+)$`)
	
	connList := strings.Split(connStr, " | ")
	for _, conn := range connList {
		connT := strings.TrimSpace(conn)
		if connT == "0" {
			continue
		}
		
		line, isPrefix, err := lineReader.ReadLine()
		for ; isPrefix && err == nil; {
			var frag []byte
			frag, isPrefix, err = lineReader.ReadLine()
			line = append(line, frag...)
		}
		if err != nil {
			break
		}
		lineS := string(line)
		res := mapConnRe.FindStringSubmatch(lineS)
		if len(res) >= 5 {
			var dest *MapConnection
			switch strings.TrimSpace(res[1]) {
			case "north":
				dest = &mapAttr.ConnectionNorth
			case "south":
				dest = &mapAttr.ConnectionSouth
			case "east":
				dest = &mapAttr.ConnectionEast
			case "west":
				dest = &mapAttr.ConnectionWest
			}
			blockOffset, err := strconv.Atoi(res[4])
			if err != nil {
				panic(err)
			}
			dest.Enabled = true
			dest.DestinationMapName = res[2]
			dest.BlockOffset = int8(blockOffset)
		}
	}
}

func ReadMapTileset(
	infile io.ReadSeekCloser,
	mapTilesetChan <-chan MapAttributes,
	mapWidthHeightChan chan<- MapAttributes,
	tilesetReqChan chan<- MapAttributes,
	tilesetRespChan <-chan MapAttributes,
) {
	defer func() {
		close(tilesetReqChan)
		close(mapWidthHeightChan)
		_ = infile.Close()
	}()
	
	infile.Seek(0, io.SeekStart)
	lineReader := bufio.NewReader(infile)
	
	var err error
	
	mapRe := regexp.MustCompile(`^\s*map (\S+), (\S+), \S+, \S+, \S+, \S+, \S+, \S+$`)
	
	mapAttrMap := make(map[string]MapAttributes)
	
	for ;; {
		mapAttr, ok := <-mapTilesetChan
		if ok {
			mapAttrMap[mapAttr.Name] = mapAttr
		} else {
			break
		}
	}
	
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
		res := mapRe.FindStringSubmatch(lineS)
		if len(res) >= 3 {
			if _, ok := mapAttrMap[res[1]]; !ok {
				continue
			}
			mapAttr := mapAttrMap[res[1]]
			mapAttr.Tileset.ConstName = res[2]
			tilesetReqChan <- mapAttr
			mapAttr = <-tilesetRespChan
			
			mapWidthHeightChan <- mapAttr
		}
	}
	
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
}

func ReadMapWidthHeight(infile io.ReadSeekCloser, mapWidthHeightChan <-chan MapAttributes) []MapAttributes {
	defer func() {
		_ = infile.Close()
	}()
	
	infile.Seek(0, io.SeekStart)
	lineReader := bufio.NewReader(infile)
	
	var err error
	
	mapConstRe := regexp.MustCompile(`^\s*map_const (\S+),\s*([0-9]+),\s*([0-9]+) ;\s*[0-9]+$`)
	
	mapAttrMap := make(map[string]MapAttributes)
	results := []MapAttributes{}
	
	for ;; {
		mapAttr, ok := <-mapWidthHeightChan
		if ok {
			mapAttrMap[mapAttr.IdName] = mapAttr
		} else {
			break
		}
	}
	
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
		res := mapConstRe.FindStringSubmatch(lineS)
		if len(res) >= 4 {
			if _, ok := mapAttrMap[res[1]]; !ok {
				continue
			}
			mapAttr := mapAttrMap[res[1]]
			width, err := strconv.ParseUint(res[2], 10, 8)
			if err != nil {
				panic(err)
			}
			height, err := strconv.ParseUint(res[3], 10, 8)
			if err != nil {
				panic(err)
			}
			mapAttr.Width = uint8(width)
			mapAttr.Height = uint8(height)
			results = append(results, mapAttr)
		}
	}
	
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
	
	return results
}
