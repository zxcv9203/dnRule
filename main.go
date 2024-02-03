package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func main() {
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
			// 해당 로직에 PR 라벨 변경 로직 추가
			fmt.Printf("#%v %v\n", *pullRequest.Number, *pullRequest.Title)
		}
		if response.NextPage == 0 {
			break
		}
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
