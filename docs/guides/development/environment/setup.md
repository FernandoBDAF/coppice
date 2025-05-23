# Local Development Environment Setup

-> IMPORTANT: Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.

## Primary Purpose and Main Goals

### Primary Purpose

This guide provides step-by-step instructions for setting up a local development environment for the Profile Service Microservices project, ensuring all developers have a consistent and functional development setup.

### Main Goals

1. Ensure consistent development environment
2. Minimize setup time and issues
3. Provide clear prerequisites
4. Document configuration steps
5. Enable quick verification

## Prerequisites

### System Requirements

1. **Operating System**

   - macOS 10.15 or later
   - Linux (Ubuntu 20.04 or later)
   - Windows 10/11 with WSL2

2. **Hardware Requirements**

   - CPU: 4 cores minimum
   - RAM: 8GB minimum
   - Storage: 20GB free space

3. **Required Software**
   - Git 2.30.0 or later
   - Docker 20.10.0 or later
   - Docker Compose 2.0.0 or later
   - Go 1.19 or later
   - Node.js 16.x or later
   - kubectl 1.22 or later
   - minikube 1.24 or later

## Installation Steps

### 1. Development Tools

#### Git Setup

```bash
# Install Git
# macOS
brew install git

# Ubuntu
sudo apt-get update
sudo apt-get install git

# Windows (with Chocolatey)
choco install git
```

#### Docker Setup

```bash
# Install Docker Desktop
# macOS
brew install --cask docker

# Ubuntu
sudo apt-get update
sudo apt-get install docker.io docker-compose

# Windows
# Download and install Docker Desktop from https://www.docker.com/products/docker-desktop
```

#### Go Setup

```bash
# Install Go
# macOS
brew install go

# Ubuntu
sudo apt-get update
sudo apt-get install golang-go

# Windows
# Download and install from https://golang.org/dl/
```

### 2. Kubernetes Setup

#### Minikube Installation

```bash
# Install Minikube
# macOS
brew install minikube

# Ubuntu
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Windows
# Download and install from https://minikube.sigs.k8s.io/docs/start/
```

#### kubectl Installation

```bash
# Install kubectl
# macOS
brew install kubectl

# Ubuntu
sudo apt-get update
sudo apt-get install kubectl

# Windows
# Download and install from https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/
```

## Configuration

### 1. Git Configuration

```bash
# Configure Git
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### 2. Docker Configuration

```bash
# Start Docker
# macOS/Windows
# Launch Docker Desktop

# Ubuntu
sudo systemctl start docker
sudo systemctl enable docker
```

### 3. Minikube Configuration

```bash
# Start Minikube
minikube start --driver=docker

# Verify installation
kubectl cluster-info
```

### 4. Project Setup

```bash
# Clone the repository
git clone https://github.com/your-org/profile-service.git
cd profile-service

# Install dependencies
go mod download
npm install
```

## Verification

### 1. Environment Check

```bash
# Verify installations
git --version
docker --version
docker-compose --version
go version
node --version
kubectl version
minikube version
```

### 2. Service Verification

```bash
# Start services
docker-compose up -d

# Verify services
docker-compose ps
```

### 3. Kubernetes Verification

```bash
# Verify Kubernetes cluster
kubectl get nodes
kubectl get pods --all-namespaces
```

## Common Issues and Solutions

### 1. Docker Issues

- **Problem**: Docker daemon not running

  - **Solution**: Start Docker Desktop or run `sudo systemctl start docker`

- **Problem**: Permission denied
  - **Solution**: Add user to docker group: `sudo usermod -aG docker $USER`

### 2. Minikube Issues

- **Problem**: Minikube fails to start

  - **Solution**: Check virtualization settings in BIOS

- **Problem**: kubectl connection refused
  - **Solution**: Ensure minikube is running: `minikube status`

### 3. Go Issues

- **Problem**: Go modules not found

  - **Solution**: Run `go mod tidy`

- **Problem**: GOPATH not set
  - **Solution**: Set GOPATH in your shell profile

## Next Steps

1. Review the [Development Process](workflow/process.md) guide
2. Set up your [IDE](ide/setup.md)
3. Read the [Testing Guide](testing/guide.md)
4. Familiarize yourself with the [Debugging Guide](debugging/guide.md)

## Notes

- Keep all tools updated
- Follow security best practices
- Document any custom configurations
- Report issues to the team

## Version History

### Current Version

- Version: To be determined
- Date: To be determined
- Changes:
  - Initial environment setup guide
  - Installation steps documented
  - Configuration instructions added
  - Verification steps included
