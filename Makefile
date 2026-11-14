.PHONY: goo interface interface-local remote

goo:
	cd goo && go run .

interface:
	cd interface && npm start

# interface dev server, talking to a local goo (make goo) for http endpoints
interface-local:
	cd interface && REACT_APP_GOO_URL=http://localhost:8089 npm start

# remote control page dev server, proxying /api to a local goo
remote:
	cd remote/frontend && npm run dev
