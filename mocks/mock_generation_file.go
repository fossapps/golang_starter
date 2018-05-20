//go:generate mockgen -destination=./mock_rate_limiter.go -package=mocks starter/adapters IRateLimiter
//go:generate mockgen -destination=./mock_db.go -package=mocks starter/db Db
//go:generate mockgen -destination=./mock_devices.go -package=mocks starter/db IDeviceManager
//go:generate mockgen -destination=./mock_migration.go -package=mocks starter/db IMigrationManager
//go:generate mockgen -destination=./mock_permissions.go -package=mocks starter/db IPermissionManager
//go:generate mockgen -destination=./mock_refresh_tokens.go -package=mocks starter/db IRefreshTokenManager
//go:generate mockgen -destination=./mock_users.go -package=mocks starter/db IUserManager
//go:generate mockgen -destination=./mock_redis.go -package=mocks starter/helpers IRedisClient
//go:generate mockgen -destination=./mock_redis.go -package=mocks starter/helpers IRedisClient
//go:generate mockgen -destination=./mock_migration_implementation.go -package=mocks starter/migrations IMigration
//go:generate mockgen -destination=./mock_logger.go -package=mocks starter ILogger
//go:generate mockgen -destination=./mock_request_helper.go -package=mocks starter IRequestHelper

package mocks
