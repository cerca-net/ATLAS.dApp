#!/bin/bash

echo "🛑 Stopping Blockchain Test Environment..."

# Kill blockchain processes
echo "🔪 Killing blockchain processes..."
pkill -f "blockchain.exe" || true
pkill -f "go run" || true

# Kill processes by PID files
for pid_file in test_nodes/*/pid; do
    if [ -f "$pid_file" ]; then
        pid=$(cat "$pid_file")
        echo "Killing process $pid"
        kill $pid 2>/dev/null || true
        rm "$pid_file"
    fi
done

# Clean up test directories
echo "🧹 Cleaning up test directories..."
rm -rf test_nodes
rm -rf test_logs

echo "✅ Test environment stopped and cleaned up!" 