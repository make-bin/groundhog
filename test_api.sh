#!/bin/bash
# API Test Script - logs all requests and responses to api_test.log

BASE_URL="http://localhost:8080/api/v1"
LOG_FILE="api_test.log"
PROVIDER="openai_compat"
MODEL="deepseek-ai/DeepSeek-V3.1-Terminus"

> "$LOG_FILE"

log() {
  echo "$@" | tee -a "$LOG_FILE"
}

call() {
  local label="$1"
  local method="$2"
  local url="$3"
  shift 3
  local extra_args=("$@")

  log ""
  log "========================================"
  log "TEST: $label"
  log "REQUEST: $method $url"

  # Print body if -d is present
  for i in "${!extra_args[@]}"; do
    if [[ "${extra_args[$i]}" == "-d" ]]; then
      log "BODY: ${extra_args[$((i+1))]}"
    fi
  done

  log "RESPONSE:"
  local resp
  resp=$(curl -s -X "$method" "$url" \
    -H "Content-Type: application/json" \
    "${extra_args[@]}")
  echo "$resp" | tee -a "$LOG_FILE"
  log ""
  echo "$resp"
}

# ─── Health ───────────────────────────────────────────────────────────────────

call "Health Check" GET "$BASE_URL/health"

# ─── Sessions ─────────────────────────────────────────────────────────────────

log ""
log "######## SESSIONS ########"

SESSION_RESP=$(call "Create Session" POST "$BASE_URL/sessions" \
  -d "{\"user_id\":\"user-001\",\"agent_id\":\"agent-001\",\"provider\":\"$PROVIDER\",\"model_name\":\"$MODEL\"}")

SESSION_ID=$(echo "$SESSION_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
log ">> Extracted SESSION_ID: $SESSION_ID"

call "List Sessions" GET "$BASE_URL/sessions"

call "List Sessions (filter by user)" GET "$BASE_URL/sessions?user_id=user-001"

if [ -n "$SESSION_ID" ]; then
  call "Get Session" GET "$BASE_URL/sessions/$SESSION_ID"

  call "Send Message" POST "$BASE_URL/sessions/$SESSION_ID/messages" \
    -d '{"user_input":"你好，请用一句话介绍你自己"}'
else
  log "SKIP: Get/SendMessage - no session ID"
fi

# ─── Channels ─────────────────────────────────────────────────────────────────

log ""
log "######## CHANNELS ########"

CHANNEL_RESP=$(call "Create Channel" POST "$BASE_URL/channels" \
  -d '{"channel_type":"discord","plugin_id":"plugin-001"}')

CHANNEL_ID=$(echo "$CHANNEL_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null || \
             echo "$CHANNEL_RESP" | grep -o '"id":"ch-[^"]*"' | head -1 | cut -d'"' -f4)
log ">> Extracted CHANNEL_ID: $CHANNEL_ID"

call "List Channels" GET "$BASE_URL/channels"

if [ -n "$CHANNEL_ID" ]; then
  call "Get Channel Status" GET "$BASE_URL/channels/$CHANNEL_ID/status"
  call "Delete Channel" DELETE "$BASE_URL/channels/$CHANNEL_ID"
else
  log "SKIP: Channel status/delete - no channel ID"
fi

# ─── Memories ─────────────────────────────────────────────────────────────────

log ""
log "######## MEMORIES ########"

# Note: user_id comes from X-User-ID header (auth context)
# If auth middleware is not wired, these will return "userID is required"
MEM_RESP=$(call "Create Memory" POST "$BASE_URL/memories" \
  -H "X-User-ID: user-001" \
  -d '{"content":"用户喜欢简洁的回答风格"}')

MEM_ID=$(echo "$MEM_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('id',''))" 2>/dev/null || \
         echo "$MEM_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
log ">> Extracted MEMORY_ID: $MEM_ID"

call "List Memories" GET "$BASE_URL/memories" \
  -H "X-User-ID: user-001"

if [ -n "$MEM_ID" ]; then
  call "Get Memory" GET "$BASE_URL/memories/$MEM_ID" \
    -H "X-User-ID: user-001"

  call "Update Memory" PUT "$BASE_URL/memories/$MEM_ID" \
    -H "X-User-ID: user-001" \
    -d '{"content":"用户喜欢简洁且带例子的回答风格"}'

  call "Search Memory" POST "$BASE_URL/memories/search" \
    -H "X-User-ID: user-001" \
    -d '{"query":"回答风格","limit":5}'

  call "Delete Memory" DELETE "$BASE_URL/memories/$MEM_ID" \
    -H "X-User-ID: user-001"
else
  log "SKIP: Memory get/update/search/delete - no memory ID (likely auth middleware not wired)"

  call "Search Memory (no auth)" POST "$BASE_URL/memories/search" \
    -d '{"query":"回答风格","limit":5}'
fi

# ─── Security / Audit ─────────────────────────────────────────────────────────

log ""
log "######## SECURITY ########"

call "Audit Logs" GET "$BASE_URL/security/audit"

# ─── Cleanup: Delete Session ──────────────────────────────────────────────────

log ""
log "######## CLEANUP ########"

if [ -n "$SESSION_ID" ]; then
  call "Delete Session" DELETE "$BASE_URL/sessions/$SESSION_ID"
fi

log ""
log "========================================"
log "Done. Full log saved to: $LOG_FILE"
