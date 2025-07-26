#!/bin/sh

SOURCE_CODE=$(cat <<'EOF'
print("hello")
EOF
)

ESCAPED_SOURCE=$(printf '%s' "$SOURCE_CODE" \
  | sed -e 's/\\/\\\\/g' \
        -e 's/"/\\"/g' \
        -e ':a;N;$!ba;s/\n/\\n/g')

JSON_PAYLOAD=$(printf '{"lang":"python","job":"ast","path":"D:\dev\kitsune\src"}' "$ESCAPED_SOURCE")

# Send it
printf '%s\n' "$JSON_PAYLOAD" | socat - UNIX-CONNECT:/tmp/kitsuned.sock
