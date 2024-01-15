# PayHere 과제

## 폴더 구조

```
.
├── api             # API 문서와 http requests 폴더
├── cmd             # 애플리케이션 명령어를 구성하기 위한 코드 폴더
├── config          # 설정 파일 폴더
├── domain          # 애플리케이션을 구성하는 도메인 코드 폴더
├── handler         # 요청을 핸들링하기 위한 코드 폴더
├── internal        
├── middleware      # 미들웨어 폴더
├── repository      # 외부 자원들과 통신을 담당하는 repository 폴더
└── usecase         # 비즈니스 로직을 담당하는 usecase 폴더
```

## 빌드

```sh
make build
```

### arm 기반일 경우
```sh
GOARCH=arm64 make build
```

## 실행

```sh
make run
```

### arm 기반일 경우
```sh
GOARCH=arm64 make run
```

## 테스트

```shell
make test
```

## API 문서

서버를 실행하고 [링크](http://localhost:1202/docs) 접속