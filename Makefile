.PHONY: goo interface remote

goo:
	cd goo && go run .

interface:
	cd interface && npm start

# remote control page dev server, proxying /api to a local goo
remote:
	cd remote/frontend && npm run dev
