# Armor Operator

A Kubernetes operator that automatically adds enterprise-grade security to your applications without code changes.

## Why Armor Operator?

Startups and teams testing market fit need to move fast, but security can't be an afterthought. Armor Operator bridges this gap by automatically adding production-ready security layers to any Kubernetes service with a single custom resource.

**The Problem:**
- Setting up WAF, authentication, and logging takes days or weeks
- Security expertise is expensive and hard to find
- Early-stage teams need to focus on product, not infrastructure
- Manual security configurations are error-prone and hard to maintain

**The Solution:**
Armor Operator turns security into a one-liner. Point it at your service, and it automatically deploys:

## Features

### 🛡️ Web Application Firewall (WAF)
- **Coraza WAF** with **OWASP Core Rule Set (CRS) 4.0**
- Protects against: SQL injection, XSS, RCE, LFI, RFI, session fixation, and 100+ attack patterns
- Zero application code changes required
- Configurable rule paths for different environments

### 🔐 Authentication & Authorization
- **OpenID Connect (OIDC)** integration
- Works with Keycloak, Auth0, Okta, and any OIDC provider
- Automatic session management
- Configurable redirect URLs and client credentials

### 📊 Centralized Logging
- **Elasticsearch** integration for all security events
- Real-time attack visibility and analytics
- Request/response logging for compliance
- Optional - enable/disable as needed

### 🚀 Zero-Configuration Deployment
- Declarative Kubernetes CRDs
- Automatic ConfigMap, Deployment, Service, and Ingress creation
- Managed by the operator - updates propagate automatically
- Works with any Kubernetes service

## Quick Example

```yaml
apiVersion: apps.armor.io/v1
kind: ArmorProxy
metadata:
  name: armorproxy-reflector
  namespace: default
  labels:
    app.kubernetes.io/name: armor-operator
    app.kubernetes.io/managed-by: kustomize
spec:
  # Name of the service to protect (the backend application)
  serviceName: reflector
  
  # Namespace where the backend service exists
  serviceNamespace: armor-reflector
  
  # Optional: Hostname for Ingress (leave empty if not using Ingress)
  ingressHost: armor.local
  
  # Optional: Ingress class name (e.g., nginx, traefik)
  ingressClassName: nginx
  
  # Enable WAF protection with Coraza + OWASP CRS
  waf: true
  
  # Enable OIDC authentication
  oidc: true
  
  # OIDC Configuration (required when oidc: true)
  oidcConfig:
    clientID: armor-proxy
    providerUrl: http://keycloak.armor-reflector.svc.cluster.local:8080/realms/Armor
    redirectApplication: http://armor-reflector.armor-reflector.svc.cluster.local:3000
  
  # Path to CRS rules inside the Armor container
  # Default: /crs4 (for containers)
  # Use ./crs4 for local development
  crsRulesPath: /crs4
  
  # Elasticsearch configuration for centralized logging
  elasticsearch:
    # Enable/disable Elasticsearch logging
    enable: true
    
    # Elasticsearch endpoint(s)
    url:
      - http://elasticsearch.armor-reflector.svc.cluster.local:9200
```

That's it. Your application is now protected by enterprise-grade security.

## Who Is This For?

- **Startups** testing product-market fit who need security without the overhead
- **Development teams** deploying internal tools that need basic auth and logging
- **Platform engineers** building secure-by-default infrastructure
- **Anyone** who wants WAF + Auth + Logging without manually configuring nginx, OAuth flows, and log shippers

## Quick Setup 
### Step 1: Install the Operator

```bash
# Clone the repository
git clone https://github.com/parth6013/Ingress-waf.git
cd Ingress-waf/armor-operator

# Install CRDs and deploy the operator
make manifests generate
make install
make deploy IMG=docker.io/parth6013/armor-operator:latest
```

**Or run locally for development:**
```bash
make manifests generate
make run
```

### Step 2: Create an ArmorProxy Resource

```bash
# Create a file: my-armor.yaml
cat <<EOF | kubectl apply -f -
apiVersion: apps.armor.io/v1
kind: ArmorProxy
metadata:
  name: my-secure-app
  namespace: default
spec:
  serviceName: my-app-service
  serviceNamespace: default
  waf: true
  oidc: false
  crsRulesPath: /crs4
  elasticsearch:
    enable: false
EOF
```

### Step 3: Verify

```bash
# Check the ArmorProxy resource
kubectl get armorproxy -n default

# Check created resources
kubectl get all,configmap -n armor -l app.kubernetes.io/managed-by=armor-operator
```

The operator automatically creates:
- ✅ ConfigMap with complete Armor configuration
- ✅ Deployment running the Armor WAF proxy
- ✅ Service exposing the protected application
- ✅ Ingress (if `ingressHost` is specified)

**Access your protected application:**
```bash
# If using Ingress
curl http://armor.local

# Or port-forward
kubectl port-forward -n armor svc/armor-<serviceName> 3000:3000
curl http://localhost:3000
```

---

## Real-World Examples

### Example 1: OIDC Authentication Protection

**Without Armor:**
Your company website is publicly accessible without authentication.

![Direct access without authentication](image.png)

**With Armor:**
Users are automatically redirected to your OIDC provider for authentication before accessing the application.

![OIDC login enforced](image-1.png)

---

### Example 2: Attack Protection

#### Shell Injection Attack

**Without Armor:**
The application is vulnerable to shell injection attacks, exposing sensitive system paths.

![Shell injection attempt](image-2.png)

![System path exposed](image-4.png)

**With Armor:**
The WAF detects and blocks the shell injection attack immediately.

![Attack blocked by Armor](image-3.png)

![Attack detection confirmation](image-5.png)

---

#### SQL Injection Attack

**Without Armor:**
The application is vulnerable to SQL injection attacks.

![SQL injection attempt](image-6.png)

![SQL injection success](image-7.png)

**With Armor:**
All SQL injection attempts are blocked by the Coraza WAF with OWASP CRS rules.

---

### Example 3: Centralized Logging & Analytics

**Elasticsearch Integration:**
All security events are automatically logged to Elasticsearch for analysis.

![Security logs in Elasticsearch](image-8.png)

**Kibana Analytics:**
Visualize attack patterns, identify vulnerabilities, and monitor your security posture in real-time.

![Kibana dashboard showing attack analytics](image-9.png)

With Armor's centralized logging, you can:
- Track all blocked attacks
- Identify attack patterns and trends
- Generate compliance reports
- Improve your security posture based on real data

---

*Built with [Kubebuilder](https://kubebuilder.io/) and powered by [Coraza WAF](https://coraza.io/)*
