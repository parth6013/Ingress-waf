## Armor - WAF Controller

Many teams within ICS develop and maintain applications that serve a web UI. While web UIs offer convenience and accessibility,
they also introduce a range of security vulnerabilities that can be exploited by malicious actors. These vulnerabilities stem
from both the complexity of web technologies and the diverse range of potential attack vectors. Attacks like SQL Injection,
Cross-Site Scripting (XSS) and Cross-Site Request Forgery (CSRF) to name a few.

Protecting applications against such vulnerabilities requires constant thought and effort from the development teams.
Teams rely on security tools to help identify new attack vectors and then mitigate issues in a "reactive" mechanism.
Although this is good, we need solutions that help us stay secure in a "proactive" mechanism. The proposed solution to
help with this problem is to utilise technologies like Web Application Firewalls (WAF) and design a way to easily integrate
it across teams in ICS.

---

## Running The Project
You must have docker compose setup on your maching (Either through podman or via rancher desktop). Run the following docker compose command to build and start the services.

```bash
docker-compose up --build
```

This will start the WAF on http://localhost:3000, which will be proxying an test http server running on http://localhost:8080.
You can test the WAF by sending requests to http://localhost:3000.

---

## Running The Reflector Application Locally (Without Docker)
You will need python version 3.10 or above to run the reflector python application locally.

Setup the local environment by running the following commands:
```bash
cd reflector
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Now you can run the reflector application by running the following command:
```bash
python3 app.py
```

---

## Running Armor WAF Locally (Without Docker)
You will need go version 1.22.0 or above to run the Armor WAF locally.

Update the configuration file to point to the reflector application running locally. The configuration file is located at `armor/config.json`.
The following is an example cofniguration that points to the reflector application running locally on port 8080 and configures armor to run on port 3000 with enable CRS set to true.

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 3000
  },
  "waf": {
    "enableCrs": true
  },
  "proxy": {
    "targetHost": "localhost",
    "targetPort": 8080
  },
  "sam": {
    "database": "telegraf",
    "measurement": "armor"
  },
  "oidc": {
    "enable": false,
  }
}
```

You will also have to download the core ruleset from here: https://github.com/coreruleset/coreruleset/archive/refs/tags/v4.7.0.tar.gz and save the rules filder in `/etc/crs4/`.
You will also have to copy the contents of `armor/coraza` into `/etc/crs4/`

Now run the following command, which will build and run the application in one go which is useful during development.

```bash
cd armor
go run main.go
```

---

## TODO
https://ics-etherpad.cisco.com/p/hackathon-waf

---
