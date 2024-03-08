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
    / ("mhmods.exe" if os_name == "nt" else "mhmods")
)


@dataclass
class Game:
    id: str
    name: str
    original_datafile: Path

    mods_folder: t.Optional[Path] = None
    mods: list[Mod] = field(default_factory=list)
    out_filename: t.Optional[str] = None

    def __post_init__(self):
        self.mods_folder = self.mods_folder or MODS_ROOT / self.id
        self.out_filename = self.out_filename or self.original_datafile.name
        self.get_mods()

    def get_mods(self) -> list[Mod]:
        mods_dir = self.mods_folder
        if not mods_dir or not mods_dir.exists():
            return []

        if not self.mods:
            self.mods = [
                Mod(id=mod_path.stem, game=self, path=mod_path)
                for mod_path in (mods_dir).iterdir()
                if mod_path.is_dir()
            ]

        return self.mods


@dataclass
class Mod:
    id: str
    game: Game
    path: Path
    readme: t.Optional[str] = None

    def __post_init__(self):
        self.get_readme()

    def get_readme(self) -> t.Optional[str]:
        if not self.readme:
            try:
                with open(self.path / "README.md", "r") as f:
                    self.readme = f.read()
            except FileNotFoundError:
                pass

        return self.readme


async def pack(
    game: Game, output_path: t.AnyStr | Path, mods: list[Mod]
) -> tuple[str, bytes, bytes, t.Optional[int]]:
    """mhmods packmod <game ID> <original data file> <output modded data file> <mod paths>... [flags]"""
    cmd = f""""{MHMODS_BINARY}" packmod {game.id.split('.', 2)[0]} "{game.original_datafile}" "{output_path}" {' '.join(f'"{mod.path / 'source'}"' for mod in mods)}"""

    process = await create_subprocess_shell(
        cmd,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )

    out, err = await process.communicate()
    code = process.returncode
    return cmd, out, err, code


GAMES = {
    "mhk_1": Game(
        id="mhk_1",
        name="Moorhuhn Kart: Extra (XXL)",
        out_filename="mhke.dat",
        original_datafile=DATAFILE_ROOT / "mhk_1.dat",
    ),
    "mhk_2.en": Game(
        id="mhk_2.en",
        name="Moorhuhn Kart 2 (English)",
        original_datafile=DATAFILE_ROOT / "mhk_2.en.dat",
        out_filename="mhk2-00.dat",
        mods_folder=MODS_ROOT / "mhk_2",
    ),
    "mhk_2.de": Game(
        id="mhk_2.de",
        name="Moorhuhn Kart 2 (German)",
        original_datafile=DATAFILE_ROOT / "mhk_2.de.dat",
        out_filename="mhk2-00.dat",
        mods_folder=MODS_ROOT / "mhk_2",
    ),
    "mhk_3": Game(
        id="mhk_3",
        name="Moorhuhn Kart 3",
        out_filename="data.sar",
        original_datafile=DATAFILE_ROOT / "mhk_3.sar",
    ),
    "mhk_4": Game(
        id="mhk_4",
        name="Moorhuhn Kart: Thunder",
        out_filename="data.sar",
        original_datafile=DATAFILE_ROOT / "mhk_4.sar",
    ),
}
