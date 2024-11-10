# Usa una imagen base de Node.js
FROM node:18-alpine AS build

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos del proyecto al contenedor
COPY package*.json ./

# Instala las dependencias
RUN npm install

# Copia el resto de los archivos del proyecto
COPY . .

# Construye la aplicación
RUN npm run build

# Usa una imagen base de Node.js para servir la aplicación
FROM node:18-alpine AS runtime

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos de construcción desde la etapa anterior
COPY --from=build /app/dist /app/dist
COPY --from=build /app/package*.json ./

# Instala las dependencias de producción
RUN npm install --only=production

# Establece las variables de entorno
ENV HOST=0.0.0.0
ENV PORT=4321

# Expone el puerto en el que se ejecutará la aplicación
EXPOSE 4321

# Comando para ejecutar la aplicación
CMD ["node", "dist/server/entry.mjs"]