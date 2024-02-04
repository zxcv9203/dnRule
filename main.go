package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func main() {
	dDayRegex := regexp.MustCompile(`D-(\d+)`)
	env, err := getEnvironment()

	if err != nil {
		log.Fatalf("토큰을 가져오는데 실패했습니다 [ %v ]", err)
	}

	ctx := context.Background()
	githubAccessToken := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env.token},
	)

	oauth2Token := oauth2.NewClient(ctx, githubAccessToken)
	client := github.NewClient(oauth2Token)

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		pullRequests, response, err := client.PullRequests.List(ctx, env.owner, env.repository, opts)
		if err != nil {
			log.Fatalf("PR 리스트를 조회하는데 실패했습니다. [ %v ]", err)
		}

		for _, pullRequest := range pullRequests {
			var newLabels []string
			var isChange = false

			for _, label := range pullRequest.Labels {
				dDayLabel := dDayRegex.FindStringSubmatch(*label.Name)
				if len(dDayLabel) == 2 {
					day, err := strconv.Atoi(dDayLabel[1])
					if err != nil {
						log.Printf("라벨 숫자 변환 실패 : %v", err)
						continue
					}
					if day > 0 {
						newLabels = append(newLabels, fmt.Sprintf("D-%d", day-1))
						isChange = true
					} else {
						newLabels = append(newLabels, *label.Name)
					}
				}
			}
			if isChange {
				_, _, err = client.Issues.ReplaceLabelsForIssue(ctx, env.owner, env.repository, *pullRequest.Number, newLabels)
				if err != nil {
					log.Printf("라벨 업데이트가 실패했습니다\nPR 번호 : %#v PR 이름 : %v\n[ %v ]", *pullRequest.Number, *pullRequest.Title, err)
				}
			}
		}
		if response.NextPage == 0 {
			break
		}
		opts.Page = response.NextPage
	}
}

func getEnvironment() (environment, error) {
	env := environment{}
	env.token = os.Getenv("GITHUB_TOKEN")
	if env.token == "" {
		return environment{}, fmt.Errorf("GITHUB_TOKEN이 존재하지 않습니다")
	}

	repositoryName := os.Getenv("GITHUB_REPOSITORY")
	if repositoryName == "" {
		return environment{}, fmt.Errorf("GITHUB_REPOSITORY가 존재하지 않습니다")
	}

	repoPart := strings.Split(repositoryName, "/")
	if len(repoPart) == 2 {
		env.owner = repoPart[0]
		env.repository = repoPart[1]
	} else {
		return environment{}, fmt.Errorf("GITHUB_REPOSITORY의 형식이 잘못되었습니다. 전달받은 값 : %s\n다음과 같은 형식을 사용해주세요 : owner/repository", repositoryName)
	}

	return env, nil
}

type environment struct {
	token      string
	owner      string
	repository string
}
