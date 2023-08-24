package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/olekukonko/tablewriter"
	"github.com/toricls/acos"
)

type CompareTo string

const (
	Yesterday CompareTo = "YESTERDAY"
	LastWeek  CompareTo = "LAST_WEEK"
)

func handler(ctx context.Context, event interface{}) ([]byte, error) {
	var err error
	var accounts acos.Accounts

	ouId := os.Getenv("OU_ID")
	if ouId == "" {
		// Try to get ROOT_OU_ID if OU_ID is not specified
		ouId = os.Getenv("ROOT_OU_ID")
	}
	if ouId != "" {
		// Get all AWS accounts in the specified OU
		if accounts, err = acos.ListAccountsByOu(ctx, ouId); err != nil {
			return nil, err
		}
	} else {
		// Get all AWS accounts in the AWS organization
		if accounts, err = acos.ListAccounts(ctx); err != nil {
			return nil, err
		}
	}

	compareTo := Yesterday
	compareToStr := os.Getenv("COMPARE_TO")
	switch compareToStr {
	case string(LastWeek):
		compareTo = LastWeek
	case string(Yesterday):
	default:
		break
	}

	asOf, err := time.Parse("2006-01-02", os.Getenv("AS_OF"))
	if err != nil {
		asOf = time.Now().UTC()
	}

	headerText := os.Getenv("HEADER_TEXT")
	if len(headerText) > 0 {
		headerText = headerText + "\n"
	}
	footerText := os.Getenv("FOOTER_TEXT")
	if len(footerText) > 0 {
		footerText = "\n" + footerText
	}

	var costs acos.Costs
	opt := acos.NewGetCostsOption(asOf)
	if costs, err = acos.GetCosts(ctx, accounts, opt); err != nil {
		return nil, err
	}

	res := print(&costs, asOf, compareTo)

	payload := slack.Payload{
		Text:      headerText + "```" + res + "```" + footerText,
		Username:  "I'm Tori's monkey",
		IconEmoji: ":monkey:",
	}
	errs := slack.Send(os.Getenv("SLACK_WEBHOOK"), "", payload)
	if len(errs) > 0 {
		fmt.Printf("error: %s\n", errs)
	}

	return nil, nil
}

func print(costs *acos.Costs, asOf time.Time, compareTo CompareTo) string {
	// Sort map keys by AWS Account ID
	keys := make([]string, 0, len(*costs))
	for k := range *costs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tableString := &strings.Builder{}
	t := tablewriter.NewWriter(tableString)
	t.SetHeader(getHeader(compareTo))
	t.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
	totalThisMonth, totalIncrease, totalLastMonth := 0.0, 0.0, 0.0
	for _, k := range keys {
		c := (*costs)[k]
		thisMonth := fmt.Sprintf("%f", c.AmountThisMonth)
		incr := getIncrease(c, compareTo)
		incrStr := fmt.Sprintf("%s %f", getAmountPrefix(incr), incr)
		lastMonth := fmt.Sprintf("%f", c.AmountLastMonth)
		t.Append([]string{c.AccountID, c.AccountName, thisMonth, incrStr, lastMonth})
		totalThisMonth += c.AmountThisMonth
		totalIncrease += incr
		totalLastMonth += c.AmountLastMonth
	}
	t.SetFooter([]string{"", "Total", fmt.Sprintf("%f", totalThisMonth), fmt.Sprintf("%s %f", getAmountPrefix(totalIncrease), totalIncrease), fmt.Sprintf("%f", totalLastMonth)})
	t.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	t.SetCaption(true, fmt.Sprintf("As of %s.", asOf.Format("2006-01-02")))
	t.Render()

	return tableString.String()
}

func getHeader(compareTo CompareTo) []string {
	if compareTo == LastWeek {
		return []string{"Account ID", "Account Name", "This Month ($)", "vs Last Week ($)", "Last Month ($)"}
	} else {
		return []string{"Account ID", "Account Name", "This Month ($)", "vs Yesterday ($)", "Last Month ($)"}
	}
}

func getIncrease(c acos.Cost, compareTo CompareTo) float64 {
	if compareTo == LastWeek {
		return c.LatestWeeklyCostIncrease
	} else {
		return c.LatestDailyCostIncrease
	}
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
