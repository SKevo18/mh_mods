from quart import Quart, render_template, request, send_file
from modder import GAMES

from tempfile import NamedTemporaryFile

WEBSERVER = Quart(__name__)


@WEBSERVER.get("/")
async def index():
    return await render_template("games.html", games=GAMES.values())


@WEBSERVER.get("/<string:game_id>")
async def game(game_id: str):
    return await render_template("game.html", game=GAMES[game_id])


@WEBSERVER.get("/<string:game_id>/packmod")
async def packmod(game_id: str):
    mods = request.args.getlist("mods")
    game = GAMES[game_id]
    mods_to_pack = [mod for mod in game.get_mods() if mod.id in mods]

    if not mods_to_pack:
        return "No mods selected/invalid IDs.", 400

    with NamedTemporaryFile(suffix=game.original_datafile.suffix) as f:
        await mods_to_pack[0].pack(f.name, mods_to_pack)
        return await send_file(
            f.name,
            as_attachment=True,
            attachment_filename=f"{game.original_datafile.name}",
        )
