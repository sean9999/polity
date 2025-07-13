#!/bin/bash

SESSION="main"
SESSION_EXISTS=$(tmux list-sessions | grep -w "$SESSION")

USER_1="blue-shadow"
USER_2="little-violet"
USER_3="quiet-bird"
USER_4="broken-hill"

if [ "$SESSION_EXISTS" = "" ]
then

	tmux new-session 	-d -s "$SESSION"
	tmux split-window	-v
	tmux split-window	-h
	tmux select-pane 	-t 1
	tmux split-window	-h

	tmux select-pane	-t "$SESSION":1.0
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_1.pem" Enter
	tmux select-pane	-t "$SESSION":1.1
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_2.pem" Enter
	tmux select-pane	-t "$SESSION":1.2
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_3.pem" Enter
	tmux select-pane	-t "$SESSION":1.3
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go -conf=testdata/$USER_4.pem" Enter

fi

tmux attach-session -t "$SESSION":1


