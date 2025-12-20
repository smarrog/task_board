#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import os
import sys
from pathlib import Path

# BOM'ы, которые будем вырезать "как есть" из начала файла (без перекодирования)
BOMS: list[bytes] = [
    b"\xEF\xBB\xBF",              # UTF-8 BOM
    b"\xFF\xFE",                  # UTF-16 LE BOM
    b"\xFE\xFF",                  # UTF-16 BE BOM
    b"\xFF\xFE\x00\x00",          # UTF-32 LE BOM
    b"\x00\x00\xFE\xFF",          # UTF-32 BE BOM
]

def find_git_root(start: Path) -> Path | None:
    """
    Рекурсивно ищет .git, поднимаясь вверх до корня файловой системы.
    (Поведение как в make_archieve.)
    """
    cur = start
    while True:
        if (cur / ".git").exists():
            return cur
        if cur.parent == cur:  # дошли до корня диска
            return None
        cur = cur.parent

def strip_bom_in_file(path: Path, *, max_probe_bytes: int = 4) -> bool:
    """
    Если файл начинается с известного BOM — удаляет его байты и перезаписывает файл.
    Возвращает True, если файл был изменён.
    """
    try:
        with path.open("rb") as f:
            head = f.read(max_probe_bytes)
            # Для корректного сравнения с BOM переменной длины нам нужен "чуть длиннее" буфер:
            # но max_probe_bytes=4 достаточно для всех BOMS выше.
            bom = next((b for b in BOMS if head.startswith(b)), None)
            if not bom:
                return False
            rest = f.read()

        # Перезаписываем файл уже без BOM
        with path.open("wb") as f:
            f.write(head[len(bom):])
            f.write(rest)

        return True
    except (PermissionError, IsADirectoryError, OSError):
        return False

def should_skip_dirname(dirname: str) -> bool:
    # минимальный скип: .git и типичные мусорные каталоги
    return dirname in {".git", ".venv", "venv", "__pycache__", "node_modules", ".idea", ".pytest_cache"}

def main() -> int:
    parser = argparse.ArgumentParser(description="Удаляет BOM (UTF-8/16/32) из всех файлов репозитория.")
    parser.add_argument(
        "--root",
        type=str,
        default=None,
        help="Корень для обхода. По умолчанию ищется git-root как в make_archieve (по .git).",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Ничего не менять, только показать какие файлы будут исправлены.",
    )
    args = parser.parse_args()

    if args.root:
        root = Path(args.root).expanduser().resolve()
    else:
        script_dir = Path(__file__).parent.resolve()
        root = find_git_root(script_dir)
        if root is None:
            print("❌ Git-репозиторий не найден (не вижу .git при подъёме вверх).", file=sys.stderr)
            return 1

    changed: list[Path] = []
    scanned = 0

    for cur_dir, dirnames, filenames in os.walk(root, topdown=True, followlinks=False):
        # фильтруем директории
        dirnames[:] = [d for d in dirnames if not should_skip_dirname(d)]

        for name in filenames:
            scanned += 1
            p = Path(cur_dir) / name

            # пропускаем очевидные бинарники по расширению (опционально, но обычно полезно)
            # BOM в них почти никогда не бывает, а трогать лишний раз не хочется.
            if p.suffix.lower() in {
                ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico",
                ".zip", ".7z", ".rar", ".gz", ".bz2", ".xz",
                ".pdf", ".mp3", ".mp4", ".mov", ".avi",
                ".woff", ".woff2", ".ttf", ".eot",
                ".exe", ".dll", ".so", ".dylib",
            }:
                continue

            try:
                with p.open("rb") as f:
                    head = f.read(4)
                if any(head.startswith(b) for b in BOMS):
                    if args.dry_run:
                        changed.append(p)
                    else:
                        if strip_bom_in_file(p):
                            changed.append(p)
            except (PermissionError, IsADirectoryError, OSError):
                continue

    print(f"📁 Root: {root}")
    print(f"🔎 Scanned files: {scanned}")
    if args.dry_run:
        print(f"🧪 DRY RUN — would fix: {len(changed)}")
    else:
        print(f"✅ Fixed: {len(changed)}")

    if changed:
        # печатаем относительные пути, чтобы было удобно
        for p in changed:
            try:
                rel = p.relative_to(root)
            except ValueError:
                rel = p
            print(f"- {rel}")

    return 0

if __name__ == "__main__":
    raise SystemExit(main())
