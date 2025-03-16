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

