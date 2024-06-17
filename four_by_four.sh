#!/bin/bash

SESSION="main"
SESSIONEXISTS=$(tmux list-sessions | grep -w "$SESSION")

if [ "$SESSIONEXISTS" = "" ]
then

	tmux new-session  -d -s "$SESSION"
	tmux split-window -v
	tmux split-window -h
	tmux select-pane  -t 0
	tmux split-window -h

	tmux select-pane	-t "$SESSION":0.0
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go --config=testdata/dawn-haze.toml" Enter
	tmux select-pane	-t "$SESSION":0.1
	tmux send-keys		-t "$SESSION" "go run ./cmd/polityd/*.go --config=testdata/holy-glade.toml" Enter
	tmux select-pane	-t "$SESSION":0.2
	tmux send-keys		-t "$SESSION" "alias dawn='go run ./cmd/polity/*.go --config=testdata/dawn-haze.toml'" Enter
	tmux select-pane	-t "$SESSION":0.3
	tmux send-keys		-t "$SESSION" "alias holy='go run ./cmd/polity/*.go --config=testdata/holy-glade.toml'" Enter

fi

tmux attach-session -t "$SESSION":0


#  tmux new-session \; \split-window -v \; \split-window -h \; \select-pane -t 0 \; \split-window -h