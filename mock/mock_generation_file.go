//go:generate mockgen -destination=./mock_rate_limiter.go -package=mock github.com/fossapps/starter/middleware RateLimiter
//go:generate mockgen -destination=./mock_db.go -package=mock github.com/fossapps/starter/db DB
//go:generate mockgen -destination=./mock_devices.go -package=mock github.com/fossapps/starter/db DeviceManager
//go:generate mockgen -destination=./mock_migration.go -package=mock github.com/fossapps/starter/db MigrationManager
//go:generate mockgen -destination=./mock_permissions.go -package=mock github.com/fossapps/starter/db PermissionManager
//go:generate mockgen -destination=./mock_refresh_tokens.go -package=mock github.com/fossapps/starter/db RefreshTokenManager
//go:generate mockgen -destination=./mock_users.go -package=mock github.com/fossapps/starter/db UserManager
//go:generate mockgen -destination=./mock_redis.go -package=mock github.com/fossapps/starter/rate RedisClient
//go:generate mockgen -destination=./mock_migration_implementation.go -package=mock github.com/fossapps/starter/migration Migration
//go:generate mockgen -destination=./mock_logger.go -mock_names "Client=MockLogger" -package=mock github.com/fossapps/starter/logger Client
//go:generate mockgen -destination=./mock_jwt.go -mock_names "Manager=MockJwtManager" -package=mock github.com/fossapps/starter/jwt Manager
//go:generate mockgen -destination=./mock_request_helper.go -package=mock github.com/fossapps/starter RequestHelper
//go:generate mockgen -destination=./mock_pushy_client.go -package=mock github.com/cyberhck/pushy IPushyClient

package mock
