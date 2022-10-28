-- \timing
-- \x

\set http_ip_address '34.107.159.115'
\set https_ip_address '34.111.236.80'
\set domain 'averyrandomdomainname.com'

WITH base AS (
    SELECT COUNT(*)
        , CASE
            WHEN user_agent ~ '^(?!.*Edge).*Chrome' THEN 'Chrome'
            WHEN user_agent ~ '^(?!.*(?:Chrome|Edge)).*Safari' THEN 'Safari'
            WHEN user_agent ~ 'MSIE ([0-9]{1,}[\.0-9]{0,})' THEN 'Internet Explorer'
            WHEN user_agent ~ 'Firefox\/(\d+(?:\.\d+)+)' THEN 'Firefox'
            WHEN user_agent ~ 'Edge' THEN 'Edge'


            WHEN user_agent ~ 'Expanse, a Palo Alto' THEN 'Expanse'
            WHEN user_agent ~ 'compatible; CensysInspect' THEN 'CensysInspect'
            WHEN user_agent ~ 'compatible; AhrefsBot' THEN 'AhrefsBot'
            WHEN user_agent ~ 'compatible; InternetMeasurement' THEN 'InternetMeasurement'
            WHEN user_agent ~ 'Cloud mapping experiment. Contact research@pdrlabs.net' THEN 'Cloud mapping experiment'

            WHEN user_agent ~ 'curl' THEN 'curl'
            WHEN user_agent ~ 'zgrab' THEN 'zgrab'
            WHEN user_agent ~ 'masscan-ng' THEN 'masscan ng'
            WHEN user_agent ~ 'masscan' THEN 'masscan'
            WHEN user_agent ~ 'okhttp' THEN 'okhttp'
            WHEN user_agent ~ 'l9tcpid' THEN 'l9tcpid'
            WHEN user_agent ~ 'l9explore' THEN 'l9explore'
            WHEN user_agent ~ 'ZmEu' THEN 'ZmEu'
            WHEN user_agent ~ 'WhatWeb' THEN 'WhatWeb'


            WHEN user_agent ~ '[Pp]ython(\-requests|\-urllib|\-httpx|.*aiohttp)' THEN 'Python (requests/urllib/httpx/aiohttp)'

            WHEN user_agent ~ 'Go-http-client' THEN 'Go http client'
            WHEN user_agent ~ 'quic-go' THEN 'Quic Go'

            WHEN user_agent ~ 'Apache-HttpClient.*Java.*' THEN 'Java Apache HttpClient'

            WHEN user_agent ~ 'libwww-perl' THEN 'libwww perl'

        END AS user_agent_client
        , user_agent
        , COUNT(DISTINCT ip) as ips
        , COUNT(*) FILTER (WHERE protocol = 'http' AND host = :'http_ip_address') AS c1
        , COUNT(*) FILTER (WHERE protocol = 'https' AND host = :'https_ip_address') AS c2
        , COUNT(*) FILTER (WHERE protocol = 'http' AND host = :'domain') AS c3
        , COUNT(*) FILTER (WHERE protocol = 'https' AND host = :'domain') AS c4
        , COUNT(*) FILTER (WHERE protocol = 'http' AND NOT host = :'http_ip_address' AND NOT host = :'domain') AS c5
        , COUNT(*) FILTER (WHERE protocol = 'https' AND NOT host = :'https_ip_address' AND NOT host = :'domain') AS c6
        , date_trunc('day', timestamp) as timestamp
    FROM logs
    GROUP BY date_trunc('day', timestamp), user_agent
    ORDER BY COUNT(*) DESC
), base2 AS (
    SELECT *
        , RANK() OVER (ORDER BY COALESCE(count_current_week, 0) DESC) AS rank_current_week
        , RANK() OVER (ORDER BY COALESCE(count_last_week, 0) DESC) AS rank_last_week
    FROM (
        SELECT coalesce(user_agent_client, user_agent) AS user_agent_group
            , sum(count) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as count_current_week
            , sum(count) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as count_last_week

            , AVG(ips) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as ips_current_week
            , AVG(ips) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as ips_last_week

            , SUM(c1) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as c1_current_week
            , SUM(c1) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as c1_last_week

            , SUM(c2) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as c2_current_week
            , SUM(c2) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as c2_last_week

            , SUM(c3) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as c3_current_week
            , SUM(c3) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as c3_last_week

            , SUM(c4) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as c4_current_week
            , SUM(c4) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as c4_last_week

            , SUM(c5) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as c5_current_week
            , SUM(c5) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as c5_last_week

            , SUM(c6) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) ) as c6_current_week
            , SUM(c6) FILTER (WHERE date_trunc('week', timestamp) = date_trunc('week', now()) - '7 days'::interval ) as c6_last_week

        FROM base
        GROUP BY coalesce(user_agent_client, user_agent)
        ORDER BY sum(count) desc
    ) _
    ORDER BY COALESCE(count_current_week, 0) DESC, COALESCE(count_last_week, 0) DESC
    LIMIT 15
)

SELECT json_agg(
        json_build_object(
            'user_agent_group', user_agent_group,
            'count', count_current_week,
            'count_delta', COALESCE(count_last_week, 0) - COALESCE(count_current_week, 0),
            'rank', rank_current_week,
            'rank_delta', COALESCE(rank_last_week, 0) - COALESCE(rank_current_week, 0),

            'daily_ips_avg', ips_current_week,
            'daily_ips_avg_delta', COALESCE(ips_last_week, 0) - COALESCE(ips_current_week, 0),

            'c1', c1_current_week,
            'c2', c2_current_week,
            'c3', c3_current_week,
            'c4', c4_current_week,
            'c5', c5_current_week,
            'c6', c6_current_week
        )
        ORDER BY count_current_week DESC
    )
-- SELECT user_agent_group, count_current_week, count_current_week
FROM base2