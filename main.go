package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type RegisterName string

const (
	AL RegisterName = "al"
	CL RegisterName = "cl"
	DL RegisterName = "dl"
	BL RegisterName = "bl"
	AH RegisterName = "ah"
	CH RegisterName = "ch"
	DH RegisterName = "dh"
	BH RegisterName = "bh"
	AX RegisterName = "ax"
	CX RegisterName = "cx"
	DX RegisterName = "dx"
	BX RegisterName = "bx"
	SP RegisterName = "sp"
	BP RegisterName = "bp"
	SI RegisterName = "si"
	DI RegisterName = "di"

	CS RegisterName = "cs"
	DS RegisterName = "ds"
	ES RegisterName = "es"
	SS RegisterName = "ss"
)

type InstructionType int

const (
	IT_Invalid InstructionType = iota
	IT_MovRegMemToFromReg
	IT_MovImToRegMem
	IT_MovImToReg
	IT_MovMemToAcc
	IT_MovAccToMem
	IT_MovRegMemToSegReg
	IT_MovSegRegToRegMem

	IT_PushRegMem
	IT_PushReg
	IT_PushSegReg

	IT_PopRegMem
	IT_PopReg
	IT_PopSegReg

	IT_ExchangeRegMemWithReg
	IT_ExchangeRegWithAcc

	IT_InFixed
	IT_InVariable
	IT_OutFixed
	IT_OutVariable

	IT_XLAT
	IT_LoadEA
	IT_LoadDS
	IT_LoadES
	IT_LoadAHWithFlags
	IT_StoreAHWithFlags
	IT_PushFlags
	IT_PopFlags

	IT_AddRegMemWithRegToEither
	IT_AddImToRegMem
	IT_AddImToAcc

	IT_AddWithCarryRegMemWithRegToEither
	IT_AddWithCarryImToRegMem
	IT_AddWithCarryImToAcc

	IT_SubRegMemWithRegToEither
	IT_SubImToRegMem
	IT_SubImFromAcc

	IT_SubWithBorrowRegMemWithRegToEither
	IT_SubWithBorrowImToRegMem
	IT_SubWithBorrowImFromAcc

	IT_CmpRegMemAndReg
	IT_CmpImWithRegMem
	IT_CmpImWithAcc

	IT_JE
	IT_JNE
	IT_JL
	IT_JLE
	IT_JB
	IT_JBE
	IT_JP
	IT_JO
	IT_JS
	IT_JNL
	IT_JNLE
	IT_JNB
	IT_JNBE
	IT_JNP
	IT_JNO
	IT_JNS
	IT_LOOP
	IT_LOOPZ
	IT_LOOPNZ
	IT_JCXZ
)

type AddressCalculationType int

const (
	ACT_Invalid AddressCalculationType = iota
	ACT_BX_SI
	ACT_BX_DI
	ACT_BP_SI
	ACT_BP_DI
	ACT_SI
	ACT_DI
	ACT_DirectAddress
	ACT_BX
	ACT_BX_SI_D8
	ACT_BX_DI_D8
	ACT_BP_SI_D8
	ACT_BP_DI_D8
	ACT_SI_D8
	ACT_DI_D8
	ACT_BP_D8
	ACT_BX_D8
	ACT_BX_SI_D16
	ACT_BX_DI_D16
	ACT_BP_SI_D16
	ACT_BP_DI_D16
	ACT_SI_D16
	ACT_DI_D16
	ACT_BP_D16
	ACT_BX_D16
)

var addressCalculationTable = [][]AddressCalculationType{
	{ACT_BX_SI, ACT_BX_DI, ACT_BP_SI, ACT_BP_DI, ACT_SI, ACT_DI, ACT_DirectAddress, ACT_BX},
	{ACT_BX_SI_D8, ACT_BX_DI_D8, ACT_BP_SI_D8, ACT_BP_DI_D8, ACT_SI_D8, ACT_DI_D8, ACT_BP_D8, ACT_BX_D8},
	{ACT_BX_SI_D16, ACT_BX_DI_D16, ACT_BP_SI_D16, ACT_BP_DI_D16, ACT_SI_D16, ACT_DI_D16, ACT_BP_D16, ACT_BX_D16},
}

type AddressCalculation struct {
	Type         AddressCalculationType
	Displacement int16
}

type DataLocationType int

const (
	DL_Invalid DataLocationType = iota
	DL_Register
	DL_Memory
	DL_Immediate
	DL_Label
)

var segmentRegisterTable = []RegisterName{ES, CS, SS, DS}
var registerTable = [][]RegisterName{
	{AL, CL, DL, BL, AH, CH, DH, BH},
	{AX, CX, DX, BX, SP, BP, SI, DI},
}

type Instruction struct {
	Type        InstructionType
	SizeInBytes int
	Destination *DataLocation
	Source      *DataLocation
}

type Label struct {
	PositionInBytes int
}

type DataLocation struct {
	Type DataLocationType

	RegisterName RegisterName

	AddressCalculation AddressCalculation

	ImmediateValue int16
	Wide           bool

	LabelPosition int

	AvoidSizeInfo bool
}

func (t InstructionType) Name() string {
	if t > IT_Invalid && t < IT_MovSegRegToRegMem {
		return "mov"
	}

	if t >= IT_PushRegMem && t <= IT_PushSegReg {
		return "push"
	}

	if t >= IT_PopRegMem && t <= IT_PopSegReg {
		return "pop"
	}

	if t >= IT_ExchangeRegMemWithReg && t <= IT_ExchangeRegWithAcc {
		return "xchg"
	}

	if t >= IT_InFixed && t <= IT_InVariable {
		return "in"
	}

	if t >= IT_OutFixed && t <= IT_OutVariable {
		return "out"
	}

	if t == IT_XLAT {
		return "xlat"
	}

	if t == IT_LoadEA {
		return "lea"
	}

	if t == IT_LoadDS {
		return "lds"
	}

	if t == IT_LoadES {
		return "les"
	}

	if t == IT_LoadAHWithFlags {
		return "lahf"
	}

	if t == IT_StoreAHWithFlags {
		return "sahf"
	}

	if t == IT_PushFlags {
		return "pushf"
	}

	if t == IT_PopFlags {
		return "popf"
	}

	if t >= IT_AddRegMemWithRegToEither && t <= IT_AddImToAcc {
		return "add"
	}

	if t >= IT_AddWithCarryRegMemWithRegToEither && t <= IT_AddWithCarryImToAcc {
		return "adc"
	}

	if t >= IT_SubRegMemWithRegToEither && t <= IT_SubImFromAcc {
		return "sub"
	}

	if t >= IT_CmpRegMemAndReg && t <= IT_CmpImWithAcc {
		return "cmp"
	}

	if t == IT_JE {
		return "je"
	}

	if t == IT_JNE {
		return "jne"
	}

	if t == IT_JL {
		return "jl"
	}

	if t == IT_JLE {
		return "jle"
	}

	if t == IT_JB {
		return "jb"
	}

	if t == IT_JBE {
		return "jbe"
	}

	if t == IT_JP {
		return "jp"
	}

	if t == IT_JO {
		return "jo"
	}

	if t == IT_JS {
		return "js"
	}

	if t == IT_JNL {
		return "jnl"
	}

	if t == IT_JNLE {
		return "jnle"
	}

	if t == IT_JNB {
		return "jnb"
	}

	if t == IT_JNBE {
		return "jnbe"
	}

	if t == IT_JNP {
		return "jnp"
	}

	if t == IT_JNO {
		return "jno"
	}

	if t == IT_JNS {
		return "jns"
	}

	if t == IT_LOOP {
		return "loop"
	}

	if t == IT_LOOPZ {
		return "loopz"
	}

	if t == IT_LOOPNZ {
		return "loopnz"
	}

	if t == IT_JCXZ {
		return "jcxz"
	}

	return ""
}

func (t InstructionType) IsImToAcc() bool {
	return t == IT_AddImToAcc || t == IT_SubImFromAcc || t == IT_CmpImWithAcc
}

func (t InstructionType) IsRegMemWithRegToEither() bool {
	return t == IT_MovRegMemToFromReg || t == IT_AddRegMemWithRegToEither || t == IT_AddWithCarryRegMemWithRegToEither || t == IT_SubRegMemWithRegToEither || t == IT_CmpRegMemAndReg || t == IT_ExchangeRegMemWithReg || t == IT_LoadEA || t == IT_LoadDS || t == IT_LoadES
}

func (t InstructionType) IsImToRegMem() bool {
	return t == IT_MovImToRegMem || t == IT_AddImToRegMem || t == IT_AddWithCarryImToRegMem || t == IT_SubImToRegMem || t == IT_CmpImWithRegMem
}

func (t InstructionType) HasSignExtension() bool {
	return t == IT_AddImToRegMem || t == IT_AddWithCarryImToRegMem || t == IT_SubImToRegMem || t == IT_CmpImWithRegMem
}

func (t InstructionType) IsJump() bool {
	return t == IT_JE || t == IT_JNE || t == IT_JL || t == IT_JLE || t == IT_JB || t == IT_JBE || t == IT_JP || t == IT_JO || t == IT_JS || t == IT_JNL || t == IT_JNLE || t == IT_JNB || t == IT_JNBE || t == IT_JNP || t == IT_JNO || t == IT_JNS || t == IT_LOOP || t == IT_LOOPZ || t == IT_LOOPNZ || t == IT_JCXZ
}

func (t InstructionType) IsInOut() bool {
	return t == IT_InFixed || t == IT_InVariable || t == IT_OutFixed || t == IT_OutVariable
}

func (t InstructionType) AlwaysToRegister() bool {
	return t == IT_ExchangeRegMemWithReg || t == IT_LoadEA || t == IT_LoadDS || t == IT_LoadES
}

func (t InstructionType) IsSingleByteInstruction() bool {
	return t == IT_XLAT || t == IT_LoadAHWithFlags || t == IT_StoreAHWithFlags || t == IT_PushFlags || t == IT_PopFlags
}

func (a AddressCalculation) String() string {
	switch a.Type {
	case ACT_BX_SI:
		return "[bx + si]"
	case ACT_BX_DI:
		return "[bx + di]"
	case ACT_BP_SI:
		return "[bp + si]"
	case ACT_BP_DI:
		return "[bp + di]"
	case ACT_SI:
		return "[si]"
	case ACT_DI:
		return "[di]"
	case ACT_DirectAddress:
		return fmt.Sprintf("[%d]", a.Displacement)
	case ACT_BX:
		return "[bx]"
	case ACT_BX_SI_D8:
		return fmt.Sprintf("[bx + si + %d]", a.Displacement)
	case ACT_BX_DI_D8:
		return fmt.Sprintf("[bx + di + %d]", a.Displacement)
	case ACT_BP_SI_D8:
		return fmt.Sprintf("[bp + si + %d]", a.Displacement)
	case ACT_BP_DI_D8:
		return fmt.Sprintf("[bp + di + %d]", a.Displacement)
	case ACT_SI_D8:
		return fmt.Sprintf("[si + %d]", a.Displacement)
	case ACT_DI_D8:
		return fmt.Sprintf("[di + %d]", a.Displacement)
	case ACT_BP_D8:
		return fmt.Sprintf("[bp + %d]", a.Displacement)
	case ACT_BX_D8:
		return fmt.Sprintf("[bx + %d]", a.Displacement)
	case ACT_BX_SI_D16:
		return fmt.Sprintf("[bx + si + %d]", a.Displacement)
	case ACT_BX_DI_D16:
		return fmt.Sprintf("[bx + di + %d]", a.Displacement)
	case ACT_BP_SI_D16:
		return fmt.Sprintf("[bp + si + %d]", a.Displacement)
	case ACT_BP_DI_D16:
		return fmt.Sprintf("[bp + di + %d]", a.Displacement)
	case ACT_SI_D16:
		return fmt.Sprintf("[si + %d]", a.Displacement)
	case ACT_DI_D16:
		return fmt.Sprintf("[di + %d]", a.Displacement)
	case ACT_BP_D16:
		return fmt.Sprintf("[bp + %d]", a.Displacement)
	case ACT_BX_D16:
		return fmt.Sprintf("[bx + %d]", a.Displacement)
	}
	return ""
}

func (d DataLocation) String() string {
	switch d.Type {
	case DL_Register:
		return string(d.RegisterName)
	case DL_Immediate:
		result := ""
		if !d.Wide {
			result += "byte "
		} else {
			result += "word "
		}
		return result + strconv.Itoa(int(d.ImmediateValue))
	case DL_Memory:
		if d.AvoidSizeInfo {
			return d.AddressCalculation.String()
		}
		result := ""
		if !d.Wide {
			result += "byte "
		} else {
			result += "word "
		}
		return result + d.AddressCalculation.String()
	case DL_Label:
		return fmt.Sprintf("label_%d", d.LabelPosition)
	}

	panic("unknown data location")
}

func (i Instruction) String() string {
	if i.Source == nil {
		if i.Destination == nil {
			return fmt.Sprintf("%s\n", i.Type.Name())
		}

		return fmt.Sprintf("%s %s\n", i.Type.Name(), i.Destination.String())
	}

	return fmt.Sprintf("%s %s, %s\n", i.Type.Name(), i.Destination.String(), i.Source.String())
}

func getInstructionType(content []byte) (InstructionType, error) {
	b := content[0]
	if b>>2 == 0b100010 {
		// Register/memory to/from register
		return IT_MovRegMemToFromReg, nil
	} else if b>>1 == 0b1100011 {
		// Immediate to register/memory
		return IT_MovImToRegMem, nil
	} else if b>>4 == 0b1011 {
		// Immediate to register
		return IT_MovImToReg, nil
	} else if b>>1 == 0b1010000 {
		// Memory to accumulator
		return IT_MovMemToAcc, nil
	} else if b>>1 == 0b1010001 {
		// Accumulator to memory
		return IT_MovAccToMem, nil
	} else if b == 0b10001110 {
		// Register/memory to segment register
		return IT_MovRegMemToSegReg, nil
	} else if b == 0b10001100 {
		// Segment register to register/memory
		return IT_MovSegRegToRegMem, nil
	}

	if b>>2 == 0b000000 {
		// Reg/memory with register to either
		return IT_AddRegMemWithRegToEither, nil
	} else if b>>2 == 0b100000 {
		// Immediate to register/memory
		b2 := content[1]
		reg := (b2 >> 3) & 0b111
		if reg == 0b000 {
			return IT_AddImToRegMem, nil
		}
		if reg == 0b010 {
			return IT_AddWithCarryImToRegMem, nil
		}
		if reg == 0b101 {
			return IT_SubImToRegMem, nil
		}
		if reg == 0b111 {
			return IT_CmpImWithRegMem, nil
		}
	} else if b>>1 == 0b0000010 {
		// Immediate to accumulator
		return IT_AddImToAcc, nil
	}

	if b>>2 == 0b001010 {
		// Register/memory with register to either
		return IT_SubRegMemWithRegToEither, nil
	} else if b>>1 == 0b0010110 {
		// Immediate from accumulator
		return IT_SubImFromAcc, nil
	}

	if b>>2 == 0b001110 {
		// Register/memory and register
		return IT_CmpRegMemAndReg, nil
	} else if b>>1 == 0b0011110 {
		// Immediate with accumulator
		return IT_CmpImWithAcc, nil
	}

	jumpInstructionsTable := []InstructionType{
		IT_JO, IT_JNO, IT_JB, IT_JNB, // 0000 - 0011
		IT_JE, IT_JNE, IT_JBE, IT_JNBE, // 0100 - 0111
		IT_JS, IT_JNS, IT_JP, IT_JNP, // 1000 - 1011
		IT_JL, IT_JNL, IT_JLE, IT_JNLE, // 1100 - 1111
	}
	if (b >> 4) == 0b0111 {
		return jumpInstructionsTable[b&0b1111], nil
	}

	if b == 0b11100010 {
		return IT_LOOP, nil
	}
	if b == 0b11100001 {
		return IT_LOOPZ, nil
	}
	if b == 0b11100000 {
		return IT_LOOPNZ, nil
	}
	if b == 0b11100011 {
		return IT_JCXZ, nil
	}

	if b == 0b11111111 {
		return IT_PushRegMem, nil
	}

	if b&0b11111000 == 0b01010000 {
		return IT_PushReg, nil
	}

	if b&0b11100111 == 0b00000110 {
		return IT_PushSegReg, nil
	}

	if b == 0b10001111 {
		return IT_PopRegMem, nil
	}

	if b&0b11111000 == 0b01011000 {
		return IT_PopReg, nil
	}

	if b&0b11100111 == 0b00000111 {
		return IT_PopSegReg, nil
	}

	if (b >> 1) == 0b1000011 {
		return IT_ExchangeRegMemWithReg, nil
	}

	if (b >> 3) == 0b10010 {
		return IT_ExchangeRegWithAcc, nil
	}

	if (b >> 1) == 0b1110010 {
		return IT_InFixed, nil
	}

	if (b >> 1) == 0b1110110 {
		return IT_InVariable, nil
	}

	if (b >> 1) == 0b1110011 {
		return IT_OutFixed, nil
	}

	if (b >> 1) == 0b1110111 {
		return IT_OutVariable, nil
	}

	if b == 0b11010111 {
		return IT_XLAT, nil
	}

	if b == 0b10001101 {
		return IT_LoadEA, nil
	}

	if b == 0b11000101 {
		return IT_LoadDS, nil
	}

	if b == 0b11000100 {
		return IT_LoadES, nil
	}

	if b == 0b10011111 {
		return IT_LoadAHWithFlags, nil
	}

	if b == 0b10011110 {
		return IT_StoreAHWithFlags, nil
	}

	if b == 0b10011100 {
		return IT_PushFlags, nil
	}

	if b == 0b10011101 {
		return IT_PopFlags, nil
	}

	if (b >> 2) == 0b000100 {
		return IT_AddWithCarryRegMemWithRegToEither, nil
	}

	if (b >> 1) == 0b0001010 {
		return IT_AddWithCarryImToAcc, nil
	}

	return IT_Invalid, fmt.Errorf("opcode %08b not implemented yet", b)
}

func assembleAndCompare(inputFileName string, inputFileContent []byte, result []byte) error {
	tmpFile, err := os.CreateTemp(".", inputFileName+"-*.asm")
	if err != nil {
		return err
	}

	_, err = tmpFile.Write(result)
	if err != nil {
		return err
	}

	err = tmpFile.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command("nasm", tmpFile.Name())
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		println(stdout.String())
		println(stderr.String())
		return err
	}

	err = os.Remove(tmpFile.Name())
	if err != nil {
		return err
	}

	trimmedTmpFile := strings.TrimSuffix(tmpFile.Name(), ".asm")
	assembled, err := os.ReadFile(trimmedTmpFile)
	if err != nil {
		return err
	}

	err = os.Remove(trimmedTmpFile)
	if err != nil {
		return err
	}

	if len(assembled) != len(inputFileContent) {
		return errors.New("length of assembled result does not match length of input")
	}
	for i, b := range assembled {
		if b != inputFileContent[i] {
			return fmt.Errorf("byte does not match, expected %08b but got %08b", inputFileContent[i], b)
		}
	}

	return nil
}

func Parse16BitValue(content []byte) int16 {
	tmp := content[0]

	highByte := int16(content[1])
	return int16(tmp) | (highByte << 8)
}

func ParseData(content []byte, wide bool) (int, int16) {
	parsedBytes := 0
	data := int16(0)
	if !wide {
		data = int16(int8(content[0]))
		parsedBytes = 1
	} else {
		data = Parse16BitValue(content[0:])
		parsedBytes = 2
	}
	return parsedBytes, data
}

func ParseAddressCalculation(content []byte, mod byte, rm byte) (int, AddressCalculation) {
	currentByte := 0
	addressCalculationType := addressCalculationTable[mod][rm]

	displacement := int16(0)
	if mod == 0b01 {
		// Memory Mode, 8 bit displacement follows
		displacement = int16(int8(content[currentByte]))
		currentByte++
	} else if mod == 0b10 || (mod == 0b00 && rm == 0b110) {
		// Memory Mode, 16 bit displacement follows
		displacement = Parse16BitValue(content[currentByte:])
		currentByte += 2
	}

	return currentByte, AddressCalculation{
		Type:         addressCalculationType,
		Displacement: displacement,
	}
}

func insert(a []Label, index int, value Label) []Label {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func insertLabel(labels []Label, position int) []Label {
	println("Inserting label", position)
	if len(labels) == 0 {
		return append(labels, Label{PositionInBytes: position})
	}

	for i, label := range labels {
		if label.PositionInBytes == position {
			return labels
		}

		if label.PositionInBytes < position {
			continue
		}

		return insert(labels, i, Label{PositionInBytes: position})
	}

	return append(labels, Label{PositionInBytes: position})
}

func stringifyInstructions(instructions []Instruction, labels []Label) string {
	result := "bits 16\n"
	currentByte := 0
	currentLabel := 0
	for _, instruction := range instructions {
		result += instruction.String()
		currentByte += instruction.SizeInBytes
		if len(labels) > currentLabel && currentByte+1 > labels[currentLabel].PositionInBytes {
			result += fmt.Sprintf("label_%d:\n", labels[currentLabel].PositionInBytes)
			currentLabel++
		}
	}
	return result
}

func disassemble(content []byte) (string, error) {
	instructions := make([]Instruction, 0)
	labels := make([]Label, 0)
	currentByte := 0
	for currentByte < len(content) {
		startByte := currentByte

		instructionType, err := getInstructionType(content[currentByte:])
		if err != nil {
			result := stringifyInstructions(instructions, labels)
			println(result)
			return "", err
		}

		b1 := content[currentByte]
		currentByte++

		if instructionType == IT_PushReg || instructionType == IT_PopReg || instructionType == IT_ExchangeRegWithAcc {
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
			println(reg)
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
			instructions = append(instructions, Instruction{
				Type:        instructionType,
				SizeInBytes: 1,
			})
			continue
		}

		if instructionType.IsJump() {
			offset := int8(content[currentByte])
			currentByte++

			labelPosition := currentByte + int(offset)
			instruction := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Destination: &DataLocation{
					Type:          DL_Label,
					LabelPosition: labelPosition,
				},
			}
			instructions = append(instructions, instruction)
			labels = insertLabel(labels, labelPosition)
			continue
		}

		if instructionType == IT_MovImToReg {
			w := (b1 >> 3) & 0b00000001
			reg := b1 & 0b00000111

			parsedBytes, data := ParseData(content[currentByte:], w == 0b1)
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
				parsedBytes, data := ParseData(content[currentByte:], false)
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

		if instructionType == IT_InVariable {
			src := DataLocation{
				Type:         DL_Register,
				RegisterName: DX,
			}
			dstRegisterName := AL
			if w == 0b1 {
				dstRegisterName = AX
			}
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: dstRegisterName,
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
			displacement := Parse16BitValue(content[currentByte:])
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
			displacement := Parse16BitValue(content[currentByte:])
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
			parsedBytes, data := ParseData(content[currentByte:], w == 0b1)
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
				src := DataLocation{
					Type:         DL_Register,
					RegisterName: registerTable[w][reg],
				}
				dst := DataLocation{
					Type:         DL_Register,
					RegisterName: registerTable[w][rm],
				}
				if instructionType == IT_ExchangeRegMemWithReg {
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

			s := (b1 >> 1) & 0b1
			wide := s == 0b0 && w == 0b1
			parsedBytes, data := ParseData(content[currentByte:], wide)
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

		parsedBytes, addressCalculation := ParseAddressCalculation(content[currentByte:], mod, rm)
		currentByte += parsedBytes

		if instructionType.IsRegMemWithRegToEither() {
			src := DataLocation{}
			dst := DataLocation{}

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

			inst := Instruction{
				Type:        instructionType,
				SizeInBytes: currentByte - startByte,
				Source:      &src,
				Destination: &dst,
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
			parsedBytes, data := ParseData(content[currentByte:], wide)
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

		result := stringifyInstructions(instructions, labels)
		print(result)
		return "", errors.New("instruction decode not implemented yet")
	}

	result := stringifyInstructions(instructions, labels)
	print(result)

	return result, nil
}

func main() {
	inputFiles := []string{
		"test.asm",
		"l_37.asm",
		"l_38.asm",
		"l_39.asm",
		"l_40.asm",
		"l_41.asm",
		"l_42.asm",
	}
	for _, inputFile := range inputFiles {
		cmd := exec.Command("nasm", inputFile)
		stdout := new(strings.Builder)
		stderr := new(strings.Builder)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
		if err != nil {
			println(stdout.String())
			println(stderr.String())
			panic(err)
		}

		assembledInputFile := strings.TrimSuffix(inputFile, ".asm")
		content, err := os.ReadFile(assembledInputFile)
		if err != nil {
			panic(err)
		}

		err = os.Remove(assembledInputFile)
		if err != nil {
			panic(err)
		}

		result, err := disassemble(content)
		if err != nil {
			panic(err)
		}

		err = assembleAndCompare(inputFile, content, []byte(result))
		if err != nil {
			panic(err)
		}

		println("Success - " + inputFile)
	}
}
