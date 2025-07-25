#!/bin/bash

SESSION="main"
SESSIONEXISTS=$(tmux list-sessions | grep -w "$SESSION")

USER_1="billowing-water"
USER_2="dark-haze"

if [ "$SESSIONEXISTS" = "" ]
then

	tmux new-session 	-d -s "$SESSION"
	tmux split-window	-v
	tmux split-window	-h
	tmux select-pane 	-t 0
	tmux split-window	-h

	tmux select-pane	-t "$SESSION":0.0
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go --config=testdata/$USER_1.json" Enter
	tmux select-pane	-t "$SESSION":0.1
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go --config=testdata/$USER_2.json" Enter
	tmux select-pane	-t "$SESSION":0.2
	tmux send-keys		-t "$SESSION" "alias $USER_1='go run ./cmd/polity/*.go --config=testdata/$USER_1.json'" Enter
	tmux select-pane	-t "$SESSION":0.3
	tmux send-keys		-t "$SESSION" "alias $USER_2='go run ./cmd/polity/*.go --config=testdata/$USER_2.json'" Enter

fi

tmux attach-session -t "$SESSION":0
