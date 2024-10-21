package service

import (
	"context"
	"scraper-first/internal/apperror"
	"scraper-first/pkg/logging"
	"strings"
	"fmt"
	"github.com/tebeka/selenium"
)

type(
	Parser interface {
		AllMarks() error
		Parsing(markName string) error
	}

	ParserImpl struct {
		wd selenium.WebDriver
		client Repository
		logger *logging.Logger
	}
)

func NewParser(wd selenium.WebDriver, client Repository, logger *logging.Logger) Parser {
	return &ParserImpl{
		wd: wd,
		client: client,
		logger: logger,
	}
}


func (p *ParserImpl) AllMarks() error {

	if err := p.wd.Get("https://ац-авто-клик.рф/"); err != nil {
		p.logger.Fatalf("Error opening the webpage: %v", err)
		return apperror.ErrorOpeningSitePage
	}

	marks, err := p.wd.FindElements(selenium.ByCSSSelector, ".list__tab--mark .list__tab-name")
	if err != nil {
		return apperror.ErrorElementHTMLSearch
	}

	var markName string
	var markNames []string

	for _, mark := range marks {
		markName, err = mark.Text()
		if err != nil {
			return apperror.ErrorRetrievingText
		}
		markName = strings.ToLower(markName)

		if strings.TrimSpace(markName) != "" {
			markNames = append(markNames, markName)
		}
	}

	for _, markName := range markNames {
		if err := p.Parsing(markName); err != nil {
			return err
		}
	}

	return nil
}

func (p *ParserImpl) Parsing(markName string) error {

	for i := 1; i < 100; i++ {
		if err := p.wd.Get(fmt.Sprintf("https://ац-авто-клик.рф/used/%s?page=%d", markName, i)); err != nil {
			p.logger.Fatalf("error opening the webpage: %v", err)
			return apperror.EndOfCatalogue
		}

		marks, err := p.wd.FindElements(selenium.ByCSSSelector, ".catalog__offer-name-mark")
		if err != nil {
			p.logger.Error(fmt.Sprintf("failed to find elements: %v", err))
			return apperror.ErrorElementHTMLSearch
		}

		models, err := p.wd.FindElements(selenium.ByCSSSelector, ".catalog__offer-name-folder")
		if err != nil {
			p.logger.Error(fmt.Sprintf("failed to find elements: %v", err))
			return apperror.ErrorElementHTMLSearch
		}

		volumes, err := p.wd.FindElements(selenium.ByCSSSelector, ".catalog__offer-name-volume")
		if err != nil {
			p.logger.Error(fmt.Sprintf("failed to find elements: %v", err))
			return apperror.ErrorElementHTMLSearch
		}

		prices, err := p.wd.FindElements(selenium.ByCSSSelector, ".catalog__offer-price")
		if err != nil {
			p.logger.Error(fmt.Sprintf("failed to find elements: %v", err))
			return apperror.ErrorElementHTMLSearch
		}

		items, err := p.wd.FindElements(selenium.ByCSSSelector,".catalog__offer-tech-item")
		if err != nil {
			p.logger.Error(fmt.Sprintf("failed to find elements: %v", err))
			return apperror.ErrorElementHTMLSearch
		}

		if len(models) == 0 || len(marks) == 0 || len(volumes) == 0 || len(prices) == 0 {
			return nil
		}

		for i := 0; i < len(marks); i++ {
			mark, err := marks[i].Text()
			if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			model, err := models[i].Text()
			if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			volume, err := volumes[i].Text()
			if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			price, err := prices[i].Text()
			if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			year, err := items[i*7].Text(); if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			km, err := items[i*7+1].Text(); if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			power, err := items[i*7+2].Text(); if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			transmission, err := items[i*7+3].Text(); if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			fuel, err := items[i*7+4].Text(); if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			owners, err := items[i*7+5].Text(); if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			drive, err := items[i*7+6].Text(); if err != nil {
				p.logger.Error(err)
				return apperror.ErrorRetrievingText
			}

			entity := CreateEntity{
				Mark: strings.ToLower(mark), 
				Model: strings.ToLower(model), 
				Volume: strings.ToLower(volume), 
				Price: strings.ToLower(price), 
				Year: strings.ToLower(year), 
				Kilometers: strings.ToLower(km), 
				Power: strings.ToLower(power), 
				Transmission: strings.ToLower(transmission), 
				Fuel: strings.ToLower(fuel), 
				Owners: strings.ToLower(owners), 
				Drive: strings.ToLower(drive),
				}

			ch, err := p.client.Check(context.Background(), &entity)
			if err != nil {
				p.logger.Error(err)
				return apperror.ErrorDB
			}

			if !ch {
				err = p.client.Create(context.Background(), &entity)

				if err != nil {
					p.logger.Error(err)
					return apperror.ErrorCreationObjectDB
				}
			}

			
		}
	}
	return nil
}