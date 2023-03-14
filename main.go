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

	IT_AddRegMemWithRegToEither
	IT_AddImToRegMem
	IT_AddImToAccumulator
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
)

var registerTable = [][]RegisterName{
	{AL, CL, DL, BL, AH, CH, DH, BH},
	{AX, CX, DX, BX, SP, BP, SI, DI},
}

type Instruction struct {
	Type InstructionType
	Src  *DataLocation
	Dst  *DataLocation
}

type DataLocation struct {
	Type DataLocationType

	RegisterName RegisterName

	AddressCalculation AddressCalculation

	ImmediateValue int16
	Wide           bool
}

func (t InstructionType) Name() string {
	if t > IT_Invalid && t < IT_MovSegRegToRegMem {
		return "mov"
	}

	if t >= IT_AddRegMemWithRegToEither && t < IT_AddImToAccumulator {
		return "add"
	}

	return ""
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
		return d.AddressCalculation.String()
	}
	return ""
}

func getInstructionType(b byte) (InstructionType, error) {
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
		return IT_AddImToRegMem, nil
	} else if b>>1 == 0b0000010 {
		// Immediate to accumulator
		return IT_AddImToAccumulator, nil
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
			return errors.New("byte %d does not match")
		}
	}

	return nil
}

func Parse16BitValue(content []byte) int16 {
	tmp := content[0]

	highByte := int16(content[1])
	return int16(tmp) | (highByte << 8)
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

func disassemble(content []byte) (string, error) {
	result := "bits 16\n"

	currentByte := 0
	for currentByte < len(content) {
		b1 := content[currentByte]
		currentByte++

		instructionType, err := getInstructionType(b1)
		if err != nil {
			println(result)
			return "", err
		}

		if instructionType == IT_MovImToReg {
			w := (b1 >> 3) & 0b00000001
			reg := b1 & 0b00000111

			data := int16(0)
			if w == 0b0 {
				data = int16(int8(content[currentByte]))
				currentByte++
			} else if w == 0b1 {
				data = Parse16BitValue(content[currentByte:])
				currentByte += 2
			}

			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: registerTable[w][reg],
			}
			result += fmt.Sprintf("%s %s, %d\n", instructionType.Name(), dst.String(), data)
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
			}
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: AX,
			}
			result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), dst.String(), src.String())
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
			}
			result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), dst.String(), src.String())
			continue
		}

		w := b1 & 0b00000001

		b2 := content[currentByte]
		currentByte++

		mod := b2 >> 6
		reg := (b2 >> 3) & 0b00000111
		rm := b2 & 0b00000111

		if mod == 0b11 {
			if instructionType == IT_MovRegMemToFromReg {
				// Register Mode (no displacement)
				src := DataLocation{
					Type:         DL_Register,
					RegisterName: registerTable[w][reg],
				}
				dst := DataLocation{
					Type:         DL_Register,
					RegisterName: registerTable[w][rm],
				}
				result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), dst.String(), src.String())
				continue
			}

			data := int16(0)
			if w == 0b0 {
				data = int16(int8(content[currentByte]))
				currentByte++
			} else if w == 0b1 {
				data = Parse16BitValue(content[currentByte:])
				currentByte += 2
			}

			src := DataLocation{
				Type:           DL_Immediate,
				ImmediateValue: data,
				Wide:           w == 0b1,
			}
			fmt.Printf("%b\n", data)
			dst := DataLocation{
				Type:         DL_Register,
				RegisterName: registerTable[w][rm],
			}
			result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), dst.String(), src.String())
			continue
		}

		parsedBytes, addressCalculation := ParseAddressCalculation(content[currentByte:], mod, rm)
		currentByte += parsedBytes

		if instructionType == IT_MovRegMemToFromReg || instructionType == IT_AddRegMemWithRegToEither {
			finalAddressCalculation := addressCalculation.String()

			d := (b1 >> 1) & 0b00000001
			if d == 0b1 {
				destinationRegisterName := registerTable[w][reg]
				result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), destinationRegisterName, finalAddressCalculation)
				continue
			} else {
				sourceRegisterName := registerTable[w][reg]
				result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), finalAddressCalculation, sourceRegisterName)
				continue
			}
		} else if instructionType == IT_MovImToRegMem || instructionType == IT_AddImToRegMem {
			data := int16(0)
			dataKind := ""
			if w == 0b0 {
				dataKind = "byte"
				data = int16(int8(content[currentByte]))
				currentByte++
			} else if w == 0b1 {
				dataKind = "word"
				data = Parse16BitValue(content[currentByte:])
				currentByte += 2
			}

			if mod == 0b11 {
				destinationRegisterName := registerTable[w][reg]
				result += fmt.Sprintf("%s %s, %s %d\n", instructionType.Name(), destinationRegisterName, dataKind, data)
				continue
			} else {
				result += fmt.Sprintf("%s %s, %s %d\n", instructionType.Name(), addressCalculation.String(), dataKind, data)
				continue
			}
		} else {
			return "", errors.New("not implemented yet")
		}
	}

	print(result)

	return result, nil
}

func main() {
	inputFiles := []string{
		"l_37.asm",
		"l_38.asm",
		"l_39.asm",
		"l_40.asm",
		"l_41.asm",
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
