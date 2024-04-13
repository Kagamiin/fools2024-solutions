package main

import (
	"log"
	"os"

	"github.com/Kagamiin/fools2024-solutions/challenge-2/cmd/maps"
)

const (
	blocksAsmPath = "extern/pokecrystal/data/maps/blocks.asm"
	attributesAsmPath = "extern/pokecrystal/data/maps/attributes.asm"
	mapsAsmPath = "extern/pokecrystal/data/maps/maps.asm"
	mapConstantsAsmPath = "extern/pokecrystal/constants/map_constants.asm"
	tilesetConstantsAsmPath = "extern/pokecrystal/constants/tileset_constants.asm"
	dataTilesetsAsmPath = "extern/pokecrystal/data/tilesets.asm"
	gfxTilesetsAsmPath = "extern/pokecrystal/gfx/tilesets.asm"
)

func main() {
	mapAttrChan := make(chan maps.MapAttributes, 0)
	mapTilesetChan := make(chan maps.MapAttributes, 0)
	mapWidthHeightChan := make(chan maps.MapAttributes, 0)
	tilesetReqChan := make(chan maps.MapAttributes, 0)
	tilesetRespChan := make(chan maps.MapAttributes, 0)

	tilesetAttrChan := make(chan maps.TilesetAttributes, 0)
	tilesetMetatileChan := make(chan maps.TilesetAttributes, 0)
	
	blocksAsmFile, err := os.Open(blocksAsmPath)
	handle(err)
	attributesAsmFile, err := os.Open(attributesAsmPath)
	handle(err)
	mapsAsmFile, err := os.Open(mapsAsmPath)
	handle(err)
	mapConstantsAsmFile, err := os.Open(mapConstantsAsmPath)
	handle(err)
	tilesetConstantsAsmFile, err := os.Open(tilesetConstantsAsmPath)
	handle(err)
	dataTilesetsAsmFile, err := os.Open(dataTilesetsAsmPath)
	handle(err)
	gfxTilesetsAsmFile, err := os.Open(gfxTilesetsAsmPath)
	handle(err)
	
	
	go maps.ReadBlocksAsm(blocksAsmFile, mapAttrChan)
	go maps.ReadMapAttributes(attributesAsmFile, mapAttrChan, mapTilesetChan)
	go maps.ReadMapTileset(mapsAsmFile, mapTilesetChan, mapWidthHeightChan, tilesetReqChan, tilesetRespChan)
	go maps.ReadTilesets(tilesetConstantsAsmFile, dataTilesetsAsmFile, tilesetAttrChan)
	go maps.ReadMetatileFileNames(gfxTilesetsAsmFile, tilesetAttrChan, tilesetMetatileChan)
	go maps.ServeTilesetAttributes(tilesetMetatileChan, tilesetReqChan, tilesetRespChan)
	mapAttrList := maps.ReadMapWidthHeight(mapConstantsAsmFile, mapWidthHeightChan)
	
	maps.InitMaps(mapAttrList)
	
	for i := range mapAttrList {
		pm := maps.LoadBlockMap(mapAttrList[i])
		maps.LoadTileMap(&pm)
		searchPattern(&pm)
	}
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func searchPattern(pm *maps.PokéMap) {
	for top := 0; top < len(pm.TileMap) - 18; top += 2 {
		for left := 0; left < len(pm.TileMap[top]) - 20; left += 2 {
			cond1 := pm.TileMap[top][left + 11] == 0x05
			cond2 := pm.TileMap[top + 6][left + 7] == 0x23
			cond3 := pm.TileMap[top + 2][left + 3] == 0x02
			cond4 := pm.TileMap[top + 2][left + 4] == 0x04
			cond5 := pm.TileMap[top + 11][left + 12] == 0x01
			if cond1 && cond2 && cond3 && cond4 && cond5 {
				log.Printf("Found match in map %s at coordinates %d, %d", pm.Attr.Name, (left + 8 - 12) / 2, (top + 8 - 12) / 2)
			}
		}
	}
}

