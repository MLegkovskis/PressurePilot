package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"intra/internal/data"
)

func synth(ts time.Time) float64 {
	min := float64(ts.Hour()*60 + ts.Minute())
	daily := 10 * math.Sin(2*math.Pi*min/1440)
	return 100 + daily + rand.Float64()*4 - 2
}

func main() {
	pool := data.NewPool(os.Getenv("PG_CONN"))
	defer pool.Close()

	ticker := time.NewTicker(500 * time.Millisecond)  // new: 1 point per 0.5 s
defer ticker.Stop()

	simTime := time.Now().UTC()
	for range ticker.C {
		if err := data.Insert(pool, simTime, synth(simTime)); err != nil {
			log.Println("insert:", err)
		}
		simTime = simTime.Add(30 * time.Minute)
	}
}
