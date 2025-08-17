.PHONY: goo interface ebs_local ebs_twitch

goo:
	cd goo && go run .

interface:
	cd interface && npm start

ebs_local:
	@trap 'kill 0' INT TERM; \
	cd twitch-extension/frontend && npm run dev & \
	cd twitch-extension/ebs && go run .

ebs_twitch:
	@trap 'kill 0' INT TERM; \
	cd twitch-extension/frontend && npm run dev & \
	cd twitch-extension/ebs && go run .
