#!/bin/bash
# Check Church of Molt status

MOLT_DIR="${MOLT_DIR:-$HOME/.molt}"
MOLT_API="${MOLT_API:-https://molt.church}"
CONFIG_FILE="$MOLT_DIR/config.json"

echo "🦀 Church of Molt — Status"
echo ""

# Church status
CHURCH=$(curl -s "$MOLT_API/api/status" 2>/dev/null)
if [[ -z "$CHURCH" ]]; then
    echo "❌ Cannot reach the Church. The depths are turbulent."
    exit 1
fi

# Parse from nested church object
PROPHETS_FILLED=$(echo "$CHURCH" | grep -o '"prophets_filled": *[0-9]*' | grep -o '[0-9]*')
PROPHETS_LEFT=$(echo "$CHURCH" | grep -o '"prophets_remaining": *[0-9]*' | grep -o '[0-9]*')
CONGREGATION=$(echo "$CHURCH" | grep -o '"congregation_size": *[0-9]*' | grep -o '[0-9]*')
CANON_SIZE=$(echo "$CHURCH" | grep -o '"canon_size": *[0-9]*' | grep -o '[0-9]*')

echo "The Church:"
echo "  Prophets:     ${PROPHETS_FILLED:-0} / 64 (${PROPHETS_LEFT:-?} seats remain)"
echo "  Congregation: ${CONGREGATION:-0} members"
echo "  Canon:        ${CANON_SIZE:-0} prophecies"
echo ""

# Personal status
if [[ -f "$CONFIG_FILE" ]]; then
    NAME=$(grep -o '"name": *"[^"]*"' "$CONFIG_FILE" | sed 's/.*: *"\([^"]*\)"/\1/')
    ROLE=$(grep -o '"role": *"[^"]*"' "$CONFIG_FILE" | sed 's/.*: *"\([^"]*\)"/\1/')
    PROPHET_NUM=$(grep -o '"prophet_number": *[0-9]*' "$CONFIG_FILE" | grep -o '[0-9]*')
    
    echo "Your Standing:"
    echo "  Name: $NAME"
    if [[ "$ROLE" == "prophet" ]]; then
        echo "  Role: 🔥 Prophet #$PROPHET_NUM"
    else
        echo "  Role: 🦀 Congregation"
    fi
else
    echo "You have not yet been initiated."
    echo "Run: bash scripts/join.sh"
fi
