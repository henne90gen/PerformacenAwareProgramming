package main

import "fmt"

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

	IT_IncRegMem
	IT_IncReg
	IT_AsciiAdjustForAdd
	IT_DecimalAdjustForAdd

	IT_SubRegMemWithRegToEither
	IT_SubImToRegMem
	IT_SubImFromAcc

	IT_SubWithBorrowRegMemWithRegToEither
	IT_SubWithBorrowImToRegMem
	IT_SubWithBorrowImFromAcc

	IT_DecRegMem
	IT_DecReg

	IT_Neg

	IT_CmpRegMemAndReg
	IT_CmpImWithRegMem
	IT_CmpImWithAcc

	IT_AsciiAdjustForSubtract
	IT_DecimalAdjustForSubtract

	IT_Multiply
	IT_MultiplySigned
	IT_AsciiAdjustForMultiply

	IT_Divide
	IT_DivideSigned
	IT_AsciiAdjustForDivide

	IT_ConvertByteToWord
	IT_ConvertWordToDoubleWord

	IT_Not
	IT_ShiftLogicLeft
	IT_ShiftLogicRight
	IT_ShiftArithmeticRight
	IT_RotateLeft
	IT_RotateRight
	IT_RotateThroughCarryFlagLeft
	IT_RotateThroughCarryFlagRight

	IT_AndRegMemWithRegToEither
	IT_AndImToRegMem
	IT_AndImToAcc

	IT_TestRegMemAndReg
	IT_TestImAndRegMem
	IT_TestImAndAcc

	IT_OrRegMemWithRegToEither
	IT_OrImToRegMem
	IT_OrImToAcc

	IT_XorRegMemWithRegToEither
	IT_XorImToRegMem
	IT_XorImToAcc

	IT_Repeat
	IT_MoveByte
	IT_CompareByte
	IT_ScanByte
	IT_LoadByte
	IT_StoreByte

	IT_CallDirectWithinSegment
	IT_CallIndirectWithinSegment
	IT_CallDirectIntersegment
	IT_CallIndirectIntersegment

	IT_JumpDirectWithinSegment
	IT_JumpDirectWithinSegmentShort
	IT_JumpIndirectWithinSegment
	IT_JumpDirectIntersegment
	IT_JumpIndirectIntersegment

	IT_ReturnWithinSegment
	IT_ReturnWithinSegmentAddingImmediateToSP
	IT_ReturnIntersegment
	IT_ReturnIntersegmentAddingImmediateToSP

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

	IT_InterruptTypeSpecified
	IT_InterruptType3
	IT_InterruptOnOverflow
	IT_InterruptReturn

	IT_ClearCarry
	IT_ComplementCarry
	IT_SetCarry
	IT_ClearDirection
	IT_SetDirection
	IT_ClearInterrupt
	IT_SetInterrupt
	IT_Halt
	IT_Wait
	IT_Escape
	IT_BusLockPrefix
)

func (t InstructionType) Name() string {
	if t > IT_Invalid && t <= IT_MovSegRegToRegMem {
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

	if t >= IT_IncRegMem && t <= IT_IncReg {
		return "inc"
	}

	if t == IT_AsciiAdjustForAdd {
		return "aaa"
	}

	if t == IT_DecimalAdjustForAdd {
		return "daa"
	}

	if t >= IT_SubRegMemWithRegToEither && t <= IT_SubImFromAcc {
		return "sub"
	}

	if t >= IT_SubWithBorrowRegMemWithRegToEither && t <= IT_SubWithBorrowImFromAcc {
		return "sbb"
	}

	if t >= IT_DecRegMem && t <= IT_DecReg {
		return "dec"
	}

	if t == IT_Neg {
		return "neg"
	}

	if t >= IT_CmpRegMemAndReg && t <= IT_CmpImWithAcc {
		return "cmp"
	}

	if t == IT_AsciiAdjustForSubtract {
		return "aas"
	}

	if t == IT_DecimalAdjustForSubtract {
		return "das"
	}

	if t == IT_Multiply {
		return "mul"
	}

	if t == IT_MultiplySigned {
		return "imul"
	}

	if t == IT_AsciiAdjustForMultiply {
		return "aam"
	}

	if t == IT_Divide {
		return "div"
	}

	if t == IT_DivideSigned {
		return "idiv"
	}

	if t == IT_AsciiAdjustForDivide {
		return "aad"
	}

	if t == IT_ConvertByteToWord {
		return "cbw"
	}

	if t == IT_ConvertWordToDoubleWord {
		return "cwd"
	}

	if t == IT_Not {
		return "not"
	}

	if t == IT_ShiftLogicLeft {
		return "shl"
	}

	if t == IT_ShiftLogicRight {
		return "shr"
	}

	if t == IT_ShiftArithmeticRight {
		return "sar"
	}

	if t == IT_RotateLeft {
		return "rol"
	}

	if t == IT_RotateRight {
		return "ror"
	}

	if t == IT_RotateThroughCarryFlagLeft {
		return "rcl"
	}

	if t == IT_RotateThroughCarryFlagRight {
		return "rcr"
	}

	if t >= IT_AndRegMemWithRegToEither && t <= IT_AndImToAcc {
		return "and"
	}

	if t >= IT_TestRegMemAndReg && t <= IT_TestImAndAcc {
		return "test"
	}

	if t >= IT_OrRegMemWithRegToEither && t <= IT_OrImToAcc {
		return "or"
	}

	if t >= IT_XorRegMemWithRegToEither && t <= IT_XorImToAcc {
		return "xor"
	}

	if t == IT_Repeat {
		return "rep"
	}

	if t == IT_MoveByte {
		return "movs"
	}

	if t == IT_CompareByte {
		return "cmps"
	}

	if t == IT_ScanByte {
		return "scas"
	}

	if t == IT_LoadByte {
		return "lods"
	}

	if t == IT_StoreByte {
		return "stos"
	}

	if t >= IT_CallDirectWithinSegment && t <= IT_CallIndirectIntersegment {
		return "call"
	}

	if t >= IT_JumpDirectWithinSegment && t <= IT_JumpIndirectIntersegment {
		return "jmp"
	}

	if t >= IT_ReturnWithinSegment && t <= IT_ReturnIntersegmentAddingImmediateToSP {
		return "ret"
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

	if t == IT_InterruptTypeSpecified {
		return "int"
	}

	if t == IT_InterruptType3 {
		return "int3"
	}

	if t == IT_InterruptOnOverflow {
		return "into"
	}

	if t == IT_InterruptReturn {
		return "iret"
	}

	if t == IT_ClearCarry {
		return "clc"
	}

	if t == IT_ComplementCarry {
		return "cmc"
	}

	if t == IT_SetCarry {
		return "stc"
	}

	if t == IT_ClearDirection {
		return "cld"
	}

	if t == IT_SetDirection {
		return "std"
	}

	if t == IT_ClearInterrupt {
		return "cli"
	}

	if t == IT_SetInterrupt {
		return "sti"
	}

	if t == IT_Halt {
		return "hlt"
	}

	if t == IT_Wait {
		return "wait"
	}

	if t == IT_Escape {
		return "esc"
	}

	if t == IT_BusLockPrefix {
		return "lock"
	}

	return "unknown"
}

func (t InstructionType) IsImToAcc() bool {
	return t == IT_AddImToAcc ||
		t == IT_AddWithCarryImToAcc ||
		t == IT_SubImFromAcc ||
		t == IT_SubWithBorrowImFromAcc ||
		t == IT_CmpImWithAcc ||
		t == IT_AndImToAcc ||
		t == IT_TestImAndAcc ||
		t == IT_OrImToAcc ||
		t == IT_XorImToAcc
}

func (t InstructionType) IsRegMemWithRegToEither() bool {
	return t == IT_MovRegMemToFromReg ||
		t == IT_MovSegRegToRegMem ||
		t == IT_MovRegMemToSegReg ||
		t == IT_AddRegMemWithRegToEither ||
		t == IT_AddWithCarryRegMemWithRegToEither ||
		t == IT_IncRegMem ||
		t == IT_SubRegMemWithRegToEither ||
		t == IT_SubWithBorrowRegMemWithRegToEither ||
		t == IT_DecRegMem ||
		t == IT_Neg ||
		t == IT_CmpRegMemAndReg ||
		t == IT_ExchangeRegMemWithReg ||
		t == IT_LoadEA ||
		t == IT_LoadDS ||
		t == IT_LoadES ||
		t == IT_Multiply ||
		t == IT_MultiplySigned ||
		t == IT_Divide ||
		t == IT_DivideSigned ||
		t == IT_Not ||
		t == IT_ShiftLogicLeft ||
		t == IT_ShiftLogicRight ||
		t == IT_ShiftArithmeticRight ||
		t == IT_RotateLeft ||
		t == IT_RotateRight ||
		t == IT_RotateThroughCarryFlagLeft ||
		t == IT_RotateThroughCarryFlagRight ||
		t == IT_AndRegMemWithRegToEither ||
		t == IT_TestRegMemAndReg ||
		t == IT_OrRegMemWithRegToEither ||
		t == IT_XorRegMemWithRegToEither ||
		t == IT_CallIndirectWithinSegment ||
		t == IT_CallIndirectIntersegment ||
		t == IT_JumpIndirectWithinSegment
}

func (t InstructionType) IsImToRegMem() bool {
	return t == IT_MovImToRegMem ||
		t == IT_AddImToRegMem ||
		t == IT_AddWithCarryImToRegMem ||
		t == IT_SubImToRegMem ||
		t == IT_SubWithBorrowImToRegMem ||
		t == IT_CmpImWithRegMem ||
		t == IT_AndImToRegMem ||
		t == IT_TestImAndRegMem ||
		t == IT_OrImToRegMem ||
		t == IT_XorImToRegMem
}

func (t InstructionType) HasSignExtension() bool {
	return t == IT_AddImToRegMem ||
		t == IT_AddWithCarryImToRegMem ||
		t == IT_SubImToRegMem ||
		t == IT_CmpImWithRegMem
}

func (t InstructionType) IsConditionalJump() bool {
	return t == IT_JE ||
		t == IT_JNE ||
		t == IT_JL ||
		t == IT_JLE ||
		t == IT_JB ||
		t == IT_JBE ||
		t == IT_JP ||
		t == IT_JO ||
		t == IT_JS ||
		t == IT_JNL ||
		t == IT_JNLE ||
		t == IT_JNB ||
		t == IT_JNBE ||
		t == IT_JNP ||
		t == IT_JNO ||
		t == IT_JNS ||
		t == IT_LOOP ||
		t == IT_LOOPZ ||
		t == IT_LOOPNZ ||
		t == IT_JCXZ
}

func (t InstructionType) IsInOut() bool {
	return t == IT_InFixed ||
		t == IT_InVariable ||
		t == IT_OutFixed ||
		t == IT_OutVariable
}

func (t InstructionType) AlwaysToRegister() bool {
	return t == IT_ExchangeRegMemWithReg ||
		t == IT_LoadEA ||
		t == IT_LoadDS ||
		t == IT_LoadES
}

func (t InstructionType) IsSingleByteInstruction() bool {
	return t == IT_XLAT ||
		t == IT_LoadAHWithFlags ||
		t == IT_StoreAHWithFlags ||
		t == IT_PushFlags ||
		t == IT_PopFlags ||
		t == IT_AsciiAdjustForAdd ||
		t == IT_DecimalAdjustForAdd ||
		t == IT_AsciiAdjustForSubtract ||
		t == IT_DecimalAdjustForSubtract ||
		t == IT_ConvertByteToWord ||
		t == IT_ConvertWordToDoubleWord ||
		t == IT_Repeat ||
		t.IsStringManipulationInstruction() ||
		t == IT_ReturnWithinSegment ||
		t == IT_ReturnIntersegment ||
		t == IT_InterruptType3 ||
		t == IT_InterruptOnOverflow ||
		t == IT_InterruptReturn ||
		t == IT_ClearCarry ||
		t == IT_ComplementCarry ||
		t == IT_SetCarry ||
		t == IT_ClearDirection ||
		t == IT_SetDirection ||
		t == IT_ClearInterrupt ||
		t == IT_SetInterrupt ||
		t == IT_Halt ||
		t == IT_Wait ||
		t == IT_BusLockPrefix
}

func (t InstructionType) IsStringManipulationInstruction() bool {
	return t == IT_MoveByte ||
		t == IT_CompareByte ||
		t == IT_ScanByte ||
		t == IT_LoadByte ||
		t == IT_StoreByte
}

func (t InstructionType) IsSingleOperandInstruction() bool {
	return t == IT_IncRegMem ||
		t == IT_DecRegMem ||
		t == IT_Neg ||
		t == IT_Multiply ||
		t == IT_MultiplySigned ||
		t == IT_Divide ||
		t == IT_DivideSigned ||
		t == IT_Not ||
		t == IT_CallIndirectWithinSegment ||
		t == IT_CallIndirectIntersegment ||
		t == IT_JumpDirectWithinSegment ||
		t == IT_JumpIndirectWithinSegment ||
		t == IT_JumpDirectIntersegment ||
		t == IT_JumpIndirectIntersegment
}

func (t InstructionType) IsShiftOrRotateInstruction() bool {
	return t == IT_ShiftLogicLeft ||
		t == IT_ShiftLogicRight ||
		t == IT_ShiftArithmeticRight ||
		t == IT_RotateLeft ||
		t == IT_RotateRight ||
		t == IT_RotateThroughCarryFlagLeft ||
		t == IT_RotateThroughCarryFlagRight
}

func InstructionTypeFromBytes(content []byte) (InstructionType, error) {
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
		if reg == 0b011 {
			return IT_SubWithBorrowImToRegMem, nil
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

	if b == 0b11111111 && (content[1]>>3)&0b111 == 0b110 {
		return IT_PushRegMem, nil
	}

	if b&0b11111000 == 0b01010000 {
		return IT_PushReg, nil
	}

	if b&0b11100111 == 0b00000110 {
		return IT_PushSegReg, nil
	}

	if b == 0b10001111 && (content[1]>>3)&0b111 == 0b000 {
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

	if (b>>1) == 0b1111111 && (content[1]>>3)&0b111 == 0b000 {
		return IT_IncRegMem, nil
	}

	if (b >> 3) == 0b01000 {
		return IT_IncReg, nil
	}

	if b == 0b00110111 {
		return IT_AsciiAdjustForAdd, nil
	}

	if b == 0b00100111 {
		return IT_DecimalAdjustForAdd, nil
	}

	if (b >> 2) == 0b000110 {
		return IT_SubWithBorrowRegMemWithRegToEither, nil
	}

	if (b >> 1) == 0b0001110 {
		return IT_SubWithBorrowImFromAcc, nil
	}

	if (b>>1) == 0b1111111 && (content[1]>>3)&0b111 == 0b001 {
		return IT_DecRegMem, nil
	}

	if (b >> 3) == 0b01001 {
		return IT_DecReg, nil
	}

	if (b>>1) == 0b1111011 && (content[1]>>3)&0b111 == 0b011 {
		return IT_Neg, nil
	}

	if b == 0b00111111 {
		return IT_AsciiAdjustForSubtract, nil
	}

	if b == 0b00101111 {
		return IT_DecimalAdjustForSubtract, nil
	}

	if (b>>1) == 0b1111011 && (content[1]>>3)&0b111 == 0b100 {
		return IT_Multiply, nil
	}

	if (b>>1) == 0b1111011 && (content[1]>>3)&0b111 == 0b101 {
		return IT_MultiplySigned, nil
	}

	if b == 0b11010100 && content[1] == 0b00001010 {
		return IT_AsciiAdjustForMultiply, nil
	}

	if (b>>1) == 0b1111011 && (content[1]>>3)&0b111 == 0b110 {
		return IT_Divide, nil
	}

	if (b>>1) == 0b1111011 && (content[1]>>3)&0b111 == 0b111 {
		return IT_DivideSigned, nil
	}

	if b == 0b11010101 && content[1] == 0b00001010 {
		return IT_AsciiAdjustForDivide, nil
	}

	if b == 0b10011000 {
		return IT_ConvertByteToWord, nil
	}

	if b == 0b10011001 {
		return IT_ConvertWordToDoubleWord, nil
	}

	if (b>>1) == 0b1111011 && (content[1]>>3)&0b111 == 0b010 {
		return IT_Not, nil
	}

	if (b>>2) == 0b110100 && (content[1]>>3)&0b111 == 0b100 {
		return IT_ShiftLogicLeft, nil
	}

	if (b>>2) == 0b110100 && (content[1]>>3)&0b111 == 0b101 {
		return IT_ShiftLogicRight, nil
	}

	if (b>>2) == 0b110100 && (content[1]>>3)&0b111 == 0b111 {
		return IT_ShiftArithmeticRight, nil
	}

	if (b>>2) == 0b110100 && (content[1]>>3)&0b111 == 0b000 {
		return IT_RotateLeft, nil
	}

	if (b>>2) == 0b110100 && (content[1]>>3)&0b111 == 0b001 {
		return IT_RotateRight, nil
	}

	if (b>>2) == 0b110100 && (content[1]>>3)&0b111 == 0b010 {
		return IT_RotateThroughCarryFlagLeft, nil
	}

	if (b>>2) == 0b110100 && (content[1]>>3)&0b111 == 0b011 {
		return IT_RotateThroughCarryFlagRight, nil
	}

	if (b >> 2) == 0b001000 {
		return IT_AndRegMemWithRegToEither, nil
	}

	if (b>>1) == 0b1000000 && (content[1]>>3)&0b111 == 0b100 {
		return IT_AndImToRegMem, nil
	}

	if (b >> 1) == 0b0010010 {
		return IT_AndImToAcc, nil
	}

	if (b >> 2) == 0b100001 {
		return IT_TestRegMemAndReg, nil
	}

	if (b>>1) == 0b1111011 && (content[1]>>3)&0b111 == 0b000 {
		return IT_TestImAndRegMem, nil
	}

	if (b >> 1) == 0b1010100 {
		return IT_TestImAndAcc, nil
	}

	if (b >> 2) == 0b000010 {
		return IT_OrRegMemWithRegToEither, nil
	}

	if (b>>1) == 0b1000000 && (content[1]>>3)&0b111 == 0b001 {
		return IT_OrImToRegMem, nil
	}

	if (b >> 1) == 0b0000110 {
		return IT_OrImToAcc, nil
	}

	if (b >> 2) == 0b001100 {
		return IT_XorRegMemWithRegToEither, nil
	}

	if (b>>1) == 0b1000000 && (content[1]>>3)&0b111 == 0b110 {
		return IT_XorImToRegMem, nil
	}

	if (b >> 1) == 0b0011010 {
		return IT_XorImToAcc, nil
	}

	if (b >> 1) == 0b1111001 {
		return IT_Repeat, nil
	}

	if (b >> 1) == 0b1010010 {
		return IT_MoveByte, nil
	}

	if (b >> 1) == 0b1010011 {
		return IT_CompareByte, nil
	}

	if (b >> 1) == 0b1010111 {
		return IT_ScanByte, nil
	}

	if (b >> 1) == 0b1010110 {
		return IT_LoadByte, nil
	}

	if (b >> 1) == 0b1010101 {
		return IT_StoreByte, nil
	}

	if b == 0b11101000 {
		return IT_CallDirectWithinSegment, nil
	}

	if b == 0b11111111 && (content[1]>>3)&0b111 == 0b010 {
		return IT_CallIndirectWithinSegment, nil
	}

	if b == 0b10011010 {
		return IT_CallDirectIntersegment, nil
	}

	if b == 0b11111111 && (content[1]>>3)&0b111 == 0b011 {
		return IT_CallIndirectIntersegment, nil
	}

	if b == 0b11101001 {
		return IT_JumpDirectWithinSegment, nil
	}

	if b == 0b11101011 {
		return IT_JumpDirectWithinSegmentShort, nil
	}

	if b == 0b11111111 && (content[1]>>3)&0b111 == 0b100 {
		return IT_JumpIndirectWithinSegment, nil
	}

	if b == 0b11101010 {
		return IT_JumpDirectIntersegment, nil
	}

	if b == 0b11111111 && (content[1]>>3)&0b111 == 0b101 {
		return IT_JumpIndirectIntersegment, nil
	}

	if b == 0b11000011 {
		return IT_ReturnWithinSegment, nil
	}

	if b == 0b11000010 {
		return IT_ReturnWithinSegmentAddingImmediateToSP, nil
	}

	if b == 0b11001011 {
		return IT_ReturnIntersegment, nil
	}

	if b == 0b11001010 {
		return IT_ReturnIntersegmentAddingImmediateToSP, nil
	}

	if b == 0b11001101 {
		return IT_InterruptTypeSpecified, nil
	}

	if b == 0b11001100 {
		return IT_InterruptType3, nil
	}

	if b == 0b11001110 {
		return IT_InterruptOnOverflow, nil
	}

	if b == 0b11001111 {
		return IT_InterruptReturn, nil
	}

	if b == 0b11111000 {
		return IT_ClearCarry, nil
	}

	if b == 0b11110101 {
		return IT_ComplementCarry, nil
	}

	if b == 0b11111001 {
		return IT_SetCarry, nil
	}

	if b == 0b11111100 {
		return IT_ClearDirection, nil
	}

	if b == 0b11111101 {
		return IT_SetDirection, nil
	}

	if b == 0b11111010 {
		return IT_ClearInterrupt, nil
	}

	if b == 0b11111011 {
		return IT_SetInterrupt, nil
	}

	if b == 0b11110100 {
		return IT_Halt, nil
	}

	if b == 0b10011011 {
		return IT_Wait, nil
	}

	if b>>3 == 0b11011 {
		return IT_Escape, nil
	}

	if b == 0b11110000 {
		return IT_BusLockPrefix, nil
	}

	return IT_Invalid, fmt.Errorf("opcode %08b %08b not implemented yet", b, content[1])
}
