Instalación y Ejecución
Para garantizar la portabilidad y facilidad de despliegue, el proyecto está completamente contenerizado.
Asegúrate de tener instalados Docker y Docker Compose.
Clona este repositorio en tu máquina local.
Desde la raíz del proyecto, ejecuta el siguiente comando para compilar y levantar los servicios:
Una vez que los contenedores estén activos:
La API de Go estará disponible en: http://localhost:8080
La API de Node.js estará disponible internamente para el procesamiento.
Uso de la API
Para iniciar el proceso, debes enviar una solicitud POST al endpoint de la API de Go con una matriz rectangular en formato JSON.
Endpoint: POST http://localhost:8080/process
Ejemplo de Payload (Matriz 4x3):
{
"matrix": [
[1.5, 2.0, 3.5],
[4.0, 5.5, 6.0],
[7.5, 8.0, 9.5],
[10.0, 11.5, 12.0]
]
}
Respuesta Esperada: El sistema devolverá un objeto JSON que integra la factorización QR realizada en Go y las estadísticas calculadas en Node.js, incluyendo la comprobación de matriz diagonal.
