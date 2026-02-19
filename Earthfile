VERSION 0.8

ARG --global TERRAFORM_VERSION=1.5.7
ARG --global TERRAFORM_PROVIDER_SOURCE=SigNoz/signoz
ARG --global TERRAFORM_PROVIDER_VERSION=0.0.12-rc7
ARG --global TERRAFORM_PROVIDER_DOWNLOAD_NAME=terraform-provider-signoz
ARG --global TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX=https://github.com/pyrex41/terraform-provider-signoz/releases/download/v${TERRAFORM_PROVIDER_VERSION}
ARG --global TERRAFORM_NATIVE_PROVIDER_BINARY=terraform-provider-signoz_v${TERRAFORM_PROVIDER_VERSION}
ARG --global REGISTRY=ghcr.io/pyrex41
ARG --global VERSION=v0.2.8

build:
    ARG BUILDPLATFORM
    ARG GOOS=linux
    ARG GOARCH
    FROM --platform=$BUILDPLATFORM golang:1.24-alpine
    WORKDIR /src
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    ENV CGO_ENABLED=0
    RUN GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags="-s -w" \
        -trimpath \
        -o /out/provider ./cmd/provider
    SAVE ARTIFACT /out/provider

# Controller runtime image — pushed separately, referenced by package/crossplane.yaml
image:
    ARG TARGETPLATFORM
    ARG TARGETOS
    ARG TARGETARCH
    FROM alpine:3.23.2
    RUN apk --no-cache add ca-certificates bash
    ENV USER_ID=65532

    COPY (+build/provider --GOOS=${TARGETOS} --GOARCH=${TARGETARCH}) /usr/local/bin/provider

    ENV PLUGIN_DIR=/terraform/provider-mirror/registry.terraform.io/${TERRAFORM_PROVIDER_SOURCE}/${TERRAFORM_PROVIDER_VERSION}/${TARGETOS}_${TARGETARCH}
    ENV TF_CLI_CONFIG_FILE=/terraform/.terraformrc
    ENV TF_FORK=0

    RUN mkdir -p ${PLUGIN_DIR}

    RUN wget -q -O /tmp/terraform.zip \
        https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${TARGETOS}_${TARGETARCH}.zip \
        && unzip /tmp/terraform.zip -d /usr/local/bin \
        && chmod +x /usr/local/bin/terraform \
        && rm /tmp/terraform.zip

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

    SAVE IMAGE --push ${REGISTRY}/provider-signoz:${VERSION}-${TARGETARCH}

# Build and push arch-specific controller images
push-images:
    ARG VERSION=$VERSION
    BUILD --platform=linux/amd64 --platform=linux/arm64 +image --VERSION=$VERSION
    FROM alpine:3.23.2
    RUN echo "Images pushed at $(date)" > /tmp/push-complete
    SAVE ARTIFACT /tmp/push-complete push-complete

# Create multi-arch manifest from pushed arch-specific images
create-manifest:
    ARG VERSION=$VERSION
    ARG GITHUB_USER=pyrex41
    FROM alpine:3.23.2
    RUN apk --no-cache add docker-cli docker-cli-buildx

    # Wait for push-images to complete
    COPY +push-images/push-complete /tmp/push-complete

    RUN --secret GITHUB_TOKEN \
        echo "$GITHUB_TOKEN" | docker login ghcr.io -u "$GITHUB_USER" --password-stdin

    RUN docker buildx imagetools create \
        -t ${REGISTRY}/provider-signoz:${VERSION} \
        -t ${REGISTRY}/provider-signoz:latest \
        ${REGISTRY}/provider-signoz:${VERSION}-amd64 \
        ${REGISTRY}/provider-signoz:${VERSION}-arm64

# Build metadata-only xpkg (CRDs + crossplane.yaml, no embedded runtime).
# Controller image is referenced via spec.controller.image in crossplane.yaml.
package-build:
    FROM golang:1.24-alpine
    RUN apk --no-cache add curl
    # Install crossplane CLI
    RUN curl -fsSL "https://releases.crossplane.io/stable/v2.1.3/bin/linux_$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')/crank" \
        -o /usr/local/bin/crossplane \
        && chmod +x /usr/local/bin/crossplane
    WORKDIR /work
    COPY package/ package/
    RUN crossplane xpkg build \
        --package-root=package \
        -o package.xpkg
    SAVE ARTIFACT package.xpkg

# Push everything: controller images + multi-arch manifest + xpkg
push:
    ARG VERSION=$VERSION
    ARG GITHUB_USER=pyrex41
    ARG XPKG_TAG=xpkg

    # Step 1: Push arch-specific controller images
    BUILD +push-images --VERSION=$VERSION

    # Step 2: Create multi-arch manifest
    BUILD +create-manifest --VERSION=$VERSION --GITHUB_USER=$GITHUB_USER

    # Step 3: Build and push metadata-only xpkg
    FROM alpine:3.23.2
    RUN apk --no-cache add docker-cli curl
    RUN curl -fsSL "https://releases.crossplane.io/stable/v2.1.3/bin/linux_$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')/crank" \
        -o /usr/local/bin/crossplane \
        && chmod +x /usr/local/bin/crossplane

    COPY +package-build/package.xpkg /tmp/provider-signoz-package.xpkg

    RUN --secret GITHUB_TOKEN \
        echo "$GITHUB_TOKEN" | docker login ghcr.io -u "$GITHUB_USER" --password-stdin

    RUN crossplane xpkg push -f /tmp/provider-signoz-package.xpkg ${REGISTRY}/provider-signoz:${XPKG_TAG}

# Push just the xpkg (when controller images are already pushed)
xpkg-push:
    ARG GITHUB_USER=pyrex41
    ARG XPKG_TAG=xpkg
    FROM alpine:3.23.2
    RUN apk --no-cache add docker-cli curl
    RUN curl -fsSL "https://releases.crossplane.io/stable/v2.1.3/bin/linux_$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')/crank" \
        -o /usr/local/bin/crossplane \
        && chmod +x /usr/local/bin/crossplane
    COPY +package-build/package.xpkg /tmp/provider-signoz-package.xpkg
    RUN --secret GITHUB_TOKEN \
        echo "$GITHUB_TOKEN" | docker login ghcr.io -u "$GITHUB_USER" --password-stdin
    RUN crossplane xpkg push -f /tmp/provider-signoz-package.xpkg ${REGISTRY}/provider-signoz:${XPKG_TAG}

# Quick local build (amd64 only, no push)
package-local:
    BUILD +image --TARGETARCH=amd64
    BUILD +package-build
