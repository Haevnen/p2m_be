ALTER TABLE links ADD COLUMN client_id VARCHAR(20) AFTER id;

UPDATE links
SET client_id = SUBSTRING_INDEX(
    SUBSTRING_INDEX(link, '\\UPLOAD\\', 1), 
    '\\', 
    -1
)
WHERE link LIKE '%\\UPLOAD\\%';