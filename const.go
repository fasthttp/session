package session

import "time"

const version = "v1.0.0"

const defaultSessionKeyName = "sessionid"
const defaultCookieLen uint32 = 32
const defaultExpires = time.Hour * 2
const defaultGCLifetime int64 = 3

const base64Table = "1234567890poiuytreqwasdfghjklmnbvcxzQWERTYUIOPLKJHGFDSAZXCVBNM-_"

const cookieCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const cookieIdxBits = 6                    // 6 bits to represent a cookie index
const cookieIdxMask = 1<<cookieIdxBits - 1 // All 1-bits, as many as cookieIdxBits
const cookieIdxMax = 63 / cookieIdxBits    // # of letter indices fitting in 63 bits
