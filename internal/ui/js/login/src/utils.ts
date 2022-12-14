/**
 * Decodes a base64url-encoded string (which differs from base64 at digits 62 and 63)
 */
export function decodeBase64url(str: string) {
	return atob(str.replace(/-/g, "+").replace(/_/g, "/"))
}

/**
 * Turns a string (its characters should have a Unicode value below 256) into a Uint8Array
 * for use with APIs like Crypto that require an ArrayBuffer
 */
export function stringToBuffer(str: string) {
	return new Uint8Array([...str].map(c => c.charCodeAt(0)))
}
