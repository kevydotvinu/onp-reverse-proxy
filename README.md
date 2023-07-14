# OpenShift Network Playground Reverse Proxy

A Golang powered reverse proxy for all the clusters that is deployed in OpenShift Network Playground.

## Installation
```
sudo git clone https://github.com/kevydotvinu/ocp-reverse-proxy /opt/openshift-network-playground/reverse-proxy
sudo cp /opt/openshift-network-playground/reverse-proxy/reverse-proxy.service /etc/systemd/system/reverse-proxy.service
sudo systemctl daemon-reload
```
