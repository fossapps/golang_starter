//go:generate mockgen -destination=./mock_rate_limiter.go -package=mocks crazy_nl_backend/adapters IRateLimiter
//go:generate mockgen -destination=./mock_db.go -package=mocks crazy_nl_backend/db Db
//go:generate mockgen -destination=./mock_devices.go -package=mocks crazy_nl_backend/db IDeviceManager
//go:generate mockgen -destination=./mock_migration.go -package=mocks crazy_nl_backend/db IMigrationManager
//go:generate mockgen -destination=./mock_permissions.go -package=mocks crazy_nl_backend/db IPermissionManager
//go:generate mockgen -destination=./mock_refresh_tokens.go -package=mocks crazy_nl_backend/db IRefreshTokenManager
//go:generate mockgen -destination=./mock_users.go -package=mocks crazy_nl_backend/db IUserManager
//go:generate mockgen -destination=./mock_redis.go -package=mocks crazy_nl_backend/helpers IRedisClient
//go:generate mockgen -destination=./mock_redis.go -package=mocks crazy_nl_backend/helpers IRedisClient
//go:generate mockgen -destination=./mock_migration_implementation.go -package=mocks crazy_nl_backend/migrations IMigration

package mocks
