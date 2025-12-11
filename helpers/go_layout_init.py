#!/usr/bin/env python3
import argparse
from pathlib import Path
from textwrap import dedent

DEFAULT_VERSION = "1.25"
DEFAULT_MODULE_PREFIX = "github.com/smarrog/"

# Дерево директорий по стандартному Go Project Layout
# https://github.com/golang-standards/project-layout
DIR_TREE = {
    "cmd": {},
    "internal": {
        "app": {},
        "pkg": {},
    },
    "pkg": {},
    "vendor": {},
    "api": {},
    "web": {},
    "configs": {},
    "init": {},
    "scripts": {},
    "build": {
        "ci": {},
        "package": {},
    },
    "deployments": {},
    "test": {
        "data": {},
        "testdata": {},
    },
    "docs": {},
    "tools": {},
    "examples": {},
    "third_party": {},
    "githooks": {},
    "assets": {},
    "website": {},
}

README_HINT = {
    "cmd": "Основные приложения проекта (каждая поддиректория = отдельный бинарник).",
    "internal": "Внутренний (неэкспортируемый) код приложения и библиотек.",
    "pkg": "Публичные библиотеки, допускающие использование в сторонних проектах.",
    "vendor": "Зависимости приложения (go mod vendor).",
    "api": "OpenAPI/Swagger спецификации, JSON Schema и определения протоколов.",
    "web": "Специфичные части веб-приложения: статика, шаблоны, SPA.",
    "configs": "Шаблоны конфигураций и значения по умолчанию.",
    "init": "Конфиги для systemd, upstart, supervisord и других init/PM систем.",
    "scripts": "Скрипты сборки, установки, анализа и прочих операций.",
    "build": "Сборка и CI: образы, пакеты, конфиги CI.",
    "deployments": "Манифесты деплоя: docker-compose, k8s/helm, terraform и пр.",
    "test": "Внешние тестовые приложения и тестовые данные.",
    "docs": "Проектная и пользовательская документация.",
    "tools": "Вспомогательные утилиты и тулзы для проекта.",
    "examples": "Примеры использования библиотек/приложения.",
    "third_party": "Сторонние утилиты, форки, помощники (например, Swagger UI).",
    "githooks": "Git hooks для репозитория.",
    "assets": "Ресурсы репозитория: изображения, логотипы и т.п.",
    "website": "Исходники сайта проекта (если не GitHub Pages).",
}


def create_tree(root: Path, tree: dict) -> None:
    for name, subtree in tree.items():
        dir_path = root / name
        dir_path.mkdir(parents=True, exist_ok=True)

        # Создаём README.md с кратким описанием (только для верхнего уровня)
        if name in README_HINT:
            readme = dir_path / "README.md"
            if not readme.exists():
                readme.write_text(
                    f"# /{name}\n\n{README_HINT[name]}\n",
                    encoding="utf-8",
                )

        # Рекурсивно создаём поддеревья
        if isinstance(subtree, dict):
            create_tree(dir_path, subtree)


def init_root_readme(root: Path, project_name: str) -> None:
    readme = root / "README.md"
    if readme.exists():
        return

    content = dedent(
        f"""\
        # {project_name}

        Скелет проекта по стандартному макету Go Project Layout:
        https://github.com/golang-standards/project-layout

        Структура директорий:
        - /cmd           — основные приложения
        - /internal      — внутренний код
        - /pkg           — публичные библиотеки
        - /api           — спецификации API
        - /web           — веб-часть (статические файлы, шаблоны, SPA)
        - /configs       — шаблоны конфигов
        - /init          — конфиги init/PM систем
        - /scripts       — скрипты сборки/анализа
        - /build         — CI и упаковка
        - /deployments   — манифесты деплоя (k8s, terraform, docker-compose, ...)
        - /test          — внешние тестовые данные и приложения
        - /docs          — документация
        - /tools         — вспомогательные утилиты
        - /examples      — примеры
        - /third_party   — сторонние тулзы и форки
        - /githooks      — git hooks
        - /assets        — ресурсы (картинки, логотипы)
        - /website       — сайт проекта

        Дальше — добавляй код и конфиги по необходимости.
        """
    )
    readme.write_text(content, encoding="utf-8")


def init_go_mod(root: Path, module_path: str | None, go_version: str | None) -> None:
    """
    Создаёт минимальный go.mod, если его нет.
    module_path по умолчанию: «example.com/{имя_директории}».
    """
    gomod = root / "go.mod"
    if gomod.exists():
        return

    if module_path is None:
        module_path = f"{DEFAULT_MODULE_PREFIX}{root.name}"

    if go_version is None:
        go_version = DEFAULT_VERSION

    gomod.write_text(
        f"module {module_path}\n\ngo {go_version}\n",
        encoding="utf-8",
    )


def ensure_target(path: Path, force: bool) -> None:
    if path.exists():
        # Если директория существует и не пустая — не портим её без --force
        if any(path.iterdir()) and not force:
            raise SystemExit(
                f"Папка '{path}' уже существует и не пуста. "
                f"Запусти с флагом --force, если хочешь использовать её."
            )
    else:
        path.mkdir(parents=True, exist_ok=True)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Создание скелета Go проекта по стандартному project-layout."
    )
    parser.add_argument(
        "path",
        help="Путь к корню проекта (будет создан, если не существует)",
    )
    parser.add_argument(
        "--version",
        help="Версия языка Go для go.mod. "
             "По умолчанию: DEFAULT_VERSION",
        default=None,
    )
    parser.add_argument(
        "--module",
        help="Путь модуля Go для go.mod (например, github.com/user/project). "
             "По умолчанию: DEFAULT_MODULE",
        default=None,
    )
    parser.add_argument(
        "-f",
        "--force",
        action="store_true",
        help="Разрешить использование уже существующей непустой директории.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    root = Path(args.path).resolve()

    ensure_target(root, force=args.force)

    create_tree(root, DIR_TREE)
    init_root_readme(root, project_name=root.name)
    init_go_mod(root, args.module, args.version)

    print(f"Готово! Скелет проекта создан в: {root}")


if __name__ == "__main__":
    main()
