package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GenerateOTP(ctx context.Context, email string) error {
	min := big.NewInt(100000)
	max := big.NewInt(999999)

	range_ := new(big.Int).Sub(max, min)
	range_ = range_.Add(range_, big.NewInt(1))

	n, err := rand.Int(rand.Reader, range_)
	if err != nil {
		return fmt.Errorf("Failed to generate random number: %v", err)
	}

	n = n.Add(n, min)
	otp := n.String()

	err = RedisClient.Set(ctx, email, otp, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("Failed to store OTP: %v", err)
	}

	htmlBody := `
	<h2>Your OTP for DEVSOC Registration</h2>
	<p>Your One Time Password is: <strong>%s</strong></p>
	<p>This OTP will expire in 5 minutes.</p>
	<p>If you did not request this OTP, please ignore this email.</p>
	`

	err = SendEmail(email, "DEVSOC Registration OTP", fmt.Sprintf(htmlBody, otp))
	if err != nil {

		RedisClient.Del(ctx, email)
		return fmt.Errorf("Failed to send OTP email: %v", err)
	}

	return nil
}
