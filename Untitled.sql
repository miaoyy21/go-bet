
SELECT LEFT(time,10),SUM(win_gold),SUM(win_gold/dx0)
FROM logs
WHERE time LIKE '2023-05-04 %'
GROUP BY LEFT(time,10);