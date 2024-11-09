FROM node:18-alpine AS runtime
WORKDIR /app

COPY . .

RUN npm install
RUN npm install @astrojs/node --save
RUN npm run build

ENV HOST=0.0.0.0
ENV PORT=4321
EXPOSE 4321
CMD node ./dist/server/entry.mjs