package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reading struct {
	ID       int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Pressure float64   `json:"pressure"`
  }
  

func Insert(pool *pgxpool.Pool, ts time.Time, p float64) error {
	_, err := pool.Exec(context.Background(),
		`INSERT INTO pressure (timestamp, pressure) VALUES ($1,$2)`, ts, p)
	return err
}

func All(pool *pgxpool.Pool) ([]Reading, error) {
	rows, err := pool.Query(context.Background(),
		`SELECT id, pressure FROM pressure ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Reading
	for rows.Next() {
		var r Reading
		if err = rows.Scan(&r.ID, &r.Pressure); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
