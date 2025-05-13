// internal/data/fourier.go  – replace LinearFit with this version
package data

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

func GenFeatures(ids []float64, period float64, harmonics int) *mat.Dense {
    rows := len(ids)
    cols := 1 + harmonics*2                 // 1 intercept + sin/cos pairs
    X := mat.NewDense(rows, cols, nil)
    for i, id := range ids {
        // intercept
        X.Set(i, 0, 1.0)

        // seasonal harmonics
        for k := 1; k <= harmonics; k++ {
            angle := 2 * math.Pi * float64(k) * id / period
            X.Set(i, 2*k-1, math.Sin(angle))  // sin(k)
            X.Set(i, 2*k,   math.Cos(angle))  // cos(k)
        }
    }
    return X
}

// LinearFit returns β̂ (as *mat.VecDense) solving ordinary least squares.
func LinearFit(X *mat.Dense, y []float64) *mat.VecDense {
	// β̂ = (XᵀX)⁻¹ Xᵀ y
	var xt mat.Dense
	xt.Mul(X.T(), X)

	var xtInv mat.Dense
	xtInv.Inverse(&xt)

	yMat := mat.NewDense(len(y), 1, y)

	var xty mat.Dense
	xty.Mul(X.T(), yMat)

	var beta mat.Dense
	beta.Mul(&xtInv, &xty)

	// ColView returns mat.Vector (interface) → assert to *mat.VecDense
	return beta.ColView(0).(*mat.VecDense)
}

func Predict(X *mat.Dense, beta *mat.VecDense) []float64 {
	var yhat mat.Dense
	yhat.Mul(X, beta)
	return yhat.RawMatrix().Data
}
