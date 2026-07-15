import React, { useState } from "react";

function App() {
  const [matrixInput, setMatrixInput] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleProcess = async () => {
    setLoading(true);
    try {
      const response = await fetch(
        "https://codingchallenge-production-674f.up.railway.app/process",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ matrix: JSON.parse(matrixInput) }),
        },
      );
      const data = await response.json();
      setResults(data);
    } catch (error) {
      alert("Error al procesar la matriz. Verifique el formato JSON.");
    }
    setLoading(false);
  };

  return (
    <div className="p-8 font-sans max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6 text-blue-800">
        Interseguro Matrix Processor
      </h1>

      <div className="mb-6">
        <label className="block mb-2 font-semibold">
          Ingrese la matriz rectangular (formato JSON):
        </label>
        <textarea
          className="w-full h-32 p-3 border border-gray-300 rounded shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
          placeholder="Ejemplo: [[4, 5], [1, 6], [3, 7]]"
          value={matrixInput}
          onChange={(e) => setMatrixInput(e.target.value)}
        />
      </div>

      <button
        onClick={handleProcess}
        className={`px-6 py-2 rounded text-white font-bold transition-colors ${
          loading ? "bg-gray-400" : "bg-blue-600 hover:bg-blue-700"
        }`}
        disabled={loading}
      >
        {loading ? "Procesando..." : "Procesar Matriz"}
      </button>

      {results && (
        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Resultados de Go: Rotación y Factorización QR */}
          <section className="bg-white p-5 rounded shadow border">
            <h2 className="text-xl font-bold mb-4 text-gray-700">
              Factorización QR (Go API)
            </h2>
            <div className="space-y-4">
              <div>
                <p className="font-semibold text-sm">Matriz Q:</p>
                <pre className="bg-gray-50 p-2 rounded text-xs overflow-auto">
                  {JSON.stringify(results.qr_factorization?.q, null, 2)}
                </pre>
              </div>
              <div>
                <p className="font-semibold text-sm">Matriz R:</p>
                <pre className="bg-gray-50 p-2 rounded text-xs overflow-auto">
                  {JSON.stringify(results.qr_factorization?.r, null, 2)}
                </pre>
              </div>
            </div>
          </section>

          {/* Estadísticas de Node.js API */}
          <section className="bg-blue-50 p-5 rounded shadow border border-blue-200">
            <h2 className="text-xl font-bold mb-4 text-blue-900">
              Estadísticas (Node.js API)
            </h2>
            <ul className="space-y-3">
              <li className="flex justify-between">
                <span>Suma Total:</span>
                <span className="font-mono font-bold">
                  {results.node_statistics?.total_sum}
                </span>
              </li>
              <li className="flex justify-between">
                <span>Promedio:</span>
                <span className="font-mono font-bold">
                  {results.node_statistics?.average?.toFixed(2)}
                </span>
              </li>
              <li className="flex justify-between">
                <span>Valor Máximo:</span>
                <span className="font-mono font-bold">
                  {results.node_statistics?.max_value}
                </span>
              </li>
              <li className="flex justify-between">
                <span>Valor Mínimo:</span>
                <span className="font-mono font-bold">
                  {results.node_statistics?.min_value}
                </span>
              </li>
              <li className="flex justify-between">
                <span>¿Es Matriz Diagonal?:</span>
                <span
                  className={`font-bold ${results.node_statistics?.is_diagonal ? "text-green-600" : "text-red-600"}`}
                >
                  {results.node_statistics?.is_diagonal ? "Sí" : "No"}
                </span>
              </li>
            </ul>
          </section>
        </div>
      )}
    </div>
  );
}

export default App;
