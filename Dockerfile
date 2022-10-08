FROM node:16-bullseye

RUN sudo apt update \
  && sudo apt install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY package.json yarn.lock tsconfig.json ./
RUN yarn install
COPY src src
RUN yarn compile
CMD [ "node", "dist" ]
