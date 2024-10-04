# Twitch EBS

Twitch "Extension Backend Service" for the twitch-extension. This is required to proxy the internal platform.

Note this is public-facing so requires auth and spam prevention techniques.

## API

see openapi/openapi.yaml

## Answers

- Architecture?
  - Internal platform (goo), running on-prem connects to the EBS if available.
	- goo listens to EBS for movement requests
  - Twitch extension running on viewer's browsers may connect to EBS if they are controlling the session.
	- Twitch extension POSTs movement requests to EBS
  - All data receipt / feedback happens through Twitch's pubsub system to aid scale.
- Why not use an mqtt bridge?
  - Because initially this was a voting system and votes didn't exist as broker topics.
  - Decoupling from the broker is nice so this can be run independently of my VPS. In a fly.io pod that is only on when machine is running.

## Questions

- What should CORS include? Twitch?