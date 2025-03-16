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

package decomp_test

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/Kagamiin/fools2024-solutions/challenge-1/cmd/decomp"
)

func xorIdentityCheck(data []byte, ops []decomp.Operation) bool {
	originalData := make([]byte, len(data))
	copy(originalData, data)
	
	for _, op := range ops {
		if op.T == decomp.DXor {
			op.DoDataXor(&data)
		}
	}
	
	for i := len(ops) - 1; i >= 0; i-- {
		op := ops[i]
		if op.T == decomp.DXor {
			op.UndoDataXor(&data)
		}
	}
	
	if !bytes.Equal(data, originalData) {
		log.Println(originalData)
		log.Println(data)
		return false
	}
	return true
}

func deltaDecIdentityCheck(data []byte, ops []decomp.Operation) bool {
	originalData := make([]byte, len(data))
	copy(originalData, data)
	
	var integrator uint8
	
	for _, op := range ops {
		if op.T == decomp.DeltaDec {
			integrator = op.DoDeltaDecode(&data, integrator)
		}
	}
	
	for i := len(ops) - 1; i >= 0; i-- {
		op := ops[i]
		if op.T == decomp.DeltaDec {
			op.UndoDeltaDecode(&data)
		}
	}
	
	if !bytes.Equal(data, originalData) {
		for i := 0; i < len(data); {
			i2 := i;
			for j := 0; j < 16 && i < len(data); j++ {
				fmt.Printf("%08b ", originalData[i])
				i++
			}
			i = i2
			fmt.Print("\n")
			for j := 0; j < 16 && i < len(data); j++ {
				if data[i] != originalData[i] {
					fmt.Printf("\033[91m%08b\033[0m ", data[i])
				} else {
					fmt.Printf("%08b ", data[i])
				}
				i++
			}
			fmt.Print("\n\n")
		}
		log.Println(originalData)
		log.Println(data)
		return false
	}
	return true
}

func Test_XorIdentity(t *testing.T) {
	conf := quick.Config{MaxCount: 10000, MaxCountScale: 1.0, Values: func(res []reflect.Value, rng *rand.Rand) {
		type args struct {
			data []byte
			ops []decomp.Operation
		}
		for i := 0; i < len(res); i++ {
			data := make([]byte, rng.Intn(63) + 2)
			ops := make([]decomp.Operation, rng.Intn(127) + 1)
			
			_, err := rng.Read(data)
			if err != nil {
				t.Fatalf("Failed to generate random data for test: %v", err)
			}
			for j := 0; j < len(ops); j++ {
				destAddr := rng.Intn(len(data))
				sourceAddr := rng.Intn(len(data))
				for ; sourceAddr == destAddr; {
					sourceAddr = rng.Intn(len(data))
				}
				ops[j] = decomp.Operation{
					T: decomp.DXor,
					DestAddr: uint16(destAddr),
					SourceAddr: uint16(sourceAddr),
					Mask: 0xff,
				}
			}
			
			res[0] = reflect.ValueOf(data)
			res[1] = reflect.ValueOf(ops)
		}
	}}
	if err := quick.Check(xorIdentityCheck, &conf); err != nil {
		t.Error(err)
	}
}

func Test_DeltaDecodeIdentity(t *testing.T) {
	conf := quick.Config{MaxCount: 1000, MaxCountScale: 1.0, Values: func(res []reflect.Value, rng *rand.Rand) {
		type args struct {
			data []byte
			ops []decomp.Operation
		}
		for i := 0; i < len(res); i++ {
			data := make([]byte, rng.Intn(128) + 16)
			ops := make([]decomp.Operation, 0)
			
			_, err := rng.Read(data)
			if err != nil {
				t.Fatalf("Failed to generate random data for test: %v", err)
			}
			numRuns := 1 // rng.Intn(4) + 1
			for j := 0; j < numRuns; j++ {
				startAddr := rng.Intn(len(data) + 1)
				stride := 2 //rng.Intn(len(data) / 16) + 15
				destAddr := startAddr + stride
				runIndicator := uint8(0)
				for ; destAddr < len(data) - 1; {
					op := decomp.Operation{
						T: decomp.DeltaDec,
						DestAddr: uint16(destAddr),
						SourceAddr: uint16(startAddr),
						Mask: 0xff,
						Value: runIndicator,
					}
					
					ops = append(ops, op)
					destAddr += stride
					startAddr += stride
					runIndicator = 1
				}
			}
				
			res[0] = reflect.ValueOf(data)
			res[1] = reflect.ValueOf(ops)
		}
	}}
	if err := quick.Check(deltaDecIdentityCheck, &conf); err != nil {
		t.Error(err)
	}
}
