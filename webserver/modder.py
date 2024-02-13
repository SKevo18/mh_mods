from __future__ import annotations
import aiofiles

from os import name as os_name
from pydantic import BaseModel
from pathlib import Path

from subprocess import run


MODS_ROOT = Path(__file__).parent.parent / "mods"
MHMODS_BINARY = (
    Path(__file__).parent.parent
    / "build"
    / ("mhmods" if os_name == "nt" else "mhmods.exe")
)


class Game(BaseModel):
    id: str
    name: str
    original_datafile: Path

    def mods(self) -> list[Mod]:
        return [
            Mod(id=mod_path.stem, game=self, path=mod_path)
            for mod_path in (MODS_ROOT / self.id).iterdir()
            if mod_path.is_dir()
        ]


class Mod(BaseModel):
    id: str
    game: Game
    path: Path

    @property
    async def get_readme(self) -> str:
        async with aiofiles.open(MODS_ROOT / self.id / "README.md", "r") as f:
            return await f.read()

    def mod(self):
        run([MHMODS_BINARY, "packmods", self.game.id, self.path])


GAMES = []
