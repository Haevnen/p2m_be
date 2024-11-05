#!/usr/bin/env python3

import os
import time
import requests
import logging
import json
from datetime import datetime
from logging.handlers import RotatingFileHandler

# Configuration
FOLDERS_TO_WATCH = [
    "/volume3/FILES STATION/CLIENTS",
    "/volume3/BRH"
]
WEBHOOK_URL = "https://api.sharpstudio.id.vn/api/v1/tickets/add_auto"
LAST_RUN_FILE = "/volume3/FILES STATION/scripts/last_run.txt"
NAS_ID = 1
LOG_FILE = "/volume3/FILES STATION/scripts/watch_folders.log"

# Set the request headers
headers = {
    "accept": "*/*",
    "Content-Type": "application/json"
}

log_handler = RotatingFileHandler(
    filename=str(LOG_FILE),
    maxBytes=10*1024*1024,  # 5 MB per log file
    backupCount=5,          # Keep up to 5 backup files
    encoding='utf-8'
)

logging.basicConfig(
    handlers=[log_handler],
    level=logging.DEBUG,
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

def find_new_folders(folders, last_run):
    new_folders = set()
    for folder in folders:
        logging.info(f"Scanning for new folders in {folder}...")
        try:
            for root, directories, _ in os.walk(folder):
                for directory_name in directories:
                    directory_path = os.path.join(root, directory_name)

                    if 'DAV' in directory_path or 'DONE' in directory_path or 'UPLOAD' not in directory_path:
                        continue

                    try:
                        if os.path.getmtime(directory_path) > last_run:
                            new_folders.add(directory_path)
                            logging.debug(f"New folder detected: {directory_path}")
                    except FileNotFoundError:
                        logging.warning(f"Directory {directory_path} was not found.")
                    except Exception as e:
                        logging.error(f"Error accessing {directory_path}: {e}")
        except Exception as e:
            logging.error(f"Error scanning folder {folder}: {e}")               
    return list(new_folders)

def send_webhook(new_folders):
    payload = {
        "nas_id": NAS_ID,
        "folders": new_folders
    }
    curl_command = log_curl_request(WEBHOOK_URL, headers, payload)
    try:
        response = requests.post(WEBHOOK_URL, headers=headers, json=payload, timeout=10)
        response.raise_for_status()
        logging.info(f"Webhook sent successfully. Status code: {response.status_code}. Payload: {payload}")
    except requests.RequestException as e:
        logging.error(f"Failed to send webhook. Error: {e}. Payload: {payload}")
        send_telegram_topic_message(curl_command, str(e))

def log_curl_request(url, headers, payload):
    # Construct cURL command equivalent
    curl_command = f"curl -X POST '{url}' "
    for key, value in headers.items():
        curl_command += f"-H '{key}: {value}' "
    curl_command += f"-d '{json.dumps(payload)}'"
    logging.info(f"Request: {curl_command}")
    return curl_command

def send_telegram_topic_message(curl_command, error):
    bot_token = 'xxxx'
    chat_id = 'yyyy'
    message_thread_id = 'zzzz'

    url = f"https://api.telegram.org/bot{bot_token}/sendMessage"
    payload = {
        'chat_id': chat_id,
        'text': f"{curl_command}\n{error}",
        'message_thread_id': message_thread_id
    }
    try:
        requests.post(url, data=payload, timeout=10)
    except Exception as e:
        logging.error(f"Failed to send Telegram message: {e}")

def main():
    logging.info("Script execution started.")
    current_time = time.time()
    last_run = get_last_run_time()

    # To prevent send a lot of existing dir in the first run
    is_first_run = last_run == 0 or (current_time - last_run) > 1200

    new_folders = find_new_folders(FOLDERS_TO_WATCH, last_run)

    if new_folders and not is_first_run:
        logging.info(f"Detected {len(new_folders)} new folder(s). Sending webhook.")
        send_webhook(new_folders)

    save_current_time(current_time)
    logging.info("Script execution finished.")

if __name__ == "__main__":
    main()