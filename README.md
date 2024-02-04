# D-n Rule

## 개요

- 리뷰시 PR의 우선순위를 파악하기 위해 사용되는 D-n룰 ([참고](https://blog.banksalad.com/tech/banksalad-code-review-culture/))을 자동화하기 위한 Actions 스크립트입니다.
- 해당 스크립트가 실행되는 경우 D-n 에서 D-(n-1)로 변경됩니다.
  - 예) D-5 인경우 D-4로 변경

## 사용방법

사용하는 레포지터리에 다음과 같이 작성합니다.

```yaml
name: dnrule Test
on:
  schedule:
    - cron: '0 15 * * *'
permissions:
  issues: write
  pull-requests: write
jobs:
  dnrule:
    runs-on: ubuntu-latest
    steps:
      - name: dnrule
        uses: zxcv9203/dnRule@v1.0.1
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          github-repository: ${{ github.repository }}
```

- name :  Actions 스크립트의 이름입니다.
- on.schedule.cron : 해당 스크립트가 실행되는 시간입니다. 한국 시간 기준 00시 00분에 실행시키기 위해서는 위와 같이 작성하면 됩니다.
- permission : Github Token의 권한을 명시합니다. 해당 스크립트의 경우 PR의 라벨을 수정하기 때문에 위와 같은 권한이 필요합니다.
- uses : owner/reponame@tag 를 작성하는 부분입니다. 해당 레포의 스크립트를 사용하기 위해서는 위와 같이 작성해주시면 됩니다.
- github-token : PR을 수정하기 위해 필요한 API Key를 지정합니다. 위와 같이 작성하면 Github에서 자동으로 토큰을 발급해줍니다.
- github-repository : PR을 수정하기 위한 레포지터리를 지정합니다. 위와 같이 작성하면 Github 환경변수를 이용해서 넘겨줄 수 있습니다.

해당 스크립트의 위치는 프로젝트 루트폴더의 `.github/workflows`에 위치시키면 됩니다.
- 예) `.github/workflows/dnrule.yml`


```text
💡 버그 발생시 이슈 등록 부탁드립니다 🙇
```