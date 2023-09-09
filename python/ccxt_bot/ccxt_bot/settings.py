from typing import Optional

import toml

from ccxt_bot.exchanges import ExchangesMap


class VariableEmptyError(Exception):
    def __init__(self, name: str):
        self.name = name


class SettingsToml:
    def __init__(self, env: dict[str, any]):
        self._env = env

    @property
    def telegram_token(self) -> str:
        if self._env["telegram"]["token"] == "":
            raise VariableEmptyError("telegram token is empty")

        return self._env["telegram"]["token"]

    @property
    def exchanges(self) -> ExchangesMap:
        exchangesMap = ExchangesMap()
        for exchange, keyMap in self._env["exchanges"].items():
            done = exchangesMap.add_exchange(
                exchange,
                keyMap["key"],
                keyMap["secret"],
            )

            if not done:
                raise ValueError(
                    f"exchange: {exchange} is empty key({keyMap['key']}) or secret({keyMap['secret']})"  # noqa: E501
                )

        return exchangesMap


def init() -> Optional[SettingsToml]:
    try:
        with open("env.toml", "r") as f:
            env = toml.load(f)
            return SettingsToml(env)
    except FileNotFoundError:
        return None
