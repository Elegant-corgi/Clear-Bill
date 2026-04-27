# Server Config

The backend now loads static startup parameters from:

- `configs/config.toml`

Supported environment overrides:

- `CLEAR_BILL_CONFIG`: custom config file path
- `CLEAR_BILL_HTTP_ADDR`: override listen address
- `CLEAR_BILL_WEB_ROOT`: override static web root

The `[GORM]` section configures the MySQL connection used by GORM.
