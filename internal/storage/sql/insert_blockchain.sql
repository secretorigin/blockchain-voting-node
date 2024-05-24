-- create blockchain and insert first block
WITH create_block_chain AS (
  INSERT INTO vb.blockchains (uuid) VALUES ($1) 
)
INSERT INTO vb.blocks (
  blockchain_uuid,
  id,
  header,
  payload
) VALUES (
  $1, 0, $2, $3
);