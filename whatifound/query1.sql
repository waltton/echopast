-- \timing
-- \x

\set http_ip_address '34.107.159.115'
\set https_ip_address '34.111.236.80'
\set domain 'averyrandomdomainname.com'

-- Where are the requests going to?
-- 1 - To the IP address via http
-- 2 - To the IP address via https
-- 3 - To the domain via http
-- 4 - To the domain via https
-- 5 - To the IP address with a different hostname via http
-- 6 - To the IP address with a different hostname via https

WITH data AS (
    select date_trunc('day', timestamp) AS "day"
        , CASE
            WHEN protocol = 'http' AND host = :'http_ip_address' THEN 1
            WHEN protocol = 'https' AND host = :'https_ip_address' THEN 2
            WHEN protocol = 'http' AND host = :'domain' THEN 3
            WHEN protocol = 'https' AND host = :'domain' THEN 4
            WHEN protocol = 'http' AND NOT host = :'http_ip_address' AND NOT host = :'domain' THEN 5
            WHEN protocol = 'https' AND NOT host = :'https_ip_address' AND NOT host = :'domain' THEN 6
        END AS "category"
    from logs
    where 1=1
    and COALESCE(protocol, '') != ''
)

SELECT json_agg(day)
FROM (
    SELECT json_build_object(
            'date', "day"::date,
            'data', json_build_object(
            1, count(*) FILTER(WHERE category = 1),
            2, count(*) FILTER(WHERE category = 2),
            3, count(*) FILTER(WHERE category = 3),
            4, count(*) FILTER(WHERE category = 4),
            5, count(*) FILTER(WHERE category = 5),
            6, count(*) FILTER(WHERE category = 6)
        ))
    FROM data
    GROUP BY "day"
    ORDER BY "day" DESC
    LIMIT 30
)_(day)