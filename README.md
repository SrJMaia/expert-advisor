# expert-advisor  

### Implementations  
- Create a function called roundIsJpy  
- Probably my annualized return is wrong  
- I could add a trade per day?  
- Delete csv  
- Put the docker inside git  and organize the files  
- Change data to read the csv and throw inside database with paramaters datehour format and the name metatrader to get only the important  
- Check if the data sended is the len of the csvLines  
- Inside the database, puts H30 H1 H4 H12? D1 hours raw and formated to m1  
- Backup database  
- Maybe change the backtest?  
- Add binance?  
- Create plot library  


### Bugs
- 001  
 There is a bug in myanalysis/save.go/SaveOtimization that in netProfit or the teenth line, after the last value will generate an extra comma  
 Solution: The append in base.go was appending to the wrong slice