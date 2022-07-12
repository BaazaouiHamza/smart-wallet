-- name: GetPolicyEvent :one
select * from policy_events WHERE nym_id=$1 AND transfer_sequence=$2;

-- name: InsertPolicyEvent :one
INSERT INTO policy_events (nym_id,transfer_sequence,policy_id) VALUES($1,$2,$3) RETURNING *;