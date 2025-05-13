package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"intra/internal/data"
)

const (
	period    = 100
	harmonics = 3
)

func main() {
	conn := data.NewPool(os.Getenv("PG_CONN")) // e.g. postgres://user:password@db:5432/pressure_db
	defer conn.Close()

	r := gin.Default()
	r.LoadHTMLGlob("internal/web/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/readings", func(c *gin.Context) {
		readings, err := data.All(conn)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, readings)
	})

	r.GET("/api/predict", func(c *gin.Context) {
		hStr := c.DefaultQuery("horizon", "10")
		horizon, _ := strconv.Atoi(hStr)

		readings, _ := data.All(conn)
		if len(readings) == 0 {
			c.JSON(200, gin.H{"pred": []float64{}})
			return
		}

		ids := make([]float64, len(readings))
		ys := make([]float64, len(readings))
		for i, r := range readings {
			ids[i] = float64(r.ID)
			ys[i] = r.Pressure
		}

		X := data.GenFeatures(ids, period, harmonics)
		beta := data.LinearFit(X, ys)

		start := ids[len(ids)-1] + 1
		fIDs := make([]float64, horizon)
		for i := range fIDs {
			fIDs[i] = start + float64(i)
		}
		Xf := data.GenFeatures(fIDs, period, harmonics)
		pred := data.Predict(Xf, beta)

		c.JSON(200, gin.H{"ids": fIDs, "pred": pred})
	})

	r.Run(":8080")
}
