# Build image
FROM golang:1.16.2-buster as build
WORKDIR /app
RUN go env -w GOPROXY=direct
# cache dependencies
ADD go.mod go.sum ./
RUN go mod download

# build
ADD *.go ./
RUN go build -o /main

# Running image
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
ADD templates /templates
ENTRYPOINT [ "/main" ]
