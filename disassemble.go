package main

import (
	"errors"
)

func parse16BitValue(content []byte) int16 {
	lowByte := int16(content[0])
	highByte := int16(content[1])
	return (highByte << 8) | lowByte
}

func parseData(content []byte, wide bool) (int, int16) {
	parsedBytes := 0
	data := int16(0)
	if !wide {
		data = int16(int8(content[0]))
		parsedBytes = 1
	} else {
		data = parse16BitValue(content)
		parsedBytes = 2
	}
	return parsedBytes, data
}

func parseAddressCalculation(content []byte, mod byte, rm byte) (int, AddressCalculation) {
	currentByte := 0
	addressCalculationType := addressCalculationTable[mod][rm]

	displacement := int16(0)
	if mod == 0b01 {
		// Memory Mode, 8 bit displacement follows
		displacement = int16(int8(content[currentByte]))
		currentByte++
	} else if mod == 0b10 || (mod == 0b00 && rm == 0b110) {
		// Memory Mode, 16 bit displacement follows
		displacement = parse16BitValue(content[currentByte:])
		currentByte += 2
	}

	return currentByte, AddressCalculation{
		Type:         addressCalculationType,
		Displacement: displacement,
	}
}

func StringifyInstructions(instructions []Instruction) string {
	result := "bits 16\n"
	for _, instruction := range instructions {
		result += instruction.String()
	}
	return result
}

func Disassemble(content []byte) ([]Instruction, error) {
	instructions := make([]Instruction, 0)
	currentByte := 0
	for currentByte < len(content) {
		startByte := currentByte

		instructionType, err := InstructionTypeFromBytes(content[currentByte:])
		if err != nil {
			return instructions, err
		}

		b1 := content[currentByte]
		currentByte++

		if instructionType == IT_PushReg || instructionType == IT_PopReg || instructionType == IT_ExchangeRegWithAcc || instructionType == IT_IncReg || instructionType == IT_DecReg {
			reg := b1 & 0b111
			var src *DataLocation
			dst := &DataLocation{
				Type:         DL_Register,
				RegisterName: registerTable[1][reg],
			}
			if instructionType == IT_ExchangeRegWithAcc {
				src = dst
				dst = &DataLocation{
					Type:         DL_Register,
					RegisterName: AX,
				}
			}
			instructions = append(instructions, Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Destination: dst,
				Source:      src,
			})
			continue
		}

		if instructionType == IT_PushSegReg || instructionType == IT_PopSegReg {
			reg := (b1 >> 3) & 0b11
			instructions = append(instructions, Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Destination: &DataLocation{
					Type:         DL_Register,
					RegisterName: segmentRegisterTable[reg],
				},
			})
			continue
		}

		if instructionType.IsSingleByteInstruction() {
			w := b1 & 0b1
			instructions = append(instructions, Instruction{
				Type:        instructionType,
				SizeInBytes: 1,
				Wide:        w == 0b1,
			})
			continue
		}

		if instructionType.IsConditionalJump() {
			offset := int8(content[currentByte])
			currentByte++

			instruction := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Destination: &DataLocation{
					Type:          DL_Label,
					LabelPosition: int(offset + 2),
				},
			}
			instructions = append(instructions, instruction)
			continue
		}

		if instructionType == IT_ReturnWithinSegmentAddingImmediateToSP || instructionType == IT_ReturnIntersegmentAddingImmediateToSP {
			parsedBytes, data := parseData(content[currentByte:], true)
			currentByte += int(parsedBytes)
			instructions = append(instructions, Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Destination: &DataLocation{
					Type:           DL_Immediate,
					ImmediateValue: data,
					AvoidSizeInfo:  true,
				},
			})
			continue
		}

		if instructionType == IT_InterruptTypeSpecified {
			parsedBytes, data := parseData(content[currentByte:], false)
			currentByte += int(parsedBytes)
			instructions = append(instructions, Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Destination: &DataLocation{
					Type:           DL_Immediate,
					ImmediateValue: data,
					AvoidSizeInfo:  true,
				},
			})
			continue
		}

		if instructionType == IT_MovImToReg {
			w := (b1 >> 3) & 0b00000001
			reg := b1 & 0b00000111

			parsedBytes, data := parseData(content[currentByte:], w == 0b1)
			currentByte += parsedBytes

			src := DataLocation{
				Type:           DL_Immediate,
				ImmediateValue: data,
				Wide:           w == 0b1,
			}
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: registerTable[w][reg],
			}
			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		w := b1 & 0b1

		if instructionType.IsInOut() {
			var src DataLocation
			if instructionType == IT_InVariable || instructionType == IT_OutVariable {
				src = DataLocation{
					Type:         DL_Register,
					RegisterName: DX,
				}
			} else {
				parsedBytes, data := parseData(content[currentByte:], false)
				currentByte += parsedBytes
				src = DataLocation{
					Type:           DL_Immediate,
					ImmediateValue: data,
				}
			}
			dstRegisterName := AL
			if w == 0b1 {
				dstRegisterName = AX
			}
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: dstRegisterName,
			}
			if instructionType == IT_OutFixed || instructionType == IT_OutVariable {
				tmp := src
				src = dst
				dst = tmp
			}
			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		if instructionType == IT_MovMemToAcc {
			displacement := parse16BitValue(content[currentByte:])
			currentByte += 2

			src := DataLocation{
				Type: DL_Memory,
				AddressCalculation: AddressCalculation{
					Type:         ACT_DirectAddress,
					Displacement: displacement,
				},
				Wide: w == 0b1,
			}
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: AX,
			}
			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		if instructionType == IT_MovAccToMem {
			displacement := parse16BitValue(content[currentByte:])
			currentByte += 2

			src := DataLocation{
				Type:         DL_Register,
				RegisterName: AX,
			}
			dst := DataLocation{
				Type: DL_Memory,
				AddressCalculation: AddressCalculation{
					Type:         ACT_DirectAddress,
					Displacement: displacement,
				},
				Wide: w == 0b1,
			}
			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		if instructionType.IsImToAcc() {
			parsedBytes, data := parseData(content[currentByte:], w == 0b1)
			currentByte += parsedBytes

			src := DataLocation{
				Type:           DL_Immediate,
				ImmediateValue: data,
				Wide:           w == 0b1,
			}
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: registerTable[w][0],
			}
			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		b2 := content[currentByte]
		currentByte++

		mod := b2 >> 6
		reg := (b2 >> 3) & 0b111
		rm := b2 & 0b111

		if mod == 0b11 {
			if instructionType.IsRegMemWithRegToEither() {
				// Register Mode (no displacement)
				src := &DataLocation{
					Type:         DL_Register,
					RegisterName: registerTable[w][reg],
				}
				dst := &DataLocation{
					Type:         DL_Register,
					RegisterName: registerTable[w][rm],
				}
				if instructionType == IT_ExchangeRegMemWithReg {
					tmp := src
					src = dst
					dst = tmp
				}
				if instructionType.IsSingleOperandInstruction() {
					src = nil
				}

				if instructionType.IsShiftOrRotateInstruction() {
					v := (b1 >> 1) & 0b1
					if v == 0b0 {
						src = &DataLocation{Type: DL_Immediate, ImmediateValue: 1, AvoidSizeInfo: true}
					} else {
						src = &DataLocation{Type: DL_Register, RegisterName: CL}
					}
				}

				inst := Instruction{
					Type:        instructionType,
					SizeInBytes: currentByte - startByte,
					Source:      src,
					Destination: dst,
				}
				instructions = append(instructions, inst)
				continue
			}

			s := (b1 >> 1) & 0b1
			wide := s == 0b0 && w == 0b1
			parsedBytes, data := parseData(content[currentByte:], wide)
			currentByte += parsedBytes

			src := DataLocation{
				Type:           DL_Immediate,
				ImmediateValue: data,
				Wide:           w == 0b1,
			}
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: registerTable[w][rm],
			}
			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		parsedBytes, addressCalculation := parseAddressCalculation(content[currentByte:], mod, rm)
		currentByte += parsedBytes

		if instructionType.IsRegMemWithRegToEither() {
			src := &DataLocation{}
			dst := &DataLocation{}

			if instructionType.IsShiftOrRotateInstruction() {
				v := (b1 >> 1) & 0b1
				if v == 0b0 {
					src = &DataLocation{Type: DL_Immediate, ImmediateValue: 1, AvoidSizeInfo: true}
				} else {
					src = &DataLocation{Type: DL_Register, RegisterName: CL}
				}
				dst.Type = DL_Memory
				dst.AddressCalculation = addressCalculation
				dst.Wide = w == 0b1
			} else {
				d := (b1 >> 1) & 0b1
				if d == 0b1 || instructionType.AlwaysToRegister() {
					if instructionType == IT_LoadES {
						w = 0b1
					}
					dst.Type = DL_Register
					dst.RegisterName = registerTable[w][reg]
					src.Type = DL_Memory
					src.AddressCalculation = addressCalculation
					src.Wide = w == 0b1
					src.AvoidSizeInfo = instructionType == IT_LoadDS || instructionType == IT_LoadES
				} else {
					src.Type = DL_Register
					src.RegisterName = registerTable[w][reg]
					dst.Type = DL_Memory
					dst.AddressCalculation = addressCalculation
					dst.Wide = w == 0b1
				}
			}

			if instructionType.IsSingleOperandInstruction() {
				dst = src
				src = nil
			}

			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      src,
				Destination: dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		if instructionType.IsImToRegMem() {
			wide := w == 0b1
			if instructionType.HasSignExtension() {
				s := (b1 >> 1) & 0b1
				wide = wide && s == 0b0
			}
			parsedBytes, data := parseData(content[currentByte:], wide)
			currentByte += parsedBytes
			src := DataLocation{
				Type:           DL_Immediate,
				ImmediateValue: data,
				Wide:           w == 0b1,
			}

			var dst DataLocation
			if mod == 0b11 {
				dst = DataLocation{
					Type:         DL_Register,
					RegisterName: registerTable[w][reg],
				}
			} else {
				dst = DataLocation{
					Type:               DL_Memory,
					AddressCalculation: addressCalculation,
					Wide:               w == 0b1,
				}
			}

			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
			}
			instructions = append(instructions, inst)
			continue
		}

		if instructionType == IT_PushRegMem || instructionType == IT_PopRegMem {
			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Destination: &DataLocation{
					Type:               DL_Memory,
					AddressCalculation: addressCalculation,
					Wide:               true,
				},
			}
			instructions = append(instructions, inst)
			continue
		}

		if instructionType == IT_AsciiAdjustForMultiply || instructionType == IT_AsciiAdjustForDivide {
			// TODO the manual says that these instructions are actually 4 bytes, but it works like this
			instructions = append(instructions, Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
			})
			continue
		}

		return instructions, errors.New("instruction decode not implemented yet")
	}

	return instructions, nil
}
