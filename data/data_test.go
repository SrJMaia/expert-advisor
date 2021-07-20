package data

import "testing"

func TestRead(t *testing.T) {
	var myData = LayoutData{}
	ReadData("C:/Users/johnk/Google Drive/Programming/Dados/EURUSD_Clean_20y.csv", &myData)
}
