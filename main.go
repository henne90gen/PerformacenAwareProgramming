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
	InvalidInstruction InstructionType = iota
	MovRegMemToFromReg
	MovImToRegMem
	MovImToReg
	MovMemToAcc
	MovAccToMem
	MovRegMemToSegReg
	MovSegRegToRegMem
)

type MemoryLocationType string

var registerTable = []RegisterName{
	AL, AX, CL, CX, DL, DX, BL, BX,
	AH, SP, CH, BP, DH, SI, BH, DI,
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
	if t > InvalidInstruction && t < MovSegRegToRegMem {
		return "mov"
	}

	return ""
}

func getInstructionType(b byte) (InstructionType, error) {
	if b>>2 == 0b100010 {
		// Register/memory to/from register
		return MovRegMemToFromReg, nil
	} else if b>>1 == 0b1100011 {
		// Immediate to register/memory
		return MovImToRegMem, nil
	} else if b>>4 == 0b1011 {
		// Immediate to register
		return MovImToReg, nil
	} else if b>>1 == 0b1010000 {
		// Memory to accumulator
		return MovMemToAcc, nil
	} else if b>>1 == 0b1010001 {
		// Accumulator to memory
		return MovAccToMem, nil
	} else if b == 0b10001110 {
		// Register/memory to segment register
		return MovRegMemToSegReg, nil
	} else if b == 0b10001100 {
		// Segment register to register/memory
		return MovSegRegToRegMem, nil
	}

	return InvalidInstruction, fmt.Errorf("opcode %08b not implemented yet", b)
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

		if instructionType == MovImToReg {
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

			destinationRegisterIndex := reg<<1 + w
			destinationRegisterName := registerTable[destinationRegisterIndex]
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
			destinationRegisterIndex := rm<<1 + w
			destinationRegisterName := registerTable[destinationRegisterIndex]
			sourceRegisterIndex := reg<<1 + w
			sourceRegisterName := registerTable[sourceRegisterIndex]
			result += fmt.Sprintf("%s %s, %s\n", instructionType.Name(), destinationRegisterName, sourceRegisterName)
		} else if mod == 0b01 {
			// Memory Mode, 8 bit displacement follows
			currentByte++
		} else if mod == 0b10 {
			// Memory Mode, 16 bit displacement follows
			currentByte += 2
		} else {
			// Memory Mode, no displacement follows (most of the time)
			if rm == 0b110 {
				return "", errors.New("(rm == 0b110) not implemented yet")
			}
		}
	}

	println(result)

	return result, nil
}

func main() {
	inputFiles := []string{"l_37", "l_38", "l_39"}
	for _, inputFile := range inputFiles {
		content, err := os.ReadFile(inputFile)
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
