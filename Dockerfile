###############
# base images #
###############
FROM golang:1.23-alpine AS build
FROM scratch AS final


########################
# build Go application #
########################
FROM build AS app-build
COPY . /go/src/github.com/rclsilver-org/k8s-volume-migration
WORKDIR /go/src/github.com/rclsilver-org/k8s-volume-migration
RUN apk add --no-cache make && \
    make k8s-volume-migration-linux-amd64


#####################
# build final image #
#####################
FROM final
WORKDIR /

COPY --from=app-build /go/src/github.com/rclsilver-org/k8s-volume-migration/dist/k8s-volume-migration-linux-amd64 /k8s-volume-migration

ENTRYPOINT [ "/k8s-volume-migration" ]
