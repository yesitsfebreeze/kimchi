FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=C.UTF-8

# Install system deps
RUN apt-get update && apt-get install -y \
  curl git make clang python3 build-essential cmake pkg-config ca-certificates gnupg && \
  rm -rf /var/lib/apt/lists/*

# Install Node.js 18 (from Nodesource)
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
  apt-get install -y nodejs && \
  npm install -g tree-sitter-cli

# Install emsdk (Emscripten)
RUN git clone https://github.com/emscripten-core/emsdk.git /emsdk
WORKDIR /emsdk
RUN ./emsdk install latest && ./emsdk activate latest

# Setup environment for future steps
ENV PATH="/emsdk:/emsdk/upstream/emscripten:$PATH"
ENV EMSDK="/emsdk"
ENV EM_CONFIG="/emsdk/.emscripten"

# Final working dir
WORKDIR /src
ENTRYPOINT ["/bin/bash"]
