import fastapi as fa

import git
import zipfile

from pathlib import Path
from tempfile import gettempdir


REPO_ROOT = Path(__file__).parent.parent
MODS_ROOT = REPO_ROOT / "mods"

APP = fa.FastAPI(title="MH Mods API", docs_url="/", redoc_url=None)
REPO = git.Repo(REPO_ROOT)


@APP.get("/mods")
async def get_mods() -> dict[str, dict[str, str]]:
    """
    Get a list of all mods for all games in the repository (key), and their last update commit hash (value).
    """

    mods = {}
    for game_folder in MODS_ROOT.iterdir():
        if game_folder.is_dir():
            mods.setdefault(game_folder.name, {})

            for mod_folder in game_folder.iterdir():
                if mod_folder.is_dir():
                    mods[game_folder.name][mod_folder.name] = REPO.git.log(
                        "-1", "--format=%h", "--", mod_folder
                    )

    return mods


@APP.get("/mods/{game}/{mod}")
async def get_mod(
    game: str, mod: str, t: fa.BackgroundTasks
) -> fa.responses.FileResponse:
    """
    Downloads the mod as a zip file.
    """

    if not all(part.isalnum() for part in (game, mod)):
        raise fa.HTTPException(
            status_code=400, detail="Game and mod names must be alphanumeric."
        )

    mod_folder = MODS_ROOT / game / mod
    if not mod_folder.exists():
        raise fa.HTTPException(status_code=404, detail="Mod not found")

    mod_zip = Path(gettempdir()) / mod_folder.with_suffix(".zip")
    if not mod_zip.exists():
        with zipfile.ZipFile(mod_zip, "w") as zf:
            zf.write(mod_folder, arcname=mod)
            for file in mod_folder.rglob("*"):
                zf.write(file, arcname=file.relative_to(mod_folder))

    t.add_task(mod_zip.unlink)
    return fa.responses.FileResponse(
        mod_zip, media_type="application/zip", filename=mod_zip.name
    )
