FROM node:18-alpine3.19 AS BUILD
WORKDIR /app/front
COPY gin-blog-front/package*.json .
RUN npm config set registry https://registry.npmmirror.com \
    && npm install -g pnpm \
    && pnpm install
COPY gin-blog-front .
RUN pnpm install && pnpm build

WORKDIR /app/admin
COPY gin-blog-admin .
RUN pnpm install && pnpm build

## 阶段２
