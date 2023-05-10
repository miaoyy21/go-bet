
-- 按照返奖率统计指定日期输赢情况
SELECT '[0, 0.850]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx <= 0.850
UNION ALL
SELECT '(0.850, 0.875]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 0.850 AND rx <= 0.875
UNION ALL
SELECT '(0.875, 0.900]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 0.875 AND rx <= 0.900
UNION ALL
SELECT '(0.900, 0.925]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 0.900 AND rx <= 0.925
UNION ALL
SELECT '(0.925, 0.950]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 0.925 AND rx <= 0.950
UNION ALL
SELECT '(0.950, 0.975]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 0.950 AND rx <= 0.975
UNION ALL
SELECT '(0.975, 1.000]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 0.975 AND rx <= 1.000
UNION ALL
SELECT '(1.000, 1.025]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 1.000 AND rx <= 1.025
UNION ALL
SELECT '(1.025, 1.050]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 1.025 AND rx <= 1.050
UNION ALL
SELECT '(1.050, 1.075]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 1.050 AND rx <= 1.075
UNION ALL
SELECT '(1.075, 1.100]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 1.075 AND rx <= 1.100
UNION ALL
SELECT '(1.100, 1.125]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 1.100 AND rx <= 1.125
UNION ALL
SELECT '(1.125, 1.150]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 1.125 AND rx <= 1.150
UNION ALL
SELECT '(1.150, 5.00]' AS scope, COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-06 %' AND bet_gold > 1000 AND rx > 1.150;

-- 按照统计指定日期的每个小时输赢情况
SELECT LEFT(time,13),COUNT(1) AS qn, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate
FROM logs
WHERE time LIKE '2023-05-07 %' AND bet_gold > 1000
GROUP BY LEFT(time,13);

-- 查询未投注的期数中奖情况
SELECT SUM(CASE WHEN rt < 1.0 THEN 1 ELSE 0 END) AS ln1,SUM(CASE WHEN rt > 1.0 THEN 1 ELSE 0 END) AS gn1
FROM logs
WHERE time LIKE '2023-05-10 %' AND win_gold = 0;

