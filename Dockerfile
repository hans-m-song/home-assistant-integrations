FROM node:18-bullseye
WORKDIR /app
COPY package.json yarn.lock tsconfig.json ./
RUN yarn install
COPY src src
RUN yarn compile
CMD [ "node", "dist" ]
