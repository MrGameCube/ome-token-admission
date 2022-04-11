# ome-token-admission
A Server which implements the Admission Webhooks used by OvenMediaEngine and provides a simple token-based authentication for Streams.
Admission Webhook specification: https://airensoft.gitbook.io/ovenmediaengine/access-control/admission-webhooks
---
## Installation Instructions (Linux)
### 1. Create a new user to run the server
> adduser ome-token --disabled-login
### 2. Create a Folder and set permissions
> mkdir /opt/ome-token

> chown ome-token:ome-token /opt/ome-token
### 3. Download and extract the latest release
> cd /opt/ome-token

> wget https://github.com/MrGameCube/ome-token-admission/releases/download/v0.1.0/linux-x86_64.zip

> unzip linux-x86_64.zip

> chmod +x ./ome-token-admission
### 4. Provide configuration values in config.ini
See config.sample.ini
### 5. Create a Systemd-Unit file
> /etc/systemd/system/ome-token.service

```
[Unit]
Description= ome-token-admission - Server for token-based admission to OvemMediaEngine
[Service]
WorkingDirectory= /opt/ome-token
Type=simple
ExecStart= /opt/ome-token/ome-token-admission
User=ome-token
Group=ome-token
[Install]
WantedBy=multi-user.target
```
> systemctl daemon-reload
> 
> systemctl enable ome-token
> 
> systemctl start ome-token
