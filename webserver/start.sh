#!/usr/bin/bash
cd "$(dirname "$0")"

if [ ! -d ../.venv ]; then
  python3 -m venv ../.venv
  source ../.venv/bin/pip install -r requirements.txt
fi

export PYTHONPATH=$(pwd)
../.venv/bin/hypercorn app:APP --bind=unix:server.sock -m 007 -w 2
