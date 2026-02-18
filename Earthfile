VERSION 0.8

ARG --global TERRAFORM_VERSION=1.5.7
ARG --global TERRAFORM_PROVIDER_SOURCE=SigNoz/signoz
ARG --global TERRAFORM_PROVIDER_VERSION=0.0.12-rc1
ARG --global TERRAFORM_PROVIDER_DOWNLOAD_NAME=terraform-provider-signoz
ARG --global TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX=https://github.com/pyrex41/terraform-provider-signoz/releases/download/v${TERRAFORM_PROVIDER_VERSION}
ARG --global TERRAFORM_NATIVE_PROVIDER_BINARY=terraform-provider-signoz_v${TERRAFORM_PROVIDER_VERSION}
ARG --global IMAGE_REPO=ghcr.io/pyrex41/provider-signoz
ARG --global IMAGE_TAG=v0.2.0

go-mod:
    FROM golang:1.24-alpine
    WORKDIR /src
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod

build:
    FROM golang:1.24-alpine
    WORKDIR /src
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    ARG GOOS=linux
    ARG GOARCH=amd64
    ENV CGO_ENABLED=0
    RUN go build -trimpath -o /out/provider ./cmd/provider
    SAVE ARTIFACT /out/provider

image:
    FROM alpine:3.23.2
    RUN apk --no-cache add ca-certificates bash
    ARG TARGETARCH
    ARG TARGETOS
    ENV USER_ID=65532

    COPY (+build/provider --GOOS=${TARGETOS} --GOARCH=${TARGETARCH}) /usr/local/bin/provider

    ENV PLUGIN_DIR=/terraform/provider-mirror/registry.terraform.io/${TERRAFORM_PROVIDER_SOURCE}/${TERRAFORM_PROVIDER_VERSION}/${TARGETOS}_${TARGETARCH}
    ENV TF_CLI_CONFIG_FILE=/terraform/.terraformrc
    ENV TF_FORK=0

    RUN mkdir -p ${PLUGIN_DIR}

    # Download terraform
    RUN wget -q -O /tmp/terraform.zip \
        https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${TARGETOS}_${TARGETARCH}.zip \
        && unzip /tmp/terraform.zip -d /usr/local/bin \
        && chmod +x /usr/local/bin/terraform \
        && rm /tmp/terraform.zip

    # Download terraform-provider-signoz from fork release
    RUN wget -q -O /tmp/provider.zip \
        ${TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX}/${TERRAFORM_PROVIDER_DOWNLOAD_NAME}_${TERRAFORM_PROVIDER_VERSION}_${TARGETOS}_${TARGETARCH}.zip \
        && unzip /tmp/provider.zip -d ${PLUGIN_DIR} \
        && chmod +x ${PLUGIN_DIR}/* \
        && rm /tmp/provider.zip

    COPY cluster/images/provider-signoz/terraformrc.hcl ${TF_CLI_CONFIG_FILE}
    RUN chown -R ${USER_ID}:${USER_ID} /terraform

    ENV TERRAFORM_VERSION=${TERRAFORM_VERSION}
    ENV TERRAFORM_PROVIDER_SOURCE=${TERRAFORM_PROVIDER_SOURCE}
    ENV TERRAFORM_PROVIDER_VERSION=${TERRAFORM_PROVIDER_VERSION}
    ENV TERRAFORM_NATIVE_PROVIDER_PATH=${PLUGIN_DIR}/${TERRAFORM_NATIVE_PROVIDER_BINARY}

    USER ${USER_ID}
    EXPOSE 8080
    ENTRYPOINT ["provider"]

    SAVE IMAGE --push ${IMAGE_REPO}:${IMAGE_TAG}

all:
    BUILD --platform=linux/amd64 --platform=linux/arm64 +image
