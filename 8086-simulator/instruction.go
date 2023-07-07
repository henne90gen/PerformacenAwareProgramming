package simulator8086

import (
	"fmt"
	"strconv"
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

var segmentRegisterTable = []RegisterName{ES, CS, SS, DS}
var registerTable = [][]RegisterName{
	{AL, CL, DL, BL, AH, CH, DH, BH},
	{AX, CX, DX, BX, SP, BP, SI, DI},
}

type Instruction struct {
	Type        InstructionType
	SizeInBytes int
	Wide        bool
	Destination *DataLocation
	Source      *DataLocation
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
		if d.AvoidSizeInfo {
			return strconv.Itoa(int(d.ImmediateValue))
		}
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
		return fmt.Sprintf("$%+d", d.LabelPosition)
	}

	panic("unknown data location")
}

func (i Instruction) String() string {
	if i.Source == nil {
		if i.Destination == nil {
			wide := ""
			if i.Type.IsStringManipulationInstruction() {
				if i.Wide {
					wide = "w"
				} else {
					wide = "b"
				}
			}
			return fmt.Sprintf("%s%s\n", i.Type.Name(), wide)
		}

		return fmt.Sprintf("%s %s\n", i.Type.Name(), i.Destination.String())
	}

	return fmt.Sprintf(
		"%s %s, %s\n",
		i.Type.Name(),
		i.Destination.String(),
		i.Source.String(),
	)
}
