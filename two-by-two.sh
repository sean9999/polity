#!/bin/bash

SESSION="dockr"
SESSIONEXISTS=$(tmux list-sessions | grep -w "$SESSION")

USER_1="falling-wave"
USER_2="dawn-haze"

if [ "$SESSIONEXISTS" = "" ]
then

    tmux set-option -g default-shell /bin/bash
    tmux set -sg escape-time 50

	tmux new-session 	-d -s "$SESSION"
	tmux split-window	-h
	#tmux select-pane 	-t 0

	tmux select-pane	-t "$SESSION":0.0
	tmux send-keys		-t "$SESSION" "docker run polity -v $PWD/testdata:/ --config=/dark-haze.json" Enter
	tmux select-pane	-t "$SESSION":0.1
	tmux send-keys		-t "$SESSION" "docker run polity -v $PWD/testdata:/ --config=/billowing-water.json" Enter

fi

tmux attach-session -t "$SESSION":0
