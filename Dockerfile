FROM node:18-alpine3.19 AS build
WORKDIR /app/front
COPY front/package.json front/pnpm-lock.yaml ./
RUN npm config set registry https://registry.npmmirror.com \
    && npm install -g pnpm@9.15.4 \
    && CI=true pnpm install --frozen-lockfile
COPY front .
RUN CI=true pnpm build
RUN rm -f /app/front/dist/images/Before.png /app/front/dist/images/After.png \
    /app/front/dist/resume/resume-ready.glb /app/front/dist/favicon.png

WORKDIR /app/admin
COPY admin/package.json admin/pnpm-lock.yaml ./
RUN CI=true pnpm install --frozen-lockfile
COPY admin .
RUN CI=true pnpm build

## 阶段２ 将静态资源部署到nginx
FROM nginx:1.24.0-alpine

RUN apk add --no-cache bash
COPY --from=build /app/front/dist /usr/share/nginx/html
COPY --from=build /app/admin/dist /usr/share/nginx/html/admin

# nginx 配置文件拷贝在容器中
COPY deploy/build/web/default.conf.template /etc/nginx/conf.d/default.conf.template
COPY deploy/build/web/default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template
COPY deploy/build/web/run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["nginx", "-g","daemon off;"]

EXPOSE 80
EXPOSE 443
