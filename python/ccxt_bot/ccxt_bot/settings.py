from typing import Optional

import toml


class VariableEmptyError(Exception):
    def __init__(self, name: str):
        self.name = name


class SettingsToml:
    def __init__(self, env: dict[str, any]):
        # if env["telegram"]["token"] == "":
        #     raise VariableEmptyError("telegram token is empty")

        self._env = env

    @property
    def telegram_token(self) -> str:
        return self._env["telegram"]["token"]


def init() -> Optional[SettingsToml]:
    try:
        with open("env.toml", "r") as f:
            env = toml.load(f)
            return SettingsToml(env)
    except FileNotFoundError:
        return None
