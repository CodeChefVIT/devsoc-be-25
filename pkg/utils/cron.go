package utils

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/robfig/cron"
)

func Cron() {
	ctx := context.Background()

	users, err := Queries.GetUsers(ctx)
	if err != nil {
		logger.Errorf("DB error: %v", err)
	}

	file, err := os.Create("users.csv")
	if err != nil {
		logger.Errorf("error creating file", err)
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)
	headers := []string{"ID", "FirstName", "LastName", "Email", "PhoneNo", "Gender", "RegNo", "TeamID", "VitEmail", "Hostel", "RoomNo", "GitHub", "Role", "IsLeader", "IsVerified", "IsBanned", "IsProfComplete"}

	if err := csvWriter.Write(headers); err != nil {
		logger.Errorf("failed to write headers", err)
		return
	}

	for _, user := range users {
		record := []string{
			user.ID.String(),
			user.FirstName,
			user.LastName,
			user.Email,
			user.PhoneNo.String,
			user.Gender,
			*user.RegNo,
			user.TeamID.UUID.String(),
			*user.VitEmail,
			user.HostelBlock,
			strconv.Itoa(int(user.RoomNo)),
			user.GithubProfile,
			user.Role,
			strconv.FormatBool(user.IsLeader),
			strconv.FormatBool(user.IsVerified),
			strconv.FormatBool(user.IsBanned),
			strconv.FormatBool(user.IsProfileComplete),
		}

		if err := csvWriter.Write(record); err != nil {
			logger.Errorf("failed to write data", err)
			return
		}
	}

	csvWriter.Flush()

	subject := fmt.Sprintf("Users Details")
	body := fmt.Sprintf("This mail contails details of all the users")

	recipients  := Config.Recipients

	if recipients == "" {
		logger.Errorf("error in getenv")
	}

	ToSendMail := strings.Split(recipients, ",")

	cr := cron.New()
	err = cr.AddFunc("@daily", func() {
		logger.Infof("Mail Sent")
		for _, user := range ToSendMail {
			SendEmail(user, subject, body, "users.csv")
		}
	})

	if err != nil {
		logger.Errorf("Error in cron", err)
	}
	cr.Start()
	logger.Infof("Cron started")
}
