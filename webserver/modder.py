from __future__ import annotations

import typing as t
import asyncio
import uvloop

from asyncio.subprocess import create_subprocess_shell
from dataclasses import dataclass, field
from os import name as os_name
from pathlib import Path

asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())

MODS_ROOT = Path(__file__).parent.parent / "mods"
DATAFILE_ROOT = Path(__file__).parent.parent / "data"
MHMODS_BINARY = (
    Path(__file__).parent.parent
    / "build"
    / ("mhmods" if os_name == "nt" else "mhmods.exe")
)


@dataclass
class Game:
    id: str
    name: str
    original_datafile: Path

    _mods: list[Mod] = field(default_factory=list)

    def get_mods(self) -> list[Mod]:
        if not self._mods:
            self._mods = [
                Mod(id=mod_path.stem, game=self, path=mod_path)
                for mod_path in (MODS_ROOT / self.id).iterdir()
                if mod_path.is_dir()
            ]

        return self._mods


@dataclass
class Mod:
    id: str
    game: Game
    path: Path
    _readme: t.Optional[str] = None

    def readme(self) -> t.Optional[str]:
        if not self._readme:
            try:
                with open(MODS_ROOT / self.id / "README.md", "r") as f:
                    self._readme = f.read()
            except FileNotFoundError:
                pass

        return self._readme

    async def pack(
        self, output_path: t.AnyStr | Path, mods: list[Mod]
    ) -> tuple[str, str]:
        process = await create_subprocess_shell(
            # mhmods packmod <game ID> <original data file> <output modded data file> <mod paths>... [flags]
            f""""{MHMODS_BINARY}" packmod {self.game.id} "{self.game.original_datafile}" "{output_path}" {' '.join(f'"{mod.path}"' for mod in mods)}"""
        )

        out, err = await process.communicate()
        return out.decode("utf-8"), err.decode("utf-8")


GAMES = {
    "mhk_1": Game(
        id="mhk_1",
        name="Moorhuhn Kart: Extra (XXL)",
        original_datafile=DATAFILE_ROOT / "mhk_1.dat",
    ),
    "mhk_2": Game(
        id="mhk_2",
        name="Moorhuhn Kart 2",
        original_datafile=DATAFILE_ROOT / "mhk_2.dat",
    ),
    "mhk_3": Game(
        id="mhk_3",
        name="Moorhuhn Kart 3",
        original_datafile=DATAFILE_ROOT / "mhk_3.sar",
    ),
    "mhk_4": Game(
        id="mhk_4",
        name="Moorhuhn Kart: Thunder",
        original_datafile=DATAFILE_ROOT / "mhk_4.sar",
    ),
}
