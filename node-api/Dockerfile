# Usa la imagen base oficial de Node.js
FROM node:20-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /usr/src/app

# Copia los archivos de paquete e instala las dependencias
COPY package*.json ./
RUN npm install --only=production

# Copia el resto del código de la aplicación
COPY . .

# Expone el puerto en el que se ejecuta el API
EXPOSE 3000

# Comando para ejecutar tu aplicación
CMD ["node", "index.js"]