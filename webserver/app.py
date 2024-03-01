from quart import Quart, abort, render_template, request, send_file
from modder import GAMES, pack

from tempfile import NamedTemporaryFile

WEBSERVER = Quart(__name__)


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

    with NamedTemporaryFile(suffix=game.original_datafile.suffix, delete=False) as f:
        out, err, code = await pack(game, f.name, mods_to_pack)
        if out is not None:
            out = out.decode("utf-8")
        if err is not None:
            err = err.decode("utf-8")
        
        if code != 0:
            return await abort(500, f"Failed to pack mods: {err}")

        return await send_file(
            f.name,
            as_attachment=True,
            attachment_filename=f"{game.out_filename}",
        )
