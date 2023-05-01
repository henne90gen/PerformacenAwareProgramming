package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	inputFiles := []string{
		"test.asm",
		"computer_enhance/perfaware/part1/listing_0037_single_register_mov.asm",
		"computer_enhance/perfaware/part1/listing_0038_many_register_mov.asm",
		"computer_enhance/perfaware/part1/listing_0039_more_movs.asm",
		"computer_enhance/perfaware/part1/listing_0040_challenge_movs.asm",
		"computer_enhance/perfaware/part1/listing_0041_add_sub_cmp_jnz.asm",
		// TODO "computer_enhance/perfaware/part1/listing_0042_completionist_decode.asm",
		"computer_enhance/perfaware/part1/listing_0043_immediate_movs.asm",
	}
	for _, inputFile := range inputFiles {
		t.Run(inputFile, func(t *testing.T) {
			content, err := assembleWithNasm(inputFile)
			require.NoError(t, err)

			result, err := disassemble(content)
			require.NoError(t, err)

			stringifiedInstructions := stringifyResult(result)
			err = assembleAndCompare(inputFile, content, []byte(stringifiedInstructions))
			require.NoError(t, err)
		})
	}
}

func TestSimulation(t *testing.T) {
	inputFiles := []string{
		"computer_enhance/perfaware/part1/listing_0043_immediate_movs.asm",
	}
	for _, inputFile := range inputFiles {
		t.Run(inputFile, func(t *testing.T) {
			content, err := assembleWithNasm(inputFile)
			require.NoError(t, err)

			result, err := disassemble(content)
			require.NoError(t, err)

			require.Len(t, result.Instructions, 8)
		})
	}
}
