FROM --platform=linux/amd64 amd64/node:16.17-alpine

RUN mkdir -p /home/app/ && chown -R node:node /home/app
WORKDIR /home/app
COPY --chown=node:node . .

USER node

ARG NEXT_PUBLIC_DOMAIN_NAME
ENV NEXT_PUBLIC_DOMAIN_NAME=$NEXT_PUBLIC_DOMAIN_NAME

ARG NEXT_PUBLIC_S3_RESOURCE_BUCKET
ENV NEXT_PUBLIC_S3_RESOURCE_BUCKET=$NEXT_PUBLIC_S3_RESOURCE_BUCKET

ARG NEXT_PUBLIC_AWS_REGION
ENV NEXT_PUBLIC_AWS_REGION=$NEXT_PUBLIC_AWS_REGION

RUN rm -rf .next

RUN yarn install --frozen-lockfile
RUN yarn build

EXPOSE 3000
CMD [ "yarn", "start" ]