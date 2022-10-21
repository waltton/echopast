-- \timing
-- \x

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

            WHEN user_agent ~ 'python-requests' THEN 'Python requests'
            WHEN user_agent ~ 'Python-urllib' THEN 'Python urllib'
            WHEN user_agent ~ 'python-httpx' THEN 'Python httpx'

            WHEN user_agent ~ 'Go-http-client' THEN 'Go http client'
            WHEN user_agent ~ 'quic-go' THEN 'Quic Go'

            WHEN user_agent ~ 'Apache-HttpClient.*Java.*' THEN 'Java Apache HttpClient'

            WHEN user_agent ~ 'libwww-perl' THEN 'libwww perl'
        END AS user_agent_client
        , user_agent
        , date_trunc('day', timestamp) as timestamp
    FROM logs
    GROUP BY date_trunc('day', timestamp), user_agent
    ORDER BY COUNT(*) DESC
), base2 AS (
    SELECT *
        , RANK() OVER (ORDER BY COALESCE(count_current_week, 0) DESC) AS rank_current_week
        , RANK() OVER (ORDER BY COALESCE(count_last_week, 0) DESC) AS rank_last_week
        , RANK() OVER (ORDER BY COALESCE(count_last_week, 0) DESC) - RANK() OVER (ORDER BY COALESCE(count_current_week, 0) DESC) as delta
    FROM (
        SELECT coalesce(user_agent_client, user_agent) AS user_agent_group
            , sum(count) FILTER (WHERE timestamp >= NOW()::DATE - '7 days'::interval ) as count_current_week
            , sum(count) FILTER (WHERE timestamp >= NOW()::DATE - '14 days'::interval AND timestamp < NOW()::DATE - '7 days'::interval ) as count_last_week
        FROM base
        GROUP BY coalesce(user_agent_client, user_agent)
        ORDER BY sum(count) desc
    ) _
    ORDER BY COALESCE(count_current_week, 0) DESC
    LIMIT 15
)

SELECT json_agg(
        json_build_object(
            'user_agent_group', user_agent_group,
            'count_current_week', count_current_week,
            'count_last_week', count_last_week,
            'rank_current_week', rank_current_week,
            'rank_last_week', rank_last_week,
            'delta', delta
        )
        ORDER BY count_current_week DESC
    )
-- SELECT *
FROM base2