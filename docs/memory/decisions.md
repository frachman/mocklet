# Decisions

1. M1 is the active delivery boundary; M2 is required before claiming the full MVP.
2. Go standard HTTP + pgx + SQL migrations are the initial implementation choices.
3. Anonymous resources expire after 24 hours, with policy values kept easy to change.
4. Public keys identify runtime resources; management tokens are separate and hashed.

