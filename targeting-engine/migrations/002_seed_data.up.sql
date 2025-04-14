INSERT INTO campaigns (id, name, image_url, cta, status) VALUES
('spotify', 'Spotify - Music for everyone', 'https://somelink', 'Download', 'ACTIVE'),
('duolingo', 'Duolingo: Best way to learn', 'https://somelink2', 'Install', 'ACTIVE'),
('subwaysurfer', 'Subway Surfer', 'https://somelink3', 'Play', 'ACTIVE');


INSERT INTO targeting_rules (campaign_id, dimension, operation, values) VALUES

('spotify', 'country', 'include', '["US", "Canada"]'),

('duolingo', 'os', 'include', '["Android", "iOS"]'),
('duolingo', 'country', 'exclude', '["US"]'),


('subwaysurfer', 'os', 'include', '["Android"]'),
('subwaysurfer', 'app', 'include', '["com.gametion.ludokinggame"]');

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app_user') THEN
        CREATE ROLE app_user WITH LOGIN PASSWORD 'securepassword';
    END IF;
END
$$;

GRANT CONNECT ON DATABASE targeting TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO app_user;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;