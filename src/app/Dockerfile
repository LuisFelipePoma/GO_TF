# Stage 1: Build
FROM node:18 AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Production
FROM node:18
WORKDIR /app
COPY --from=build /app/dist /app/dist
COPY --from=build /app/package*.json ./

# Install production dependencies
RUN npm install --only=production

# Set environment variables
ENV HOST=0.0.0.0
ENV PORT=4321

# Expose the port
EXPOSE 4321

# Command to run the application
CMD ["node", "dist/server/entry.mjs"]