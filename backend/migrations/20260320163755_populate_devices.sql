-- +goose Up
INSERT INTO devices (hostname, ip, is_active, is_deleted, created_at) VALUES
    ('edge-router-01',     '192.168.10.1',    true,  false, '2025-11-15 08:30:00'),
    ('core-sw-ldn-1',      '10.240.0.11',     true,  false, '2025-10-03 14:15:00'),
    ('access-ap-warehouse','172.16.45.120',   true,  false, '2026-01-22 09:45:00'),
    ('fw-dmz-prod',        '203.0.113.45',    true,  false, '2025-09-12 11:20:00'),
    ('cam-front-gate-01',  '192.168.200.81',  false, false, '2026-02-08 16:10:00'),
    ('printer-mfp-3f',     '10.30.5.203',     true,  false, '2025-12-19 10:00:00'),
    ('iot-temp-sensor-07', 'fd12:3456::1a:7', true,  false, '2026-03-01 07:25:00'),
    ('nas-backup-01',      '192.168.88.250',  true,  false, '2025-08-30 13:40:00'),
    ('test-laptop-dev01',  '10.99.1.105',     false, false, '2026-03-10 15:55:00'),
    ('legacy-voip-pbx',    '172.31.254.10',   false, true,  '2024-06-05 09:15:00');

-- +goose Down
DELETE FROM devices
WHERE hostname IN (
    'edge-router-01',
    'core-sw-ldn-1',
    'access-ap-warehouse',
    'fw-dmz-prod',
    'cam-front-gate-01',
    'printer-mfp-3f',
    'iot-temp-sensor-07',
    'nas-backup-01',
    'test-laptop-dev01',
    'legacy-voip-pbx'
);
