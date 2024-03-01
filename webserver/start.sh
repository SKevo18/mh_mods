#!/usr/bin/bash
cd "$(dirname "$0")"

handle_sigint() {
  echo "Killing..."
  kill $pid1 $pid2
}
trap handle_sigint SIGINT

if [ ! -d .venv ]; then
  python3 -m venv .venv
  source .venv/bin/activate
  pip install -r requirements.txt
fi

export PYTHONPATH=$(pwd)
../.venv/bin/hypercorn app:WEBSERVER --bind=unix:./mhmods_flask.sock -k uvloop -m 007 -w 4 --timeout=300 &
pid1=$!


wait $pid1 $pid2
