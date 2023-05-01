package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func TestDisassemble(t *testing.T) {
	inputFiles := []string{
		"test.asm",
		"computer_enhance/perfaware/part1/listing_0037_single_register_mov.asm",
		"computer_enhance/perfaware/part1/listing_0038_many_register_mov.asm",
		"computer_enhance/perfaware/part1/listing_0039_more_movs.asm",
		"computer_enhance/perfaware/part1/listing_0040_challenge_movs.asm",
		"computer_enhance/perfaware/part1/listing_0041_add_sub_cmp_jnz.asm",
		// TODO "computer_enhance/perfaware/part1/listing_0042_completionist_decode.asm",
		"computer_enhance/perfaware/part1/listing_0043_immediate_movs.asm",
		"computer_enhance/perfaware/part1/listing_0044_register_movs.asm",
		"computer_enhance/perfaware/part1/listing_0045_challenge_register_movs.asm",
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
