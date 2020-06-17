# piper
A command-line tool to synchronise your workspace with a development server in real time and debug as you go.

# Project goals
The goal of piper is to be a pipe between your local code, and a development server that synchronizes your local files as you go with what's on the server and restarts the service.
Example workflow for a generic boring NodeJS web app
 - You have piper running on your host server and computer that you are working on
 - You edit one of the files (lets say `index.js`) in your favorite IDE
 - Piper notices the file change and replaces the one on the server
 - Piper restarts the nodejs service again and reports errors/stdout to your machine so you can see if it worked
 
 One of the main goals is to make Piper completely cross platform (so being able to use a Mac editor on a Linux host)
 
 mooi toch?