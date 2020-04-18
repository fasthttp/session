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
