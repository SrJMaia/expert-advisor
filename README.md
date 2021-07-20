# expert-advisor  

### Implementations  
- Create a function called roundIsJpy  
- Probably my annualized return is wrong  
- I could add a trade per day?  


### Bugs
- 001  
 There is a bug in myanalysis/save.go/SaveOtimization that in netProfit or the teenth line, after the last value will generate an extra comma  
 Solution: The append in base.go was appending to the wrong slice