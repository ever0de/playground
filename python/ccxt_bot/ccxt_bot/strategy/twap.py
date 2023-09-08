import ccxt
import pandas as pd


def calculate_twap(
    exchange: ccxt.Exchange, symbol: str, timeframe: str, start_date: str, end_date: str
):
    start_date = exchange.parse8601(start_date)

    ohlcv = exchange.fetch_ohlcv(
        symbol, timeframe, exchange.parse8601(start_date), exchange.parse8601(end_date)
    )

    df = pd.DataFrame(
        ohlcv, columns=["timestamp", "open", "high", "low", "close", "volume"]
    )

    df["timestamp"] = pd.to_datetime(df["timestamp"], unit="ms")
    df["twap"] = (df["close"] * df["volume"]).cumsum() / df["volume"].cumsum()

    return df


# symbol = "BTC/USDT"
# timeframe = "1h"
# start_date = "2023-01-01T00:00:00Z"
# end_date = "2023-01-02T00:00:00Z"
