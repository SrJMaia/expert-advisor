# expert-advisor  

### Implementations  
- In backtestHear change to return a pointer 
- Instead returning three slices, return a struct in backtest  
- Unit test finance calculation improve  
Mudar backtest?
    - Após checar se o candle é do tf ou nao, ou seja, o primeiro if else
    - Dividir o codigo posterior em pequenas funções, e chamar cada uma no qual queira
- Delete read csv ?  
- I could add a trade per day?  
- Change the backtest
- Add binance?  
- Create plot library  

### Bugs
- 001  
 There is a bug in myanalysis/save.go/SaveOtimization that in netProfit or the teenth line, after the last value will generate an extra comma  
 Solution: The append in base.go was appending to the wrong slice