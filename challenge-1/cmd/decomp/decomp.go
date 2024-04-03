package decomp

import (
	"log"
)

func (r *RecordedDecompression) ApplyRecording(destMemory *[]byte) {
	
	log.Println("Applying recorded decompression journal over saved data...")
	
	var integrator uint8
	for _, operation := range r.Operations {
		switch (operation.T) {
		case Fill:
			operation.DoFill(destMemory)
		case Or:
			operation.DoOr(destMemory)
		case DCopy:
			operation.DoDataCopy(destMemory)
		case DXor:
			operation.DoDataXor(destMemory)
		case DeltaDec:
			integrator = operation.DoDeltaDecode(destMemory, integrator)
		}
	}
}

func (o Operation) DoFill(destMemory *[]byte) {
	(*destMemory)[o.DestAddr] = o.Value & o.Mask
}

func (o Operation) DoOr(destMemory *[]byte) {
	(*destMemory)[o.DestAddr] |= o.Value & o.Mask
}

func (o Operation) DoDataCopy(destMemory *[]byte) {
	(*destMemory)[o.DestAddr] = (*destMemory)[o.SourceAddr] & o.Mask
}

func (o Operation) DoDataXor(destMemory *[]byte) {
	(*destMemory)[o.DestAddr] = (((*destMemory)[o.SourceAddr] ^ (*destMemory)[o.DestAddr]) & o.Mask) | ((*destMemory)[o.DestAddr] & ^o.Mask)
}

func (o Operation) DoDeltaDecode(destMemory *[]byte, state uint8) uint8 {
	originalVal := (*destMemory)[o.DestAddr]
	(*destMemory)[o.DestAddr] = 0
	for bit := 0; bit < 8; bit++ {
		if (originalVal & (0x80 >> bit)) != 0 {
			if state != 0 {
				state = 0
			} else {
				state = 1
			}
		}
		if state != 0 {
			(*destMemory)[o.DestAddr] |= 0x80 >> bit
		}
	}
	
	return state
}

