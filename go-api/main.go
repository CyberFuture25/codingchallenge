package main

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

// RotateMatrix: Maneja la rotación de M x N a N x M de forma matemática y correcta.
func RotateMatrix(m [][]float64) [][]float64 {
	if len(m) == 0 || len(m[0]) == 0 {
		return m
	}
	
	numRows := len(m)    
	numCols := len(m[0]) // lee el número real de columnas de la matriz rectangular

	result := make([][]float64, numCols)
	for i := range result {
		result[i] = make([]float64, numRows)
	}

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			// Rotación de 90 grados en sentido de las agujas del reloj
			result[j][numRows-1-i] = m[i][j]
		}
	}
	return result
}

// QRFactorization: Lógica de Gram-Schmidt que soporta correctamente matrices rectangulares M x N
func QRFactorization(a [][]float64) ([][]float64, [][]float64) {
	if len(a) == 0 || len(a[0]) == 0 {
		return nil, nil
	}
	
	M := len(a)    
	N := len(a[0]) //  lee el número real de columnas de la matriz rectangular

	// Q es de tamaño M x N, R es de tamaño N x N
	q := make([][]float64, M)
	for i := range q {
		q[i] = make([]float64, N)
	}
	r := make([][]float64, N)
	for i := range r {
		r[i] = make([]float64, N)
	}

	for j := 0; j < N; j++ {
		v := make([]float64, M)
		for i := 0; i < M; i++ {
			v[i] = a[i][j]
		}

		for i := 0; i < j; i++ {
			var dot float64
			for k := 0; k < M; k++ {
				dot += q[k][i] * a[k][j]
			}
			r[i][j] = dot
			for k := 0; k < M; k++ {
				v[k] -= r[i][j] * q[k][i]
			}
		}

		var norm float64
		for i := 0; i < M; i++ {
			norm += v[i] * v[i]
		}
		r[j][j] = math.Sqrt(norm)

		for i := 0; i < M; i++ {
			if r[j][j] > 1e-12 {
				q[i][j] = v[i] / r[j][j]
			}
		}
	}
	return q, r
}

func main() {
	app := fiber.New()
	
	// Middleware esencial de recuperación ante pánicos, logger e inyección de CORS
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/process", func(c *fiber.Ctx) error {
		var req MatrixRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
		}
		if len(req.Matrix) == 0 || len(req.Matrix[0]) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Empty matrix"})
		}

		// Realizar operaciones en Go
		rotated := RotateMatrix(req.Matrix)
		q, r := QRFactorization(req.Matrix)

		payload := QRResult{RotatedMatrix: rotated, Q: q, R: r}
		jsonData, _ := json.Marshal(payload)

		
		nodeURL := os.Getenv("NODE_SERVICE_URL")
		if nodeURL == "" {
			
			nodeURL = "http://node-service.railway.internal:3000"
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(nodeURL+"/statistics", "application/json", bytes.NewBuffer(jsonData))
		
		if err != nil {
			
			resp, err = client.Post("http://localhost:3000/statistics", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "El microservicio de estadísticas (Node.js) no está disponible.",
				})
			}
		}
		defer resp.Body.Close()

		var stats interface{}
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse Node.js stats"})
		}

		return c.JSON(fiber.Map{
			"status":           "success",
			"qr_factorization": payload,
			"node_statistics":  stats,
		})
	})

	// El puerto de escucha del backend
	app.Listen(":8080")
}