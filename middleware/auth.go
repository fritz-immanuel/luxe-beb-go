package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"luxe-beb-go/configs"
	"luxe-beb-go/library"
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

func Auth(c *gin.Context) {
	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}

	CheckIPClientIP(c, config)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})
	if err != nil {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token Invalid"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}
	claimJWT, ok := library.GetJWTClaims(c, tokenString)
	if !ok {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token Invalid"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	val, errRedis := redisClient.Get(tokenString).Result()
	if errRedis != nil {
		log.Printf(`
		======================================================================
		Error Collecting Caching in "Auth":
		Error: %v
		======================================================================
		`, errRedis)
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token is Expired"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}
	if val == "" {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token is Expired"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	c.Set("SessionID", token)
	c.Set("KitchenTypeID", claimJWT["KitchenTypeID"])
	c.Set("BusinessID", claimJWT["BusinessID"])
	c.Set("SupervisorUserID", claimJWT["SupervisorUserID"])
	c.Set("UserID", claimJWT["ID"])
	c.Set("Email", claimJWT["Email"])
	c.Set("Type", claimJWT["Type"])

	if errRedis := redisClient.Set(
		tokenString,
		fmt.Sprintf("{\"id\":%s}", claimJWT["ID"]),
		time.Second*time.Duration(config.RedisTimeOut),
	).Err(); errRedis != nil {
		log.Printf(`
		======================================================================
		Error Storing Caching in "Auth":
		Error: %v,
		======================================================================
		`, err)
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token is Expired"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	// check hak akses
	route := c.Request.RequestURI
	routeIndex := strings.Index(route, "?")
	var fixRoute string

	if routeIndex == -1 {
		fixRoute = route
	} else {
		fixRoute = string([]rune(route)[0:routeIndex])
	}

	db, err := sqlx.Open("mysql", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT
	user_permission.permission_id,
	permission.package AS permission_package,
	permission.name AS permission_name,
	permission.action AS permission_action,
	permission.type AS permission_type,
	permission.route AS permission_route
	FROM user_permission
	JOIN permission on permission.id = user_permission.permission_id
	WHERE package = 'Website' AND user_permission.user_id = ?
	`, claimJWT["ID"])
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	permissionObjArr := []models.Permission{}

	method := c.Request.Method
	hasAccess := false
	for rows.Next() {
		var (
			id           int64
			package_name string
			name         string
			action       string
			typeMethod   string
			route        string
		)
		if err := rows.Scan(&id, &package_name, &name, &action, &typeMethod, &route); err != nil {
			log.Fatal(err)
		}

		data := models.Permission{
			ID:      uint(id),
			Package: package_name,
			Name:    name,
			Action:  action,
			Type:    typeMethod,
			Route:   route,
		}
		permissionObjArr = append(permissionObjArr, data)

		checkRoute := true
		arrRoutes := strings.Split(data.Route, "/")

		arrFixRoutes := strings.Split(fixRoute, "/")

		if len(arrRoutes) == len(arrFixRoutes) {
			for i, arrRoute := range arrRoutes {

				if arrRoute != arrFixRoutes[i] && !strings.HasPrefix(arrRoute, ":") {

					checkRoute = false
					break
				}
			}
		} else {
			checkRoute = false
		}

		if checkRoute && data.Type == method {
			hasAccess = true
			break
		}
	}

	if hasAccess == false {
		response := types.Result{Status: "Warning", StatusCode: http.StatusForbidden, Message: "No Permission Access"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}
}

func AuthPOS(c *gin.Context) {
	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}

	CheckIPClientIP(c, config)

	CheckApplicationVersionPOS(c)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secretmobile"), nil
	})
	if err != nil {

		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token Invalid"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}
	claimJWT, ok := library.GetJWTMobileClaims(c, tokenString)
	if !ok {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token Invalid"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	val, errRedis := redisClient.Get(tokenString).Result()
	if errRedis != nil {
		log.Printf(`
		======================================================================
		Error Collecting Caching in "AuthMobile":
		Error: %v
		======================================================================
		`, errRedis)
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token is Expired"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	if val == "" {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token is Expired"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	c.Set("SessionID", token)
	c.Set("BusinessID", claimJWT["BusinessID"])
	c.Set("Type", claimJWT["Type"])
	c.Set("UserID", claimJWT["ID"])
	c.Set("Email", claimJWT["Email"])
	c.Set("BusinessShiftID", claimJWT["BusinessShiftID"])
	c.Set("IsCaptain", claimJWT["IsCaptain"])
	c.Set("IsDisabledChangeBusinessPOS", claimJWT["IsDisabledChangeBusinessPOS"])

	if errRedis := redisClient.Set(
		tokenString,
		fmt.Sprintf("{\"id\":%s}", claimJWT["ID"]),
		time.Second*time.Duration(config.RedisTimeOut),
	).Err(); errRedis != nil {
		log.Printf(`
		======================================================================
		Error Storing Caching in "Auth":
		Error: %v,
		======================================================================
		`, err)
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token is Expired"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	if (claimJWT["BusinessShiftID"]) != 0.0 {
		db, err := sqlx.Open("mysql", config.DBConnectionString)
		if err != nil {
			log.Fatalln("failed to open database x: ", err)
		}
		defer db.Close()

		rows, err := db.Query(`SELECT open_shift.*
		FROM open_shift
		INNER JOIN business_shift on business_shift.id = open_shift.shift_id
		WHERE open_shift.business_id = ? AND DATE(open_shift.created_at) = DATE(UTC_TIMESTAMP + INTERVAL 7 hour) 
			AND open_shift.shift_id = ? AND open_shift.closed_at IS NULL`, claimJWT["BusinessID"], claimJWT["BusinessShiftID"])
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		if !rows.Next() {
			rowsShift, err := db.Query(`SELECT id
			FROM business_shift bs
			WHERE bs.business_id = ? AND bs.deleted_at IS NULL
			ORDER BY CAST(bs.end_hour AS TIME) DESC`, claimJWT["BusinessID"])
			if err != nil {
				log.Fatal(err)
			}
			defer rowsShift.Close()

			rowsShift.Next()

			var lastID int
			rowsShift.Scan(lastID)

			if lastID == claimJWT["BusinessShiftID"] {
				rowsOrder, err := db.Query(`SELECT id
				FROM orders o
				WHERE business_id = ? AND status_id IN (1,2) AND order_type_id = 2`, claimJWT["BusinessID"])
				if err != nil {
					log.Fatal(err)
				}
				defer rowsOrder.Close()

				if rowsOrder.Next() {
					response := types.Result{Status: "Warning", StatusCode: http.StatusForbidden, Message: "Shift sudah habis! Masih ada order yang aktif"}
					result := gin.H{
						"result": response,
					}
					c.JSON(http.StatusTeapot, result)
					c.Abort()
					return
				}
			}

			response := types.Result{Status: "Warning", StatusCode: http.StatusForbidden, Message: "Shift sudah berubah"}
			result := gin.H{
				"result": response,
			}
			c.JSON(http.StatusForbidden, result)
			c.Abort()
			return
		}
	}
}

func AuthExternal(c *gin.Context) {
	CheckSecretTokenWebApp(c)

	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}

	CheckIPClientIP(c, config)

	var token string
	tokenString := c.Request.Header.Get("Authorization")
	_, err = fmt.Sscanf(tokenString, "Bearer %s", &token)
	if err != nil {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token Format Wrong"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	db, err := sqlx.Open("mysql", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT
	api_client.id,
	api_client.name
	FROM api_client
	WHERE api_client.token = ? and name = 'Account'
	`, token)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	hasAccess := false
	if rows.Next() {
		hasAccess = true
	}

	if hasAccess == false {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Token Not Found"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}
}

func AuthCheckIP(c *gin.Context) {
	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}

	CheckIPClientIP(c, config)
}

func CheckApplicationVersionPOS(c *gin.Context) {
	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}
	if err != nil {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "failed to get configuration"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	// Version format xx.xx.xx (Major.Minor.Bugfix)
	minimumVersionStr := "1.0.0"
	minimumVersion := strings.Split(minimumVersionStr, ".")

	requestAppVersionStr := c.Request.Header.Get("AppVersion")
	requestAppVersion := strings.Split(requestAppVersionStr, ".")

	requestAndroidVersionStr := c.Request.Header.Get("AndroidVersion")
	requestAndroidVersion := strings.Split(requestAndroidVersionStr, ".")

	requestIOSVersionStr := c.Request.Header.Get("IOSVersion")
	requestIOSVersion := strings.Split(requestIOSVersionStr, ".")

	if strings.Compare(requestAndroidVersionStr, "") == 0 && strings.Compare(requestIOSVersionStr, "") == 0 {
		//response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Application Need To Be Updated. Please Update your application on playstore"}
		result := gin.H{
			"code":    "warning",
			"Status":  "Warning",
			"Message": "Application Need To Be Updated. Please Update your application on Playstore/ App Store",
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	if strings.Compare(requestAndroidVersionStr, "") != 0 {
		requestAppVersion = requestAndroidVersion

		minimumVersionStr = config.AndroidPOSAppMinimumVersion
		minimumVersion = strings.Split(minimumVersionStr, ".")
	}

	if strings.Compare(requestIOSVersionStr, "") != 0 {
		requestAppVersion = requestIOSVersion

		minimumVersionStr = config.IosPOSAppMinimumVersion
		minimumVersion = strings.Split(minimumVersionStr, ".")
	}

	for i := 0; i < len(minimumVersion); i++ {
		minimumVer, errConversion := strconv.Atoi(minimumVersion[i])
		if errConversion != nil {
			response := types.Result{Status: "Warning", StatusCode: http.StatusInternalServerError, Message: "Server Mobile App Minimum Version String Conversion Error"}
			result := gin.H{
				"result": response,
			}
			c.JSON(http.StatusUnauthorized, result)
			c.Abort()
			return
		}

		requestVer, errConversion := strconv.Atoi(requestAppVersion[i])
		if errConversion != nil {
			response := types.Result{Status: "Warning", StatusCode: http.StatusInternalServerError, Message: "Request Mobile App Version String Conversion Error"}
			result := gin.H{
				"result": response,
			}
			c.JSON(http.StatusUnauthorized, result)
			c.Abort()
			return
		}

		if requestVer < minimumVer {
			response := types.Result{Status: "Warning", StatusCode: http.StatusUpgradeRequired, Message: "Application Need To Be Updated. Please Update your application on Playstore/ App Store"}
			result := gin.H{
				"result": response,
			}
			c.JSON(http.StatusUpgradeRequired, result)
			c.Abort()
			return

		} else if requestVer > minimumVer {
			break
		}
	}
}

func CheckSecretTokenWebApp(c *gin.Context) {
	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}
	if err != nil {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "failed to get configuration"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	// CHECK SECRET TOKEN
	var secretToken string
	secretTokenString := c.Request.Header.Get("Access-Token")
	_, err = fmt.Sscanf(secretTokenString, "Bearer %s", &secretToken)
	if err != nil {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Access Token Format Wrong"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	db, err := sqlx.Open("mysql", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT
	api_client.id,
	api_client.name
	FROM api_client
	WHERE api_client.token = ? and name = 'External'
	`, secretToken)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	hasAccess := false
	if rows.Next() {
		hasAccess = true
	}

	if hasAccess == false {
		response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Access Token Not Found"}
		result := gin.H{
			"result": response,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}
}

func CheckIPClientIP(c *gin.Context, config *configs.Config) {
	clientIP := c.ClientIP()

	if clientIP != "::1" {
		clientIPSplit := strings.Split(clientIP, ".")

		whitelistSplit := strings.Split(config.WhitelistedIps, ",")

		var first []string
		var second []string
		var third []string
		var fourth []string

		counter := 0

		for _, wl := range whitelistSplit {
			wl = strings.TrimSpace(wl)
			splited := strings.Split(wl, ".")

			first = append(first, splited[0])
			second = append(second, splited[1])
			third = append(third, splited[2])
			fourth = append(fourth, splited[3])

			if clientIPSplit[0] == splited[0] || splited[0] == "0" {
				counter += 1
			}

			if clientIPSplit[1] == splited[1] || splited[1] == "0" {
				counter += 1
			}

			if clientIPSplit[2] == splited[2] || splited[2] == "0" {
				counter += 1
			}

			if clientIPSplit[3] == splited[3] || splited[3] == "0" {
				counter += 1
			}

			if counter == 4 {
				break
			}

			counter = 0
		}

		if counter != 4 {
			response := types.Result{Status: "Warning", StatusCode: http.StatusUnauthorized, Message: "Unauthorized Access"}
			result := gin.H{
				"result": response,
			}
			c.JSON(http.StatusUnauthorized, result)
			c.Abort()
			return
		}
	}
}
