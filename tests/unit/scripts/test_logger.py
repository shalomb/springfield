import json
import os
import shutil
import sys
import subprocess
from unittest.mock import patch
from scripts.logger import log, main

def test_log_creates_files():
    if os.path.exists("logs"):
        shutil.move("logs", "logs_backup")
    
    try:
        test_agent = "test-agent"
        log_message = "test message"
        log(log_message, agent=test_agent)
        
        assert os.path.exists("logs")
        assert os.path.exists(f"logs/{test_agent}.log")
        assert os.path.exists("logs/springfield.log")
        
        with open(f"logs/{test_agent}.log", "r") as f:
            entry = json.loads(f.read().strip())
            assert entry["message"] == log_message
            assert entry["agent"] == test_agent
            assert "timestamp" in entry
    finally:
        if os.path.exists("logs"):
            shutil.rmtree("logs")
        if os.path.exists("logs_backup"):
            shutil.move("logs_backup", "logs")

def test_log_with_extra_fields():
    if os.path.exists("logs"):
        shutil.move("logs", "logs_backup")
        
    try:
        log("msg", level="DEBUG", agent="agent1", epic="EPIC-1", task="Task 1", data={"key": "value"})
        
        with open("logs/agent1.log", "r") as f:
            entry = json.loads(f.read().strip())
            assert entry["level"] == "DEBUG"
            assert entry["epic"] == "EPIC-1"
            assert entry["task"] == "Task 1"
            assert entry["data"] == {"key": "value"}
    finally:
        if os.path.exists("logs"):
            shutil.rmtree("logs")
        if os.path.exists("logs_backup"):
            shutil.move("logs_backup", "logs")

def test_cli_execution():
    if os.path.exists("logs"):
        shutil.move("logs", "logs_backup")
        
    try:
        test_args = [
            "scripts/logger.py", 
            "cli message", 
            "--level", "WARNING", 
            "--agent", "cli-agent",
            "--epic", "EPIC-CLI",
            "--task", "Task-CLI",
            "--data", '{"foo": "bar"}'
        ]
        with patch.object(sys, 'argv', test_args):
            main()
        
        assert os.path.exists("logs/cli-agent.log")
        with open("logs/cli-agent.log", "r") as f:
            entry = json.loads(f.read().strip())
            assert entry["message"] == "cli message"
            assert entry["level"] == "WARNING"
            assert entry["epic"] == "EPIC-CLI"
            assert entry["task"] == "Task-CLI"
            assert entry["data"] == {"foo": "bar"}
    finally:
        if os.path.exists("logs"):
            shutil.rmtree("logs")
        if os.path.exists("logs_backup"):
            shutil.move("logs_backup", "logs")

def test_cli_execution_invalid_json():
    if os.path.exists("logs"):
        shutil.move("logs", "logs_backup")
        
    try:
        test_args = [
            "scripts/logger.py", 
            "bad json", 
            "--agent", "bad-agent",
            "--data", "not-json"
        ]
        with patch.object(sys, 'argv', test_args):
            main()
        
        with open("logs/bad-agent.log", "r") as f:
            entry = json.loads(f.read().strip())
            assert entry["data"] == {"raw": "not-json"}
    finally:
        if os.path.exists("logs"):
            shutil.rmtree("logs")
        if os.path.exists("logs_backup"):
            shutil.move("logs_backup", "logs")

def test_main_execution_as_script():
    if os.path.exists("logs"):
        shutil.move("logs", "logs_backup")
        
    try:
        # Run the script as a subprocess to hit the 'if __name__ == "__main__":' block
        subprocess.run([
            "python3", "scripts/logger.py", 
            "main execution", 
            "--agent", "main-agent"
        ], check=True)
        
        assert os.path.exists("logs/main-agent.log")
    finally:
        if os.path.exists("logs"):
            shutil.rmtree("logs")
        if os.path.exists("logs_backup"):
            shutil.move("logs_backup", "logs")
