package main

import "fmt"

type FlagIndex int

const (
	Flag_Zero FlagIndex = iota
	Flag_Sign
	Flag_Carry
	Flag_AuxilliaryCarry
	Flag_Parity
	Flag_Overflow
)

type Context struct {
	Registers          [24]byte
	Flags              [6]bool
	InstructionPointer int16
	Memory             [1024 * 1024]byte
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

func (c *Context) SetFlag(index FlagIndex, value bool) {
	c.Flags[index] = value
}

func (c *Context) GetFlag(index FlagIndex) bool {
	return c.Flags[index]
}

func (c *Context) ResetFlags() {
	for i := range c.Flags {
		c.Flags[i] = false
	}
}

func (c *Context) GetValue(location *DataLocation) int16 {
	switch location.Type {
	case DL_Invalid:
		panic("Cannot get value of invalid data location")
	case DL_Immediate:
		return location.ImmediateValue
	case DL_Label:
		panic("Cannot get value of label")
	case DL_Register:
		return c.GetRegister(location.RegisterName)
	case DL_Memory:
		panic("TODO")
	}
	return 0
}

func (c *Context) SetValue(destination *DataLocation, value int16, updateFlags bool) {
	switch destination.Type {
	case DL_Invalid:
		panic("Cannot set value of invalid data location")
	case DL_Immediate:
		panic("Cannot set value of immediate")
	case DL_Label:
		panic("Cannot set value of label")
	case DL_Register:
		c.SetRegister(destination.RegisterName, value)
	case DL_Memory:
		panic("TODO")
	}

	if !updateFlags {
		return
	}

	if value == 0 {
		c.SetFlag(Flag_Zero, true)
	}

	parity := 0
	for i := 0; i < 16; i++ {
		parity += int((value >> i) & 0b1)
	}
	if parity%2 == 0 {
		c.SetFlag(Flag_Parity, true)
	}

	if value < 0 {
		c.SetFlag(Flag_Sign, true)
	}
}

func SimulateInstruction(context *Context, instruction Instruction) error {
	switch instruction.Type {
	case IT_MovImToReg:
		context.SetRegister(instruction.Destination.RegisterName, instruction.Source.ImmediateValue)
	case IT_MovRegMemToFromReg:
		fallthrough
	case IT_MovSegRegToRegMem:
		fallthrough
	case IT_MovRegMemToSegReg:
		value := context.GetValue(instruction.Source)
		context.SetValue(instruction.Destination, value, false)
	case IT_SubImToRegMem:
		fallthrough
	case IT_SubRegMemWithRegToEither:
		context.ResetFlags()
		srcValue := context.GetValue(instruction.Source)
		dstValue := context.GetValue(instruction.Destination)
		value := dstValue - srcValue
		context.SetValue(instruction.Destination, value, true)
	case IT_AddImToRegMem:
		fallthrough
	case IT_AddRegMemWithRegToEither:
		srcValue := context.GetValue(instruction.Source)
		dstValue := context.GetValue(instruction.Destination)
		value := srcValue + dstValue
		context.SetValue(instruction.Destination, value, true)
	case IT_CmpRegMemAndReg:
		// srcValue := context.GetValue(instruction.Source)
		// dstValue := context.GetValue(instruction.Destination)
	default:
		return fmt.Errorf("simulation not implemented for instruction %s (%d)", instruction.Type.Name(), instruction.Type)
	}
	return nil
}

func Simulate(context *Context, instructions []Instruction) error {
	for _, instruction := range instructions {
		err := SimulateInstruction(context, instruction)
		if err != nil {
			return err
		}
	}
	return nil
}
