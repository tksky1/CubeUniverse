
FROM node:12

WORKDIR /

COPY dev-frontend/ /

RUN cd /dev-frontend/ && npm install

EXPOSE 3000
CMD [ "cd /dev-frontend/ && npm run dev" ]