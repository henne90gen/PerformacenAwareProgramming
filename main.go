package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	{ACT_BX_SI, ACT_BX_DI, ACT_BP_SI, ACT_BP_DI, ACT_SI, ACT_DI, ACT_Invalid, ACT_BX},
	{ACT_BX_SI_D8, ACT_BX_DI_D8, ACT_BP_SI_D8, ACT_BP_DI_D8, ACT_SI_D8, ACT_DI_D8, ACT_BP_D8, ACT_BX_D8},
	{ACT_BX_SI_D16, ACT_BX_DI_D16, ACT_BP_SI_D16, ACT_BP_DI_D16, ACT_SI_D16, ACT_DI_D16, ACT_BP_D16, ACT_BX_D16},
}

type MemoryLocationType string

var registerTable = [][]RegisterName{
	{AL, CL, DL, BL, AH, CH, DH, BH},
	{AX, CX, DX, BX, SP, BP, SI, DI},
}

type Instruction struct {
	Type InstructionType
	Src  *MemoryLocation
	Dst  *MemoryLocation
}

type MemoryLocation struct {
	Type MemoryLocationType
}

func (t InstructionType) Name() string {
	if t > IT_Invalid && t < IT_MovSegRegToRegMem {
		return "mov"
	}

	return ""
}

func (a AddressCalculationType) String(displacement int16) string {
	switch a {
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
	case ACT_BX:
		return "[bx]"
	case ACT_BX_SI_D8:
		return fmt.Sprintf("[bx + si + %d]", displacement)
	case ACT_BX_DI_D8:
		return fmt.Sprintf("[bx + di + %d]", displacement)
	case ACT_BP_SI_D8:
		return fmt.Sprintf("[bp + si + %d]", displacement)
	case ACT_BP_DI_D8:
		return fmt.Sprintf("[bp + di + %d]", displacement)
	case ACT_SI_D8:
		return fmt.Sprintf("[si + %d]", displacement)
	case ACT_DI_D8:
		return fmt.Sprintf("[di + %d]", displacement)
	case ACT_BP_D8:
		return fmt.Sprintf("[bp + %d]", displacement)
	case ACT_BX_D8:
		return fmt.Sprintf("[bx + %d]", displacement)
	case ACT_BX_SI_D16:
		return fmt.Sprintf("[bx + si + %d]", displacement)
	case ACT_BX_DI_D16:
		return fmt.Sprintf("[bx + di + %d]", displacement)
	case ACT_BP_SI_D16:
		return fmt.Sprintf("[bp + si + %d]", displacement)
	case ACT_BP_DI_D16:
		return fmt.Sprintf("[bp + di + %d]", displacement)
	case ACT_SI_D16:
		return fmt.Sprintf("[si + %d]", displacement)
	case ACT_DI_D16:
		return fmt.Sprintf("[di + %d]", displacement)
	case ACT_BP_D16:
		return fmt.Sprintf("[bp + %d]", displacement)
	case ACT_BX_D16:
		return fmt.Sprintf("[bx + %d]", displacement)
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

			tmp := content[currentByte]
			currentByte++
			data := int16(0)
			if w == 0b0 {
				data = int16(int8(tmp))
			} else if w == 0b1 {
				highByte := int16(content[currentByte])
				data = int16(tmp) | (highByte << 8)
				currentByte++
			}

			destinationRegisterName := registerTable[w][reg]
			result += fmt.Sprintf("%s %s, %d\n", instructionType.Name(), destinationRegisterName, data)
			continue
		}

		w := b1 & 0b00000001

		b2 := content[currentByte]
		currentByte++

		mod := b2 >> 6
		reg := (b2 >> 3) & 0b00000111
		rm := b2 & 0b00000111

		if mod == 0b11 {
			// Register Mode (no displacement)
			destinationRegisterName := registerTable[w][rm]
			sourceRegisterName := registerTable[w][reg]
			result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), destinationRegisterName, sourceRegisterName)
			continue
		} else {
			// Memory Mode
			if rm == 0b110 && mod == 0b11 {
				print(result)
				return "", errors.New("(rm == 0b110) not implemented yet")
			}

			addressCalculationType := addressCalculationTable[mod][rm]

			displacement := int16(0)
			if mod == 0b01 {
				// Memory Mode, 8 bit displacement follows
				displacement = int16(int8(content[currentByte]))
				currentByte++
			} else if mod == 0b10 {
				// Memory Mode, 16 bit displacement follows
				tmp := content[currentByte]
				currentByte++

				highByte := int16(content[currentByte])
				displacement = int16(tmp) | (highByte << 8)
				currentByte++
			}

			d := (b1 >> 1) & 0b00000001
			if d == 0b1 {
				destinationRegisterName := registerTable[w][reg]
				result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), destinationRegisterName, addressCalculationType.String(displacement))
				continue
			} else {
				sourceRegisterName := registerTable[w][reg]
				result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), addressCalculationType.String(displacement), sourceRegisterName)
				continue
			}
		}
	}

	print(result)

	return result, nil
}

func main() {
	inputFiles := []string{"l_37.asm", "l_38.asm", "l_39.asm"}
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
