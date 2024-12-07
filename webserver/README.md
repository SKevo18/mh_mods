# MH Mods Webserver

This is the (very simple) Flask MHK mods webserver. It is used for packing existing mods present in this repository into a single data file - useful for people that do not want to use the CLI and want to simply play the mods.

You can access it live at [https://mhmods.svit.ac](https://mhmods.svit.ac).

> **Note:** Mod folders starting with dot (`.`) are marked as WIP/temporary, and will be excluded from the webpage's mod list.

## Running it locally (for development)

1. Install Python 3.10+ virtualenv in `./.venv`: `python3 -m venv .venv`;
    - Python 3.12 is recommended;
2. Activate the virtualenv: `source .venv/bin/activate`;
3. Install dependencies: `pip install -r requirements.txt`;
4. `service mhmods start` if running via systemd, or simply `./start.sh`;
    - The webserver operates as a socket file. You can use `systemd` to manage it (see `mhmods.service` for example), or simply run `./start.sh` to start the Flask app.
5. Place your original data files in `../data/` (see `../data/README.md` for more info);
