package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/olekukonko/tablewriter"
	"github.com/toricls/acos"
)

func handler(ctx context.Context, event interface{}) ([]byte, error) {
	var err error
	var accounts acos.Accounts

	ouId := os.Getenv("OU_ID")
	if ouId != "" {
		// Get all AWS accounts in the specified OU
		if accounts, err = acos.ListAccountsByOu(ctx, os.Getenv("OU_ID")); err != nil {
			return nil, err
		}
	} else {
		// Get all AWS accounts in the AWS organization
		if accounts, err = acos.ListAccounts(ctx); err != nil {
			return nil, err
		}
	}

	var costs acos.Costs
	opt := acos.AcosGetCostsOption{
		ExcludeCredit:  true,
		ExcludeUpfront: true,
		ExcludeRefund:  false,
		ExcludeSupport: false,
	}
	if costs, err = acos.GetCosts(ctx, accounts, opt); err != nil {
		return nil, err
	}
	res := print(&costs)

	payload := slack.Payload{
		Text:      "Here's our AWS bills:```\n" + res + "```",
		Username:  "I'm Tori's monkey",
		IconEmoji: ":monkey:",
	}
	errs := slack.Send(os.Getenv("SLACK_WEBHOOK"), "", payload)
	if len(errs) > 0 {
		fmt.Printf("error: %s\n", errs)
	}

	return nil, nil
}

func print(costs *acos.Costs) string {
	tableString := &strings.Builder{}
	t := tablewriter.NewWriter(tableString)
	t.SetHeader([]string{"Account ID", "Account Name", "This Month ($)", "vs Yesterday ($)", "Last Month ($)"})
	t.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
	totalThisMonth, totalYesterday, totalLastMonth := 0.0, 0.0, 0.0
	for _, c := range *costs {
		thisMonth := fmt.Sprintf("%f", c.AmountThisMonth)
		vsYesterday := fmt.Sprintf("%s %f", getAmountPrefix(c.LatestDailyCostIncrease), c.LatestDailyCostIncrease)
		lastMonth := fmt.Sprintf("%f", c.AmountLastMonth)
		t.Append([]string{c.AccountID, c.AccountName, thisMonth, vsYesterday, lastMonth})
		totalThisMonth += c.AmountThisMonth
		totalYesterday += c.LatestDailyCostIncrease
		totalLastMonth += c.AmountLastMonth
	}
	t.SetFooter([]string{"", "Total", fmt.Sprintf("%f", totalThisMonth), fmt.Sprintf("%s %f", getAmountPrefix(totalYesterday), totalYesterday), fmt.Sprintf("%f", totalLastMonth)})
	t.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	t.SetCaption(true, fmt.Sprintf("As of %s.", time.Now().Format("2006-01-02")))
	t.Render()

	return tableString.String()
}

func getAmountPrefix(amount float64) string {
	if amount > 0.0 {
		return "+"
	} else if amount < 0.0 {
		return "-"
	}
	return ""
}

func main() {
	lambda.Start(handler)
}
