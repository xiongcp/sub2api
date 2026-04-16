#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
按邮箱通过 Admin API 重置用户密码。

用法:
  tools/reset_user_password.sh \
    --base-url https://example.com \
    --email user@example.com \
    --password 'NewStrongPass123!' \
    --token '<admin-jwt>'

或:
  tools/reset_user_password.sh \
    --base-url https://example.com \
    --email user@example.com \
    --password 'NewStrongPass123!' \
    --api-key 'admin-xxxxxxxx'

参数:
  --base-url   服务地址，例如 https://cc.taylor-link.xyz
  --email      要重置密码的用户邮箱
  --password   新密码
  --token      管理员 JWT
  --api-key    管理员 API Key

说明:
  - 脚本会先调用 GET /api/v1/admin/users?search=... 查找精确邮箱，再调用 PUT /api/v1/admin/users/:id 更新 password。
  - 需要系统中已存在管理员认证能力；脚本不会直接访问数据库。
EOF
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "缺少依赖命令: $1" >&2
    exit 1
  }
}

BASE_URL=""
EMAIL=""
NEW_PASSWORD=""
ADMIN_TOKEN=""
ADMIN_API_KEY=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --base-url)
      BASE_URL="${2:-}"
      shift 2
      ;;
    --email)
      EMAIL="${2:-}"
      shift 2
      ;;
    --password)
      NEW_PASSWORD="${2:-}"
      shift 2
      ;;
    --token)
      ADMIN_TOKEN="${2:-}"
      shift 2
      ;;
    --api-key)
      ADMIN_API_KEY="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "未知参数: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$BASE_URL" || -z "$EMAIL" || -z "$NEW_PASSWORD" ]]; then
  echo "缺少必填参数。" >&2
  usage
  exit 1
fi

if [[ -n "$ADMIN_TOKEN" && -n "$ADMIN_API_KEY" ]]; then
  echo "--token 和 --api-key 只能二选一。" >&2
  exit 1
fi

if [[ -z "$ADMIN_TOKEN" && -z "$ADMIN_API_KEY" ]]; then
  echo "必须提供 --token 或 --api-key 其中之一。" >&2
  exit 1
fi

require_cmd curl
require_cmd python3

BASE_URL="${BASE_URL%/}"

AUTH_HEADER_NAME="Authorization"
AUTH_HEADER_VALUE="Bearer ${ADMIN_TOKEN}"
if [[ -n "$ADMIN_API_KEY" ]]; then
  AUTH_HEADER_NAME="x-api-key"
  AUTH_HEADER_VALUE="${ADMIN_API_KEY}"
fi

EMAIL_ENCODED="$(python3 - <<'PY' "$EMAIL"
import sys, urllib.parse
print(urllib.parse.quote(sys.argv[1], safe=''))
PY
)"

lookup_response_file="$(mktemp)"
update_response_file="$(mktemp)"
trap 'rm -f "$lookup_response_file" "$update_response_file"' EXIT

lookup_status="$(
  curl -sS \
    -o "$lookup_response_file" \
    -w '%{http_code}' \
    -H "${AUTH_HEADER_NAME}: ${AUTH_HEADER_VALUE}" \
    "${BASE_URL}/api/v1/admin/users?page=1&page_size=100&search=${EMAIL_ENCODED}"
)"

if [[ "$lookup_status" != "200" ]]; then
  echo "查询用户失败，HTTP ${lookup_status}" >&2
  cat "$lookup_response_file" >&2
  exit 1
fi

user_id="$(python3 - <<'PY' "$lookup_response_file" "$EMAIL"
import json, sys

path = sys.argv[1]
target = sys.argv[2].strip().lower()

with open(path, 'r', encoding='utf-8') as fh:
    payload = json.load(fh)

items = (((payload or {}).get('data') or {}).get('items') or [])
matched = [item for item in items if str((item or {}).get('email', '')).strip().lower() == target]

if len(matched) != 1:
    sys.exit(2)

print(matched[0]['id'])
PY
)" || {
  echo "没有找到唯一匹配的用户邮箱: ${EMAIL}" >&2
  echo "请先用后台确认该邮箱存在且唯一。" >&2
  cat "$lookup_response_file" >&2
  exit 1
}

update_status="$(
  python3 - <<'PY' "$BASE_URL" "$user_id" "$NEW_PASSWORD" "$AUTH_HEADER_NAME" "$AUTH_HEADER_VALUE" "$update_response_file"
import json, subprocess, sys

base_url, user_id, password, header_name, header_value, output_file = sys.argv[1:]
payload = json.dumps({"password": password})
cmd = [
    "curl", "-sS",
    "-o", output_file,
    "-w", "%{http_code}",
    "-X", "PUT",
    "-H", f"{header_name}: {header_value}",
    "-H", "Content-Type: application/json",
    "--data", payload,
    f"{base_url}/api/v1/admin/users/{user_id}",
]
result = subprocess.run(cmd, check=True, capture_output=True, text=True)
print(result.stdout, end="")
PY
)"

if [[ "$update_status" != "200" ]]; then
  echo "重置密码失败，HTTP ${update_status}" >&2
  cat "$update_response_file" >&2
  exit 1
fi

python3 - <<'PY' "$update_response_file" "$EMAIL" "$user_id"
import json, sys

path, email, user_id = sys.argv[1:]
with open(path, 'r', encoding='utf-8') as fh:
    payload = json.load(fh)

if payload.get("code") != 0:
    print("接口返回失败:", json.dumps(payload, ensure_ascii=False), file=sys.stderr)
    sys.exit(1)

print(f"已重置用户密码: email={email} user_id={user_id}")
print("管理员改密已通过 token_version 失效旧会话；用户需要使用新密码重新登录。")
PY
