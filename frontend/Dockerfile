FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm build

FROM node:22-alpine

RUN npm install -g serve

COPY --from=builder /app/dist /app/dist

EXPOSE 5173

CMD ["serve", "-s", "/app/dist", "-l", "5173"]
