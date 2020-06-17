# piper ![Go](https://github.com/Mindgamesnl/piper/workflows/Go/badge.svg)
Piper is a service that deploys your code as you write it with a clear overview of what changed as well as its output. Useful for development where you have to run/test your code on a specific environment.
It watches your files for changes and uploads them to your server and automatically restarts the service, so you can continue debugging/testing.

You can also disable automatic uploads and only push files manually (by hitting enter when Piper detected your changes) or use it to maintain/setup production environments based on a local setup (you don't have to constantly reload code, you can also upload your entire package, restart the service and call it a day)

Cross platform with support for
- Windows
- MacOS
- Linux (+arm)

![screenshot](https://i.imgur.com/P76mffQ.png)

# Piper client
The Piper client is what you run on your local machine. Place the Piper binary in the root directory of your project and create a configuration YAML file. It should look like this
```yaml
# Server details for the host to use
server: localhost
port: 4723
password: nice

# If you want it to synchronize automatically, timeout is the interval where files will be uploaded (if changed) in seconds
# You can also chose to only do this manually, then set it to false and just use the ENTER key in the CUI
auto-sync: true
auto-sync-timeout: 5

# Commands and routines. Use `[]` if you want to disable them
pre-update-commands:
  - "echo STOPPING SERVICE"

post-update-commands:
  - "echo STARTING SERVICE! WELCOME BACK"

# The service that will be stopped before a update and resumed when all the files are in place
# You can just empty a small echo command if you dont want a service (like python, nodejs or go) to start/stop
service-command: "/usr/local/bin/node test.js"

# Folders that will be ignored
ignored-directories:
  - .idea
  - .git

# File extensions that will be synced
watched-extensions:
  - go
  - txt
  - js
  - yml
  - md
  - sh
  - php
  - py
  - jar
  - java
  - xml
```
To start the client, simply run `piper client ./my-config.yml` and you are ready to go. If you wish to set up a new environment and push all your code at once by adding the `--upload-all` flag (piper will look for all applicable files, upload them and then close again which is super useful to push to production)

# Piper server/host
Chose a directory where you want your code to be placed and executed from and place the Piper binary.
You need to create a server configuration file (just like for the client) with the following values
```yaml
port: 4723
password: nice
```
and can then start the server using `piper server ./my-config.yml`

