deploy:
	go build
	cp potato_killer dist/potato_killer
	cp install.sh dist/install.sh
	cp upgrade.sh dist/upgrade.sh
	cp potato_killer.service dist/potato_killer.service
	cp config.json.example dist/config.json.example
	chmod +x dist/install.sh
	chmod +x dist/upgrade.sh