package main

import (
	"encoding/json"
	"fmt"
	"math"
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
	Pairs []PointPair `json:"pairs"`
}

func Square(A float64) float64 {
	return A * A
}

func RadiansFromDegrees(Degrees float64) float64 {
	Result := 0.01745329251994329577 * Degrees
	return Result
}

// NOTE EarthRadius is generally expected to be 6372.8
func Haversine(X0, Y0, X1, Y1, EarthRadius float64) float64 {
	/* NOTE: This is not meant to be a "good" way to calculate the Haversine distance.
	   Instead, it attempts to follow, as closely as possible, the formula used in the real-world
	   question on which these homework exercises are loosely based.
	*/

	lat1 := Y0
	lat2 := Y1
	lon1 := X0
	lon2 := X1

	dLat := RadiansFromDegrees(lat2 - lat1)
	dLon := RadiansFromDegrees(lon2 - lon1)
	lat1 = RadiansFromDegrees(lat1)
	lat2 = RadiansFromDegrees(lat2)

	a := Square(math.Sin(dLat/2.0)) + math.Cos(lat1)*math.Cos(lat2)*Square(math.Sin(dLon/2))
	c := 2.0 * math.Asin(math.Sqrt(a))

	return EarthRadius * c
}

type Cluster struct {
	X  float64
	Y  float64
	DX float64
	DY float64
}

func GenerateClusters(random *rand.Rand, numClusters int) []Cluster {
	result := make([]Cluster, numClusters)

	for i := 0; i < numClusters; i++ {
		x := random.Float64()*220.0 - 110.0
		y := random.Float64()*100.0 - 50.0
		dx := random.Float64()*60.0 + 10.0
		dy := random.Float64()*30.0 + 10.0
		result[i] = Cluster{
			X:  x,
			Y:  y,
			DX: dx,
			DY: dy,
		}
	}

	return result
}

func RandomXCoordinateWithinCluster(random *rand.Rand, cluster Cluster) float64 {
	f := random.Float64()*2.0 - 1.0
	return cluster.X + f*cluster.DX
}

func RandomYCoordinateWithinCluster(random *rand.Rand, cluster Cluster) float64 {
	f := random.Float64()*2.0 - 1.0
	return cluster.Y + f*cluster.DY
}

func GeneratePointPairs(numPointPairs int, seed int64) []PointPair {
	randomSource := rand.NewSource(seed)
	random := rand.New(randomSource)

	numClusters := 64
	clusters := GenerateClusters(random, numClusters)

	result := make([]PointPair, numPointPairs)
	distanceSum := 0.0
	for i := 0; i < numPointPairs; i++ {
		clusterIndex := i % numClusters
		x0 := RandomXCoordinateWithinCluster(random, clusters[clusterIndex])
		y0 := RandomYCoordinateWithinCluster(random, clusters[clusterIndex])
		x1 := RandomXCoordinateWithinCluster(random, clusters[clusterIndex])
		y1 := RandomYCoordinateWithinCluster(random, clusters[clusterIndex])
		result[i] = PointPair{
			X0: x0,
			Y0: y0,
			X1: x1,
			Y1: y1,
		}
		distance := Haversine(x0, y0, x1, y1, 6372.8)
		distanceSum += distance
	}

	distanceAvg := distanceSum / float64(numPointPairs)
	os.Stderr.WriteString(fmt.Sprintf("Average distance: %f\n", distanceAvg))

	return result
}

func Main(ctx *cli.Context) error {
	numPointPairs := ctx.Int("total")
	if numPointPairs == 0 {
		numPointPairs = 1000
	}
	seed := ctx.Int64("seed")
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	outFilePath := ctx.String("out")

	var result Result
	result.Pairs = GeneratePointPairs(numPointPairs, seed)

	buf, err := json.Marshal(&result)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	if outFilePath == "" {
		fmt.Printf("%s", buf)
		return nil
	}

	return os.WriteFile(outFilePath, buf, os.ModePerm)
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
			&cli.StringFlag{
				Name: "out",
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
