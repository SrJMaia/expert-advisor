from binance.client import Client
from time import sleep
from binance import ThreadedWebsocketManager
from binance.enums import *
from binance.exceptions import BinanceAPIException, BinanceOrderException

api_key_real    = "4DxDa3VoWiDDLbzE8rsJ9OEOrQdXOjB5vKBAMPtJEHkkuKfcZLKKqEqN37Dxhgm3"
api_secret_real = "pdBDYEvzV2t1IWaoBl3XYNu5pzp2QUvAfMMhnm0pROYuzCLrHgSf1J42wpgyuHXs"

client = Client(api_key_real, api_secret_real)

print("-"*90) 

def expert_advisor(msg):
    flag_open_trade = True
    balance = 100
    should_exit = False
    buyId = 0
    sellId = 0
    trades = client.get_open_orders(symbol="USDCUSDT")
    if trades != []:
        try:
            if trades[0]["side"] == "BUY":
                buyId = trades[0]["orderId"]
                flag_open_trade = False
        except:
            if trades[0]["side"] == "SELL":
                sellId = trades[0]["orderId"]
                flag_open_trade = False

    should_exit = check_account(client, should_exit)
    should_exit = check_api(client, should_exit)
    
    capital = get_balance(client, balance)

    price = get_prices(msg)

    print(f"Teste {price}")

    if buyId != 0 and price["bid"] <= 0.9997 and not flag_open_trade:
        cancel_order(client, buyId)
        sellId = sell(client, price["bid"], capital)
        return
    elif sellId != 0 and price["ask"] >= 1.0003 and not flag_open_trade:
        cancel_order(client, sellId)
        buyId = buy(client, price['ask'], capital)
        return

    if price["ask"] >= 1.0003 and flag_open_trade:
        buyId = buy(client, price['ask'], capital)
        flag_open_trade = False
    elif price["bid"] <= 0.9997 and flag_open_trade:
        sellId = sell(client, price["bid"], capital)
        flag_open_trade = False

    if should_exit:
        bsm.stop()

bsm = ThreadedWebsocketManager()
bsm.start()

bsm.start_symbol_ticker_socket(callback=expert_advisor, symbol="USDCUSDT")

def get_balance(client, balance):
    response = client.get_asset_balance(asset="USDC")
    if response != None:
        balance = response
        return balance
    return balance

def check_account(client , should_exit):
    status = client.get_account_status()
    if status["data"] != "Normal":
        print("System Maintenance")
        bsm.stop()
        should_exit = True
        return should_exit

def check_api(client, should_exit):
    status = client.get_system_status()
    if status["status"] != 0:
        print("System Maintenance")
        bsm.stop()
        should_exit = True
        return should_exit

def get_prices(msg):
    price = {'error':False}
    if msg['e'] != 'error':
        price['vol'] = float(msg['v'])
        price['open'] = float(msg['o'])
        price['high'] = float(msg['h'])
        price['low'] = float(msg['l'])
        price['close'] = float(msg['c'])
        price['bid'] = float(msg['b']) # Sell orders
        price['ask'] = float(msg['a']) # Buy orders
        price['error'] = False
    else:
        price['error'] = True
    return price

def buy(client, buy_price, volume, symbol="USDCUSDT"):
    buy_limit = {}
    try:
        buy_limit = client.create_order(
            symbol=symbol,
            side="BUY",
            type="LIMIT",
            timeInForce="GTC",
            quantity=volume,
            price=buy_price)
    except BinanceAPIException as e:
        print(e)
    except BinanceOrderException as e:
        print(e)
    else:
        return buy_limit['orderId']

def sell(client, sell_price, volume, symbol="USDCUSDT"):
    sell_limit = {}
    try:
        sell_limit = client.create_order(
            symbol=symbol,
            side="SELL",
            type="LIMIT",
            timeInForce="GTC",
            quantity=volume,
            price=sell_price)
    except BinanceAPIException as e:
        print(e)
    except BinanceOrderException as e:
        print(e)
    else:
        return sell_limit['orderId']

def cancel_order(client, id, symbol="USDCUSDT"):
    try:
        cancel = client.cancel_order(
        symbol=symbol, 
        orderId=id)
    except BinanceAPIException as e:
        print(e)
    except BinanceOrderException as e:
        print(e)
    else:
        print("Order cancelled")