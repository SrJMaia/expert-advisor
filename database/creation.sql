CREATE TABLE prices.eurusd (
	index_table INT NOT NULL AUTO_INCREMENT,
	date_hour TIMESTAMP NOT NULL,
	open_prices DECIMAL(6,5) NOT NULL,
	high_prices DECIMAL(6,5) NOT NULL,
	low_prices DECIMAL(6,5) NOT NULL,
	close_prices DECIMAL(6,5) NOT NULL,
	PRIMARY KEY (index_table)
);

DROP TABLE prices.eurusd;

SELECT * FROM prices.eurusd LIMIT 100;

SELECT COUNT(eurusd.index_table) FROM prices.eurusd;

-- H4
SELECT 
	MIN(eurusd.date_hour)
FROM 
	prices.eurusd
GROUP BY 
	DATE(eurusd.date_hour), HOUR(eurusd.date_hour) 
LIKE
	"% 04:_____"
LIMIT
	100;

-- D1
SELECT 
	*
FROM
	prices.eurusd
WHERE
	eurusd.index_table 
IN (
	SELECT 
		MIN(eurusd.index_table)
	FROM 
		prices.eurusd
	GROUP BY 
		DATE(eurusd.date_hour), DAY(eurusd.date_hour) 
)
LIMIT  
	100;

-- H1
SELECT 
	COUNT(eurusd.index_table)
FROM
	prices.eurusd
WHERE
	eurusd.index_table 
IN (
	SELECT 
		MIN(eurusd.index_table)
	FROM 
		prices.eurusd
	GROUP BY 
		DATE(eurusd.date_hour), HOUR(eurusd.date_hour) 
);
	