package cassandra

const insertByMerchant = `
INSERT INTO auth_events_by_merchant
(merchant_id, event_time, auth_id, event_id, event_type, decision, payload)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

const insertByCard = `
INSERT INTO auth_events_by_card
(card_hash, event_time, auth_id, merchant_id, event_id, event_type, decision, payload)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

const insertByTime = `
INSERT INTO auth_events_by_time
(day_bucket, event_time, auth_id, merchant_id, card_hash, event_type, decision, payload, event_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`
