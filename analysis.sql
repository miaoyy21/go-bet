
-- 按照返奖率统计指定日期输赢情况
SELECT '[0, 0.95]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx <= 0.95
UNION ALL
SELECT '(0.95, 0.96]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 0.95 AND rx <= 0.96
UNION ALL
SELECT '(0.96, 0.97]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 0.96 AND rx <= 0.97
UNION ALL
SELECT '(0.97, 0.98]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 0.97 AND rx <= 0.98
UNION ALL
SELECT '(0.98, 0.99]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 0.98 AND rx <= 0.99
UNION ALL
SELECT '(0.99, 1.00]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 0.99 AND rx <= 1.00
UNION ALL
SELECT '(1.00, 1.01]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 1.00 AND rx <= 1.01
UNION ALL
SELECT '(1.01, 1.02]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 1.01 AND rx <= 1.02
UNION ALL
SELECT '(1.02, 1.03]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 1.02 AND rx <= 1.03
UNION ALL
SELECT '(1.03, 1.04]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 1.03 AND rx <= 1.04
UNION ALL
SELECT '(1.04, 1.05]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 1.04 AND rx <= 1.05
UNION ALL
SELECT '(1.05, 5.00]' AS scope, CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000 AND rx > 1.05;

-- 按照统计指定日期的每个小时输赢情况
SELECT LEFT(time,13),CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate
FROM logs
WHERE time LIKE '2023-05-05 %' AND bet_gold <> 1000
GROUP BY LEFT(time,13);

