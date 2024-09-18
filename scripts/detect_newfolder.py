#!/usr/bin/env python3

import os
import time
import requests
from datetime import datetime

# Configuration
FOLDER_TO_WATCH = "/volume5/FOR DEVELOPER/CLIENTS"
WEBHOOK_URL = "https://webhook.site/7852d57c-0594-4525-a795-e4bd5a20625e"
LAST_RUN_FILE = "/volume5/FOR DEVELOPER/scripts/last_run.txt"

def get_last_run_time():
    if os.path.exists(LAST_RUN_FILE):
        with open(LAST_RUN_FILE, 'r') as f:
            return float(f.read().strip())
    return 0

def save_current_time(timestamp):
    with open(LAST_RUN_FILE, 'w') as f:
        f.write(str(timestamp))

def find_new_folders(folder, last_run):
    new_folders = []
    for root, dirs, _ in os.walk(folder):
        for dir_name in dirs:
            dir_path = os.path.join(root, dir_name)
            if os.path.getmtime(dir_path) > last_run:
                new_folders.append(dir_path)
    return new_folders

def send_webhook(new_folders):
    payload = {
        "message": "New folders detected",
        "folders": "\n".join(new_folders)
    }
    try:
        response = requests.post(WEBHOOK_URL, json=payload)
        response.raise_for_status()
        print(f"Webhook sent successfully. Status code: {response.status_code}")
    except requests.RequestException as e:
        print(f"Failed to send webhook: {e}")

def main():
    current_time = time.time()
    last_run = get_last_run_time()

    new_folders = find_new_folders(FOLDER_TO_WATCH, last_run)

    if new_folders:
        send_webhook(new_folders)

    save_current_time(current_time)

if __name__ == "__main__":
    main()"