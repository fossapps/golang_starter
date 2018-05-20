//go:generate mockgen -destination=./mock_rate_limiter.go -package=mocks github.com/fossapps/starter/adapters IRateLimiter
//go:generate mockgen -destination=./mock_db.go -package=mocks github.com/fossapps/starter/db Db
//go:generate mockgen -destination=./mock_devices.go -package=mocks github.com/fossapps/starter/db IDeviceManager
//go:generate mockgen -destination=./mock_migration.go -package=mocks github.com/fossapps/starter/db IMigrationManager
//go:generate mockgen -destination=./mock_permissions.go -package=mocks github.com/fossapps/starter/db IPermissionManager
//go:generate mockgen -destination=./mock_refresh_tokens.go -package=mocks github.com/fossapps/starter/db IRefreshTokenManager
//go:generate mockgen -destination=./mock_users.go -package=mocks github.com/fossapps/starter/db IUserManager
//go:generate mockgen -destination=./mock_redis.go -package=mocks github.com/fossapps/starter/helpers IRedisClient
//go:generate mockgen -destination=./mock_redis.go -package=mocks github.com/fossapps/starter/helpers IRedisClient
//go:generate mockgen -destination=./mock_migration_implementation.go -package=mocks github.com/fossapps/starter/migrations IMigration
//go:generate mockgen -destination=./mock_logger.go -package=mocks github.com/fossapps/starter ILogger
//go:generate mockgen -destination=./mock_request_helper.go -package=mocks github.com/fossapps/starter IRequestHelper

package mocks
