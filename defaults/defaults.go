// Default package implements versioning primitives.
package defaults

import "time"

// Major version number.
const Major = "0"

// Minor version number.
const Minor = "1"

// Default Content-Type used for web calls.
const ContentType = "application/json; charset=utf8"

// Release number. It will be overwritten during build. Do not try to manage it here.
var Release = "0-dev"

// Version compiled into a string.
var Version = Major + "." + Minor + "." + Release

// Telegram API base URL
const TelegramAPIBase string = "https://api.telegram.org/bot"

// User data TTL in database - 2 weeks
const TTL = 2 * 168 * time.Hour

// Debug messages
const Debug = false
