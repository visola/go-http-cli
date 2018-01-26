FROM ubuntu

RUN apt-get update

# Install latest certificates
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

# Install add-apt-repository
RUN apt-get install -y software-properties-common

# Install Go
RUN apt-get update
RUN apt-get install -y golang-1.9-go

# Install Java OpenJDK 8
RUN apt-get install -y openjdk-8-jdk

# Install Git
RUN apt-get install -y git

# Create appuser and switch to it
RUN useradd -c 'Application user' -m -d /home/appuser -s /bin/bash appuser
USER appuser

ENV USER=appuser
ENV HOME=/home/appuser
WORKDIR /home/appuser

# Clone repo
RUN mkdir -pv /home/appuser/goworkspace/src/github.com/visola
WORKDIR /home/appuser/goworkspace/src/github.com/visola
RUN git clone https://github.com/visola/go-http-cli.git

# Run build
WORKDIR /home/appuser/goworkspace/src/github.com/visola/go-http-cli
COPY --chown=appuser:appuser build.sh /home/appuser/goworkspace/src/github.com/visola/go-http-cli/build.sh
RUN ./build.sh
