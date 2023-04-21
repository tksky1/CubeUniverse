# docker build -t cubeuniverse-dashboard -f dockerfiles/dashboard.Dockerfile dev-frontend
# docker tag cubeuniverse-dashboard tksky1/cubeuniverse-dashboard:0.1alpha
# docker push tksky1/cubeuniverse-dashboard:0.1alpha
FROM node:16
MAINTAINER tk_sky

COPY . .

RUN cd / && npm install && npm run build

EXPOSE 3000
CMD [ "npm" , "run" , "start" ]