package session

import "time"

const defaultSessionKeyName = "sessionid"
const defaultDomain = ""
const defaultExpiration = 2 * time.Hour
const defaultGCLifetime = 1 * time.Minute
const defaultSecure = true
const defaultSessionIDInURLQuery = false
const defaultSessionIDInHTTPHeader = false
const defaultCookieLen uint32 = 32

// If set the cookie expiration when the browser is closed (-1), set the expiration as a keep alive (2 days)
// so as not to keep dead sessions for a long time
const keepAliveExpiration = 2 * 24 * time.Hour
const expirationAttrKey = "__store:expiration__"
