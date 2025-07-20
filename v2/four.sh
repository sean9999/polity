#!/bin/bash

SESSION="polityd"
SESSION_EXISTS=$(tmux list-sessions | grep -w "$SESSION")

USER_1="blue-shadow"
USER_2="little-violet"
USER_3="quiet-bird"
USER_4="broken-hill"

if [ "$SESSION_EXISTS" = "" ]
then

	tmux new-session 	-d -s "$SESSION"

	tmux split-window -h -t "$SESSION"        # Split pane 0 horizontally -> pane 1
	tmux split-window -v -t "$SESSION:0.0"    # Split pane 0 vertically -> pane 2
	tmux split-window -v -t "$SESSION:0.1"    # Split pane 1 vertically -> pane 3

	tmux select-layout tiled

	tmux select-pane	-t "$SESSION":0.0
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_1.pem" Enter
	tmux select-pane	-t "$SESSION":0.1
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_2.pem" Enter
	tmux select-pane	-t "$SESSION":0.2
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_3.pem" Enter
	tmux select-pane	-t "$SESSION":0.3
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_4.pem" Enter

fi

tmux attach-session -t "$SESSION":0


