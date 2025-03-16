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

package decomp

import (
	"errors"
	// "fmt"
	"log"
)

type opType int

const (
	Fill opType = iota
	Or
	DeltaDec
	DCopy
	DXor
)

type Operation struct {
	T opType
	DestAddr uint16
	Mask uint8
	Value uint8
	SourceAddr uint16
}

type RecordedDecompression struct {
	Operations []Operation
}

type BitstreamReader struct {
	bitstream []byte
	currentByte uint8
	bytePosition int
	bitPosition int
}

func (b *BitstreamReader) readBit() (uint8, error) {
	if b.bitPosition == 0 {
		if b.bytePosition < len(b.bitstream) {
			b.currentByte = b.bitstream[b.bytePosition]
		} else {
			return 0, errors.New("premature end of stream")
		}
	}
	
	if b.bitPosition < 8 {
		bit := (b.currentByte & 0x80) >> 7
		b.currentByte <<= 1
		b.bitPosition++
		//log.Printf("previousByte = 0x%02x, bit = %d, b.currentByte = 0x%02x, b.bitPosition = %d", previousByte, bit, b.currentByte, b.bitPosition)
		return bit, nil
	} else {
		b.bitPosition = 0;
		b.bytePosition++
		return b.readBit()
	}
}

func (b *BitstreamReader) readBits(numBits int) (int, error) {
	var result int
	for i := 0; i < numBits; i++ {
		bit, err := b.readBit()
		if err != nil {
			return 0, err
		}
		result <<= 1
		result |= int(bit)
		//log.Printf("result = %b %d", result, result)
	}
	//log.Printf("---")
	return result, nil
}

func RecordDecompressSprite(rom []byte, spritePtr, baseDataWidth, baseDataHeight int) *RecordedDecompression {
	recording := RecordedDecompression{
		Operations: make([]Operation, 0),
	}
	
	spriteReader := BitstreamReader{
		bitstream: rom,
		bytePosition: spritePtr,
	}
	
	log.Println("Clearing buffers")
	recording.fillBuffer(1);
	recording.fillBuffer(2);
	
	widthTiles, err := spriteReader.readBits(4)
	handle(err)
	heightTiles, err := spriteReader.readBits(4)
	handle(err)
	log.Printf("Sprite size is %dx%d\n", widthTiles, heightTiles)
	
	if baseDataWidth < 0 {
		baseDataWidth = widthTiles
	}
	if baseDataHeight < 0 {
		baseDataHeight = heightTiles
	}
	
	firstBuffer, err := spriteReader.readBit()
	handle(err)
	firstBuffer++
	var secondBuffer uint8
	
	if firstBuffer == 1 {
		secondBuffer = 2
	} else {
		secondBuffer = 1
	}
	
	bufferOrder := []uint8{firstBuffer, secondBuffer}
	log.Printf("Starting with BP%d, then BP%d\n", firstBuffer, secondBuffer)
	
	var decodeMode uint8
	
	for i := 0; i < 2; i++ {
		if i == 1 {
			decodeMode, err = spriteReader.readBit()
			handle(err)
			if decodeMode == 1 {
				decodeMode <<= 1
				bit, err := spriteReader.readBit()
				handle(err)
				decodeMode |= bit
			}
		}
		
		log.Printf("Decompressing plane %d into BP%d...\n", i, int(bufferOrder[i]))
		recording.decompressPlane(&spriteReader, heightTiles, widthTiles, int(bufferOrder[i]))
	}
	
	log.Printf("Using decode mode %d\n", decodeMode)
	
	
	switch decodeMode {
	case 0:
		recording.deltaDecode(heightTiles, widthTiles, 1)
		recording.deltaDecode(heightTiles, widthTiles, 2)
	case 2:
		recording.deltaDecode(heightTiles, widthTiles, int(firstBuffer))
		recording.xorBuffers(heightTiles, widthTiles, int(firstBuffer), int(secondBuffer))
	case 3:
		recording.deltaDecode(heightTiles, widthTiles, int(secondBuffer))
		recording.deltaDecode(heightTiles, widthTiles, int(firstBuffer))
		recording.xorBuffers(heightTiles, widthTiles, int(firstBuffer), int(secondBuffer))
	}
	
	recording.copyAlignSpriteData(baseDataHeight, baseDataWidth)
	recording.interlaceBuffers()
	
	return &recording
}

func (r *RecordedDecompression) interlaceBuffers() {
	log.Println("Interlacing buffers...")
	for offset := 0x187; offset >= 0; offset-- {
		srcAddr2 := getBufferBaseAddr(1) + uint16(offset)
		srcAddr1 := getBufferBaseAddr(0) + uint16(offset)
		destAddr2 := getBufferBaseAddr(1) + uint16(offset) * 2 + 1
		destAddr1 := getBufferBaseAddr(1) + uint16(offset) * 2
		
		r.Operations = append(r.Operations, Operation{
			T: DCopy,
			DestAddr: destAddr2,
			Mask: 0xff,
			SourceAddr: srcAddr2,
		})
		r.Operations = append(r.Operations, Operation{
			T: DCopy,
			DestAddr: destAddr1,
			Mask: 0xff,
			SourceAddr: srcAddr1,
		})
	}
}

func (r *RecordedDecompression) copyAlignSpriteData(heightTiles, widthTiles int) {
	log.Printf("Copying/aligning sprite data with size %dx%d...\n", widthTiles, heightTiles)
	startOffset := (7 * ((8 - widthTiles) / 2)) & 0xff
	startOffset = (startOffset + (7 - heightTiles)) & 0xff
	startOffset = (8 * startOffset) & 0xff

	rowCountForProcessing := (heightTiles * 8)
	if rowCountForProcessing == 0 {
		rowCountForProcessing = 256
	}
	if widthTiles == 0 {
		widthTiles = 256
	}
	
	r.fillBuffer(0)
	
	for column := 0; column < widthTiles; column++ {
		for row := 0; row < rowCountForProcessing; row++ {
			destAddr := getBufferBaseAddr(0) + uint16(startOffset) + uint16(column * 7 * 8) + uint16(row)
			srcAddr := getBufferBaseAddr(1) + uint16(column * rowCountForProcessing) + uint16(row)
			r.Operations = append(r.Operations, Operation {
				T: DCopy,
				DestAddr: destAddr,
				Mask: 0xff,
				SourceAddr: srcAddr,
			})
		}
	}
	
	r.fillBuffer(1)
	
	for column := 0; column < widthTiles; column++ {
		for row := 0; row < rowCountForProcessing; row++ {
			destAddr := getBufferBaseAddr(1) + uint16(startOffset) + uint16(column * 7 * 8) + uint16(row)
			srcAddr := getBufferBaseAddr(2) + uint16(column * rowCountForProcessing) + uint16(row)
			r.Operations = append(r.Operations, Operation {
				T: DCopy,
				DestAddr: destAddr,
				Mask: 0xff,
				SourceAddr: srcAddr,
			})
		}
	}
}

func (r *RecordedDecompression) xorBuffers(heightTiles, widthTiles, firstBuffer, secondBuffer int) {
	log.Printf("Applying XOR from BP%d to BP%d\n", firstBuffer, secondBuffer)
	rowCount := uint16(heightTiles * 8)
	rowCountForProcessing := rowCount
	if rowCountForProcessing == 0 {
		rowCountForProcessing = 256
	}
	if widthTiles == 0 {
		widthTiles = 32
	}
	
	for column := uint16(0); column < uint16(widthTiles); column++ {
		for row := uint16(0); row < rowCountForProcessing; row++ {
			sourceAddr := getBufferBaseAddr(firstBuffer) + rowCount * column + row
			destAddr := getBufferBaseAddr(secondBuffer) + rowCount * column + row
			r.Operations = append(r.Operations, Operation{
				T: DXor,
				DestAddr: destAddr,
				Mask: 0xff,
				SourceAddr: sourceAddr,
			})
		}
	}
}

func (r *RecordedDecompression) deltaDecode(heightTiles, widthTiles, bufferIdx int) {
	log.Printf("Performing delta decode on BP%d\n", bufferIdx)
	rowCount := uint16(heightTiles * 8)
	rowCountForProcessing := rowCount
	if rowCountForProcessing == 0 {
		rowCountForProcessing = 256
	}
	if widthTiles == 0 {
		widthTiles = 32
	}
	
	for row := uint16(0); row < rowCountForProcessing; row++ {
		for column := uint16(0); column < uint16(widthTiles); column++ {
			addr := getBufferBaseAddr(bufferIdx) + rowCount * column + row
			var prevAddr uint16
			if column > 0 {
				prevAddr = getBufferBaseAddr(bufferIdx) + rowCount * (column - 1) + row
			} else {
				prevAddr = addr
			}
			
			var startValueMask uint8
			if column == 0 {
				startValueMask = 0
			} else {
				startValueMask = 1
			}
			r.Operations = append(r.Operations, Operation{
				T: DeltaDec,
				DestAddr: addr,
				Mask: 0xff,
				Value: startValueMask,
				SourceAddr: prevAddr,
			})
		}
	}
}

func (r *RecordedDecompression) decompressPlane(spriteReader *BitstreamReader, 
                                                heightTiles, widthTiles, bufferIdx int) {
	rowCount := uint16(heightTiles * 8)
	if heightTiles == 0 {
		rowCount = 256
	}
	// NOTE: each column is 2 pixels wide
	columnCount := uint16(widthTiles * 4)
	if widthTiles == 0 {
		columnCount = 128
	}
	totalOffset := rowCount * columnCount
	outputRowIdx := uint16(0)
	outputColumnIdx := uint16(0)
	outputOffset := uint16(0)
	
	currentMode, err := spriteReader.readBit()
	handle(err)
	if currentMode == 0 {
		readRLEPacket(spriteReader, &outputOffset, &outputRowIdx, &outputColumnIdx, rowCount)
		currentMode = 1
	}
	
	for ; outputOffset < totalOffset; {
		if currentMode == 0 {
			readRLEPacket(spriteReader, &outputOffset, &outputRowIdx, &outputColumnIdx, rowCount)
			currentMode = 1
		} else {
			code, err := spriteReader.readBits(2)
			handle(err)
			if code == 0 {
				currentMode = 0
				continue
			}
			r.Operations = append(r.Operations, Operation {
				T: Or,
				DestAddr: getBufferPixelPairAddr(bufferIdx, outputRowIdx, outputColumnIdx, rowCount),
				Mask: 0xc0 >> ((outputColumnIdx % 4) * 2),
				Value: uint8(code) << (6 - ((outputColumnIdx % 4) * 2)),
			})
			//log.Println(r.Operations[len(r.Operations) - 1], outputRowIdx, outputColumnIdx, rowCount)
			outputOffset++
			outputRowIdx, outputColumnIdx = recalcRowColumnIdx(outputOffset, rowCount)
		}
	}
}

func (r *RecordedDecompression) fillBuffer(bufIndex int) {
	log.Printf("- clearing BP%d...", bufIndex)
	var initAddr uint16
	initAddr = getBufferBaseAddr(bufIndex)
	var addr uint16
	for addr = initAddr; addr < initAddr + 0x188; addr++ {
		r.Operations = append(r.Operations,
			Operation {
				T: Fill,
				DestAddr: addr,
				Mask: 0xff,
				Value: 0,
			},
		)
	}
}

func readRLEPacket(reader *BitstreamReader,
                   outputOffset, rowIdx, columnIdx *uint16,
                   rowCount uint16) {

	offset, err := readExpGolombNumber(reader)
	handle(err)
	*outputOffset += uint16(offset)
	*rowIdx, *columnIdx = recalcRowColumnIdx(*outputOffset, rowCount)
}

func recalcRowColumnIdx(offset, rowCount uint16) (rowIdx, columnIdx uint16) {
	columnIdx = offset / rowCount
	rowIdx = offset % rowCount
	return
}

func readExpGolombNumber(reader *BitstreamReader) (int, error) {
	// log.Println("Reading Exp-Golomb number...")
	var numBitsToRead int
	var result int
	var bit uint8
	var err error
	// assume
	bit = 1
	result = 0
	for ; bit == 1; {
		bit, err = reader.readBit()
		if err != nil {
			return 0, err
		}
		result <<= 1
		result |= int(bit)
		numBitsToRead++
	}
	result += 1
	offset, err := reader.readBits(numBitsToRead)
	if err != nil {
		return 0, err
	}
	result += offset

	// fmt.Print("raw: ")
	// for i := 0; i < numBitsToRead - 1; i++ {
	// 	fmt.Print("1")
	// }
	// fmt.Print("0")
	// for i := numBitsToRead - 1; i >= 0; i-- {
	// 	if ((1 << i) & offset) != 0 {
	// 		fmt.Print("1")
	// 	} else {
	// 		fmt.Print("0")
	// 	}
	// }
	// fmt.Printf("  numbits: %d  result: %d\n", numBitsToRead, result)

	return result, nil
}

func getBufferBaseAddr(bufIndex int) uint16 {
	return 0xa000 + uint16(bufIndex) * 0x188
}

func getBufferPixelPairAddr(bufIndex int,
                            row, column, rowCount uint16) uint16 {
	baseAddr := getBufferBaseAddr(bufIndex)
	columnAddr := baseAddr + (column / 4) * rowCount
	pixelPairAddr := columnAddr + row
	return pixelPairAddr
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
