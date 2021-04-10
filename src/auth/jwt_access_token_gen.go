package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/osins/osin-simple/simple/model/entity"
	"github.com/osins/osin-simple/simple/model/face"
)

var (
	privatekeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`)

	publickeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----`)
)

func NewJwt() (face.AccessToken, error) {
	var err error
	gen := &jwtAccessTokenGen{}
	if gen.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privatekeyPEM); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	if gen.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publickeyPEM); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	return gen, nil
}

type jwtAccessTokenGen struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func (c *jwtAccessTokenGen) GenerateAccessToken(data face.Access, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	// generate JWT access token

	var claims jwt.Claims
	fmt.Printf("user data: %v\n", data.GetUser())
	if data.GetUser() != nil && structs.IsStruct(data.GetUser()) {
		u := structs.Map(data.GetUser())
		claims = jwt.MapClaims{
			"cid":    data.GetClient().GetId(),
			"uid":    u["Id"],
			"name":   u["Username"],
			"email":  u["EMail"],
			"mobile": u["Mobile"],
			"exp":    data.ExpireAt().Unix(),
		}
	} else {
		claims = jwt.MapClaims{
			"cid": data.GetClient().GetId(),
			"exp": data.ExpireAt().Unix(),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	accesstoken, err = token.SignedString(c.PrivateKey)
	if err != nil {
		return "", "", err
	}

	fmt.Printf("generaterefresh: %v\n", generaterefresh)
	if !generaterefresh {
		return
	}

	if data.GetUser() != nil {
		claims = jwt.MapClaims{
			"cid":    data.GetClient().GetId(),
			"uid":    data.GetUser().GetId(),
			"name":   data.GetUser().GetUsername(),
			"email":  data.GetUser().GetEmail(),
			"mobile": data.GetUser().GetMobile(),
		}
	} else {
		claims = jwt.MapClaims{
			"cid": data.GetClient().GetId(),
		}
	}

	// generate JWT refresh token
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	refreshtoken, err = token.SignedString(c.PrivateKey)
	if err != nil {
		return "", "", err
	}
	return
}

func (c *jwtAccessTokenGen) VerifyToken(accessToken string) (face.Access, error) {
	parts := strings.Split(accessToken, ".")
	header, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		fmt.Printf("token header info base64 to string faild: %v\n", parts[0])
		return nil, fmt.Errorf("token verify faild.")
	}

	fmt.Printf("token header: %s\n", header)

	var h map[string]string
	if err := json.Unmarshal([]byte(header), &h); err != nil {
		fmt.Printf("token type info parse faild: %v\n", parts[0])
		return nil, fmt.Errorf("token verify faild.")
	}

	if h["typ"] != "JWT" {
		fmt.Printf("token type error: %v\n", h["typ"])
		return nil, fmt.Errorf("token verify faild.")
	}

	fmt.Printf("token header map: %v\n", h)

	sigMethod := jwt.GetSigningMethod(h["alg"])
	if err := sigMethod.Verify(strings.Join(parts[0:2], "."), parts[2], c.PublicKey); err != nil {
		fmt.Printf("\n\ntoken verify faild, jwt parse faild, error: %s\n\nstrings.Join(parts[0:2], \".\"): %s\n\nparts[2]: %s\n\n", err, strings.Join(parts[0:2], "."), parts[2])
		return nil, fmt.Errorf("token verify faild.")
	}

	fmt.Printf("jwt verify type: %s, ok!\n", sigMethod.Alg())

	pl, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Printf("token payload info base64 to string faild: %v\nerror: %s\n\n", parts[1], err)
		return nil, fmt.Errorf("token verify faild.")
	}

	fmt.Printf("token payload: %s\n", pl)

	var claims jwt.MapClaims
	if err := json.Unmarshal([]byte(pl), &claims); err != nil {
		fmt.Printf("token claims: %v, parse error: %s\n", pl, err)
		return nil, fmt.Errorf("token verify faild.")
	}

	cid, ok := claims["cid"].(string)
	if !ok {
		fmt.Printf("token verify faild, cid get faild, claims: %v\n", claims)
		return nil, fmt.Errorf("token verify faild.")
	}

	access := &entity.Access{
		ClientId: cid,
		Client: &entity.Client{
			Id: cid,
		},
	}

	if exp, ok := claims["exp"].(string); ok {
		i64, err := strconv.ParseInt(exp, 10, 32)
		if err == nil {
			access.ExpiresIn = int32(i64)
		}
	}

	if uid, ok := claims["uid"].(string); ok {
		access.UserId = uid
		u := &entity.User{
			Id: uid,
		}
		if name, ok := claims["name"].(string); ok {
			u.Username = name
		}
		if email, ok := claims["email"].(string); ok {
			u.EMail = email
		}
		if mobile, ok := claims["mobile"].(string); ok {
			u.Mobile = mobile
		}

		access.User = u
	}

	return access, nil
}
