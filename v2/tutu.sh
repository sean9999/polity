#!/bin/bash

SESSION="polityctl"
SESSION_EXISTS=$(tmux list-sessions | grep -w "$SESSION")

USER_1="quiet-bird"
USER_2="little-violet"

if [ "$SESSION_EXISTS" = "" ]
then

	tmux new-session 	-d -s "$SESSION"
	tmux split-window -v -t "$SESSION"


	tmux select-pane	-t "$SESSION":0.0
	tmux send-keys		-t "$SESSION" "sleep 1 && go run ./cmd/polityd/*.go -conf=testdata/$USER_1.pem" Enter
	tmux select-pane	-t "$SESSION":0.1
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_2.pem" Enter

fi

tmux attach-session -t "$SESSION":0


