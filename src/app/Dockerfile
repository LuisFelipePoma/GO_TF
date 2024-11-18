# This is the newer version alpine
FROM node:20.10-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

RUN npm i -g serve

COPY . .

RUN npm run build

EXPOSE 4321

CMD ["serve", "-s", "dist", "-p", "4321"]