SET @RESULT :=18;

SELECT TT.*
FROM (
	SELECT T.result,COUNT(1) AS NN,
		CASE 
			WHEN T.result = 0 OR T.result = 27 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/1000)
			WHEN T.result = 1 OR T.result = 26 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/333.33) 
			WHEN T.result = 2 OR T.result = 25 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/166.67) 
			WHEN T.result = 3 OR T.result = 24 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/100.00)
			WHEN T.result = 4 OR T.result = 23 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/66.67) 
			WHEN T.result = 5 OR T.result = 22 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/47.62) 
			WHEN T.result = 6 OR T.result = 21 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/35.71) 
			WHEN T.result = 7 OR T.result = 20 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/27.78)
			WHEN T.result = 8 OR T.result = 19 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/22.22) 
			WHEN T.result = 9 OR T.result = 18 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/18.18)
			WHEN T.result = 10 OR T.result = 17 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/15.87) 
			WHEN T.result = 11 OR T.result = 16 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/14.49)
			WHEN T.result = 12 OR T.result = 15 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/13.70) 
			WHEN T.result = 13 OR T.result = 14 THEN COUNT(1)/((SELECT COUNT(1) FROM histories WHERE result = @RESULT)/13.33)
			ELSE 1 
		END AS RATE
	FROM histories T
	WHERE EXISTS (SELECT 1 FROM histories X WHERE X.result = @RESULT AND X.issue + 1 = T.issue)
	GROUP BY T.result
) TT
ORDER BY TT.result ASC