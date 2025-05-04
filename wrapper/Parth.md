Undertsanding How Middleware works

1) We create one New RP ( Reverse Proxy what is rp when request is recieved to server intead going to main service(target)) it goes through proxy logic

Creating Oidc Server
1) mkdir -p ./keycloak-data
Authorization-> Checks authorization policies

User is redirected to Keycloak for login.
After login, Keycloak sends an authorization code to the client.
The client exchanges this code for an access token (securely).
2) 

Run this command
podman run -p 3003:8080 \                
 -e KEYCLOAK_ADMIN=admin \  
 -e KEYCLOAK_ADMIN_PASSWORD=admin \                                         
 -v $(pwd)/keycloak-data:/opt/keycloak/data \
 quay.io/keycloak/keycloak:latest \
 start-dev

podman run -p 3003:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin -v $(pwd)/keycloak-data:/opt/keycloak/data quay.io/keycloak/keycloak:latest start-dev

QnA
Why we do state check

An attacker can generate an OAuth authorization code, but they can’t access the victim’s browser cookies, so they can’t forge a correct state value.
The victim’s browser only sends the state value to legitimate requests initiated from their session, preventing CSRF attacks.
Even if an attacker tricks a victim into clicking a malicious link, the state validation stops the attack.

Use Of Nonce
✅ Ensures that the authentication request and response belong to the same session.
✅ Prevents replay attacks where an attacker tries to reuse an old authentication response.


The session is stored in a cookie, but the script itself does not explicitly set an expiration time.
The expiration is inherited from the ID token issued by the OIDC provider.(keycloak-data in this case)

partbhar
test1234

____________________________
What are we trying to Achieve
1) So we arec reating Middlewares for other software to use on the cluster
2) We want to enable middleware to track security defaults logging and other testing things for this we want to have a operator/controller that will do this thing inject sidecars on the pods or setup a proxy


3) I will create a controller which will decide to create these things as a sidecar or proxy

______________

Let's get back to all facilities

Time Logging failure
1) Dashboard for Api calls and failures at something

__
1) Logging as a Service what we'll do
-> In progress

1) Run elasticsearch then
podman exec -it elasticsearch bin/elasticsearch-service-tokens create kibana kibana-token

then use that in kibana
1) Open Kibana http://localhost:5601
2) Go to the Kibana "Discover" tab:
3) Click on "Create index pattern logs*
4) 

_____________

Coraza What it does
Intercept all HTTP requests
Apply WAF rules to detect security threats
Log any violations to a monitoring system
Either block malicious requests or allow safe ones to proceed to the next handler

_____________

Coraza Rules

coraza.NewWAFConfig() - Creates a new configuration object for the WAF.
.WithErrorCallback(logError) - Sets up a callback function (logError) that will be called whenever the WAF detects a rule violation. This is how the system knows what to do when it detects an attack - in this case, it logs the error and sends information to a monitoring system.
.WithDirectivesFromFile("/etc/crs4/coraza.conf") - Loads the main WAF configuration from a file. This file contains basic settings and rules for the WAF.
The conditional block if config.Waf.EnableCrs { ... } checks if the Core Rule Set (CRS) is enabled in the application configuration.
If CRS is enabled, it loads two additional sets of files:
/etc/crs4/crs-setup.conf: The main setup file for the Core Rule Set
/etc/crs4/rules/*.conf: All rule files in the rules directory
The Core Rule Set (CRS) is a set of generic attack detection rules that protect web applications from a wide range of attacks, including the OWASP Top Ten. By making this conditional, the application allows for flexibility:

Basic mode: Just use the main Coraza configuration
Enhanced mode: Use the main config plus the comprehensive CRS rules
This pattern uses method chaining (each method returns the config object) to build up the complete WAF configuration before it's used to create the actual WAF instance in the next part of the code.



1) sql injection
1 or 1=1 ( the or 1=1 is always true and ) expose the whole databse


Learning I was sending to big of a message hence trimmed the message


2) attack 

http://localhost:8080/echo?key=%3Cscript%3Ealert(%27hacked%27)%3C/script%3E

3) shell injection employee id 
data=foo; ls


Kubernetes

1) minikube start --driver=podman --container-runtime=containerd --network=host

2) create images through dockerfile and push

3) create resources

Create a nginx controller

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install nginx-ingress ingress-nginx/ingress-nginx --set controller.publishService.enabled=true

2) create ingress resources
You have an NGINX Ingress Controller running in Minikube.

You have created Ingress rules mapping armor.local → your armor service (port 3000).

You port-forwarded the Ingress Controller service to your local localhost:80.

sudo kubectl port-forward svc/nginx-ingress-ingress-nginx-controller 80:80 -n default

armor.local should point to 127.0.0.1 — your local machine loopback

127.0.0.1   armor.local
127.0.0.1   reflector.local

in your /etc/hosts

3) We will create a controller which will
1) If user creates a crd of this type 
deployment name and the web service port so the controller will create armor deployment and service and a ingress resource for it