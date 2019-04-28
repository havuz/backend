## Havuz TC Reverse Proxy [![Codacy Badge](https://api.codacy.com/project/badge/Grade/6c8f37ee6d1b4e23982b13be8568c94c)](https://www.codacy.com/app/0xbkt/backend?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=havuz/backend&amp;utm_campaign=Badge_Grade)

An app that verifies Havuz users' license tokens and
reverse proxies requests to [Vekilio's Tunnel Collector](https://github.com/Vekilio/tunnel-collector).

---

### Configuration

#### Environment Variables

```json
// Google Sheets URL, exported as CSV. Must be public.
// Must have CSV headers corresponding to type User.
// Location headers are followed by HTTP client.
// Example: https://bit.ly/2IndTaa
"SHEET_URL": "",

// NaCl crypto_sign public key to verify and decrypt the digest.
// Must be encoded with base64std, not url variant.
"PUB_KEY": "",

// HTTP URL of TC for reverse proxying.
// User-Pass may be omitted. Only relative root dir ("/") will be hit.
// Example: https://usern:passw@collector.example.com
//          https://collector.another.com/another-path
"TC_URL": "",

// The TCP port to listen on. Default is 8090, as given below.
"PORT": "8090"
```

---

### Generating Keypair

NaCl's [crypto_sign](https://nacl.cr.yp.to/sign.html) utility is used for generating keys, creating license codes and verifying them on the server-side. I myself use [Step CLI](https://github.com/smallstep/cli) tool to generate crypto_sign keypairs and signing data (in this app, data = user ID). Verification is done by [x/crypto/nacl/sign](https://godoc.org/golang.org/x/crypto/nacl/sign) package of Go.

1. Generate a keypair.
   - `step crypto nacl sign keypair pub-file priv-file`
     - Remember to keep the private pair secure and offline.
       The public pair can be distributed anywhere.
2. Give the public key to app.
   - `cat pub-file | base64`
     - This will print the base64std encoding of `pub-file`.
       This final output can be set to `PUB_KEY` env var.
     - With the public key, the app will be able to verify and
       decrypt the digested license code and extract the ID to
       compare it against the Sheet.
3. Create a license code.
   - `step crypto nacl sign sign priv-file`
     - Enter the intended user ID as message. Afterwards, you will be given a message digest encoded in Base64.
       See the usage below.

---

### Authenticating to Reverse Proxy

Now that the client has their license code, either they or a mediator can access the reverse proxy by using [HTTP Basic Auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization).

```bash
$ curl -I -XGET --user "_:<DIGEST/LICENSE_CODE>" https://tc-reverse-proxy:8090
$ # curl will take care of Authorization header when used with `--user` flag.
$ # now we do it ourselves:
$ echo "_:<DIGEST/LICENSE_CODE>" | base64
Xzo8RElHRVNUL0xJQ0VOU0U+Cg==
$ curl -I -XGET -H "Authorization: Xzo8RElHRVNUL0xJQ0VOU0U+Cg==" https://tc-reverse-proxy:8090
```
