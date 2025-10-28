#!/usr/bin/env zsh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
JWT_SECRET="a-very-super-secret-key-for-testing"
DB_FILE="scrypts_test1.db"
LOG_FILE="/tmp/scrypts_test.log"
SERVER_PID=""
BASE_URL="http://localhost:8080"

# use a unique username per run to avoid "User already exists"
USERNAME="testuser_$(date +%s%N)"
# --- Helper Functions ---

# Function to print colored headers
print_header() {
  echo "\n\033[1;34m--- $1 ---\033[0m"
}

# Cleanup function to be called on script exit
cleanup() {
  print_header "CLEANUP"
  if [[ -n "$SERVER_PID" ]]; then
    echo "Stopping server (PID: $SERVER_PID)..."
    # Kill the process group to ensure any child processes are also terminated
    kill -9 -- "-$SERVER_PID" 2>/dev/null || echo "Server already stopped."
  fi
  echo "Removing test database: $DB_FILE"
  rm -f "$DB_FILE" "$DB_FILE-shm" "$DB_FILE-wal"
  echo "Test script finished."
}

# Trap EXIT signal to run cleanup function
trap cleanup EXIT

# --- Main Script ---

# 1. Initial Setup
print_header "SETUP"
echo "Using database: $DB_FILE"
echo "Using log file: $LOG_FILE"

# Check for jq for cleaner JSON parsing
if ! command -v jq &> /dev/null; then
  echo "\033[1;33mWarning: 'jq' is not installed. JSON output will be less readable. Using 'sed' for parsing.\033[0m"
fi

# Clean up previous runs
rm -f "$DB_FILE" "$DB_FILE-shm" "$DB_FILE-wal" "$LOG_FILE"

# 2. Start Server
print_header "STARTING SERVER"
export JWT_SECRET
export SCRYPTS_DB_PATH="$DB_FILE"

# If any process is listening on :8080, try to kill it to avoid "address already in use" errors
if command -v lsof >/dev/null 2>&1; then
  OLD_PID=$(lsof -ti:8080 || true)
  if [[ -n "$OLD_PID" ]]; then
    echo "Found existing process on :8080 (PID: $OLD_PID) — killing it"
    kill -9 $OLD_PID 2>/dev/null || true
    sleep 0.2
  fi
fi

# Start the server in background and capture logs
go run ./cmd/scrypts > "$LOG_FILE" 2>&1 &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# Wait for the server to be ready
echo "Waiting for server to respond..."
for i in {1..10}; do
  if curl -s -o /dev/null "$BASE_URL"; then
    echo "Server is up!"
    break
  fi
  if (( i == 10 )); then
    echo "\033[1;31mServer failed to start. Check logs:\033[0m"
    cat "$LOG_FILE"
    exit 1
  fi
  sleep 0.5
done

# 3. API Test Flow
print_header "REGISTER USER"
# use unique username and save raw response for debugging
curl -i -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"ValidPassw0rd!\"}" \
  -o /tmp/register_resp.txt
cat /tmp/register_resp.txt

print_header "LOGIN USER"
# save raw login response for easier debugging
curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"ValidPassw0rd!\"}" \
  -o /tmp/login.json

echo "Raw login response:"
cat /tmp/login.json

if command -v jq &> /dev/null; then
  TOKEN=$(jq -r '.token' /tmp/login.json 2>/dev/null || true)
else
  TOKEN=$(sed -n 's/.*"token":"\([^\"]*\)".*/\1/p' /tmp/login.json || true)
fi

if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "\033[1;31mFailed to get JWT token. Login response:\033[0m"
  echo "$LOGIN_RESP"
  exit 1
fi
echo "Token acquired."

print_header "CREATE NOTE"
CREATE_RESP=$(curl -s -X POST "$BASE_URL/notes" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"hello persistent world"}')

if command -v jq &> /dev/null; then
  NOTE_ID=$(echo "$CREATE_RESP" | jq -r '.id')
else
  NOTE_ID=$(echo "$CREATE_RESP" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
fi

if [[ -z "$NOTE_ID" || "$NOTE_ID" == "null" ]]; then
  echo "\033[1;31mFailed to create note. Create response:\033[0m"
  echo "$CREATE_RESP"
  exit 1
fi
echo "Note created with ID: $NOTE_ID"

print_header "GET NOTES (after create)"
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/notes" | (command -v jq &> /dev/null && jq . || cat)

print_header "UPDATE NOTE"
curl -i -X PUT "$BASE_URL/notes" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"id\":\"$NOTE_ID\",\"content\":\"updated content\"}" \
  -o /tmp/update.json
cat /tmp/update.json

print_header "GET NOTES (after update)"
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/notes" | (command -v jq &> /dev/null && jq . || cat)

print_header "DELETE NOTE"
curl -i -X DELETE "$BASE_URL/notes" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"id\":\"$NOTE_ID\"}"

print_header "GET NOTES (after delete)"
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/notes" | (command -v jq &> /dev/null && jq . || cat)

# 4. Final Log Output
print_header "SERVER LOGS"
cat "$LOG_FILE"