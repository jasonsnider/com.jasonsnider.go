FROM node:14-buster

# Install Python 2 and other necessary build tools
RUN apt-get update && apt-get install -y python2 make g++ && ln -sf python2 /usr/bin/python

# Install Gulp globally
RUN npm install -g gulp-cli

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

CMD ["sh", "-c", "gulp"]