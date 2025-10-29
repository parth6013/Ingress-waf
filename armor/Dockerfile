# Stage 1: Build the binary
FROM golang:1.22 AS build

# Set Go proxy (optional: use "direct" if behind a firewall or corporate network)
ENV GOPROXY=https://proxy.golang.org,direct

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 go build -ldflags='-s -w -X "main.version=0.0.1"' -o /bin/armor .

# Uncomment this block if needed for CRS setup
# RUN wget https://github.com/coreruleset/coreruleset/archive/refs/tags/v4.7.0.tar.gz && \
#     echo "13b652768a0ea3c3f5e1cae6bc25d520ef22943faf87564340751374e5df238b17526466cafd452846edfc436f6ad98cdd5ed5496531981c1ebebd2c179d22d7 v4.7.0.tar.gz" | sha512sum -c && \
#     mkdir /etc/crs4 && \
#     tar -xzvf v4.7.0.tar.gz --strip-components 1 -C /etc/crs4

# Stage 2: Minimal runtime image
FROM alpine:3.19

WORKDIR /app

# Copy the built binary
COPY --from=build /bin/armor /bin/armor

# Copy configuration and rules
COPY crs4 /crs4
COPY config.json .

# Run the binary
CMD ["/bin/armor"]
