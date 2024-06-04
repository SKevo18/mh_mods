# MH Mods API

This is a simple Python FastAPI webserver to query and download mods from the GitHub repository, for the GUI application to use.

## Running

1. Install Python 3.10+ virtualenv in `./.venv`: `python3 -m venv .venv`;
    - Python 3.12 is recommended;
2. Activate the virtualenv: `source .venv/bin/activate`;
3. Install dependencies: `pip install -r requirements.txt`;
4. `service mhmods start` if running via systemd, or simply `./start.sh`;
    - The webserver operates as a socket file. You can use `systemd` to manage it (see `mhmods.service` for example), or simply run `./start.sh` to start the Flask app.
5. Place your original data files in `../data/` (see `../data/README.md` for more info);
