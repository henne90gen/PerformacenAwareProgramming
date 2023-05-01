package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func assembleWithNasm(inputFile string) ([]byte, error) {
	cmd := exec.Command("nasm", inputFile)
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		println(stdout.String())
		println(stderr.String())
		return nil, err
	}

	assembledInputFile := strings.TrimSuffix(inputFile, ".asm")
	content, err := os.ReadFile(assembledInputFile)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(inputFile, "computer_enhance") {
		err = os.Remove(assembledInputFile)
		if err != nil {
			return nil, err
		}
	}

	return content, nil
}

func assembleAndCompare(inputFileName string, inputFileContent []byte, result []byte) error {
	dir, fileName := filepath.Split(inputFileName)
	tmpFile, err := os.CreateTemp(dir, fileName+"-*.asm")
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
		return fmt.Errorf("length of assembled result (%d) does not match length of input (%d)", len(assembled), len(inputFileContent))
	}

	for i, b := range assembled {
		if b != inputFileContent[i] {
			return fmt.Errorf("[%s] byte %d does not match, expected %08b but got %08b", inputFileName, i, inputFileContent[i], b)
		}
	}

	return nil
}

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

			instructions, err := Disassemble(content)
			stringifiedInstructions := StringifyInstructions(instructions)
			require.NoError(t, err, stringifiedInstructions)

			err = assembleAndCompare(inputFile, content, []byte(stringifiedInstructions))
			require.NoError(t, err, stringifiedInstructions)
		})
	}
}

func createExpectedContext(inputFile string) (*Context, error) {
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
		if len(parts) != 3 {
			return nil, errors.New("unexpected line in 'expected state'")
		}

		registerName := strings.TrimRight(parts[0], ":")
		value, err := strconv.ParseInt(parts[1], 0, 16)
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
}

func TestSimulation(t *testing.T) {
	inputFiles := []string{
		"computer_enhance/perfaware/part1/listing_0043_immediate_movs.asm",
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
			err = simulate(context, instructions)
			require.NoError(t, err)

			requireContextsToBeEqual(t, expectedContext, context)
		})
	}
}
