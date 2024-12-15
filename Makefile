.PHONY: goo interface ebs ebs_backend

goo:
	cd goo && go run .

interface:
	cd interface && npm start

ebs:
	@trap 'kill 0' INT TERM; \
	cd twitch-extension/frontend && npm start & \
	cd twitch-extension/ebs && go run .
