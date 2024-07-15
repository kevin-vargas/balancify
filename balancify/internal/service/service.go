package service

import (
	"balancify/internal/server"
	"balancify/internal/service/email"
	"balancify/internal/transaction"
	"bytes"
	"context"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"net/smtp"
	"strconv"
	"text/template"
	"time"
)

//go:embed resources/template.html
var EMAIL_TEMPLATE string

type service struct {
	manager     transaction.Manager
	smtpAddress string
}

func (s *service) ProcessTrxs(ctx context.Context, receiver string, trxs []transaction.Trx) error {
	err := s.manager.Insert(receiver, trxs)
	if err != nil {
		return err
	}
	return nil
}

type DataResults struct {
	Total     float64
	AvgDebt   float64
	AvgCredit float64
	Periods   []transaction.Period
}

func (s *service) fetchDataConcurrently(e string) (DataResults, error) {
	totalCh := make(chan float64)
	avgDebtCh := make(chan float64)
	avgCreditCh := make(chan float64)
	periodsCh := make(chan []transaction.Period)
	errCh := make(chan error, 4)

	go func() {
		total, err := s.manager.Total(e)
		if err != nil {
			errCh <- err
			return
		}
		totalCh <- total
	}()

	go func() {
		avgDebt, err := s.manager.AvgDebt(e)
		if err != nil {
			errCh <- err
			return
		}
		avgDebtCh <- avgDebt
	}()

	go func() {
		avgCredit, err := s.manager.AvgCredit(e)
		if err != nil {
			errCh <- err
			return
		}
		avgCreditCh <- avgCredit
	}()

	go func() {
		periods, err := s.manager.CountByPeriod(e)
		if err != nil {
			errCh <- err
			return
		}
		periodsCh <- periods
	}()

	var results DataResults

	for i := 0; i < 4; i++ {
		select {
		case total := <-totalCh:
			results.Total = total
		case avgDebt := <-avgDebtCh:
			results.AvgDebt = avgDebt
		case avgCredit := <-avgCreditCh:
			results.AvgCredit = avgCredit
		case periods := <-periodsCh:
			results.Periods = periods
		case err := <-errCh:
			return results, err
		}
	}
	return results, nil
}

func (s *service) SendEmail(ctx context.Context, from string, receiver string) error {
	const (
		subject       = "Balance"
		template_name = "mailTemplate"
	)
	to := []string{
		receiver,
	}
	type TemplateData struct {
		Title       string
		LogoURL     string
		CompanyName string
		Website     string
		Receiver    string
		DataResults
	}
	data := TemplateData{
		LogoURL:     "https://play-lh.googleusercontent.com/oXTAgpljdbV5LuAOt1NP9_JafUZe9BNl7pwQ01ndl4blYL4N4IQh4-n456P5l_hc1A",
		CompanyName: "Balancify",
		Website:     "balancify.fast.ar",
		Title:       "Balance",
		Receiver:    receiver,
	}

	res, err := s.fetchDataConcurrently(receiver)
	if err != nil {
		return err
	}
	data.DataResults = res

	buf := new(bytes.Buffer)
	t, err := template.New(template_name).Parse(EMAIL_TEMPLATE)
	if err != nil {
		return err
	}

	t.Execute(buf, data)

	msg := email.Format(
		email.WithFrom(from),
		email.WithReceiver(receiver),
		email.WithSubject(subject),
		email.WithHtmlBody(buf.Bytes()),
	)

	err = smtp.SendMail(s.smtpAddress, nil, from, to, msg)
	if err != nil {
		return fmt.Errorf("error al enviar email, %v", err)
	}
	return nil
}

func (s *service) ParseCSV(ctx context.Context, r io.Reader) ([]transaction.Trx, error) {
	reader := csv.NewReader(r)
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	trxs := make([]transaction.Trx, len(records))
	for i, record := range records {
		date, err := time.Parse("1/2/2006", record[0])
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		trxs[i] = transaction.Trx{
			Date:   date,
			Amount: amount,
		}
	}
	return trxs, nil
}

func New(a string, m transaction.Manager) server.Service {
	return &service{
		manager:     m,
		smtpAddress: a,
	}
}
