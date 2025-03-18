# 🚀 ARM64 Docker Image Build for Raspberry Pi

This is the Docker build configuration to compile and package a Go application for the **linux/arm64** architecture, making it ready to run on a Raspberry Pi.

## ✅ Requirements
- Docker with **BuildKit** and **Buildx** enabled
- QEMU for cross-architecture builds
- Make (if using the provided Makefile)

## 🛠 QEMU Setup (One-time)
Register QEMU handlers to enable ARM emulation during the build:

```bash
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

## 🛠 Create and Use a Buildx Builder (One-time)
```bash
docker buildx create --name crossbuilder --use
docker buildx inspect --bootstrap
```

## 🔨 Build the ARM64 Docker Image
Use the following `buildx` command or the provided Makefile target to build the image for Raspberry Pi:

```bash
docker buildx build --platform=linux/arm64 -t <image-name>:<version> --load .
```

For debugging and detailed output:
```bash
docker buildx build --platform=linux/arm64 -t <image-name>:<version> --progress=plain --load .
```

## 📦 Save the Image to a Tar File (Optional)
You can export the image to a tarball for easier transfer to the Raspberry Pi:

```bash
docker save -o <image-name>_<version>_arm64.tar <image-name>:<version>
```

## 💡 Alternative: Native Go Cross-Compilation (Optional)
If you want to avoid QEMU during the build process, you can cross-compile the binary with Go and create a minimal Docker image:

### 1. Cross-compile the binary:
```bash
GOOS=linux GOARCH=arm64 go build -o my-app ./cmd
```

### 2. Use a minimal Dockerfile:
```dockerfile
FROM alpine:latest
COPY my-app /usr/local/bin/my-app
CMD ["my-app"]
```

### 3. Build the image:
```bash
docker buildx build --platform=linux/arm64 -t <image-name>:<version> --load .
```

## ✅ Notes
- `--load`: Loads the image into the local Docker daemon (useful for testing or saving).
- `--push`: Pushes the image directly to a container registry (optional).
- QEMU is **required** when `RUN` commands execute during the build and depend on the target architecture.
- Cross-compiling with Go (`GOOS` and `GOARCH`) is often faster if your build does not require architecture-specific dependencies.

---

## 🐳 Example Makefile Target
```makefile
dockerize:
	docker buildx build --platform=linux/arm64 -t $(CONTAINER_NAME):$(VERSION) --load .
	docker save -o $(CONTAINER_NAME)_$(VERSION)_arm64.tar $(CONTAINER_NAME):$(VERSION)
```
Run with:
```bash
make dockerize
```

