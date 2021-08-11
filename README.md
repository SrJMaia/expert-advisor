# expert-advisor  

### Implementations  
- I need to have a flag called non eur and when i use a non eur currency pair, i need to pass eur price   
- Instead returning three slices, return a struct in backtest  ?
- Delete read csv ?  
- I could add a trade per day?  
- Add binance?  
- Create plot library  

### Bugs
- 001  
 There is a bug in myanalysis/save.go/SaveOtimization that in netProfit or the teenth line, after the last value will generate an extra comma  
 Solution: The append in base.go was appending to the wrong slice