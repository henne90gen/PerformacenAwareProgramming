package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

// x [-180, 180]
// y [-90, 90]

type PointPair struct {
	X0 float64 `json:"x0"`
	Y0 float64 `json:"y0"`
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
}

type Result struct {
	Pairs []PointPair
}

func RandomXCoordinate(r *rand.Rand) float64 {
	return r.Float64()*360.0 - 180.0
}

func RandomYCoordinate(r *rand.Rand) float64 {
	return r.Float64()*180.0 - 90.0
}

func GeneratePointPairs(numPointPairs int, seed int64) []PointPair {
	randomSource := rand.NewSource(seed)
	random := rand.New(randomSource)

	result := make([]PointPair, numPointPairs)

	for i := 0; i < numPointPairs; i++ {
		x0 := RandomXCoordinate(random)
		y0 := RandomYCoordinate(random)
		x1 := RandomXCoordinate(random)
		y1 := RandomYCoordinate(random)
		result[i] = PointPair{
			X0: x0,
			Y0: y0,
			X1: x1,
			Y1: y1,
		}
	}

	return result
}

func Main(ctx *cli.Context) error {
	numPointPairs := ctx.Int("total")
	seed := ctx.Int64("seed")
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	var result Result
	result.Pairs = GeneratePointPairs(numPointPairs, seed)

	buf, err := json.Marshal(&result)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	fmt.Printf("%s", buf)
	return nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name: "total",
			},
			&cli.Int64Flag{
				Name: "seed",
			},
		},
		DefaultCommand: "default",
		Commands: []*cli.Command{
			{
				Name:   "default",
				Action: Main,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s", err))
	}
}
