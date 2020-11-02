#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi


systemctl stop potato_killer

cp potato_killer /opt/potato_killer/potato_killer

chown pufferpanel:pufferpanel -R /opt/potato_killer

systemctl start potato_killer