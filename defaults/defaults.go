// Default package implements versioning primitives.
package defaults

// Major version number.
const Major = "0"

// Minor version number.
const Minor = "2"

// Release number. It will be overwritten during build. Do not try to manage it here.
var Release = "0"

// Version compiled into a string.
var Version = Major + "." + Minor + "." + Release

// Default Content-Type used for web calls.
const ContentType = "application/json; charset=utf8"

// Telegram API base URL
const TelegramAPIBase string = "https://api.telegram.org/bot"

// Debug messages
const Debug = false
