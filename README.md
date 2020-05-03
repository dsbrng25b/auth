# Auth Overview
This repository tries to give an overview over several authentication mechanisms.

## HTTP Authorization header
The client sends the credentials in the `Authorization` header.
* WWW-Authenticate and Authorization header: https://tools.ietf.org/html/rfc7235#section-4.1
```
curl -H 'Authorization: credentials'
```

### Basic Authentication
In the basic auth scheme we send the username and password base64 encoded in the auth header: `Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=`.
```
curl -u john:tes http://localhost:8080/handle-basic
```
If a browser does not send the Authorization header, the server can request the credentials by sending the Header `WWW-Authenticate: Basic realm="Your Realm"` and HTTP code 401.

* HTTP Authentication: Basic and Digest Access Authentication: https://tools.ietf.org/html/rfc2617
* The 'Basic' HTTP Authentication Scheme: https://tools.ietf.org/html/rfc7617

## Bearer Token
Credentials in the auth header other then username and password are often called bearer tokens and are used like this: `Authorization: Bearer apfdo2x50mn6f1qa`. The token doesn't have to conform to a special encoding or format. Tokens which do not have any special meaning to the client are called opaque tokens:
```
curl -H 'Authorization: Bearer apfdo2x50mn6f1qa' http://localhost:8080/handle-bearer
```
The most common tokens which follow a standard format are JWT tokens.

## HTTP POST / Cookies / Sessions
Credentials can be sent to the server in a HTTP POST. The following are common formats:
* `application/x-www-form-urlencoded`
```
curl --data-urlencode 'user=myuser' --data-urlencode 'password=foobar$' http://localhost/login
```

* `multipart/form-data`
```
curl -F user=myuser -F foo='foobar$' http://localhost/login
```

* `application/json`
```
curl -d '{"username":"myuser","password":"foobar$"}' -H 'Content-Type: application/json' http://localhost/login
```

The server checks the credentials and if they are valid he either sets a cookie with the `Set-Cookie` header or returns a token in the payload, which the client has to store and send in subsequent requests.

### Cookies
* https://tools.ietf.org/html/rfc6265
Cookies are set by using the `Set-Cookie` HTTP header.
```
Set-Cookie: <cookie-name>=<cookie-value>; Domain=<domain-value>; Secure; HttpOnly
```


## JSON Web Token (JWT)
JWT tokens are digitally signed tokens which consist of three base64 encoded parts sperated by  dots:
```
base64( header ) + "." + base64( payload ) + "." + base64(signature)
```

* https://tools.ietf.org/html/rfc7519
* https://jwt.io/introduction/

### Header
The header specifies that the token is a JWT token and which signing algorithm is used:
```json
{
  "typ": "JWT",
  "alg": "HS256"
}
```

### Payload
The payload contains a number of values called [claims](https://tools.ietf.org/html/rfc7519#section-4). There are three kind of claims: registerd, public and private claims. According to the RFC no claim is requreid but some names (registerd, public) have a special meaning:
```json
{
  "iss": "the issuer",
  "aud": "the receiver",
  "sub": "subject", //user id
  "exp": 1588417578, //not after
  "nbf": 1588417278, //not before
  "iat": 1588417278, /issued at
  "jti": "jwt id"
}
```
A client can not change the values in the payload since the token is digitally signed.

### Signing
JWT tokens are digitally signed either using a symetric or asymetric algorithm (RSA, ECDSA).

#### Symetric
If a symetric signature algorithm (e.g. `HS256`) is used there is only one key. All servers which want to issue or verify JWT tokens need that key. The signature is created with the header and the payload:
```
hmac_sha_256( base64(header) + "." + base64(payload), key)
```

#### Asymetric
If an asymetric signature algorithm (e.g. `RS256`) is used we can issue/sign JWT tokens with the private key but to verify we only need the public key.

### Example
* Create a JWT token from the shell
```
header='{"typ":"JWT","alg":"HS256"}'
payload='{"sub":"john.doe@foo.local","iat":'$(date +%s)'}'
headerAndPayload=$( echo -n $header | base64 -w 0 | sed 's/[-= ]*$//' ).$( echo -n $payload | base64 -w 0 | sed 's/[-= ]*$//' )
key=mysupersecurekey
signature=$( echo -n $headerAndPayload | openssl sha256 -hmac $key -binary | base64 -w 0 | sed 's/[-= ]*$//' )
token=${headerAndPayload}.${signature}
```
If you need a small example on how to achieve the same with Go take a look at [brianvoe/sjwt](https://github.com/brianvoe/sjwt).
If you look for a more complete solution: https://github.com/square/go-jose

## TLS / x509 Certificates
A TLS client can send a certificate (client certificate) to the server. To authenticate the client the server verifies that the certificate is signed from a certian CA.
The common name in the certificate can be used as username and the organizational units are often used as groups for authZ.

```
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=foobar'
```

for testing we use the same certificate as client certificate
```
curl -k --cert cert.pem --key key.pem https://localhost:8443/handle-tls
```

* [No, don't enable revocation checking](https://www.imperialviolet.org/2014/04/19/revchecking.html)

## OAuth2 / OpenID connect
* Authorization Code Flow
* Implicit
* Client Credential
* User Password
