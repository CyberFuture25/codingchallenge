const express = require("express");
const app = express();
app.use(express.json({ limit: "10mb" }));

// Verifica de forma segura si una matriz es diagonal (debe ser cuadrada)
function isDiagonal(m) {
  if (!m || m.length === 0 || m.length !== m[0].length) return false;

  for (let i = 0; i < m.length; i++) {
    for (let j = 0; j < m[i].length; j++) {
      // Si no es la diagonal principal y el valor no es cero (con tolerancia decimal)
      if (i !== j && Math.abs(m[i][j]) > 1e-10) {
        return false;
      }
    }
  }
  return true;
}

app.post("/statistics", (req, res) => {
  const { rotated_matrix, q, r } = req.body;

  // Reunimos todos los valores para calcular las estadísticas globales
  const all = [...rotated_matrix.flat(), ...q.flat(), ...r.flat()];

  const sumTotal = all.reduce((a, b) => a + b, 0);

  const diagonalChecks = {
    rotated: isDiagonal(rotated_matrix),
    q_matrix: isDiagonal(q),
    r_matrix: isDiagonal(r),
  };

  //
  res.json({
    max_value: Math.max(...all),
    min_value: Math.min(...all),
    average: sumTotal / all.length,
    total_sum: sumTotal,
    is_diagonal: Object.values(diagonalChecks).some((v) => v === true),
    detailed_checks: diagonalChecks,
  });
});

app.listen(3000, () => console.log("Node.js listening on port 3000"));
