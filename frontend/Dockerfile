FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
ENV DISABLE_ESLINT_PLUGIN=true
RUN npm run build
RUN npm install -g serve
CMD ["serve", "-s", "build", "-l", "3000"]