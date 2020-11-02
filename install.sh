#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

mkdir -p /opt/potato_killer
cp potato_killer /opt/potato_killer/potato_killer
cp config.json.example /opt/potato_killer/config.json.example
cp potato_killer.service /opt/potato_killer/potato_killer.service
cp /opt/potato_killer/potato_killer.service /etc/systemd/system/potato_killer.service

chown pufferpanel:pufferpanel -R /opt/potato_killer

systemctl enable potato_killer.service

