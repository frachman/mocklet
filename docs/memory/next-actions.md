# Next actions

1. Create `mocklet.mikrolyt.com` DNS for the production VPS.
2. Confirm the SFTP transport used by Hookbin can receive a separate `mocklet-*.dump.gpg` namespace; direct VPS access to homelab `192.168.1.5:22` is currently unreachable.
3. Rehearse encrypted backup and isolated restore for the Mocklet database.
4. Run target-specific preflight and obtain explicit approval before Caddy/Compose production changes.
