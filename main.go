package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/havuz/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kr/pretty"
	"golang.org/x/crypto/nacl/sign"
)

var (
	pubKey [32]byte

	sheetURL = os.Getenv("SHEET_URL")
	tcURL    = os.Getenv("TC_URL")

	isOnCloudfoundry = flag.Bool("cloudfoundry", false, "indicates that the runtime is deployed on cloudfoundry")

	errUnauthorized = &httpError{http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)}
)

func init() {
	flag.Parse()

	strPubKey := os.Getenv("PUB_KEY")
	pv, _ := base64.StdEncoding.DecodeString(strPubKey)
	copy(pubKey[:], pv)
}

func main() {
	ADDR := ":8090"
	if portFromEnv := os.Getenv("PORT"); portFromEnv != "" {
		ADDR = ":" + portFromEnv
	}

	http.Handle("/", http.HandlerFunc(handler))

	log.Printf("Reverse Proxy is listening at: %+q\n", ADDR)
	log.Fatal(http.ListenAndServe(ADDR, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var finalDecision bool

	pretty.Logln("")

	var remoteIP string
	{
		xForwardedFor := strings.Split(
			strings.ReplaceAll(r.Header.Get("X-Forwarded-For"), " ", ""),
			",",
		)

		switch {
		case *isOnCloudfoundry:
			remoteIP = xForwardedFor[len(xForwardedFor)-2]
		default:
			remoteIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
	}

	pretty.Logln("IP:", remoteIP)

	defer func() {
		pretty.Logln("DECISION:", finalDecision)

		if err := recover(); err != nil && err != http.ErrAbortHandler {
			httpErr, ok := err.(*httpError)
			if !ok {
				httpErr = &httpError{Message: fmt.Sprint(err)}
			}

			http.Error(w, httpErr.Error(), httpErr.Code)
		}
	}()

	// obtain the encrypted message digest from
	// auth header and base64 decode it into raw bytes
	var digest []byte
	{
		_, rawDigest, ok := r.BasicAuth()
		if !ok {
			panic(errUnauthorized)
		}

		var err error
		if digest, err = base64.RawURLEncoding.DecodeString(rawDigest); err != nil {
			panic(errUnauthorized)
		}

		pretty.Logln("DIGEST:", rawDigest)
	}

	// verify and decrypt the signed digest and
	// extract user id into UID
	var UID string
	{
		uid, ok := sign.Open(nil, digest, &pubKey)
		if !ok {
			panic(errUnauthorized)
		}

		UID = string(uid)

		pretty.Logln("UID:", UID)
	}

	// use UID to match a user in sheet
	// and validate them
	var user *types.User
	{
		resp, err := http.Get(sheetURL)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var allUsers []*types.User
		if err := gocsv.Unmarshal(resp.Body, &allUsers); err != nil {
			panic(err)
		}

		// traverse all users in sheet
		for _, u := range allUsers {
			if u.ID == UID {
				user = u
				break
			}
		}

		pretty.Logln("USER:", user)

		if user == nil {
			panic(errUnauthorized)
		}

		serUser, _ := json.Marshal(user)
		w.Header().Set("X-User", string(serUser))

		var isAuthorizedIP bool
		{
			IPs := strings.Split(
				strings.ReplaceAll(user.AllowedIPs, " ", ""),
				",",
			)

			for _, IP := range IPs {
				if IP == "*" || IP == remoteIP {
					isAuthorizedIP = true
					break
				}
			}
		}

		var isActive = user.Status == "ACTIVE"
		var isExpired bool
		{
			expiredAt, _ := time.Parse("2006-01-02", user.ExpiredAt)
			isExpired = time.Now().UTC().After(expiredAt)
		}

		switch {
		case isExpired, !isActive, !isAuthorizedIP:
			panic(errUnauthorized)
		}

		finalDecision = true
	}

	proxyTC(w, r)
}

func proxyTC(w http.ResponseWriter, r *http.Request) {
	var hopByHop = [...]string{
		"CF-Ray", "Expect-CT", "Set-Cookie", "Via", "Server",
	}

	u, err := url.Parse(tcURL)
	if err != nil {
		panic(err)
	}

	r.Host = u.Host
	r.URL.Path = "/"

	if tcAuth := u.User; tcAuth != nil {
		usern := tcAuth.Username()
		passw, _ := tcAuth.Password()
		r.SetBasicAuth(usern, passw)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ModifyResponse = func(resp *http.Response) error {
		for _, key := range hopByHop {
			resp.Header.Del(key)
		}

		return nil
	}

	proxy.ServeHTTP(w, r)
}

type httpError struct {
	Code    int
	Message string
}

func (herr *httpError) Error() string {
	if herr.Code == 0 {
		herr.Code = http.StatusInternalServerError
	}

	if herr.Message == "" {
		herr.Message = http.StatusText(herr.Code)
	}

	return herr.Message
}
