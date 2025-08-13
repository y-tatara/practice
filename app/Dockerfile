FROM node:20

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

EXPOSE 5174

CMD ["npm", "run", "dev", "--", "--host", "--port", "5174"]
