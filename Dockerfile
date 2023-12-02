FROM node:18-bullseye

RUN apt update \
  && apt install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY package.json package-lock.json tsconfig.json ./
RUN npm install
COPY src src
RUN npm run compile
CMD [ "node", "dist" ]
