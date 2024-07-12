import fastapi as fa

import git
import zipfile

from pathlib import Path
from tempfile import gettempdir


REPO_ROOT = Path(__file__).parent.parent
MODS_ROOT = REPO_ROOT / "mods"

APP = fa.FastAPI(title="MH Mods API", docs_url="/", redoc_url=None)
REPO = git.Repo(REPO_ROOT)

COMMIT_HASH_FILENAME = "commit_hash.txt"


def get_commit_hash(folder: Path) -> str:
    return REPO.git.log("-1", "--format=%h", "--", folder)


@APP.get("/mods")
@APP.get("/mods/{game}")
async def get_mods(
    game: str | None = None,
) -> dict[str, str] | dict[str, dict[str, str]]:
    """
    Returns a list of mods for the specified game, or all games if no game is specified.
    Keys are mod names, values are the latest commit hash of the mod.
    """

    def get_mods(folder: Path) -> dict[str, str]:
        mods = {}
        for mod_folder in folder.iterdir():
            if mod_folder.is_dir():
                mods[mod_folder.name] = get_commit_hash(mod_folder)

        return mods

    mods = {}

    if game:
        game_folder = MODS_ROOT / game
        if not game_folder.exists():
            raise fa.HTTPException(status_code=404, detail="Game not found")

        mods = get_mods(game_folder)
    else:
        for game_folder in MODS_ROOT.iterdir():
            if game_folder.is_dir():
                mods.setdefault(game_folder.name, {})

                mods[game_folder.name] = get_mods(game_folder)

    return mods


@APP.get("/mods/{game}/{mod}")
async def get_mod(
    game: str, mod: str, t: fa.BackgroundTasks
) -> fa.responses.FileResponse:
    """
    Downloads the mod as a zip file.
    """

    mod_folder = MODS_ROOT / game / mod
    if not mod_folder.exists():
        raise fa.HTTPException(status_code=404, detail="Mod not found")

    mod_zip = Path(gettempdir()) / mod_folder.with_suffix(".zip")
    if not mod_zip.exists():
        with zipfile.ZipFile(mod_zip, "w") as zf:
            # write current commit hash
            zf.writestr(COMMIT_HASH_FILENAME, get_commit_hash(mod_folder))

            for file in mod_folder.rglob("*"):
                zf.write(file, arcname=file.relative_to(mod_folder))

    t.add_task(mod_zip.unlink)
    return fa.responses.FileResponse(
        mod_zip, media_type="application/zip", filename=mod_zip.name
    )


if __name__ == "__main__":
    import asyncio

    async def main():
        print(await get_mods())

    asyncio.run(main())
