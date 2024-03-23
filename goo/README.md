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

## Youtube auth

### Errors

- `oauth2: "invalid_grant" "Bad Request"`
	- Delete the `youtube-credentials-cache.json` file
	- Locally, **With /mnt/md0/light-stores mounted**, run `cd goo && go run . -yt`
	- Follow the link and the instructions
	- After the redirect, copy the `code=` value and paste it in terminal.

Not sure if we should be updating the credentials cache file after the token gets refreshed.
If the above error happens every week, this needs more work.

