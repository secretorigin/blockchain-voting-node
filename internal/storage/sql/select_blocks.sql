-- $1 - blockchain uuid, $2 - first id (inclusive), $3 last id (inclusive)
SELECT 
  id,
  header,
  payload
FROM vb.blocks
WHERE blockchain_uuid = $1 AND id >= $2 AND id <= $3;