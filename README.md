# Auto Publisher
This is a server that listens to the GitHub webhook for any change. Whenever a change is detected it runs following actions:
- Updates repository pointed by the environment variable
- Executes bash script in the specified directory

The intention for this is to update content of my blog whenever I submit anything to the blog repository and to publish new version. This was previously manual process for me so I thought it would be great to make it automatic.

The server is running as a systemctl daemon on my VPS.

I am very well aware that using go server to perform two simple actions is a slight overkill however there is a plan to extend this server to do following:
- Create a Discord bot that will communicate with this server and whenever new content is added it will post message on Discord with a commit message as a title