#!/bin/bash

SESSION="polityctl"
SESSION_EXISTS=$(tmux list-sessions | grep -w "$SESSION")

USER_1="quiet-bird"


if [ "$SESSION_EXISTS" = "" ]
then

	tmux new-session 	-d -s "$SESSION"
	tmux split-window -v -t "$SESSION"

	tmux select-pane	-t "$SESSION":0.0
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_1.pem" Enter
	tmux select-pane	-t "$SESSION":0.1
	tmux send-keys		-t "$SESSION" "\
	  sleep 1 && \
	  go run ./cmd/polityctl/*.go \
	  -conf=testdata/$USER_1.pem \
	  -join udp://e282a7ab4ed4854baa3fd30a4e87da86e69dcc0e945aff08906fbbb81b5b1b2a863a5bff18d7520962a903761582e6e3c2af042f457eb694fc3cf38304ae87e2@127.0.0.1:50971\
	" Enter

fi

tmux attach-session -t "$SESSION":0


