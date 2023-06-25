package middleware

import (
	"net/http"
	"net/url"

	"a21hc3NpZ25tZW50/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session_token")

		if cookie == "" {
			head := ctx.Request.Header.Get("Content-Type")

			if head == "application/json" {
				ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
				return
			}
			location := url.URL{Path: "/client/login"}
			ctx.Redirect(http.StatusSeeOther, location.RequestURI())
			return
		}

		// Ambil value dari cookie token
		tknStr := cookie

		// Deklarasi variable claims yang akan kita isi dengan data hasil parsing JWT
		claim := &model.Claims{}

		// parse JWT token ke dalam claims
		tkn, err := jwt.ParseWithClaims(tknStr, claim, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// return unauthorized ketika ada kesalahan saat parsing token
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
				return
			}
			// return bad request ketika field token tidak ada
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}

		//return unauthorized ketika token sudah tidak valid (biasanya karena token expired)
		if !tkn.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}

		ctx.Set("email", claim.Email)

		ctx.Next()
	})
}
