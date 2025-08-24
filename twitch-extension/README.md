# Milkelangelo Twitch extension

## Running locally

0. Stop Goo on milkelangelo to avoid interference
1. Run the ebs frontend and backend with `make ebs_local` in project root.
2. Run goo locally with `cd goo && go run .`. It will connect to the local ebs
3. Go to `localhost:3000` to view the extension.

You could test with the deployed goo by setting its EBS_HOST to your tailscale hostname so it connects to your ebs instance.