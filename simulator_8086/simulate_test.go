package simulator8086

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func parseFlags(flagsStr string) [6]bool {
	flags := [6]bool{}
	if strings.Contains(flagsStr, "Z") {
		flags[Flag_Zero] = true
	}
	if strings.Contains(flagsStr, "S") {
		flags[Flag_Sign] = true
	}
	if strings.Contains(flagsStr, "P") {
		flags[Flag_Parity] = true
	}
	return flags
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

type ContextTransition struct {
	HasRegisterUpdate bool
	RegisterName      RegisterName
	NewValue          int16

	HasFlagsUpdate bool
	Flags          [6]bool
}

func createExpectedContext(inputFile string) ([]ContextTransition, error) {
	outputFile := strings.TrimSuffix(inputFile, ".asm") + ".txt"
	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(expectedOutput), "\n")
	contextTransitions := make([]ContextTransition, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "Final registers:") {
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--- ") {
			continue
		}

		lineParts := strings.Split(line, " ; ")
		if len(lineParts) != 2 {
			continue
		}

		contextTransition := ContextTransition{
			HasRegisterUpdate: false,
			HasFlagsUpdate:    false,
		}

		contextUpdate := lineParts[1]
		contextUpdateParts := strings.Split(contextUpdate, "flags:")
		if strings.Contains(contextUpdate, "flags:") {
			flagsStr := contextUpdateParts[len(contextUpdateParts)-1]

			contextTransition.HasFlagsUpdate = true
			contextTransition.Flags = parseFlags(flagsStr)
		}

		if (len(contextUpdateParts) > 1 && contextUpdateParts[0] != "") || !strings.Contains(contextUpdate, "flags:") {
			stateUpdate := strings.TrimSpace(contextUpdateParts[0])
			stateUpdateParts := strings.Split(stateUpdate, ":")
			registerName := stateUpdateParts[0]
			valueUpdate := stateUpdateParts[1]
			valueUpdateParts := strings.Split(valueUpdate, "->")

			toValueStr := valueUpdateParts[1]
			toValue, err := strconv.ParseInt(toValueStr, 0, 17)
			if err != nil {
				return nil, err
			}

			contextTransition.HasRegisterUpdate = true
			contextTransition.RegisterName = RegisterName(registerName)
			contextTransition.NewValue = int16(toValue)
		}

		contextTransitions = append(contextTransitions, contextTransition)
	}

	return contextTransitions, nil
}

func requireContextsToBeEqual(t *testing.T, expected *Context, actual *Context) {
	for i := range expected.Registers {
		require.Equalf(t, expected.Registers[i], actual.Registers[i], "mismatch in register state at position %d", i)
	}
	for i := range expected.Flags {
		require.Equalf(t, expected.Flags[i], actual.Flags[i], "mismatch in %s", FlagIndex(i).Name())
	}
}

func TestSimulation(t *testing.T) {
	inputFiles := []string{
		COMPUTER_ENHANCE_PATH + "/perfaware/part1/listing_0043_immediate_movs.asm",
		COMPUTER_ENHANCE_PATH + "/perfaware/part1/listing_0044_register_movs.asm",
		COMPUTER_ENHANCE_PATH + "/perfaware/part1/listing_0045_challenge_register_movs.asm",
		// COMPUTER_ENHANCE_PATH + "/perfaware/part1/listing_0046_add_sub_cmp.asm",
	}
	for _, inputFile := range inputFiles {
		t.Run(inputFile, func(t *testing.T) {
			content, err := assembleWithNasm(inputFile)
			require.NoError(t, err)

			stateTransitions, err := createExpectedContext(inputFile)
			require.NoError(t, err)

			instructions, err := Disassemble(content)
			require.NoError(t, err)

			require.Equal(t, len(stateTransitions), len(instructions))

			context := &Context{}
			expectedContext := &Context{}
			for i, instruction := range instructions {
				transition := stateTransitions[i]
				if transition.HasRegisterUpdate {
					expectedContext.SetRegister(transition.RegisterName, transition.NewValue)
				}
				if transition.HasFlagsUpdate {
					expectedContext.Flags = transition.Flags
				}

				err := SimulateInstruction(context, instruction)
				require.NoError(t, err)

				fmt.Printf("%+v", instruction)
				requireContextsToBeEqual(t, expectedContext, context)
			}
		})
	}
}
