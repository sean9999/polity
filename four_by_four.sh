#!/bin/bash

SESSION="main"
SESSIONEXISTS=$(tmux list-sessions | grep -w "$SESSION")

USER_1="blue-shadow"
USER_2="little-violet"
USER_3="quiet-bird"
USER_4="broken-hill"

if [ "$SESSIONEXISTS" = "" ]
then

	tmux new-session 	-d -s "$SESSION"
	tmux split-window	-v
	tmux split-window	-h
	tmux select-pane 	-t 0
	tmux split-window	-h

	tmux select-pane	-t "$SESSION":0.0
	tmux send-keys		-t "$SESSION" "cd v2 && go run ./cmd/polityd/*.go -conf=testdata/$USER_1.pem" Enter
	tmux select-pane	-t "$SESSION":0.1
	tmux send-keys		-t "$SESSION" "cd v2 && go run ./cmd/polityd/*.go -conf=testdata/$USER_2.pem" Enter
	tmux select-pane	-t "$SESSION":0.2
	tmux send-keys		-t "$SESSION" "cd v2 && go run ./cmd/polityd/*.go -conf=testdata/$USER_3.pem" Enter
	tmux select-pane	-t "$SESSION":0.3
	tmux send-keys		-t "$SESSION" "cd v2 && go run ./cmd/polityd/*.go -conf=testdata/$USER_4.pem" Enter

fi

tmux attach-session -t "$SESSION":0


