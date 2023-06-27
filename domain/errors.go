package domain

import "errors"

var ErrResourceNotFound = errors.New("resource not found")

var ErrUserNotFound = errors.New("user not found")

var ErrCouponNotFound = errors.New("coupon not found")

var ErrMenuNotFound = errors.New("menu not found")

var ErrOrderNotFound = errors.New("order not found")

var ErrCartItemNotFound = errors.New("cart item not found")

var ErrInternalServer = errors.New("internal server error")

var ErrUpdateCart = errors.New("error update cart")

var ErrUploadImage = errors.New("error upload image")

var ErrDeleteImage = errors.New("error delete image")

var ErrUnsyncedData = errors.New("data is not synced")

var ErrInvalidIdentifier = errors.New("invalid identifier, please insert valid username/email/phone number")

var ErrInvalidEmail = errors.New("invalid email")

var ErrInvalidEmailFormat = errors.New("invalid email format")

var ErrInvalidUsername = errors.New("invalid username")

var ErrInvalidUsernameFormat = errors.New("invalid username format")

var ErrInvalidPhone = errors.New("invalid phone number")

var ErrInvalidPhoneFormat = errors.New("invalid phone number format")

var ErrDuplicateEmail = errors.New("duplicate email")

var ErrDuplicateUsername = errors.New("duplicate username")

var ErrDuplicatePhone = errors.New("duplicate phone number")

var ErrInvalidPassword = errors.New("invalid password")

var ErrInvalidToken = errors.New("invalid token")

var ErrInvalidCategory = errors.New("invalid menu category")

var ErrInvalidDeliveryStatus = errors.New("invalid delivery status")

var ErrUpdateMenu = errors.New("error update menu")

var ErrForbiddenAccess = errors.New("forbidden access")

var ErrInvalidRequest = errors.New("invalid request")

var ErrUnauthorized = errors.New("unauthorized access")

var ErrInvalidBody = errors.New("invalid body request")

var ErrInvalidParams = errors.New("invalid params")

var ErrMalformedRequest = errors.New("malformed request")

var ErrInvalidQuery = errors.New("invalid query")

var ErrHasNotOrdered = errors.New("you have not ordered yet")

var ErrDuplicateReview = errors.New("you have already reviewed this menu")

var ErrNoGamesAttempt = errors.New("you have no games attempt")

var ErrGameAlreadyAnswered = errors.New("you have already answered this game")

var ErrGameNotFound = errors.New("game not found")

var ErrPaymentOptionNotFound = errors.New("payment option not found")

var ErrPromotionNotFound = errors.New("promotion not found")

var ErrUserCouponNotFound = errors.New("user coupon not found")
