import typing as t

import subprocess
import aiofiles
import shutil

from pathlib import Path

from fastapi import FastAPI, HTTPException, Request, Query, BackgroundTasks
from fastapi.responses import StreamingResponse
from starlette.background import BackgroundTask


MHMODS_FASTAPI_APP = FastAPI(root_path="/api", docs_url="/")

GITHUB_ROOT = "https://github.com/SKevo18/mhmods/tree/main/mods"



async def _async_download(path: Path, *args, **kwargs) -> StreamingResponse:
    async def _iter_file() -> t.AsyncGenerator[bytes, None]:
        async with aiofiles.open(path, "rb") as f:
            while chunk := await f.read((1024 * 1024) * 25):
                yield chunk

    return StreamingResponse(
        _iter_file(),
        headers={
            "Content-Length": str(path.stat().st_size),
            "Content-Disposition": f'attachment; filename="{path.name}"',
        },
        *args,
        **kwargs,
    )


@MHMODS_FASTAPI_APP.get("/info")
async def get_game_ids() -> t.Dict[str, str]:
    return {game_id: game.name for game_id, game in MHK_GAMES.items()}


@MHMODS_FASTAPI_APP.get("/info/{game_id}")
async def get_mod_ids_for_game(game_id: str) -> t.List[str]:
    game = _get_game(game_id)
    root_path = game.mod_root_path()

    if root_path.exists():
        return [
            path.stem
            for path in root_path.iterdir()
            if path.is_dir() and not path.stem.startswith(".")
        ]

    else:
        raise HTTPException(404, f"Game `{game_id}` has no mods.")


@MHMODS_FASTAPI_APP.get("/info/{game_id}/{mod_id}")
async def get_mod_data(game_id: str, mod_id: str) -> t.Dict[str, t.Any]:
    if mod_id.startswith("."):
        raise HTTPException(404, f"Mod file for `{game_id}/{mod_id}` does not exist!")

    game = _get_game(game_id)
    root = game.mod_root_path(mod_id)
    data = {
        "game_id": game_id,
        "mod_id": mod_id,
        "download_url": "/api"
        + MHMODS_FASTAPI_APP.url_path_for("download_mod", game_id=game_id, mod_id=mod_id),
        "source": GITHUB_ROOT + f"/{game_id}/{mod_id}",
        "meta": {"title": None, "description": None, "readme": None},
    }

    readme = root / "README.md"

    if readme.exists():
        readme_text = readme.read_text()

        title = next(
            line for line in readme_text.splitlines() if line.startswith("# ")
        ).removeprefix("# ")
        description = next(
            line for line in readme_text.splitlines() if line.startswith("> ")
        ).removeprefix("> ")

        data["meta"]["title"] = title
        data["meta"]["description"] = description
        data["meta"]["readme"] = readme_text

    return data


@MHMODS_FASTAPI_APP.get("/download/{game_id}/{mod_id}")
async def download_mod(game_id: str, mod_id: str):
    game = _get_game(game_id)
    mod_file = game.mod_root_path(mod_id) / game.data_filename

    if mod_id.startswith("."):
        raise HTTPException(404, f"Mod `{game_id}/{mod_id}` is WIP!")

    if not mod_file.exists():
        try:
            compile(game_id=game_id, mod_id=mod_id, force=True)

        except Exception:
            raise HTTPException(
                502,
                f"An error ocurred while downloading mod `{game_id}/{mod_id}`. Please, try again later.",
            )

    return await _async_download(mod_file)


@MHMODS_FASTAPI_APP.get("/merge/{game_id}")
async def merge_mods(game_id: str, mod_ids: t.List[str] = Query(alias="mod_id")):
    game = _get_game(game_id)

    merged_mod_id = f".merged.{'_'.join(mod_ids)}"
    merged_mod_path = game.mod_root_path(merged_mod_id)
    merged_mod_file = merged_mod_path / game.data_filename

    if len(mod_ids) < 2:
        raise HTTPException(400, "Please, specify at least 2 mod IDs to merge.")

    def mod_dir_cleanup():
        return shutil.rmtree(merged_mod_file.parent)

    if not merged_mod_file.exists():
        try:
            merge(game_id=game_id, merged_mod_id=merged_mod_id, mod_ids=mod_ids)

        except Exception:
            mod_dir_cleanup()
            raise HTTPException(500, "An error ocurred. Are mod IDs correct?")

    return await _async_download(
        merged_mod_file, background=BackgroundTask(mod_dir_cleanup)
    )


@MHMODS_FASTAPI_APP.post("/github_webhook", include_in_schema=False)
async def github_webhook(request: Request, background_tasks: BackgroundTasks):
    payload = await request.json()

    if payload.get("ref", None) != "refs/heads/main":
        raise HTTPException(200, "Not main branch.")

    def pull():
        try:
            subprocess.run(
                "git pull origin main",
                cwd=ROOT_PATH,
                shell=True,
                capture_output=True,
                check=True,
            )

        except subprocess.CalledProcessError as e:
            raise RuntimeError(e.returncode, ROOT_PATH, e.stderr) from e

        return compile_all(game_id=None, force=True)

    background_tasks.add_task(pull)
    return {"status": "success"}