FROM --platform=linux/amd64 amd64/node:16.17-alpine

RUN mkdir -p /home/app/ && chown -R node:node /home/app
WORKDIR /home/app
COPY --chown=node:node . .

USER node

RUN yarn install --frozen-lockfile
RUN yarn build

EXPOSE 3000
CMD [ "yarn", "start" ]