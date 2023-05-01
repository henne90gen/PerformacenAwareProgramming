package main

import "fmt"

type Context struct {
	Registers [24]byte
	Memory    [1024 * 1024]byte
}

func getPositionAndWide(registerName RegisterName) (int, bool) {
	position := 0
	wide := false
	switch registerName {
	case AL:
		position = 1
		wide = false
	case BL:
		position = 3
		wide = false
	case CL:
		position = 5
		wide = false
	case DL:
		position = 7
		wide = false
	case AH:
		position = 0
		wide = false
	case BH:
		position = 2
		wide = false
	case CH:
		position = 4
		wide = false
	case DH:
		position = 6
		wide = false
	case AX:
		position = 0
		wide = true
	case BX:
		position = 2
		wide = true
	case CX:
		position = 4
		wide = true
	case DX:
		position = 6
		wide = true
	case SP:
		position = 8
		wide = true
	case BP:
		position = 10
		wide = true
	case SI:
		position = 12
		wide = true
	case DI:
		position = 14
		wide = true
	case CS:
		position = 16
		wide = true
	case DS:
		position = 18
		wide = true
	case ES:
		position = 20
		wide = true
	case SS:
		position = 22
		wide = true
	default:
		panic(fmt.Sprintf("unknown register name: %s", registerName))
	}
	return position, wide
}

func (c *Context) SetRegister(registerName RegisterName, value int16) {
	position, wide := getPositionAndWide(registerName)

	if wide {
		c.Registers[position] = byte(value >> 8)
		c.Registers[position+1] = byte(value & 0xff)
	} else {
		c.Registers[position] = byte(value)
	}
}

func (c *Context) GetRegister(registerName RegisterName) int16 {
	position, wide := getPositionAndWide(registerName)
	if !wide {
		return int16(c.Registers[position])
	}

	value := int16(c.Registers[position]) << 8
	value |= int16(c.Registers[position+1])
	return value
}

func Simulate(context *Context, instructions []Instruction) error {
	for _, instruction := range instructions {
		switch instruction.Type {
		case IT_MovImToReg:
			context.SetRegister(instruction.Destination.RegisterName, instruction.Source.ImmediateValue)
		case IT_MovRegMemToFromReg:
			fallthrough
		case IT_MovSegRegToRegMem:
			fallthrough
		case IT_MovRegMemToSegReg:
			if instruction.Destination.Type == DL_Register && instruction.Source.Type == DL_Register {
				value := context.GetRegister(instruction.Source.RegisterName)
				context.SetRegister(instruction.Destination.RegisterName, value)
			}
		default:
			return fmt.Errorf("instruction simulation not implemented for %s", instruction.Type.Name())
		}
	}
	return nil
}
