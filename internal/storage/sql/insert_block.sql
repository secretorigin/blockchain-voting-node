-- insert new block and delete all payloads in other blocks
WITH insert_block AS (
  INSERT INTO vb.blocks (
    blockchain_uuid,
    id,
    header,
    payload
  ) VALUES (
    $1, $2, $3, $4
  )
)
UPDATE vb.blocks
SET payload = NULL
WHERE blockchain_uuid = $1 AND id != $2 AND payload IS NOT NULL;