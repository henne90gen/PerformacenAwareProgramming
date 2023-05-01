package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func parseFlags(context *Context, flags string) {
	if strings.Contains(flags, "Z") {
		context.SetFlag(Flag_Zero, true)
	}
	if strings.Contains(flags, "S") {
		context.SetFlag(Flag_Sign, true)
	}
	if strings.Contains(flags, "P") {
		context.SetFlag(Flag_Parity, true)
	}

	// if strings.Contains(flags, "O") {
	// 	context.SetFlag(Flag_Overflow, true)
	// }
	// if strings.Contains(flags, "C") {
	// 	context.SetFlag(Flag_Carry, true)
	// }
	// if strings.Contains(flags, "A") {
	// 	context.SetFlag(Flag_AuxilliaryCarry, true)
	// }
}

func createExpectedContext(inputFile string) (*Context, error) {
	// TODO parse state changes after every instruction as well

	outputFile := strings.TrimSuffix(inputFile, ".asm") + ".txt"
	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(expectedOutput), "\n")
	context := &Context{}
	foundRegisters := false
	for _, line := range lines {
		if strings.HasPrefix(line, "Final registers:") {
			foundRegisters = true
			continue
		}

		if !foundRegisters {
			continue
		}

		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) > 3 {
			return nil, fmt.Errorf("unexpected line in 'expected state': %s", line)
		}

		registerName := strings.TrimRight(parts[0], ":")
		if registerName == "flags" {
			parseFlags(context, parts[1])
			continue
		}

		value, err := strconv.ParseInt(parts[1], 0, 17)
		if err != nil {
			return nil, err
		}

		context.SetRegister(RegisterName(registerName), int16(value))
	}

	return context, nil
}

func requireContextsToBeEqual(t *testing.T, expected *Context, actual *Context) {
	for i := range expected.Registers {
		require.Equalf(t, expected.Registers[i], actual.Registers[i], "mismatch in register state at position %d", i)
	}
	for i := range expected.Flags {
		require.Equalf(t, expected.Flags[i], actual.Flags[i], "mismatch in flag at position %d", i)
	}
}

func TestSimulation(t *testing.T) {
	inputFiles := []string{
		"computer_enhance/perfaware/part1/listing_0043_immediate_movs.asm",
		"computer_enhance/perfaware/part1/listing_0044_register_movs.asm",
		"computer_enhance/perfaware/part1/listing_0045_challenge_register_movs.asm",
		"computer_enhance/perfaware/part1/listing_0046_add_sub_cmp.asm",
	}
	for _, inputFile := range inputFiles {
		t.Run(inputFile, func(t *testing.T) {
			content, err := assembleWithNasm(inputFile)
			require.NoError(t, err)

			expectedContext, err := createExpectedContext(inputFile)
			require.NoError(t, err)

			instructions, err := Disassemble(content)
			require.NoError(t, err)

			context := &Context{}
			err = Simulate(context, instructions)
			require.NoError(t, err)

			requireContextsToBeEqual(t, expectedContext, context)
		})
	}
}
