import os
import sys
import zipfile
import subprocess
from pathlib import Path

def read_text(p: Path) -> str:
    try:
        return p.read_text(encoding="utf-8")
    except FileNotFoundError:
        return ""
    except UnicodeDecodeError:
        # .gitignore обычно UTF-8; в случае сбоя читаем как latin-1
        return p.read_text(encoding="latin-1", errors="ignore")

def find_git_root(start: Path) -> Path | None:
    """
    Рекурсивно ищет .git, поднимаясь вверх до корня файловой системы.
    """
    cur = start
    while True:
        if (cur / ".git").exists():
            return cur
        if cur.parent == cur:  # дошли до корня диска
            return None
        cur = cur.parent

def get_global_excludes_path() -> Path | None:
    try:
        out = subprocess.check_output(
            ["git", "config", "--path", "core.excludesfile"],
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        return Path(out) if out else None
    except Exception:
        return None

def list_ignored_via_git(repo_root: Path, subdir: Path | None = None) -> set[str] | None:
    """
    Возвращает множество относительных путей, игнорируемых git'ом.

    Если subdir указан и это подкаталог repo_root, то возвращает пути
    уже относительные к subdir (чтобы совпадали с rel-путями в архиваторе).
    """
    try:
        cmd = ["git", "ls-files", "--ignored", "--exclude-standard", "--others", "-z"]
        out = subprocess.check_output(cmd, cwd=str(repo_root))
        items = out.split(b"\x00")
        ignored: set[str] = set()

        sub_prefix = ""
        if subdir is not None:
            sub_prefix = subdir.relative_to(repo_root).as_posix().rstrip("/")

        for b in items:
            if not b:
                continue
            rel = b.decode("utf-8", errors="ignore").replace("\\", "/")

            # Если работаем только с подкаталогом — отфильтруем лишнее
            if sub_prefix:
                if not rel.startswith(sub_prefix + "/"):
                    continue
                # обрезаем префикс subdir/, чтобы путь стал относительным к subdir
                rel = rel[len(sub_prefix) + 1 :]

            ignored.add(rel)

        return ignored
    except Exception:
        return None

def build_pathspec_from_gitignores(root: Path):
    """
    Пытается собрать pathspec из .gitignore, .git/info/exclude и глобального excludes.
    Вернёт (spec|None, used_pathspec: bool)
    """
    try:
        import pathspec
        from pathspec.patterns import GitWildMatchPattern
    except ImportError:
        return None, False

    patterns_text = []

    # .gitignore в корне + любые .gitignore ниже по дереву
    gitignores = [root / ".gitignore", *root.rglob(".gitignore")]

    # 👉 добавляем .gitignore на один уровень выше
    parent_gi = root.parent / ".gitignore"
    gitignores.append(parent_gi)

    for gi in gitignores:
        patterns_text.append(read_text(gi))

    # .git/info/exclude (если репо всё-таки в root)
    info_exclude = root / ".git" / "info" / "exclude"
    patterns_text.append(read_text(info_exclude))

    # глобальный excludesfile
    ge = get_global_excludes_path()
    if ge:
        patterns_text.append(read_text(ge))

    all_text = "\n".join(patterns_text)
    spec = pathspec.PathSpec.from_lines(GitWildMatchPattern, all_text.splitlines())
    return spec, True

def should_skip(rel_path: str, ignored_set: set[str] | None, spec) -> bool:
    # Приводим к виду с прямыми слэшами (git-стиль)
    p = rel_path.replace("\\", "/")
    if p in (".gitignore",):
        return True
    if p == ".git" or p.startswith(".git/"):
        return True
    if ignored_set is not None:
        # Быстрый путь: точные пути из git ls-files.
        # Для директорий git может вернуть их с завершающим слэшем; проверим префикс.
        if p in ignored_set:
            return True
        # Если какой-то предок в ignored_set с завершающим /, то файл внутри тоже игнорируется.
        parts = p.split("/")
        prefix = ""
        for part in parts[:-1]:
            prefix = f"{prefix}{part}/" if prefix else f"{part}/"
            if prefix in ignored_set:
                return True

    if spec is not None:
        # PathSpec ожидает путь относительный без ведущего './'
        return spec.match_file(p)
    return False

def main():
    script_dir = Path(__file__).parent.resolve()
    repo_root = find_git_root(script_dir)
    if repo_root is None:
        print("❌ Git-репозиторий не найден")
        sys.exit(1)

    root = repo_root
    project_name = root.name
    zip_name = f"{project_name}.zip"
    zip_path = root / zip_name  # 👉 архив всегда в корне репы

    ignored_set = None
    spec = None

    repo_root = find_git_root(root)
    if repo_root is not None:
        ignored_set = list_ignored_via_git(repo_root, subdir=root)

    if spec is None:
        spec, used_pathspec = build_pathspec_from_gitignores(root)
        if ignored_set is None and not used_pathspec:
            print("⚠️ Не удалось учесть правила .gitignore: установите один из вариантов:")
            print("   • Git в системе (чтобы работала команда 'git')")
            print("   • Библиотеку pathspec:  pip install pathspec")
            print("Продолжаю без исключений по .gitignore…")

    # Создание архива, не заходя в игнорируемые директории
    with zipfile.ZipFile(zip_path, "w", compression=zipfile.ZIP_DEFLATED) as zf:
        for cur_dir, dirnames, filenames in os.walk(root, topdown=True, followlinks=False):
            cur_dir_path = Path(cur_dir)
            rel_dir = cur_dir_path.relative_to(root)
            rel_dir_str = "" if str(rel_dir) == "." else str(rel_dir).replace("\\", "/")

            keep_dirs = []
            for d in dirnames:
                rel = f"{rel_dir_str}/{d}" if rel_dir_str else d
                if not should_skip(rel, ignored_set, spec):
                    keep_dirs.append(d)
            dirnames[:] = keep_dirs

            for fname in filenames:
                abs_path = cur_dir_path / fname
                rel = f"{rel_dir_str}/{fname}" if rel_dir_str else fname

                if should_skip(rel, ignored_set, spec):
                    continue

                # не кладём сам архив внутрь архива
                if abs_path == zip_path:
                    continue

                zf.write(abs_path, arcname=rel)

    print(f"✅ Архив создан: {zip_path}")

if __name__ == "__main__":
    main()
