#!/bin/sh

##  this is for apple silicon. you may wish to swap out "container" for "docker" or "podman", etc
##  or just start bare-metal redis locally
##  the main objective is to have a redis server running on port 6379

container run -p 6379:6379/tcp redis:8.4