tmux new-window -c ~
tmux split-window -h -c /tmp
tmux split-window -v -c /
tmux select-pane -t 1
tmux split-window -v -c /home
tmux select-pane -t 1

