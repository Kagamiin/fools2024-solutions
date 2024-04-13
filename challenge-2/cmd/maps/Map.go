package maps

import (
	"fmt"
	"io"
	"os"
)

type PokéMap struct {
	Attr MapAttributes
	Tileset [][]byte
	BlockMap [][]byte
	TileMap [][]byte
}

func readBlockMapFile(filepath string, width, height uint8) [][]byte {
	infile, err := os.Open(filepath)
	handle(err)
	
	defer func() {
		_ = infile.Close()
	}()
	
	result := make([][]byte, height)
	for i := range result {
		row := make([]byte, width)
		n, err := io.ReadFull(infile, row)
		handle(err)
		if n != int(width) {
			panic("io.ReadFull returned less bytes than expected")
		}
		result[i] = row
	}
	
	return result
}

func readMetatileFile(filepath string) [][]byte {
	infile, err := os.Open(filepath)
	handle(err)
	
	defer func() {
		_ = infile.Close()
	}()
	
	result := make([][]byte, 0)
	for ; err == nil; {
		row := make([]byte, 16)
		n, err := io.ReadFull(infile, row)
		if err != nil {
			break
		}
		if n != 16 {
			panic(fmt.Sprintf("io.ReadFull returned less bytes than expected (expected 16, got %d)", n))
		}
		result = append(result, row)
	}
	
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		panic(err)
	}
	
	return result
}

var mapsByName map[string]MapAttributes

func InitMaps(mapAttrs []MapAttributes) {
	if mapsByName == nil {
		mapsByName = make(map[string]MapAttributes)
	}
	for _, m := range mapAttrs {
		mapsByName[m.Name] = m
	}
}

func LoadBlockMap(mapAttr MapAttributes) PokéMap {
	result := PokéMap{
		Attr: mapAttr,
	}

	result.BlockMap = make([][]byte, int(mapAttr.Height) + 6)
	for i := range result.BlockMap {
		row := make([]byte, int(mapAttr.Width) + 6)
		for j := range row {
			row[j] = mapAttr.BorderBlock
		}
		result.BlockMap[i] = row
	}

	mapBlocks := readBlockMapFile(mapAttr.BlockFilePath, mapAttr.Width, mapAttr.Height)
	for i := range mapBlocks {
		for j := range mapBlocks[i] {
			result.BlockMap[i + 3][j + 3] = mapBlocks[i][j]
		}
	}
	
	loadNorthMapConnection(&result)
	loadSouthMapConnection(&result)
	loadWestMapConnection(&result)
	loadEastMapConnection(&result)

	return result
}

func LoadTileMap(pm *PokéMap) {
	pm.Tileset = readMetatileFile(pm.Attr.Tileset.MetatileFileName)

	result := make([][]byte, len(pm.BlockMap) * 4)
	for i := range result {
		row := make([]byte, len(pm.BlockMap[0]) * 4)
		result[i] = row
	}

	pm.TileMap = result

	for i := range pm.BlockMap {
		for j := range pm.BlockMap[i] {
			block := pm.BlockMap[i][j]
			for m := 0; m < 4; m++ {
				for n := 0; n < 4; n++ {
					tile := pm.Tileset[block][m * 4 + n]
					pm.TileMap[i * 4 + m][j * 4 + n] = tile
				}
			}
		}
	}
}

func loadNorthMapConnection(m *PokéMap) {
	if m.Attr.ConnectionNorth.Enabled {
		if destMapAttr, ok := mapsByName[m.Attr.ConnectionNorth.DestinationMapName]; ok {
			mapBlocks := readBlockMapFile(destMapAttr.BlockFilePath, destMapAttr.Width, destMapAttr.Height)
			for i := 0; i < 3; i++ {
				j := -(m.Attr.ConnectionNorth.BlockOffset + 3)
				for x := 0; x < int(m.Attr.Width) + 6; x++ {
					if j >= 0 && j < int8(destMapAttr.Width) {
						m.BlockMap[i][x] = mapBlocks[destMapAttr.Height - 3 + uint8(i)][j]
					}
					j++
				}
			}
		}
	}
}

func loadSouthMapConnection(m *PokéMap) {
	if m.Attr.ConnectionSouth.Enabled {
		if destMapAttr, ok := mapsByName[m.Attr.ConnectionSouth.DestinationMapName]; ok {
			mapBlocks := readBlockMapFile(destMapAttr.BlockFilePath, destMapAttr.Width, destMapAttr.Height)
			for i := 0; i < 3; i++ {
				j := -(m.Attr.ConnectionSouth.BlockOffset + 3)
				for x := 0; x < int(m.Attr.Width) + 6; x++ {
					if j >= 0 && j < int8(destMapAttr.Width) {
						m.BlockMap[m.Attr.Height + 3 + uint8(i)][x] = mapBlocks[i][j]
					}
					j++
				}
			}
		}
	}
}

func loadEastMapConnection(m *PokéMap) {
	if m.Attr.ConnectionEast.Enabled {
		if destMapAttr, ok := mapsByName[m.Attr.ConnectionEast.DestinationMapName]; ok {
			mapBlocks := readBlockMapFile(destMapAttr.BlockFilePath, destMapAttr.Width, destMapAttr.Height)
			i := -(m.Attr.ConnectionEast.BlockOffset + 3)
			for y := 0; y < int(m.Attr.Height) + 6; y++ {
				for j := 0; j < 3; j++ {
					if i >= 0 && i < int8(destMapAttr.Height) {
						m.BlockMap[y][m.Attr.Width + 3 + uint8(j)] = mapBlocks[i][j]
					}
				}
				i++
			}
		}
	}
}

func loadWestMapConnection(m *PokéMap) {
	if m.Attr.ConnectionWest.Enabled {
		if destMapAttr, ok := mapsByName[m.Attr.ConnectionWest.DestinationMapName]; ok {
			mapBlocks := readBlockMapFile(destMapAttr.BlockFilePath, destMapAttr.Width, destMapAttr.Height)
			i := -(m.Attr.ConnectionWest.BlockOffset + 3)
			for y := 0; y < int(m.Attr.Height) + 6; y++ {
				for j := 0; j < 3; j++ {
					if i >= 0 && i < int8(destMapAttr.Height) {
						m.BlockMap[y][j] = mapBlocks[i][destMapAttr.Width - 3 + uint8(j)]
					}
				}
				i++
			}
		}
	}
}
