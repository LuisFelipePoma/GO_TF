# This is the newer version
FROM node:20 as build-stage

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

RUN npm i -g serve

COPY . .

RUN npm run build

CMD [ "npm", "run", "dev" ]