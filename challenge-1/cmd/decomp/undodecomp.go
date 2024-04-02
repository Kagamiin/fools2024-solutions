package decomp

import (
	"log"
)

func (r *RecordedDecompression) UndoRecording(destMemory *[]byte) (unknownBitMap *[]byte) {
	data := make([]byte, 65536)
	unknownBitMap = &data
	
	log.Println("Undoing recorded decompression journal over saved data...")

	for i := len(r.Operations) - 1; i >= 0; i-- {
		operation := r.Operations[i]
		switch(operation.T) {
		case Fill:
			operation.MarkFill(unknownBitMap)
		case Or:
			operation.MarkOr(unknownBitMap)
		case DCopy:
			operation.UndoDataCopy(destMemory)
		case DXor:
			operation.UndoDataXor(destMemory)
		case DeltaDec:
			operation.UndoDeltaDecode(destMemory)
		}
	}
	
	return
}

func (o Operation) MarkOr(unknownBitMap *[]byte) {
	(*unknownBitMap)[o.DestAddr] |= o.Value & o.Mask
}

func (o Operation) MarkFill(unknownBitMap *[]byte) {
	(*unknownBitMap)[o.DestAddr] |= o.Mask
}

func (o Operation) UndoDataCopy(destMemory *[]byte) {
	(*destMemory)[o.SourceAddr] = (*destMemory)[o.DestAddr] & o.Mask
}

func (o Operation) UndoDataXor(destMemory *[]byte) {
	(*destMemory)[o.DestAddr] = (((*destMemory)[o.SourceAddr] ^ (*destMemory)[o.DestAddr]) & o.Mask) | ((*destMemory)[o.DestAddr] & ^o.Mask)
}

func (o Operation) UndoDeltaDecode(destMemory *[]byte) {
	originalVal := (*destMemory)[o.DestAddr]
	valInPrevPosition := (*destMemory)[o.SourceAddr]
	(*destMemory)[o.DestAddr] = 0
	for bit := 0; bit < 8; bit++ {
		if bit < 7 {
			if ((originalVal & (1 << bit)) != 0 && (originalVal & (1 << (bit + 1))) == 0) ||
				((originalVal & (1 << bit)) == 0 && (originalVal & (1 << (bit + 1))) != 0) {
				
				(*destMemory)[o.DestAddr] |= (1 << bit)
			}
		} else if o.Value != 0 {
			if ((originalVal & (1 << bit)) != 0 && (valInPrevPosition & 1) == 0) ||
				((originalVal & (1 << bit)) == 0 && (valInPrevPosition & 1) != 0) {
				
				(*destMemory)[o.DestAddr] |= (1 << bit)
			}
		} else {
			if ((originalVal & (1 << bit)) != 0) {
				(*destMemory)[o.DestAddr] |= (1 << bit)
			}
		}
	}
}
