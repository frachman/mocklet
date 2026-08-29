# Architecture overview

Mocklet has a management plane and a runtime plane. The management plane creates and inspects disposable mock configuration. The runtime plane resolves a public key and returns the configured response. PostgreSQL is the independent Mocklet domain store.

The initial service is a single Go process with clear handler/repository boundaries. `owner_user_id` is nullable for future shared identity and no Mocklet credential store is introduced.

