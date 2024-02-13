# MH Mods Webserver

This is the (very simple) FastAPI + Flask MHK mods webserver.

You can access it live at [https://mhmods.svit.ac](https://mhmods.svit.ac).

API documentation is available at [https://mhmods.svit.ac/api/docs](https://mhmods.svit.ac/api).

> **Note:** Mod folders starting with dot (`.`) are marked as WIP/temporary, and will be excluded from the webpage's mod list.

## Running it locally (for development)

1. Install Python 3.10+ virtualenv in `./.venv`: `python3 -m venv .venv`;
    - Python 3.12 is recommended (likely `python3.12 -m venv .venv` on Linux with APT)
2. Activate the virtualenv: `source .venv/bin/activate`;
3. Install dependencies: `pip install -r requirements.txt`;
4. `service mhkmods start` if running via systemd, or simply `./start.sh`

The webserver runs as a Unix `.sock` socket file. You can use `systemd` to manage it (see `mhmods.service` for example), or simply run `./start.sh` to start it via terminal.
