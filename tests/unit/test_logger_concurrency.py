import json
import subprocess
import threading
import os
import shutil

def run_logger(i):
    subprocess.run(["./bin/log", "--agent", "test-agent", f"message {i}"])

def test_concurrency():
    if os.path.exists("logs"):
        shutil.rmtree("logs")
    
    threads = []
    for i in range(50):
        t = threading.Thread(target=run_logger, args=(i,))
        threads.append(t)
        t.start()
    
    for t in threads:
        t.join()
    
    # Check if logs/test-agent.log is valid JSON
    with open("logs/test-agent.log", "r") as f:
        lines = f.readlines()
        assert len(lines) == 50
        for line in lines:
            json.loads(line)
    
    print("Concurrency test passed!")

if __name__ == "__main__":
    test_concurrency()
