# Stage 1: Build the binary
FROM containers.cisco.com/ngcs/hcps-go:develop AS build

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -ldflags='-s -w -X "main.version=0.0.1"' -o /bin/armor .

RUN wget https://github.com/coreruleset/coreruleset/archive/refs/tags/v4.7.0.tar.gz && \
    echo "13b652768a0ea3c3f5e1cae6bc25d520ef22943faf87564340751374e5df238b17526466cafd452846edfc436f6ad98cdd5ed5496531981c1ebebd2c179d22d7 v4.7.0.tar.gz" | sha512sum -c && \
    mkdir /etc/crs4 && \
    tar -xzvf v4.7.0.tar.gz --strip-components 1 -C /etc/crs4

# Stage 2: Copy the binary to a new image
FROM containers.cisco.com/sto-ccc-cloud9/hardened_alpine:latest

WORKDIR /app

COPY --from=build /bin/armor /bin/armor
COPY --from=build /etc/crs4/rules /etc/crs4/rules
COPY coraza /etc/crs4/
COPY config.json .

CMD ["/bin/armor"]
