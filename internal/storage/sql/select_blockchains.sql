WITH max_ids AS (
  SELECT 
    blockchain_uuid,
    MAX(id) AS max_id
  FROM vb.blocks
  GROUP BY blockchain_uuid
)
SELECT b.blockchain_uuid, b.header AS last_block_header
FROM vb.blocks AS b
JOIN max_ids mi
ON b.blockchain_uuid = mi.blockchain_uuid AND b.id = mi.max_id;