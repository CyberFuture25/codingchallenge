const express = require("express");
const app = express();
app.use(express.json({ limit: "10mb" }));

function isDiagonal(m) {
  // A diagonal matrix must be square [3]
  if (!m || m.length === 0 || m.length !== m.length) return false;
  for (let i = 0; i < m.length; i++) {
    for (let j = 0; j < m[i].length; j++) {
      if (i !== j && Math.abs(m[i][j]) > 1e-10) return false;
    }
  }
  return true;
}

app.post("/statistics", (req, res) => {
  const { rotated_matrix, q, r } = req.body;
  const all = [...rotated_matrix.flat(), ...q.flat(), ...r.flat()];

  const sumTotal = all.reduce((a, b) => a + b, 0);
  const diagonalChecks = {
    rotated: isDiagonal(rotated_matrix),
    q_matrix: isDiagonal(q),
    r_matrix: isDiagonal(r),
  };

  res.json({
    maxValue: Math.max(...all),
    minValue: Math.min(...all),
    average: sumTotal / all.length,
    totalSum: sumTotal,
    isDiagonalAny: Object.values(diagonalChecks).some((v) => v === true),
    detailedChecks: diagonalChecks,
  });
});

app.listen(3000, () => console.log("Node.js listening on 3000"));
