CREATE OR REPLACE VIEW user_stats AS
SELECT
  u.unid,

  COALESCE(SUM(CASE
    WHEN tl.link_type = 'output' AND c.spent = FALSE THEN c.value
    ELSE 0
  END), 0) AS balance,

  COALESCE(SUM(CASE
    WHEN tl.link_type = 'output' THEN c.value
    ELSE 0
  END), 0) AS total_earned,

  COALESCE(SUM(CASE
    WHEN tl.link_type = 'output' AND c.spent = TRUE THEN c.value
    ELSE 0
  END), 0) AS total_spent,

  COALESCE(SUM(CASE
    WHEN tl.link_type = 'output' THEN p.message_count
    ELSE 0
  END), 0) AS total_messages,

  CASE
    WHEN COALESCE(SUM(CASE WHEN tl.link_type = 'output' THEN p.message_count ELSE 0 END), 0) = 0
    THEN NULL
    ELSE ROUND(
      SUM(CASE WHEN tl.link_type = 'output' THEN c.value ELSE 0 END) /
      SUM(CASE WHEN tl.link_type = 'output' THEN p.message_count ELSE 0 END),
    4)
  END AS avg_coins_per_message,

  COUNT(DISTINCT CASE WHEN tl.link_type = 'output' AND p.id IS NOT NULL THEN p.id END) AS payout_count,
  COUNT(DISTINCT tx_init.txid) AS transactions_initiated,
  COUNT(DISTINCT CASE WHEN a.winner_unid = u.unid THEN a.auction_id END) AS auctions_won,
  COUNT(DISTINCT CASE WHEN a.seller_unid = u.unid THEN a.auction_id END) AS auctions_created,
  DATEDIFF(NOW(), u.timestamp) AS account_age_days,
  MIN(CASE WHEN tl.link_type = 'output' THEN tl.timestamp END) AS first_earned_at,
  MAX(CASE WHEN tl.link_type = 'output' THEN tl.timestamp END) AS last_earned_at

FROM users u
LEFT JOIN coins c ON c.owner_unid = u.unid
LEFT JOIN transaction_links tl ON tl.coin_id = c.coin_id
LEFT JOIN payouts p ON p.txlink = tl.id
LEFT JOIN transactions tx_init ON tx_init.initiated_by = u.unid
LEFT JOIN auctions a ON a.winner_unid = u.unid OR a.seller_unid = u.unid

GROUP BY u.unid, u.timestamp
ORDER BY total_earned DESC;
