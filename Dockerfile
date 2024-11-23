FROM node:18-alpine3.19 AS BUILD
WORKDIR /app/front
COPY front/package*.json .
RUN npm config set registry https://registry.npmmirror.com \
    && npm install -g pnpm \
    && pnpm install
COPY front .
RUN pnpm install && pnpm build

WORKDIR /app/admin
COPY admin .
RUN pnpm install && pnpm build

## 阶段２ 将静态资源部署到nginx
FROM nginx:1.24.0-alpine

RUN apk add --no-cache bash
COPY --from=BUILD /app/front/dist /usr/share/nginx/html
COPY --from=BUILD /app/admin/dist /usr/share/nginx/html/admin

# nginx 配置文件拷贝在容器中
COPY deploy/build/web/default.conf.template /etc/nginx/conf.d/default.conf.template
COPY deploy/build/web/default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template
COPY deploy/build/web/run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["nginx", "-g","daemon off;"]

EXPOSE 80
EXPOSE 443
