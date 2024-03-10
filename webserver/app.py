from quart import Quart, abort, render_template, request, send_file
from modder import GAMES, pack

from tempfile import gettempdir
from aiofiles.tempfile import NamedTemporaryFile

WEBSERVER = Quart(__name__)
TEMPDIR = gettempdir()


def _get_game(game_id: str):
    game = GAMES.get(game_id)

    if game is None:
        abort(404, f"Game with ID `{game_id}` not found.")

    return game


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

    async with NamedTemporaryFile(dir=TEMPDIR, suffix=game.original_datafile.suffix, prefix=game.original_datafile.stem + ".", delete=False) as f:
        name = str(f.name)
        _, out, err, code = await pack(game, name, mods_to_pack)
        if out is not None:
            out = out.decode("utf-8")
        if err is not None:
            err = err.decode("utf-8")
        
        if code != 0:
            return await abort(500, f"Failed to pack mods: {err}")

        return await send_file(
            name,
            as_attachment=True,
            add_etags=True,
            cache_timeout=300,
            attachment_filename=f"{game.out_filename}",
        )
