#!/usr/bin/env python3

import os
import time
import requests
import logging
import json
from datetime import datetime
from logging.handlers import RotatingFileHandler

# Configuration
# TODO: Update env to each NAS device
FOLDER_TO_WATCH = "/Users/phan.van.thanh/Desktop/CLIENTS"
WEBHOOK_URL = "http://localhost:38081/api/v1/tickets/add_auto"
LAST_RUN_FILE = "/Users/phan.van.thanh/Desktop/scripts/last_run.txt"
NAS_ID = 1
LOG_FILE = "/Users/phan.van.thanh/Desktop/scripts/watch_folders.log"

# Set the request headers
headers = {
    "accept": "*/*",
    "Content-Type": "application/json"
}

log_handler = RotatingFileHandler(
    filename=str(LOG_FILE),
    maxBytes=5*1024*1024,  # 5 MB per log file
    backupCount=5,          # Keep up to 5 backup files
    encoding='utf-8'
)

logging.basicConfig(
    handlers=[log_handler],
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)

def get_last_run_time():
    if os.path.exists(LAST_RUN_FILE):
        try:
            with open(LAST_RUN_FILE, 'r') as file:
                return float(file.read().strip())
        except ValueError as e:
            logging.error(f"Invalid timestamp in {LAST_RUN_FILE}: {e}")
            return 0
    logging.info(f"{LAST_RUN_FILE} does not exist. Assuming first run.")
    return 0

def save_current_time(timestamp):
    try:
        with open(LAST_RUN_FILE, 'w') as file:
            file.write(str(timestamp))
        logging.info(f"Updated last run time to {timestamp}.")
    except Exception as e:
        logging.error(f"Failed to write to {LAST_RUN_FILE}: {e}")


def find_new_folders(folder, last_run):
    new_folders = []
    for root, directories, _ in os.walk(folder):
        for directory_name in directories:
            directory_path = os.path.join(root, directory_name)
            try:
                if os.path.getmtime(directory_path) > last_run_time:
                    new_folders.append(directory_path)
                    logging.debug(f"New folder detected: {directory_path}")
            except FileNotFoundError:
                logging.warning(f"Directory {directory_path} was not found.")
            except Exception as e:
                logging.error(f"Error accessing {directory_path}: {e}")
    return new_folders

def send_webhook(new_folders):
    payload = {
        "nas_id": NAS_ID,
        "folders": new_folders
    }
    try:
        response = requests.post(WEBHOOK_URL, headers=headers, json=payload, timeout=10)
        response.raise_for_status()
        logging.info(f"Webhook sent successfully. Status code: {response.status_code}. Payload: {payload}")
    except requests.RequestException as e:
        logging.error(f"Failed to send webhook. Error: {e}. Payload: {payload}")
        log_curl_request(WEBHOOK_URL, headers, payload)

def log_curl_request(url, headers, payload):
    # Construct cURL command equivalent
    curl_command = f"curl -X POST '{url}' "
    for key, value in headers.items():
        curl_command += f"-H '{key}: {value}' "
    curl_command += f"-d '{json.dumps(payload)}'"
    logging.error(f"Equivalent cURL command for failed request: {curl_command}")

def main():
    logging.info("Script execution started.")
    current_time = time.time()
    last_run = get_last_run_time()

    # To prevent send a lot of existing dir in the first run
    is_first_run = last_run == 0

    new_folders = find_new_folders(FOLDER_TO_WATCH, last_run)

    if new_folders and not is_first_run:
        logging.info(f"Detected {len(new_folders)} new folder(s). Sending webhook.")
        send_webhook(new_folders)

    save_current_time(current_time)
    logging.info("Script execution finished.")

if __name__ == "__main__":
    main()