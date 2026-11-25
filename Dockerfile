FROM node:20-alpine
WORKDIR /app
COPY server-node ./server-node
COPY web ./web
WORKDIR /app/server-node
RUN npm install
ENV PORT=8080
EXPOSE 8080
CMD ["npm", "start"]
