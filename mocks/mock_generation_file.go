//go:generate mockgen -destination=./mock_rate_limiter.go -package=mocks golang_starter/adapters IRateLimiter
//go:generate mockgen -destination=./mock_db.go -package=mocks golang_starter/db Db
//go:generate mockgen -destination=./mock_devices.go -package=mocks golang_starter/db IDeviceManager
//go:generate mockgen -destination=./mock_migration.go -package=mocks golang_starter/db IMigrationManager
//go:generate mockgen -destination=./mock_permissions.go -package=mocks golang_starter/db IPermissionManager
//go:generate mockgen -destination=./mock_refresh_tokens.go -package=mocks golang_starter/db IRefreshTokenManager
//go:generate mockgen -destination=./mock_users.go -package=mocks golang_starter/db IUserManager
//go:generate mockgen -destination=./mock_redis.go -package=mocks golang_starter/helpers IRedisClient
//go:generate mockgen -destination=./mock_redis.go -package=mocks golang_starter/helpers IRedisClient
//go:generate mockgen -destination=./mock_migration_implementation.go -package=mocks golang_starter/migrations IMigration
//go:generate mockgen -destination=./mock_logger.go -package=mocks golang_starter ILogger
//go:generate mockgen -destination=./mock_request_helper.go -package=mocks golang_starter IRequestHelper

package mocks
