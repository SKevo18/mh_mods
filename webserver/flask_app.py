from flask import Flask, render_template, abort


MHMODS_FLASK_APP = Flask(__name__)


@MHMODS_FLASK_APP.get('/')
def index():
    return render_template("games.html")


@MHMODS_FLASK_APP.get('/<string:game_id>')
def game(game_id: str):
    return render_template("game.html", game=_get_game(game_id))