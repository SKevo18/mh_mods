import hashlib

from pathlib import Path
from tempfile import gettempdir
from aiofiles.tempfile import NamedTemporaryFile
from collections import OrderedDict

from quart import Quart, abort, render_template, request, send_file
from modder import GAMES, pack


WEBSERVER = Quart(__name__)
TEMPDIR = gettempdir()
CACHE_DIR = Path(TEMPDIR) / "idlemod_cache"
CACHE_DIR.mkdir(parents=True, exist_ok=True)
CACHE_LIMIT = 10
cache = OrderedDict()


def _get_game(game_id: str):
    game = GAMES.get(game_id)

    if game is None:
        abort(404, f"Game with ID `{game_id}` not found.")

    return game


def _get_cache_key(mods: list[str]):
    return hashlib.md5("".join(sorted(mods)).encode()).hexdigest()


def _add_to_cache(key: str, filepath: Path):
    if key in cache:
        cache.move_to_end(key)
    else:
        cache[key] = filepath
        if len(cache) > CACHE_LIMIT:
            _, old_file = cache.popitem(last=False)
            old_file.unlink()


@WEBSERVER.get("/")
@WEBSERVER.get("/game")
async def index():
    return await render_template("games.html", games=GAMES.values())


@WEBSERVER.get("/game/<string:game_id>")
async def game(game_id: str):
    return await render_template("game.html", game=_get_game(game_id))


@WEBSERVER.get("/game/<string:game_id>/packmod")
async def packmods(game_id: str):
    mods = request.args.getlist("mod")
    game = _get_game(game_id)
    mods_to_pack = [mod for mod in game.get_mods() if mod.id in mods]
    if not mods_to_pack:
        return await abort(400, "No valid mods selected to pack.")

    cache_key = _get_cache_key(mods)
    cached_file = cache.get(cache_key)

    if cached_file and cached_file.exists():
        return await send_file(
            cached_file,
            as_attachment=True,
            add_etags=True,
            mimetype="application/octet-stream",
            cache_timeout=300,
            attachment_filename=f"{game.out_filename}",
        )

    async with NamedTemporaryFile(dir=CACHE_DIR, suffix=game.original_datafile.suffix, prefix=game.original_datafile.stem + ".", delete=False) as f:
        path = Path(str(f.name)).resolve()
        _, out, err, code = await pack(game, path, mods_to_pack)
        if out is not None:
            out = out.decode("utf-8")
        if err is not None:
            err = err.decode("utf-8")
        
        if code != 0:
            return await abort(500, f"Failed to pack mods: {err}")

        _add_to_cache(cache_key, path)

        return await send_file(
            path,
            as_attachment=True,
            add_etags=True,
            mimetype="application/octet-stream",
            attachment_filename=f"{game.out_filename}",
        )
