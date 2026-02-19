#!/usr/bin/env python3
import json
import sys
import os
import fcntl
from datetime import datetime, timezone

def log(message, level="INFO", agent="orchestrator", epic=None, task=None, data=None):
    log_entry = {
        "timestamp": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z",
        "level": level,
        "agent": agent,
        "message": message
    }
    
    if epic:
        log_entry["epic"] = epic
    if task:
        log_entry["task"] = task
    if data:
        log_entry["data"] = data

    # Write to agent-specific log
    log_dir = "logs"
    os.makedirs(log_dir, exist_ok=True)
    write_locked_log(os.path.join(log_dir, f"{agent}.log"), log_entry)
    
    # Also write to combined log
    write_locked_log(os.path.join(log_dir, "springfield.log"), log_entry)

def write_locked_log(filename, entry):
    line = json.dumps(entry) + "\n"
    with open(filename, "a") as f:
        try:
            fcntl.flock(f, fcntl.LOCK_EX)
            f.write(line)
        finally:
            fcntl.flock(f, fcntl.LOCK_UN)

    # Output to stdout for tmux capture (optional, but good for visibility)
    # print(json.dumps(log_entry))

def main():
    import argparse
    parser = argparse.ArgumentParser(description="Springfield Structured Logger")
    parser.add_argument("message", help="Log message")
    parser.add_argument("--level", default="INFO", help="Log level (default: INFO)")
    parser.add_argument("--agent", default="orchestrator", help="Agent name")
    parser.add_argument("--epic", help="Epic ID")
    parser.add_argument("--task", help="Task description")
    parser.add_argument("--data", help="Extra JSON data")

    args = parser.parse_args()
    
    extra_data = None
    if args.data:
        try:
            extra_data = json.loads(args.data)
        except json.JSONDecodeError:
            extra_data = {"raw": args.data}

    log(args.message, args.level, args.agent, args.epic, args.task, extra_data)

if __name__ == "__main__":
    main()
