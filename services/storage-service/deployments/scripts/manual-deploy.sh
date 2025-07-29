#!/bin/bash

# Storage Service - Manual Step-by-Step Deployment Script
# This script provides an interactive, educational deployment experience

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="default"
TIMEOUT=300
INTERACTIVE=true

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${CYAN}[STEP]${NC} $1"
}

# Interactive prompt function
prompt_continue() {
    if [ "$INTERACTIVE" = true ]; then
        echo
        read -p "Press Enter to continue to the next step, or Ctrl+C to exit..."
        echo
    fi
}

# Show step header
show_step_header() {
    echo
    echo "==============================================="
    echo -e "${CYAN}$1${NC}"
    echo "==============================================="
    echo
}

# Wait for user confirmation
wait_for_confirmation() {
    local message="$1"
    if [ "$INTERACTIVE" = true ]; then
        echo -e "${YELLOW}$message${NC}"
        read -p "Type 'yes' to continue: " confirmation
        if [ "$confirmation" != "yes" ]; then
            log_error "Deployment cancelled by user"
            exit 1
        fi
    fi
}

# Verify prerequisites
verify_prerequisites() {
    show_step_header "Step 1: Verifying Prerequisites"
    
    log_step "Checking required tools..."
    
    # Check kubectl
    if command -v kubectl &> /dev/null; then
        log_success "kubectl is installed: $(kubectl version --client --short)"
    else
        log_error "kubectl is not installed. Please install kubectl and try again."
        exit 1
    fi
    
    # Check cluster connectivity
    log_step "Testing cluster connectivity..."
    if kubectl cluster-info &> /dev/null; then
        log_success "Connected to Kubernetes cluster: $(kubectl config current-context)"
    else
        log_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    # Check permissions
    log_step "Verifying permissions..."
    local permissions=(
        "create deployments"
        "create services"
        "create configmaps"
        "create secrets"
        "create horizontalpodautoscalers"
    )
    
    for perm in "${permissions[@]}"; do
        if kubectl auth can-i $perm &> /dev/null; then
            log_success "Permission verified: $perm"
        else
            log_warning "Permission missing: $perm"
        fi
    done
    
    # Check cluster resources
    log_step "Checking cluster resources..."
    log_info "Node status:"
    kubectl get nodes
    
    log_info "Available storage classes:"
    kubectl get storageclass 2>/dev/null || log_warning "No storage classes found"
    
    prompt_continue
}

# Deploy dependencies
deploy_dependencies() {
    show_step_header "Step 2: Deploying Dependencies (PostgreSQL & RabbitMQ)"
    
    log_step "Creating PostgreSQL deployment..."
    
    cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: "profiles"
        - name: POSTGRES_USER
          value: "profile_user"
        - name: POSTGRES_PASSWORD
          value: "profile_password"
        - name: PGDATA
          value: "/var/lib/postgresql/data/pgdata"
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: postgres-storage
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  selector:
    app: postgres
  ports:
  - name: postgres
    port: 5432
    targetPort: 5432
  type: ClusterIP
EOF
    
    log_success "PostgreSQL deployment created"
    
    log_step "Waiting for PostgreSQL to be ready..."
    kubectl wait --for=condition=available deployment/postgres --timeout=${TIMEOUT}s
    log_success "PostgreSQL is ready"
    
    log_step "Testing PostgreSQL connection..."
    if kubectl run postgres-test --rm -it --restart=Never --image=postgres:15-alpine -- psql -h postgres -U profile_user -d profiles -c "SELECT version();" 2>/dev/null; then
        log_success "PostgreSQL connection test passed"
    else
        log_warning "PostgreSQL connection test failed, but deployment continues"
    fi
    
    prompt_continue
    
    log_step "Creating RabbitMQ deployment..."
    
    cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rabbitmq
  labels:
    app: rabbitmq
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rabbitmq
  template:
    metadata:
      labels:
        app: rabbitmq
    spec:
      containers:
      - name: rabbitmq
        image: rabbitmq:3.11-management-alpine
        ports:
        - containerPort: 5672
          name: amqp
        - containerPort: 15672
          name: management
        env:
        - name: RABBITMQ_DEFAULT_USER
          value: "admin"
        - name: RABBITMQ_DEFAULT_PASS
          value: "password"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        volumeMounts:
        - name: rabbitmq-data
          mountPath: /var/lib/rabbitmq
      volumes:
      - name: rabbitmq-data
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: rabbitmq
  labels:
    app: rabbitmq
spec:
  selector:
    app: rabbitmq
  ports:
  - name: amqp
    port: 5672
    targetPort: 5672
  - name: management
    port: 15672
    targetPort: 15672
  type: ClusterIP
EOF
    
    log_success "RabbitMQ deployment created"
    
    log_step "Waiting for RabbitMQ to be ready..."
    kubectl wait --for=condition=available deployment/rabbitmq --timeout=${TIMEOUT}s
    log_success "RabbitMQ is ready"
    
    log_info "Dependencies deployed successfully!"
    log_info "PostgreSQL: Available at postgres:5432"
    log_info "RabbitMQ: Available at rabbitmq:5672 (Management: rabbitmq:15672)"
    
    prompt_continue
}

# Create configuration
create_configuration() {
    show_step_header "Step 3: Creating Configuration (ConfigMaps & Secrets)"
    
    log_step "Creating ConfigMap..."
    
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: storage-service-config
  labels:
    app: storage-service
data:
  config.yaml: |
    server:
      host: "0.0.0.0"
      port: 8080
      grpc_port: 9090
      read_timeout: 30s
      write_timeout: 30s
      idle_timeout: 120s

    database:
      max_connections: 100
      idle_connections: 20
      max_lifetime: 3600s
      connection_timeout: 30s

    queue:
      prefetch_count: 5
      process_timeout: 30s
      max_retries: 3
      reconnect_delay: 5s

    metrics:
      enabled: true
      port: 8080
      path: "/metrics"

    logging:
      level: "info"
      format: "json"

    circuit_breaker:
      enabled: true
      timeout: 10s
      max_requests: 100
      interval: 30s
      ratio: 0.6
EOF
    
    log_success "ConfigMap created"
    
    prompt_continue
    
    log_step "Creating Secrets..."
    
    wait_for_confirmation "This will create secrets with default development credentials. Is this acceptable for your environment?"
    
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: storage-service-secrets
  labels:
    app: storage-service
type: Opaque
data:
  # postgresql://profile_user:profile_password@postgres:5432/profiles
  database-url: cG9zdGdyZXNxbDovL3Byb2ZpbGVfdXNlcjpwcm9maWxlX3Bhc3N3b3JkQHBvc3RncmVzOjU0MzIvcHJvZmlsZXM=
  # amqp://admin:password@rabbitmq:5672/
  rabbitmq-url: YW1xcDovL2FkbWluOnBhc3N3b3JkQHJhYmJpdG1xOjU2NzIv
EOF
    
    log_success "Secrets created"
    
    log_step "Verifying configuration..."
    log_info "ConfigMap status:"
    kubectl describe configmap storage-service-config
    
    echo
    log_info "Secret status:"
    kubectl describe secret storage-service-secrets
    
    prompt_continue
}

# Deploy storage service
deploy_storage_service() {
    show_step_header "Step 4: Deploying Storage Service"
    
    wait_for_confirmation "This will deploy the Storage Service with 3 replicas. Continue?"
    
    log_step "Creating Storage Service deployment..."
    
    # Use the base Kubernetes manifests
    log_info "Applying Storage Service manifests..."
    kubectl apply -f ../kubernetes/deployment.yaml
    kubectl apply -f ../kubernetes/service.yaml
    
    log_success "Storage Service deployment created"
    
    log_step "Waiting for deployment to be ready..."
    kubectl wait --for=condition=available deployment/storage-service --timeout=${TIMEOUT}s
    
    log_success "Storage Service is ready!"
    
    log_step "Checking deployment status..."
    log_info "Pod status:"
    kubectl get pods -l app=storage-service
    
    echo
    log_info "Service status:"
    kubectl get services -l app=storage-service
    
    prompt_continue
}

# Verify deployment
verify_deployment() {
    show_step_header "Step 5: Verifying Deployment"
    
    log_step "Performing health checks..."
    
    # Port forward for testing
    log_info "Setting up port forwarding for testing..."
    kubectl port-forward svc/storage-service 8080:8080 &
    PORT_FORWARD_PID=$!
    
    # Wait for port forward to establish
    sleep 5
    
    # Test health endpoint
    log_step "Testing health endpoint..."
    if curl -f -s http://localhost:8080/health >/dev/null 2>&1; then
        log_success "Health endpoint is accessible"
        
        log_info "Health check response:"
        curl -s http://localhost:8080/health | jq . 2>/dev/null || curl -s http://localhost:8080/health
    else
        log_warning "Health endpoint test failed - this might be expected if the service is still starting"
    fi
    
    echo
    log_step "Testing metrics endpoint..."
    if curl -f -s http://localhost:8080/metrics >/dev/null 2>&1; then
        log_success "Metrics endpoint is accessible"
        
        log_info "Sample metrics:"
        curl -s http://localhost:8080/metrics | head -10
    else
        log_warning "Metrics endpoint test failed"
    fi
    
    # Clean up port forward
    kill $PORT_FORWARD_PID 2>/dev/null || true
    
    prompt_continue
    
    log_step "Checking logs..."
    log_info "Recent Storage Service logs:"
    kubectl logs -l app=storage-service --tail=20
    
    prompt_continue
}

# Setup monitoring (optional)
setup_monitoring() {
    show_step_header "Step 6: Setup Monitoring (Optional)"
    
    if [ "$INTERACTIVE" = true ]; then
        echo -e "${YELLOW}Would you like to set up Prometheus monitoring?${NC}"
        read -p "Type 'yes' to setup monitoring, or 'no' to skip: " setup_monitoring_choice
        
        if [ "$setup_monitoring_choice" = "yes" ]; then
            log_step "Creating ServiceMonitor for Prometheus..."
            
            if kubectl apply -f ../monitoring/servicemonitor.yaml 2>/dev/null; then
                log_success "ServiceMonitor created successfully"
                log_info "Prometheus will now scrape metrics from the Storage Service"
            else
                log_warning "ServiceMonitor creation failed - Prometheus Operator might not be installed"
                log_info "You can manually install it later using: kubectl apply -f ../monitoring/servicemonitor.yaml"
            fi
        else
            log_info "Skipping monitoring setup"
        fi
    else
        log_info "Skipping monitoring setup (non-interactive mode)"
    fi
    
    prompt_continue
}

# Show final information
show_final_info() {
    show_step_header "Deployment Complete!"
    
    log_success "Storage Service has been successfully deployed!"
    
    echo
    echo "=== Deployment Summary ==="
    echo "Service: Storage Service"
    echo "Namespace: $NAMESPACE"
    echo "Replicas: 3"
    echo "Dependencies: PostgreSQL, RabbitMQ"
    echo
    
    echo "=== Access Information ==="
    echo
    echo "1. Check pod status:"
    echo "   kubectl get pods -l app=storage-service"
    echo
    echo "2. View logs:"
    echo "   kubectl logs -l app=storage-service -f"
    echo
    echo "3. Access service locally:"
    echo "   kubectl port-forward svc/storage-service 8080:8080"
    echo "   curl http://localhost:8080/health"
    echo
    echo "4. Access PostgreSQL (debugging):"
    echo "   kubectl port-forward svc/postgres 5432:5432"
    echo
    echo "5. Access RabbitMQ Management (debugging):"
    echo "   kubectl port-forward svc/rabbitmq 15672:15672"
    echo "   http://localhost:15672 (admin/password)"
    echo
    echo "6. Scale the service:"
    echo "   kubectl scale deployment storage-service --replicas=5"
    echo
    echo "7. Update the service:"
    echo "   kubectl set image deployment/storage-service storage-service=storage-service:new-version"
    echo
    echo "=== Troubleshooting ==="
    echo
    echo "- Check events: kubectl get events --sort-by=.metadata.creationTimestamp"
    echo "- Describe pods: kubectl describe pods -l app=storage-service"
    echo "- Check resources: kubectl top pods -l app=storage-service"
    echo "- View all resources: kubectl get all -l app=storage-service"
    echo
    
    echo "=== Cleanup (when needed) ==="
    echo
    echo "To remove the deployment:"
    echo "   ./manual-cleanup.sh"
    echo
    
    log_success "Manual deployment completed successfully!"
}

# Main execution
main() {
    log_info "Starting Storage Service manual deployment..."
    echo
    echo "This script will guide you through deploying the Storage Service step by step."
    echo "Each step will be explained and you can review the results before continuing."
    echo
    
    if [ "$INTERACTIVE" = true ]; then
        read -p "Press Enter to begin, or Ctrl+C to exit..."
    fi
    
    verify_prerequisites
    deploy_dependencies
    create_configuration
    deploy_storage_service
    verify_deployment
    setup_monitoring
    show_final_info
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        --timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        --non-interactive)
            INTERACTIVE=false
            shift
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo
            echo "Options:"
            echo "  --namespace NAMESPACE    Kubernetes namespace (default: default)"
            echo "  --timeout SECONDS        Timeout for operations (default: 300)"
            echo "  --non-interactive        Run without user prompts"
            echo "  --help                   Show this help message"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Run main function
main "$@" 