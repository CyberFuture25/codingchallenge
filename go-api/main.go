package main

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type QRResult struct {
	RotatedMatrix [][]float64 `json:"rotated_matrix"`
	Q             [][]float64 `json:"q"`
	R             [][]float64 `json:"r"`
}

// RotateMatrix: Specifically handles M x N to N x M transformation
func RotateMatrix(m [][]float64) [][]float64 {
	if len(m) == 0 || len(m) == 0 { return m }
	
	numRows := len(m)    // 4
	numCols := len(m) // 3

	// Result must be 3 rows x 4 columns
	result := make([][]float64, numCols)
	for i := range result {
		result[i] = make([]float64, numRows)
	}

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			// Mathematical 90-degree clockwise rotation
			result[j][numRows-1-i] = m[i][j]
		}
	}
	return result
}

// QRFactorization: Gram-Schmidt logic for M x N matrices
func QRFactorization(a [][]float64) ([][]float64, [][]float64) {
	if len(a) == 0 || len(a) == 0 { return nil, nil }
	
	M := len(a)    // 4 rows
	N := len(a) // 3 columns

	// Q is M x N (4x3), R is N x N (3x3)
	q := make([][]float64, M)
	for i := range q { q[i] = make([]float64, N) }
	r := make([][]float64, N)
	for i := range r { r[i] = make([]float64, N) }

	for j := 0; j < N; j++ {
		v := make([]float64, M)
		for i := 0; i < M; i++ { v[i] = a[i][j] }

		for i := 0; i < j; i++ {
			var dot float64
			for k := 0; k < M; k++ { dot += q[k][i] * a[k][j] }
			r[i][j] = dot
			for k := 0; k < M; k++ { v[k] -= r[i][j] * q[k][i] }
		}

		var norm float64
		for i := 0; i < M; i++ { norm += v[i] * v[i] }
		r[j][j] = math.Sqrt(norm)

		for i := 0; i < M; i++ {
			if r[j][j] > 1e-12 { q[i][j] = v[i] / r[j][j] }
		}
	}
	return q, r
}

func main() {
	app := fiber.New()
	app.Use(recover.New()) // This will catch any panic and prevent "socket hang up"
	app.Use(logger.New())

	app.Post("/process", func(c *fiber.Ctx) error {
		var req MatrixRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
		}
		if len(req.Matrix) == 0 || len(req.Matrix) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Empty matrix"})
		}

		// Perform operations required by the challenge [1]
		rotated := RotateMatrix(req.Matrix)
		q, r := QRFactorization(req.Matrix)

		payload := QRResult{RotatedMatrix: rotated, Q: q, R: r}
		jsonData, _ := json.Marshal(payload)

		// Forwarding to Node.js service for additional operations [2]
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post("http://node-service:3000/statistics", "application/json", bytes.NewBuffer(jsonData))
		
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Node service unreachable"})
		}
		defer resp.Body.Close()

		var stats interface{}
		json.NewDecoder(resp.Body).Decode(&stats)

		return c.JSON(fiber.Map{
			"status": "success",
			"qr_factorization": payload,
			"node_statistics": stats,
		})
	})

	app.Listen(":8080")
}