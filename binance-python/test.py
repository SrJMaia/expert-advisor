from binance.client import Client
from time import sleep
from binance import ThreadedWebsocketManager
from binance.enums import *
from binance.exceptions import BinanceAPIException, BinanceOrderException
import time

api_key_real    = "X"
api_secret_real = "X"

client = Client(api_key_real, api_secret_real)

print("-"*180)

start = time.time()
book = client.get_order_book(symbol='USDCUSDT')
end = time.time()
print("Time Price: " ,end - start)

start = time.time()
balance_usdc = client.get_asset_balance(asset="USDC")
balance_usdt = client.get_asset_balance(asset="USDT")
end = time.time()
print("Time Balance: " ,end - start)

start = time.time()
orders = client.get_open_orders(symbol='USDCUSDT')
end = time.time()
print("Time Orders: " ,end - start)

print("-"*180)

print("Bid Price: ", book["bids"][0][0]) # The first price to buy, but its used to sell
print("Ask Price: ", book["asks"][0][0]) # The first price to sell, but its used to buy
print("USDC Balance: ", balance_usdc) # Balance to sell
print("USDT Balance: ", balance_usdt) # Balance to buy
print("Orders: ", orders)

print("-"*180)
