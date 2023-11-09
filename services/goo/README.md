# goo - go orchestrator

This isn't really an orchestrator. It's a place for the following and other
related features. Hopefully doesn't become too gooey.

- [ ] Store / make available the Piece data.
- [ ] Store / make available the Events.
- [ ] Connect with the cloud system over gRPC, facilitating communication.
- Perhaps, eventually replace pygateway and have controller<>goo be proto-based.

## Twitch auth

Set up the twitch cli and do

```bash
twitch token -u -s 'chat:read chat:edit moderator:manage:announcements'
```

Save the refresh token in `kv/TWITCH_REFRESH_TOKEN`