package session

import "time"

const version = "v1.0.0"

const defaultSessionKeyName = "sessionid"
const defaultCookieLen uint32 = 32
const defaultExpires = time.Hour * 2
const defaultGCLifetime int64 = 3

const base64Table = "1234567890poiuytreqwasdfghjklmnbvcxzQWERTYUIOPLKJHGFDSAZXCVBNM-_"

const cookieCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
