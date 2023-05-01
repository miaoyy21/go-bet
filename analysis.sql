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

-- 按照返奖率统计指定日期输赢情况
SELECT '[0, 0.85]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx <= 0.85
UNION ALL
SELECT '(0.85, 0.90]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 0.85 AND rx <= 0.90
UNION ALL
SELECT '(0.90, 0.95]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 0.90 AND rx <= 0.95
UNION ALL
SELECT '(0.95, 1.00]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 0.95 AND rx <= 1.00
UNION ALL
SELECT '(1.00, 1.05]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 1.00 AND rx <= 1.05
UNION ALL
SELECT '(1.05, 1.10]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 1.05 AND rx <= 1.10
UNION ALL
SELECT '(1.10, 1.15]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 1.10 AND rx <= 1.15
UNION ALL
SELECT '(1.15, 1.20]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 1.15 AND rx <= 1.20
UNION ALL
SELECT '(1.20, 1.25]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 1.20 AND rx <= 1.25
UNION ALL
SELECT '(1.25, 5]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-01 %' AND rx > 1.25

-- 按照统计指定日期的每个小时输赢情况
SELECT LEFT(time,13),CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate
FROM logs
WHERE time LIKE '2023-05-01 %'
GROUP BY LEFT(time,13);

