package session

import "time"

const version = "v1.0.0"

const defaultSessionKeyName = "sessionid"
const defaultCookieLen uint32 = 32
const defaultExpires = 2 * time.Hour
const defaultGCLifetime = 30 * time.Second
const defaultSessionLifetime int64 = 60
const defaultSecure = true
const defaultSessionIDInURLQuery = false
const defaultSessionIDInHTTPHeader = false

const base64Table = "1234567890poiuytreqwasdfghjklmnbvcxzQWERTYUIOPLKJHGFDSAZXCVBNM-_"

const cookieCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
