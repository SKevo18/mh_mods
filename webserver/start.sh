#!/usr/bin/bash
cd "$(dirname "$0")"

if [ ! -d ../.venv ]; then
  python3 -m venv ../.venv
  source ../.venv/bin/pip install -r requirements.txt
fi

export PYTHONPATH=$(pwd)

if [[ $* == *--local* ]]; then
  ../.venv/bin/hypercorn app:APP -m 007 -w 2 --bind=localhost:8080 --reload
else
  ../.venv/bin/hypercorn app:APP -m 007 -w 2 --bind=unix:server.sock
fi
