SESSION="four_panes"

# Create a new detached tmux session
tmux new-session -d -s "$SESSION"

# Split the window into 4 equal-sized panes (2x2 grid)
tmux split-window -h -t "$SESSION:0"      # Split pane 0 horizontally -> pane 1
tmux split-window -v -t "$SESSION:0.0"    # Split pane 0 vertically -> pane 2
tmux split-window -v -t "$SESSION:0.1"    # Split pane 1 vertically -> pane 3

# Send a unique command to each pane
tmux send-keys -t "$SESSION:0.0" "echo 'Pane 0: Monitoring logs'; tail -f /var/log/syslog" C-m
tmux send-keys -t "$SESSION:0.1" "echo 'Pane 1: Running top'; top" C-m
tmux send-keys -t "$SESSION:0.2" "echo 'Pane 2: Pinging google'; ping google.com" C-m
tmux send-keys -t "$SESSION:0.3" "echo 'Pane 3: Watching disk usage'; watch df -h" C-m

tmux select-layout tiled

# Attach to the session
tmux attach-session -t "$SESSION"

